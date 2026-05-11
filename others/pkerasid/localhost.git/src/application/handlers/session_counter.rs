//! Session-counter demo handler.
//!
//! Called by the pipeline after it has loaded (or created) the session and
//! incremented the `visits` key. Builds an HTML response showing the count.

use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

pub fn handle(visits: u64) -> Response {
    let body = format!(
        "<!DOCTYPE html>\n\
         <html lang=\"en\"><head><meta charset=\"utf-8\">\
         <title>Session Counter</title></head>\n\
         <body><h1>Session Counter</h1>\
         <p>You have visited this page <strong>{visits}</strong> time(s).</p>\
         </body></html>\n"
    )
    .into_bytes();

    let mut builder = Response::builder(Status::OK);
    if let (Ok(n), Ok(v)) = (
        HeaderName::new("content-type"),
        HeaderValue::new("text/html; charset=utf-8"),
    ) {
        builder = builder.header(n, v);
    }
    builder.body(body).build()
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used)]

    use super::*;

    #[test]
    fn status_ok_and_html_body() {
        let r = handle(3);
        assert_eq!(r.status(), Status::OK);
        assert_eq!(r.headers().get("content-type"), Some("text/html; charset=utf-8"));
        let body = String::from_utf8_lossy(r.body());
        assert!(body.contains("3"));
        assert!(body.contains("Session Counter"));
    }

    #[test]
    fn first_visit_shows_one() {
        let r = handle(1);
        let body = String::from_utf8_lossy(r.body());
        assert!(body.contains("<strong>1</strong>"));
    }
}
