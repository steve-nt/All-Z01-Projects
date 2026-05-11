use crate::domain::http::headers::{HeaderName, HeaderValue, Headers};
use crate::domain::http::method::Method;
use crate::domain::http::url::{Query, RequestPath};
use crate::domain::http::version::HttpVersion;

/// Inbound HTTP request, fully validated.
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct Request {
    method: Method,
    path: RequestPath,
    query: Option<Query>,
    version: HttpVersion,
    headers: Headers,
    body: Vec<u8>,
}

impl Request {
    pub fn builder(method: Method, path: RequestPath, version: HttpVersion) -> RequestBuilder {
        RequestBuilder {
            method,
            path,
            query: None,
            version,
            headers: Headers::new(),
            body: Vec::new(),
        }
    }

    pub fn method(&self) -> &Method {
        &self.method
    }

    pub fn path(&self) -> &RequestPath {
        &self.path
    }

    pub fn query(&self) -> Option<&Query> {
        self.query.as_ref()
    }

    pub fn version(&self) -> HttpVersion {
        self.version
    }

    pub fn headers(&self) -> &Headers {
        &self.headers
    }

    pub fn body(&self) -> &[u8] {
        &self.body
    }

    /// Returns the requested host (without port), case-insensitive.
    pub fn host(&self) -> Option<&str> {
        let raw = self.headers.get("host")?;
        Some(raw.split(':').next().unwrap_or(raw))
    }
}

#[derive(Debug)]
pub struct RequestBuilder {
    method: Method,
    path: RequestPath,
    query: Option<Query>,
    version: HttpVersion,
    headers: Headers,
    body: Vec<u8>,
}

impl RequestBuilder {
    #[must_use]
    pub fn query(mut self, q: Query) -> Self {
        self.query = Some(q);
        self
    }

    #[must_use]
    pub fn header(mut self, name: HeaderName, value: HeaderValue) -> Self {
        self.headers.append(name, value);
        self
    }

    #[must_use]
    pub fn headers(mut self, headers: Headers) -> Self {
        self.headers = headers;
        self
    }

    #[must_use]
    pub fn body(mut self, body: Vec<u8>) -> Self {
        self.body = body;
        self
    }

    pub fn build(self) -> Request {
        Request {
            method: self.method,
            path: self.path,
            query: self.query,
            version: self.version,
            headers: self.headers,
            body: self.body,
        }
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn req() -> Request {
        Request::builder(
            Method::Get,
            RequestPath::parse("/index.html").unwrap(),
            HttpVersion::Http11,
        )
        .header(
            HeaderName::new("Host").unwrap(),
            HeaderValue::new("example.com:8080").unwrap(),
        )
        .build()
    }

    #[test]
    fn accessors() {
        let r = req();
        assert_eq!(r.method().as_str(), "GET");
        assert_eq!(r.path().as_str(), "/index.html");
        assert_eq!(r.version(), HttpVersion::Http11);
        assert!(r.body().is_empty());
    }

    #[test]
    fn host_extracts_without_port() {
        let r = req();
        assert_eq!(r.host(), Some("example.com"));
    }

    #[test]
    fn host_missing_returns_none() {
        let r = Request::builder(
            Method::Get,
            RequestPath::parse("/").unwrap(),
            HttpVersion::Http10,
        )
        .build();
        assert_eq!(r.host(), None);
    }
}
