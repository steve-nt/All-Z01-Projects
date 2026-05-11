//! Error page resolution.
//!
//! Tries the host's configured custom error page file; on any I/O failure
//! (file missing, unreadable) falls back to the built-in HTML template.
//! This two-level fallback means a misconfigured error page never causes
//! a 500-loop.

use crate::application::ports::filesystem::FileSystem;
use crate::domain::config::server::HostConfig;
use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

pub fn error_response(
    status: Status,
    host_config: Option<&HostConfig>,
    fs: &dyn FileSystem,
) -> Response {
    if let Some(host) = host_config
        && let Some(path) = host.error_page(status.code())
        && let Ok(body) = fs.read_file(path)
    {
        return build_html_response(status, body);
    }
    build_html_response(status, default_body(status).into_bytes())
}

fn build_html_response(status: Status, body: Vec<u8>) -> Response {
    let mut builder = Response::builder(status);
    if let (Ok(n), Ok(v)) = (
        HeaderName::new("content-type"),
        HeaderValue::new("text/html; charset=utf-8"),
    ) {
        builder = builder.header(n, v);
    }
    builder.body(body).build()
}

fn default_body(status: Status) -> String {
    let code = status.code();
    let reason = status.reason();
    let description = error_description(status);
    format!(
        "<!DOCTYPE html>\n\
         <html lang=\"en\">\n\
         <head>\n\
         <meta charset=\"utf-8\">\n\
         <title>{code} {reason}</title>\n\
         <style>body{{font-family:sans-serif;max-width:600px;margin:4rem auto;color:#333}}\
         h1{{border-bottom:1px solid #ccc;padding-bottom:.5rem}}\
         p{{color:#555}}</style>\n\
         </head>\n\
         <body>\n\
         <h1>{code} {reason}</h1>\n\
         <p>{description}</p>\n\
         </body>\n\
         </html>\n"
    )
}

fn error_description(status: Status) -> &'static str {
    match status.code() {
        400 => "The request could not be understood by the server due to malformed syntax.",
        403 => "You do not have permission to access the requested resource.",
        404 => "The requested resource could not be found on this server.",
        405 => "The HTTP method used is not allowed for the requested resource.",
        408 => "The server timed out waiting for the request.",
        411 => "A Content-Length header is required for this request.",
        413 => "The request body exceeds the maximum allowed size.",
        414 => "The request URI is too long for the server to process.",
        500 => "The server encountered an internal error and could not complete the request.",
        501 => "The server does not support the functionality required to fulfil this request.",
        505 => "The HTTP protocol version used in the request is not supported.",
        _ => "An error occurred while processing your request.",
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::collections::HashMap;
    use std::path::PathBuf;

    use super::*;
    use crate::application::ports::filesystem::fake::FakeFileSystem;
    use crate::domain::config::route::{RouteConfig, RouteKind};
    use crate::domain::config::server::HostConfig;
    use crate::domain::http::method::Method;

    fn host_with_error_page(code: u16, path: PathBuf) -> HostConfig {
        let mut pages = HashMap::new();
        pages.insert(code, path);
        let route = RouteConfig::new(
            "/",
            vec![Method::Get],
            RouteKind::Static {
                root: PathBuf::from("/var/www"),
                index: None,
                autoindex: false,
                upload_dir: None,
            },
            None,
        )
        .unwrap();
        HostConfig::new(
            vec!["0.0.0.0:80".parse().unwrap()],
            vec![],
            1024 * 1024,
            pages,
            vec![route],
        )
        .unwrap()
    }

    #[test]
    fn builtin_template_on_no_config() {
        let fs = FakeFileSystem::new();
        let resp = error_response(Status::NOT_FOUND, None, &fs);
        assert_eq!(resp.status(), Status::NOT_FOUND);
        let body = std::str::from_utf8(resp.body()).unwrap();
        assert!(body.contains("404"));
        assert!(body.contains("Not Found"));
        assert_eq!(
            resp.headers().get("content-type"),
            Some("text/html; charset=utf-8")
        );
    }

    #[test]
    fn custom_page_served_when_present() {
        let fs = FakeFileSystem::new();
        fs.add_file("/pages/404.html", b"<h1>Custom 404</h1>");
        let host = host_with_error_page(404, PathBuf::from("/pages/404.html"));
        let resp = error_response(Status::NOT_FOUND, Some(&host), &fs);
        assert_eq!(resp.body(), b"<h1>Custom 404</h1>");
    }

    #[test]
    fn falls_back_to_builtin_when_custom_missing() {
        let fs = FakeFileSystem::new();
        let host = host_with_error_page(404, PathBuf::from("/pages/missing.html"));
        let resp = error_response(Status::NOT_FOUND, Some(&host), &fs);
        let body = std::str::from_utf8(resp.body()).unwrap();
        assert!(body.contains("404"));
    }

    fn check_default(status: Status) {
        let fs = FakeFileSystem::new();
        let resp = error_response(status, None, &fs);
        assert_eq!(resp.status(), status);
        let body = std::str::from_utf8(resp.body()).unwrap();
        assert!(
            body.contains(&status.code().to_string()),
            "body missing status code for {status}"
        );
        assert!(
            body.contains(status.reason()),
            "body missing reason phrase for {status}"
        );
        assert_eq!(resp.headers().get("content-type"), Some("text/html; charset=utf-8"));
    }

    #[test]
    fn default_page_400() {
        check_default(Status::BAD_REQUEST);
    }

    #[test]
    fn default_page_403() {
        check_default(Status::FORBIDDEN);
    }

    #[test]
    fn default_page_404() {
        check_default(Status::NOT_FOUND);
    }

    #[test]
    fn default_page_405() {
        check_default(Status::METHOD_NOT_ALLOWED);
    }

    #[test]
    fn default_page_408() {
        check_default(Status::REQUEST_TIMEOUT);
    }

    #[test]
    fn default_page_413() {
        check_default(Status::PAYLOAD_TOO_LARGE);
    }

    #[test]
    fn default_page_500() {
        check_default(Status::INTERNAL_SERVER_ERROR);
    }

    #[test]
    fn template_is_valid_html5() {
        let fs = FakeFileSystem::new();
        let resp = error_response(Status::NOT_FOUND, None, &fs);
        let body = std::str::from_utf8(resp.body()).unwrap();
        assert!(body.starts_with("<!DOCTYPE html>"));
        assert!(body.contains("<html"));
        assert!(body.contains("<head>"));
        assert!(body.contains("<body>"));
    }
}
