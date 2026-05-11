use crate::domain::http::headers::{HeaderName, HeaderValue, Headers};
use crate::domain::http::status::Status;
use crate::domain::http::version::HttpVersion;

/// Outbound HTTP response, fully validated.
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct Response {
    version: HttpVersion,
    status: Status,
    headers: Headers,
    body: Vec<u8>,
}

impl Response {
    pub fn builder(status: Status) -> ResponseBuilder {
        ResponseBuilder {
            version: HttpVersion::Http11,
            status,
            headers: Headers::new(),
            body: Vec::new(),
        }
    }

    pub fn version(&self) -> HttpVersion {
        self.version
    }

    pub fn status(&self) -> Status {
        self.status
    }

    pub fn headers(&self) -> &Headers {
        &self.headers
    }

    pub fn body(&self) -> &[u8] {
        &self.body
    }

    /// Append a header to an already-built response (used by middleware layers).
    pub fn add_header(&mut self, name: HeaderName, value: HeaderValue) {
        self.headers.append(name, value);
    }
}

#[derive(Debug)]
pub struct ResponseBuilder {
    version: HttpVersion,
    status: Status,
    headers: Headers,
    body: Vec<u8>,
}

impl ResponseBuilder {
    #[must_use]
    pub fn version(mut self, v: HttpVersion) -> Self {
        self.version = v;
        self
    }

    #[must_use]
    pub fn header(mut self, name: HeaderName, value: HeaderValue) -> Self {
        self.headers.append(name, value);
        self
    }

    #[must_use]
    pub fn set_header(mut self, name: HeaderName, value: HeaderValue) -> Self {
        self.headers.set(name, value);
        self
    }

    #[must_use]
    pub fn body(mut self, body: Vec<u8>) -> Self {
        self.body = body;
        self
    }

    pub fn build(self) -> Response {
        Response {
            version: self.version,
            status: self.status,
            headers: self.headers,
            body: self.body,
        }
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn defaults_to_http_1_1() {
        let r = Response::builder(Status::OK).build();
        assert_eq!(r.version(), HttpVersion::Http11);
        assert_eq!(r.status(), Status::OK);
        assert!(r.body().is_empty());
    }

    #[test]
    fn header_append_vs_set() {
        let r = Response::builder(Status::OK)
            .header(
                HeaderName::new("X").unwrap(),
                HeaderValue::new("1").unwrap(),
            )
            .header(
                HeaderName::new("X").unwrap(),
                HeaderValue::new("2").unwrap(),
            )
            .set_header(
                HeaderName::new("Y").unwrap(),
                HeaderValue::new("z").unwrap(),
            )
            .build();
        let xs: Vec<&str> = r.headers().get_all("x").collect();
        assert_eq!(xs, vec!["1", "2"]);
        assert_eq!(r.headers().get("y"), Some("z"));
    }

    #[test]
    fn body_round_trip() {
        let r = Response::builder(Status::CREATED)
            .body(b"hello".to_vec())
            .build();
        assert_eq!(r.body(), b"hello");
        assert_eq!(r.status(), Status::CREATED);
    }
}
