//! Static file and directory handler.
//!
//! Given a resolved filesystem path:
//! - **File** → serves content with `Content-Type`, `Content-Length`,
//!   `Last-Modified`.
//! - **Directory** → tries `index` file, then `autoindex` listing, then 403.
//! - **Not found** → 404.

use std::fmt::Write as _;
use std::path::Path;

use crate::application::ports::filesystem::FileSystem;
use crate::domain::http::date::format_http_date;
use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

/// Serve the filesystem path resolved by the request pipeline.
///
/// `dir_index` is the optional index filename (e.g. `"index.html"`).
/// `autoindex` enables automatic directory listing.
pub fn handle(
    file_path: &Path,
    dir_index: Option<&str>,
    autoindex: bool,
    fs: &dyn FileSystem,
) -> Response {
    match fs.stat(file_path) {
        Err(e) if e.kind() == std::io::ErrorKind::NotFound => not_found(),
        Err(_) => internal_error(),
        Ok(meta) if meta.is_file => serve_file(file_path, meta.modified_secs, fs),
        Ok(meta) if meta.is_dir => serve_dir(file_path, dir_index, autoindex, fs),
        Ok(_) => not_found(),
    }
}

fn serve_file(path: &Path, modified_secs: u64, fs: &dyn FileSystem) -> Response {
    match fs.read_file(path) {
        Ok(body) => {
            let mime = mime_for_path(path);
            let mut builder = Response::builder(Status::OK);
            if let (Ok(n), Ok(v)) = (HeaderName::new("content-type"), HeaderValue::new(mime)) {
                builder = builder.header(n, v);
            }
            if let (Ok(n), Ok(v)) = (
                HeaderName::new("last-modified"),
                HeaderValue::new(format_http_date(modified_secs)),
            ) {
                builder = builder.header(n, v);
            }
            builder.body(body).build()
        }
        Err(_) => internal_error(),
    }
}

fn serve_dir(
    dir_path: &Path,
    dir_index: Option<&str>,
    autoindex: bool,
    fs: &dyn FileSystem,
) -> Response {
    // Try the index file first.
    if let Some(idx) = dir_index {
        let index_path = dir_path.join(idx);
        match fs.stat(&index_path) {
            Ok(meta) if meta.is_file => {
                return serve_file(&index_path, meta.modified_secs, fs);
            }
            _ => {}
        }
    }

    if autoindex {
        return autoindex_listing(dir_path, fs);
    }

    forbidden()
}

fn autoindex_listing(dir_path: &Path, fs: &dyn FileSystem) -> Response {
    let Ok(entries) = fs.read_dir(dir_path) else {
        return internal_error();
    };

    let dir_display = dir_path.display().to_string();
    let mut html = format!(
        "<!DOCTYPE html><html><head><title>Index of {dir}</title></head>\
<body><h1>Index of {dir}</h1><hr><pre>",
        dir = html_escape(&dir_display),
    );

    // Parent link (always useful).
    html.push_str("<a href=\"../\">../</a>\n");

    for entry in &entries {
        let suffix = if entry.is_dir { "/" } else { "" };
        let link = format!("{}{}", html_escape(&entry.name), suffix);
        let display = format!("{}{suffix}", entry.name);
        // String::write_fmt never fails; the error is infallible.
        let _ = writeln!(html, "<a href=\"{link}\">{display}</a>");
    }

    html.push_str("</pre><hr></body></html>");

    let mut builder = Response::builder(Status::OK);
    if let (Ok(n), Ok(v)) = (
        HeaderName::new("content-type"),
        HeaderValue::new("text/html; charset=utf-8"),
    ) {
        builder = builder.header(n, v);
    }
    builder.body(html.into_bytes()).build()
}

/// Minimal HTML escaping for paths/filenames in listings.
fn html_escape(s: &str) -> String {
    let mut out = String::with_capacity(s.len());
    for ch in s.chars() {
        match ch {
            '&' => out.push_str("&amp;"),
            '<' => out.push_str("&lt;"),
            '>' => out.push_str("&gt;"),
            '"' => out.push_str("&quot;"),
            _ => out.push(ch),
        }
    }
    out
}

fn mime_for_path(path: &Path) -> &'static str {
    let ext = path.extension().and_then(|e| e.to_str()).unwrap_or("");
    mime_for_extension(ext)
}

fn mime_for_extension(ext: &str) -> &'static str {
    match ext {
        "html" | "htm" => "text/html; charset=utf-8",
        "css" => "text/css; charset=utf-8",
        "js" | "mjs" => "application/javascript",
        "json" => "application/json",
        "png" => "image/png",
        "jpg" | "jpeg" => "image/jpeg",
        "gif" => "image/gif",
        "svg" => "image/svg+xml",
        "ico" => "image/x-icon",
        "txt" | "text" => "text/plain; charset=utf-8",
        "xml" => "application/xml",
        "pdf" => "application/pdf",
        "zip" => "application/zip",
        "woff" => "font/woff",
        "woff2" => "font/woff2",
        "ttf" => "font/ttf",
        "otf" => "font/otf",
        "webp" => "image/webp",
        "mp4" => "video/mp4",
        "webm" => "video/webm",
        "ogg" => "audio/ogg",
        "mp3" => "audio/mpeg",
        _ => "application/octet-stream",
    }
}

fn not_found() -> Response {
    plain_response(Status::NOT_FOUND)
}

fn forbidden() -> Response {
    plain_response(Status::FORBIDDEN)
}

fn internal_error() -> Response {
    plain_response(Status::INTERNAL_SERVER_ERROR)
}

fn plain_response(status: Status) -> Response {
    let body = format!("{}\n", status.reason()).into_bytes();
    let mut builder = Response::builder(status);
    if let (Ok(n), Ok(v)) = (
        HeaderName::new("content-type"),
        HeaderValue::new("text/plain; charset=utf-8"),
    ) {
        builder = builder.header(n, v);
    }
    builder.body(body).build()
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::path::PathBuf;

    use super::*;
    use crate::application::ports::filesystem::fake::FakeFileSystem;

    #[test]
    fn serves_html_file_with_correct_mime() {
        let fs = FakeFileSystem::new();
        fs.add_file("/www/index.html", b"<h1>Hello</h1>");
        let resp = handle(&PathBuf::from("/www/index.html"), None, false, &fs);
        assert_eq!(resp.status(), Status::OK);
        assert_eq!(
            resp.headers().get("content-type"),
            Some("text/html; charset=utf-8")
        );
        assert_eq!(resp.body(), b"<h1>Hello</h1>");
    }

    #[test]
    fn not_found_when_file_missing() {
        let fs = FakeFileSystem::new();
        let resp = handle(&PathBuf::from("/www/missing.html"), None, false, &fs);
        assert_eq!(resp.status(), Status::NOT_FOUND);
    }

    #[test]
    fn directory_with_index_serves_index() {
        let fs = FakeFileSystem::new();
        fs.add_dir("/www/");
        fs.add_file("/www/index.html", b"<h1>Home</h1>");
        let resp = handle(&PathBuf::from("/www/"), Some("index.html"), false, &fs);
        assert_eq!(resp.status(), Status::OK);
        assert_eq!(resp.body(), b"<h1>Home</h1>");
    }

    #[test]
    fn directory_autoindex_lists_entries() {
        let fs = FakeFileSystem::new();
        fs.add_dir("/www");
        fs.add_file("/www/a.txt", b"hello");
        fs.add_file("/www/b.html", b"world");
        let resp = handle(&PathBuf::from("/www"), None, true, &fs);
        assert_eq!(resp.status(), Status::OK);
        let body = std::str::from_utf8(resp.body()).unwrap();
        assert!(body.contains("a.txt"));
        assert!(body.contains("b.html"));
    }

    #[test]
    fn directory_without_index_or_autoindex_returns_403() {
        let fs = FakeFileSystem::new();
        fs.add_dir("/www");
        let resp = handle(&PathBuf::from("/www"), None, false, &fs);
        assert_eq!(resp.status(), Status::FORBIDDEN);
    }

    #[test]
    fn mime_table_spot_checks() {
        assert_eq!(mime_for_extension("html"), "text/html; charset=utf-8");
        assert_eq!(mime_for_extension("png"), "image/png");
        assert_eq!(mime_for_extension("js"), "application/javascript");
        assert_eq!(mime_for_extension("xyz"), "application/octet-stream");
    }
}
