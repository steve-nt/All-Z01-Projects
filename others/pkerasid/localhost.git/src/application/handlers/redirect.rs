//! Redirect handler: emits a 3xx response with a `Location` header.

use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

pub fn handle(location: &str, status: Status) -> Response {
    let mut builder = Response::builder(status);
    if let (Ok(n), Ok(v)) = (HeaderName::new("location"), HeaderValue::new(location)) {
        builder = builder.header(n, v);
    }
    builder.build()
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn sets_location_and_status() {
        let resp = handle("/new/path", Status::MOVED_PERMANENTLY);
        assert_eq!(resp.status(), Status::MOVED_PERMANENTLY);
        assert_eq!(resp.headers().get("location"), Some("/new/path"));
        assert!(resp.body().is_empty());
    }

    #[test]
    fn found_redirect() {
        let resp = handle("https://example.com/", Status::FOUND);
        assert_eq!(resp.status(), Status::FOUND);
        assert_eq!(resp.headers().get("location"), Some("https://example.com/"));
    }
}
