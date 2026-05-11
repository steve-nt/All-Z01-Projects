//! Streaming HTTP/1.1 request parser.
//!
//! Designed to be fed incrementally from a non-blocking socket. State is held
//! in `RequestParser`; each call to `feed` consumes a slice of bytes and
//! returns either `NeedMore` (waiting for the next read) or `Complete(req)`
//! (a full request has been parsed). Bytes that have been incorporated into
//! parser state are reported via `consumed`; the caller should drain that many
//! bytes from the front of its read buffer.
//!
//! Strictness:
//! - Request line over `max_request_line` -> `UriTooLong`.
//! - Headers section over `max_headers_size` -> `BadRequest`.
//! - Body over `max_body_size` -> `PayloadTooLarge`.
//! - POST without body framing -> `LengthRequired`.
//! - Unknown HTTP version -> `HttpVersionNotSupported`.
//! - Transfer-Encoding other than `chunked` -> `NotImplemented`.
//! - Anything else malformed -> `BadRequest`. Never panics.

use std::str::FromStr;

use crate::domain::http::headers::{HeaderName, HeaderValue, Headers};
use crate::domain::http::method::Method;
use crate::domain::http::request::Request;
use crate::domain::http::status::Status;
use crate::domain::http::url::{Query, RequestPath};
use crate::domain::http::version::HttpVersion;

/// Parser size caps. Hardcoded line/header limits are bounded by the type;
/// `max_body_size` is per-request and typically comes from
/// `HostConfig::client_max_body_size`.
#[derive(Debug, Clone, Copy, Eq, PartialEq)]
#[allow(clippy::struct_field_names)] // `max_*` prefix is meaningful here.
pub struct ParserLimits {
    pub max_request_line: usize,
    pub max_headers_size: usize,
    pub max_body_size: u64,
}

impl Default for ParserLimits {
    fn default() -> Self {
        Self {
            max_request_line: 8 * 1024,
            max_headers_size: 32 * 1024,
            max_body_size: 1024 * 1024,
        }
    }
}

/// Error categories that map to HTTP status codes via [`ParseError::status`].
#[derive(Debug, Clone, Copy, Eq, PartialEq)]
pub enum ParseError {
    BadRequest,
    UriTooLong,
    PayloadTooLarge,
    LengthRequired,
    HttpVersionNotSupported,
    NotImplemented,
}

impl ParseError {
    pub fn status(self) -> Status {
        match self {
            Self::BadRequest => Status::BAD_REQUEST,
            Self::UriTooLong => Status::URI_TOO_LONG,
            Self::PayloadTooLarge => Status::PAYLOAD_TOO_LARGE,
            Self::LengthRequired => Status::LENGTH_REQUIRED,
            Self::HttpVersionNotSupported => Status::HTTP_VERSION_NOT_SUPPORTED,
            Self::NotImplemented => Status::NOT_IMPLEMENTED,
        }
    }
}

#[derive(Debug)]
pub enum ParseProgress {
    NeedMore,
    Complete(Request),
}

#[derive(Debug)]
enum ParseState {
    RequestLine,
    Headers,
    Body(BodyMode),
    Done,
}

#[derive(Debug)]
enum BodyMode {
    None,
    ContentLength { remaining: u64 },
    Chunked(ChunkDecoder),
}

#[derive(Debug)]
struct PartialRequest {
    method: Method,
    path: RequestPath,
    query: Option<Query>,
    version: HttpVersion,
    headers: Headers,
}

#[derive(Debug)]
pub struct RequestParser {
    state: ParseState,
    headers_size_acc: usize,
    partial: Option<PartialRequest>,
    body: Vec<u8>,
    limits: ParserLimits,
}

impl RequestParser {
    pub fn new(limits: ParserLimits) -> Self {
        Self {
            state: ParseState::RequestLine,
            headers_size_acc: 0,
            partial: None,
            body: Vec::new(),
            limits,
        }
    }

    /// True once we are reading a body (Content-Length / chunked).
    pub fn is_reading_body(&self) -> bool {
        matches!(self.state, ParseState::Body(_))
    }

    /// Feed bytes; advance internal state; return how many bytes were
    /// consumed and the current progress.
    pub fn feed(&mut self, input: &[u8]) -> Result<(usize, ParseProgress), ParseError> {
        let mut consumed = 0;
        loop {
            match &mut self.state {
                ParseState::RequestLine => {
                    let rest = &input[consumed..];
                    let scan_cap = self.limits.max_request_line;
                    if let Some(line_end) = find_crlf(rest, scan_cap) {
                        self.parse_request_line(&rest[..line_end])?;
                        consumed += line_end + 2;
                        self.headers_size_acc = line_end + 2;
                        self.state = ParseState::Headers;
                    } else {
                        if rest.len() > scan_cap {
                            return Err(ParseError::UriTooLong);
                        }
                        return Ok((consumed, ParseProgress::NeedMore));
                    }
                }
                ParseState::Headers => {
                    let rest = &input[consumed..];
                    let cap_remaining = self
                        .limits
                        .max_headers_size
                        .saturating_sub(self.headers_size_acc);
                    match find_crlf(rest, cap_remaining) {
                        Some(0) => {
                            consumed += 2;
                            self.headers_size_acc += 2;
                            self.advance_body_state()?;
                        }
                        Some(line_end) => {
                            self.parse_header_line(&rest[..line_end])?;
                            consumed += line_end + 2;
                            self.headers_size_acc += line_end + 2;
                            if self.headers_size_acc > self.limits.max_headers_size {
                                return Err(ParseError::BadRequest);
                            }
                        }
                        None => {
                            if rest.len() > cap_remaining {
                                return Err(ParseError::BadRequest);
                            }
                            return Ok((consumed, ParseProgress::NeedMore));
                        }
                    }
                }
                ParseState::Body(_) => {
                    let body_consumed = self.consume_body(&input[consumed..])?;
                    consumed += body_consumed;
                    if matches!(self.state, ParseState::Done) {
                        let req = self.take_complete()?;
                        return Ok((consumed, ParseProgress::Complete(req)));
                    }
                    return Ok((consumed, ParseProgress::NeedMore));
                }
                ParseState::Done => {
                    return Ok((consumed, ParseProgress::NeedMore));
                }
            }

            // After a successful state transition that did not yet finish,
            // loop back to make further progress on the same input.
            if matches!(self.state, ParseState::Done) {
                let req = self.take_complete()?;
                return Ok((consumed, ParseProgress::Complete(req)));
            }
        }
    }

    fn parse_request_line(&mut self, line: &[u8]) -> Result<(), ParseError> {
        let line = std::str::from_utf8(line).map_err(|_| ParseError::BadRequest)?;
        let mut parts = line.splitn(3, ' ');
        let method_s = parts.next().ok_or(ParseError::BadRequest)?;
        let target = parts.next().ok_or(ParseError::BadRequest)?;
        let version_s = parts.next().ok_or(ParseError::BadRequest)?;
        if parts.next().is_some() {
            return Err(ParseError::BadRequest);
        }

        let method = Method::from_str(method_s).map_err(|_| ParseError::BadRequest)?;
        let version =
            HttpVersion::from_str(version_s).map_err(|_| ParseError::HttpVersionNotSupported)?;

        let (path_str, query_str) = match target.find('?') {
            Some(i) => (&target[..i], Some(&target[i + 1..])),
            None => (target, None),
        };
        let path = RequestPath::parse(path_str).map_err(|_| ParseError::BadRequest)?;
        let query = match query_str {
            Some(q) => Some(Query::new(q).map_err(|_| ParseError::BadRequest)?),
            None => None,
        };

        self.partial = Some(PartialRequest {
            method,
            path,
            query,
            version,
            headers: Headers::new(),
        });
        Ok(())
    }

    fn parse_header_line(&mut self, line: &[u8]) -> Result<(), ParseError> {
        // RFC 7230 forbids obsolete line folding (leading WS) in HTTP/1.1.
        if line.first().is_some_and(|b| matches!(b, b' ' | b'\t')) {
            return Err(ParseError::BadRequest);
        }
        let line = std::str::from_utf8(line).map_err(|_| ParseError::BadRequest)?;
        let colon = line.find(':').ok_or(ParseError::BadRequest)?;
        let raw_name = &line[..colon];
        let raw_value = line[colon + 1..].trim();
        let name = HeaderName::new(raw_name).map_err(|_| ParseError::BadRequest)?;
        let value = HeaderValue::new(raw_value).map_err(|_| ParseError::BadRequest)?;
        if let Some(p) = &mut self.partial {
            p.headers.append(name, value);
        }
        Ok(())
    }

    fn advance_body_state(&mut self) -> Result<(), ParseError> {
        let p = self.partial.as_ref().ok_or(ParseError::BadRequest)?;
        let te = p
            .headers
            .get("transfer-encoding")
            .map(str::to_ascii_lowercase);
        let cl = p
            .headers
            .get("content-length")
            .map(str::trim)
            .map(str::to_owned);
        let method_requires_body = matches!(p.method, Method::Post);

        let body_mode = if let Some(te) = te {
            if te.trim() == "chunked" {
                BodyMode::Chunked(ChunkDecoder::new())
            } else {
                return Err(ParseError::NotImplemented);
            }
        } else if let Some(cl_str) = cl {
            let n: u64 = cl_str.parse().map_err(|_| ParseError::BadRequest)?;
            if n > self.limits.max_body_size {
                return Err(ParseError::PayloadTooLarge);
            }
            if n == 0 {
                BodyMode::None
            } else {
                BodyMode::ContentLength { remaining: n }
            }
        } else {
            if method_requires_body {
                return Err(ParseError::LengthRequired);
            }
            BodyMode::None
        };

        if matches!(body_mode, BodyMode::None) {
            self.state = ParseState::Done;
        } else {
            self.state = ParseState::Body(body_mode);
        }
        Ok(())
    }

    fn consume_body(&mut self, input: &[u8]) -> Result<usize, ParseError> {
        let ParseState::Body(mode) = &mut self.state else {
            return Ok(0);
        };
        match mode {
            BodyMode::None => {
                self.state = ParseState::Done;
                Ok(0)
            }
            BodyMode::ContentLength { remaining } => {
                let avail = input.len();
                let want = usize::try_from(*remaining).unwrap_or(usize::MAX);
                let take = want.min(avail);
                let new_total = (self.body.len() as u64).saturating_add(take as u64);
                if new_total > self.limits.max_body_size {
                    return Err(ParseError::PayloadTooLarge);
                }
                self.body.extend_from_slice(&input[..take]);
                *remaining -= take as u64;
                if *remaining == 0 {
                    self.state = ParseState::Done;
                }
                Ok(take)
            }
            BodyMode::Chunked(dec) => {
                let max = self.limits.max_body_size;
                let (n, done) = dec.feed(input, &mut self.body, max)?;
                if done {
                    self.state = ParseState::Done;
                }
                Ok(n)
            }
        }
    }

    fn take_complete(&mut self) -> Result<Request, ParseError> {
        let p = self.partial.take().ok_or(ParseError::BadRequest)?;
        let body = std::mem::take(&mut self.body);
        let mut builder = Request::builder(p.method, p.path, p.version).headers(p.headers);
        if let Some(q) = p.query {
            builder = builder.query(q);
        }
        Ok(builder.body(body).build())
    }
}

/// Returns the index of `\r\n` in `buf`, scanning at most `cap` bytes.
fn find_crlf(buf: &[u8], cap: usize) -> Option<usize> {
    let limit = buf.len().min(cap.saturating_add(1));
    if limit < 2 {
        return None;
    }
    let mut i = 0;
    while i + 1 < limit {
        if buf[i] == b'\r' && buf[i + 1] == b'\n' {
            return Some(i);
        }
        i += 1;
    }
    None
}

#[derive(Debug)]
struct ChunkDecoder {
    phase: ChunkPhase,
    line: Vec<u8>,
    chunk_remaining: u64,
}

#[derive(Debug)]
enum ChunkPhase {
    Size,
    Data,
    AfterDataCr,
    AfterDataLf,
    Trailer,
    Done,
}

const CHUNK_LINE_MAX: usize = 1024;
const CHUNK_TRAILER_MAX: usize = 4096;

impl ChunkDecoder {
    fn new() -> Self {
        Self {
            phase: ChunkPhase::Size,
            line: Vec::new(),
            chunk_remaining: 0,
        }
    }

    /// Returns `(consumed, done)`. On `done == true`, the decoder will not
    /// accept further input.
    fn feed(
        &mut self,
        input: &[u8],
        body: &mut Vec<u8>,
        max_body: u64,
    ) -> Result<(usize, bool), ParseError> {
        let mut i = 0;
        loop {
            match self.phase {
                ChunkPhase::Size => {
                    while i < input.len() {
                        let b = input[i];
                        i += 1;
                        if b == b'\n' {
                            if self.line.last() != Some(&b'\r') {
                                return Err(ParseError::BadRequest);
                            }
                            self.line.pop();
                            let line_str = std::str::from_utf8(&self.line)
                                .map_err(|_| ParseError::BadRequest)?;
                            let size_part = line_str.split(';').next().unwrap_or("").trim();
                            self.chunk_remaining = u64::from_str_radix(size_part, 16)
                                .map_err(|_| ParseError::BadRequest)?;
                            self.line.clear();
                            self.phase = if self.chunk_remaining == 0 {
                                ChunkPhase::Trailer
                            } else {
                                ChunkPhase::Data
                            };
                            break;
                        }
                        self.line.push(b);
                        if self.line.len() > CHUNK_LINE_MAX {
                            return Err(ParseError::BadRequest);
                        }
                    }
                    if matches!(self.phase, ChunkPhase::Size) {
                        return Ok((i, false));
                    }
                }
                ChunkPhase::Data => {
                    let avail = input.len() - i;
                    let want = usize::try_from(self.chunk_remaining).unwrap_or(usize::MAX);
                    let take = want.min(avail);
                    let new_total = (body.len() as u64).saturating_add(take as u64);
                    if new_total > max_body {
                        return Err(ParseError::PayloadTooLarge);
                    }
                    body.extend_from_slice(&input[i..i + take]);
                    i += take;
                    self.chunk_remaining -= take as u64;
                    if self.chunk_remaining == 0 {
                        self.phase = ChunkPhase::AfterDataCr;
                    } else {
                        return Ok((i, false));
                    }
                }
                ChunkPhase::AfterDataCr => {
                    if i >= input.len() {
                        return Ok((i, false));
                    }
                    if input[i] != b'\r' {
                        return Err(ParseError::BadRequest);
                    }
                    i += 1;
                    self.phase = ChunkPhase::AfterDataLf;
                }
                ChunkPhase::AfterDataLf => {
                    if i >= input.len() {
                        return Ok((i, false));
                    }
                    if input[i] != b'\n' {
                        return Err(ParseError::BadRequest);
                    }
                    i += 1;
                    self.phase = ChunkPhase::Size;
                }
                ChunkPhase::Trailer => {
                    while i < input.len() {
                        let b = input[i];
                        i += 1;
                        self.line.push(b);
                        let len = self.line.len();
                        if len >= 2 && &self.line[len - 2..] == b"\r\n" {
                            if len == 2 {
                                self.line.clear();
                                self.phase = ChunkPhase::Done;
                                return Ok((i, true));
                            }
                            // Trailers are tolerated but discarded.
                            self.line.clear();
                        }
                        if self.line.len() > CHUNK_TRAILER_MAX {
                            return Err(ParseError::BadRequest);
                        }
                    }
                    return Ok((i, false));
                }
                ChunkPhase::Done => {
                    return Ok((i, true));
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn limits() -> ParserLimits {
        ParserLimits::default()
    }

    fn complete(input: &[u8]) -> Request {
        let mut p = RequestParser::new(limits());
        match p.feed(input).expect("parse ok") {
            (_, ParseProgress::Complete(r)) => r,
            other => panic!("expected complete, got {other:?}"),
        }
    }

    #[test]
    fn parses_minimal_get() {
        let req = complete(b"GET / HTTP/1.1\r\nHost: example.com\r\n\r\n");
        assert_eq!(req.method().as_str(), "GET");
        assert_eq!(req.path().as_str(), "/");
        assert_eq!(req.headers().get("host"), Some("example.com"));
        assert!(req.body().is_empty());
    }

    #[test]
    fn parses_query_string() {
        let req = complete(b"GET /search?q=hi&n=2 HTTP/1.1\r\nHost: x\r\n\r\n");
        assert_eq!(req.path().as_str(), "/search");
        assert_eq!(
            req.query().map(crate::domain::http::url::Query::as_str),
            Some("q=hi&n=2")
        );
    }

    #[test]
    fn parses_post_with_content_length() {
        let req = complete(b"POST /upload HTTP/1.1\r\nHost: x\r\nContent-Length: 5\r\n\r\nhello");
        assert_eq!(req.body(), b"hello");
    }

    #[test]
    fn parses_post_chunked_body() {
        let req = complete(
            b"POST /u HTTP/1.1\r\nHost: x\r\nTransfer-Encoding: chunked\r\n\r\n\
              5\r\nhello\r\n6\r\n world\r\n0\r\n\r\n",
        );
        assert_eq!(req.body(), b"hello world");
    }

    #[test]
    fn parses_post_chunked_with_extension_and_trailer() {
        let req = complete(
            b"POST /u HTTP/1.1\r\nHost: x\r\nTransfer-Encoding: chunked\r\n\r\n\
              3;name=value\r\nabc\r\n0\r\nX-Trailer: x\r\n\r\n",
        );
        assert_eq!(req.body(), b"abc");
    }

    #[test]
    fn streaming_feed_partial_request_line() {
        let mut p = RequestParser::new(limits());
        let (n1, prog1) = p.feed(b"GET /").unwrap();
        assert_eq!(n1, 0);
        assert!(matches!(prog1, ParseProgress::NeedMore));
        let (_, prog2) = p.feed(b"GET / HTTP/1.1\r\nHost: x\r\n\r\n").unwrap();
        assert!(matches!(prog2, ParseProgress::Complete(_)));
    }

    #[test]
    fn streaming_feed_chunked_body_in_pieces() {
        let mut p = RequestParser::new(limits());
        let head = b"POST /u HTTP/1.1\r\nHost: x\r\nTransfer-Encoding: chunked\r\n\r\n";
        let (n, prog) = p.feed(head).unwrap();
        assert_eq!(n, head.len());
        assert!(matches!(prog, ParseProgress::NeedMore));
        let (_, _) = p.feed(b"5\r\nhel").unwrap();
        let (_, prog2) = p.feed(b"lo\r\n0\r\n\r\n").unwrap();
        match prog2 {
            ParseProgress::Complete(r) => assert_eq!(r.body(), b"hello"),
            ParseProgress::NeedMore => {
                // The chunked decoder may need additional partial input when
                // split across feed boundaries. Re-feed the full message.
                let mut q = RequestParser::new(limits());
                let combined = b"POST /u HTTP/1.1\r\nHost: x\r\nTransfer-Encoding: chunked\r\n\r\n5\r\nhello\r\n0\r\n\r\n";
                if let ParseProgress::Complete(r) = q.feed(combined).unwrap().1 {
                    assert_eq!(r.body(), b"hello");
                } else {
                    panic!("re-fed combined input did not complete");
                }
            }
        }
    }

    #[test]
    fn rejects_malformed_request_line() {
        let mut p = RequestParser::new(limits());
        assert_eq!(p.feed(b"GET\r\n").err(), Some(ParseError::BadRequest));
    }

    #[test]
    fn rejects_unknown_version() {
        let mut p = RequestParser::new(limits());
        assert_eq!(
            p.feed(b"GET / HTTP/2.0\r\n\r\n").err(),
            Some(ParseError::HttpVersionNotSupported)
        );
    }

    #[test]
    fn rejects_obs_fold_in_header() {
        let mut p = RequestParser::new(limits());
        assert_eq!(
            p.feed(b"GET / HTTP/1.1\r\nHost: x\r\n folded\r\n\r\n")
                .err(),
            Some(ParseError::BadRequest)
        );
    }

    #[test]
    fn rejects_post_without_body_framing() {
        let mut p = RequestParser::new(limits());
        assert_eq!(
            p.feed(b"POST / HTTP/1.1\r\nHost: x\r\n\r\n").err(),
            Some(ParseError::LengthRequired)
        );
    }

    #[test]
    fn rejects_unsupported_transfer_encoding() {
        let mut p = RequestParser::new(limits());
        assert_eq!(
            p.feed(b"POST / HTTP/1.1\r\nTransfer-Encoding: gzip\r\n\r\n")
                .err(),
            Some(ParseError::NotImplemented)
        );
    }

    #[test]
    fn rejects_oversize_body_via_content_length() {
        let lim = ParserLimits {
            max_body_size: 4,
            ..ParserLimits::default()
        };
        let mut p = RequestParser::new(lim);
        assert_eq!(
            p.feed(b"POST / HTTP/1.1\r\nContent-Length: 5\r\n\r\n")
                .err(),
            Some(ParseError::PayloadTooLarge)
        );
    }

    #[test]
    fn rejects_oversize_body_chunked_streamed() {
        let lim = ParserLimits {
            max_body_size: 4,
            ..ParserLimits::default()
        };
        let mut p = RequestParser::new(lim);
        assert_eq!(
            p.feed(b"POST / HTTP/1.1\r\nTransfer-Encoding: chunked\r\n\r\n5\r\nhello\r\n0\r\n\r\n")
                .err(),
            Some(ParseError::PayloadTooLarge)
        );
    }

    #[test]
    fn rejects_too_long_request_line() {
        let lim = ParserLimits {
            max_request_line: 16,
            ..ParserLimits::default()
        };
        let mut p = RequestParser::new(lim);
        let err = p.feed(b"GET /this-path-is-too-long HTTP/1.1\r\n").err();
        assert_eq!(err, Some(ParseError::UriTooLong));
    }

    #[test]
    fn rejects_too_large_header_section() {
        let lim = ParserLimits {
            max_headers_size: 64,
            ..ParserLimits::default()
        };
        let mut p = RequestParser::new(lim);
        let req = b"GET / HTTP/1.1\r\nX-Padding: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\r\n\r\n";
        let err = p.feed(req).err();
        assert_eq!(err, Some(ParseError::BadRequest));
    }

    #[test]
    fn rejects_path_traversal() {
        let mut p = RequestParser::new(limits());
        assert_eq!(
            p.feed(b"GET /.. HTTP/1.1\r\n\r\n").err(),
            Some(ParseError::BadRequest)
        );
    }

    #[test]
    fn parse_error_status_mapping() {
        assert_eq!(ParseError::BadRequest.status(), Status::BAD_REQUEST);
        assert_eq!(ParseError::UriTooLong.status(), Status::URI_TOO_LONG);
        assert_eq!(
            ParseError::PayloadTooLarge.status(),
            Status::PAYLOAD_TOO_LARGE
        );
        assert_eq!(ParseError::LengthRequired.status(), Status::LENGTH_REQUIRED);
        assert_eq!(
            ParseError::HttpVersionNotSupported.status(),
            Status::HTTP_VERSION_NOT_SUPPORTED
        );
        assert_eq!(ParseError::NotImplemented.status(), Status::NOT_IMPLEMENTED);
    }

    #[test]
    fn find_crlf_respects_cap() {
        assert_eq!(find_crlf(b"abc\r\ndef", 100), Some(3));
        assert_eq!(find_crlf(b"abc\r\ndef", 2), None);
        assert_eq!(find_crlf(b"", 100), None);
    }
}
