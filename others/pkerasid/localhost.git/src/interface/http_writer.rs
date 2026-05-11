//! Response serializer (`Response` -> bytes).
//!
//! Adds the standard framing headers when not already set:
//! - `Date` (RFC 7231 IMF-fixdate, computed locally without `chrono`).
//! - `Server: localhost/<VERSION>`.
//! - `Content-Length` when no `Transfer-Encoding` is present.
//! - `Connection: keep-alive | close` based on the caller's decision.
//!
//! The writer never panics on a well-formed `Response` and produces a
//! deterministic byte order: status line, framing headers, user headers
//! (excluding ones we override), blank line, body.

use std::io::Write as _;
use std::time::{SystemTime, UNIX_EPOCH};

use crate::VERSION;
use crate::domain::http::date::format_http_date;
use crate::domain::http::response::Response;

/// Serialize `resp` to wire bytes, supplying default framing headers.
///
/// `keep_alive` is the server-side decision; client requests via
/// `Connection: keep-alive` are subject to this override (server may downgrade
/// to close at any time).
pub fn serialize(resp: &Response, keep_alive: bool) -> Vec<u8> {
    serialize_at(resp, keep_alive, current_unix_secs())
}

/// Test-friendly variant that takes a fixed unix timestamp for `Date`.
pub fn serialize_at(resp: &Response, keep_alive: bool, unix_secs: u64) -> Vec<u8> {
    let mut buf: Vec<u8> = Vec::with_capacity(256 + resp.body().len());
    write_status_line(&mut buf, resp);

    let mut have_date = false;
    let mut have_server = false;
    let mut have_cl = false;
    let mut have_te = false;

    for (name, value) in resp.headers().iter() {
        match name.as_str() {
            "date" => have_date = true,
            "server" => have_server = true,
            "content-length" => have_cl = true,
            "transfer-encoding" => have_te = true,
            "connection" => continue, // server controls this
            _ => {}
        }
        write_header(&mut buf, name.as_str(), value.as_str());
    }

    if !have_date {
        write_header(&mut buf, "date", &format_http_date(unix_secs));
    }
    if !have_server {
        let server = format!("localhost/{VERSION}");
        write_header(&mut buf, "server", &server);
    }
    if !have_cl && !have_te {
        let cl = resp.body().len().to_string();
        write_header(&mut buf, "content-length", &cl);
    }
    write_header(
        &mut buf,
        "connection",
        if keep_alive { "keep-alive" } else { "close" },
    );

    buf.extend_from_slice(b"\r\n");
    buf.extend_from_slice(resp.body());
    buf
}

fn write_status_line(buf: &mut Vec<u8>, resp: &Response) {
    let _ = write!(
        buf,
        "{} {} {}\r\n",
        resp.version().as_str(),
        resp.status().code(),
        resp.status().reason()
    );
}

fn write_header(buf: &mut Vec<u8>, name: &str, value: &str) {
    buf.extend_from_slice(name.as_bytes());
    buf.extend_from_slice(b": ");
    buf.extend_from_slice(value.as_bytes());
    buf.extend_from_slice(b"\r\n");
}

fn current_unix_secs() -> u64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_secs())
        .unwrap_or(0)
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;
    use crate::domain::http::headers::{HeaderName, HeaderValue};
    use crate::domain::http::status::Status;
    use crate::domain::http::version::HttpVersion;

    fn n(s: &str) -> HeaderName {
        HeaderName::new(s).unwrap()
    }
    fn v(s: &str) -> HeaderValue {
        HeaderValue::new(s).unwrap()
    }

    #[test]
    fn serialize_minimal_200() {
        let resp = Response::builder(Status::OK).body(b"ok\n".to_vec()).build();
        let bytes = serialize_at(&resp, true, 1_704_067_200);
        let s = std::str::from_utf8(&bytes).unwrap();
        assert!(s.starts_with("HTTP/1.1 200 OK\r\n"));
        assert!(s.contains("date: Mon, 01 Jan 2024 00:00:00 GMT\r\n"));
        assert!(s.contains("server: localhost/"));
        assert!(s.contains("content-length: 3\r\n"));
        assert!(s.contains("connection: keep-alive\r\n"));
        assert!(s.ends_with("\r\nok\n"));
    }

    #[test]
    fn caller_supplied_date_and_server_are_kept() {
        let resp = Response::builder(Status::OK)
            .header(n("Date"), v("Wed, 21 Oct 2015 07:28:00 GMT"))
            .header(n("Server"), v("custom/1.0"))
            .body(b"hi".to_vec())
            .build();
        let s = String::from_utf8(serialize_at(&resp, false, 0)).unwrap();
        assert!(s.contains("Wed, 21 Oct 2015 07:28:00 GMT"));
        assert!(s.contains("custom/1.0"));
        assert!(!s.contains("localhost/"));
    }

    #[test]
    fn server_overrides_caller_connection_header() {
        let resp = Response::builder(Status::OK)
            .header(n("Connection"), v("keep-alive"))
            .body(b"x".to_vec())
            .build();
        let s = String::from_utf8(serialize_at(&resp, false, 0)).unwrap();
        // Should appear exactly once and reflect the server's decision.
        assert_eq!(s.matches("connection:").count(), 1);
        assert!(s.contains("connection: close\r\n"));
    }

    #[test]
    fn transfer_encoding_suppresses_content_length() {
        let resp = Response::builder(Status::OK)
            .header(n("Transfer-Encoding"), v("chunked"))
            .body(b"".to_vec())
            .build();
        let s = String::from_utf8(serialize_at(&resp, true, 0)).unwrap();
        assert!(s.contains("transfer-encoding: chunked\r\n"));
        assert!(!s.contains("content-length:"));
    }

    #[test]
    fn empty_body_for_204() {
        let resp = Response::builder(Status::NO_CONTENT).build();
        let s = String::from_utf8(serialize_at(&resp, true, 0)).unwrap();
        assert!(s.starts_with("HTTP/1.1 204 No Content\r\n"));
        assert!(s.contains("content-length: 0\r\n"));
    }

    #[test]
    fn status_line_uses_response_version() {
        let resp = Response::builder(Status::OK)
            .version(HttpVersion::Http10)
            .build();
        let s = String::from_utf8(serialize_at(&resp, false, 0)).unwrap();
        assert!(s.starts_with("HTTP/1.0 200 OK\r\n"));
    }
}
