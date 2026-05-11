//! Upload handler (POST).
//!
//! Supports two modes:
//! - **Multipart** (`Content-Type: multipart/form-data; boundary=…`): parses
//!   the first file part, extracts the filename from `Content-Disposition`, and
//!   writes the part body to `upload_dir/<filename>`.
//! - **Raw**: treats the entire request body as the file content and uses the
//!   last path segment of the request URL as the filename.
//!
//! Returns 201 Created on success (with a `Location` header pointing to the
//! uploaded file), 400 on parse/validation errors, 500 on I/O failure.

use std::path::{Path, PathBuf};

use crate::application::ports::filesystem::FileSystem;
use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::request::Request;
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

pub fn handle(req: &Request, upload_dir: &Path, fs: &dyn FileSystem) -> Response {
    let content_type = req.headers().get("content-type").unwrap_or("").to_owned();

    if let Some(boundary) = multipart_boundary(&content_type) {
        handle_multipart(req.body(), &boundary, upload_dir, fs)
    } else {
        handle_raw(req, upload_dir, fs)
    }
}

// --- multipart ---

fn multipart_boundary(content_type: &str) -> Option<String> {
    // content-type: multipart/form-data; boundary=<token>
    if !content_type
        .to_ascii_lowercase()
        .starts_with("multipart/form-data")
    {
        return None;
    }
    for part in content_type.split(';') {
        let part = part.trim();
        if let Some(b) = part.strip_prefix("boundary=") {
            return Some(b.trim_matches('"').to_owned());
        }
    }
    None
}

fn handle_multipart(
    body: &[u8],
    boundary: &str,
    upload_dir: &Path,
    fs: &dyn FileSystem,
) -> Response {
    let delim = format!("--{boundary}");
    let delim_bytes = delim.as_bytes();

    // Find first part after the opening delimiter.
    let Some(start) = find_subsequence(body, delim_bytes) else {
        return bad_request("missing multipart boundary");
    };
    let after_delim = start + delim_bytes.len();
    // Skip \r\n after delimiter.
    let after_crlf = if body.get(after_delim..after_delim + 2) == Some(b"\r\n") {
        after_delim + 2
    } else {
        after_delim
    };

    // Headers end at \r\n\r\n.
    let Some(headers_end) = find_subsequence(&body[after_crlf..], b"\r\n\r\n") else {
        return bad_request("missing part headers");
    };
    let headers_bytes = &body[after_crlf..after_crlf + headers_end];
    let part_body_start = after_crlf + headers_end + 4; // skip \r\n\r\n

    // End delimiter.
    let end_delim = format!("\r\n--{boundary}");
    let part_body =
        if let Some(end) = find_subsequence(&body[part_body_start..], end_delim.as_bytes()) {
            &body[part_body_start..part_body_start + end]
        } else {
            &body[part_body_start..]
        };

    // Extract filename from Content-Disposition.
    let headers_str = std::str::from_utf8(headers_bytes).unwrap_or("");
    let Some(filename) = extract_filename(headers_str) else {
        return bad_request("missing filename in Content-Disposition");
    };

    write_upload(&sanitize_filename(&filename), part_body, upload_dir, fs)
}

fn extract_filename(headers: &str) -> Option<String> {
    for line in headers.lines() {
        let lower = line.to_ascii_lowercase();
        if lower.starts_with("content-disposition:") {
            for segment in line.split(';') {
                let seg = segment.trim();
                if let Some(val) = seg.strip_prefix("filename=") {
                    return Some(val.trim_matches('"').to_owned());
                }
                // case-insensitive
                if let Some(val) = seg.to_ascii_lowercase().strip_prefix("filename=") {
                    // re-extract original casing
                    let offset = seg.len() - val.len() - "filename=".len();
                    let _ = offset; // suppress warning
                    return Some(seg["filename=".len()..].trim_matches('"').to_owned());
                }
            }
        }
    }
    None
}

// --- raw ---

fn handle_raw(req: &Request, upload_dir: &Path, fs: &dyn FileSystem) -> Response {
    let path_str = req.path().as_str();
    let filename = path_str
        .rsplit('/')
        .find(|s| !s.is_empty())
        .unwrap_or("upload");
    write_upload(filename, req.body(), upload_dir, fs)
}

// --- shared ---

fn write_upload(filename: &str, data: &[u8], upload_dir: &Path, fs: &dyn FileSystem) -> Response {
    if filename.is_empty() {
        return bad_request("empty filename");
    }
    let dest: PathBuf = upload_dir.join(filename);
    match fs.write_file(&dest, data) {
        Ok(()) => {
            let location = format!("/{filename}");
            let mut builder = Response::builder(Status::CREATED);
            if let (Ok(n), Ok(v)) = (HeaderName::new("location"), HeaderValue::new(location)) {
                builder = builder.header(n, v);
            }
            builder.build()
        }
        Err(_) => Response::builder(Status::INTERNAL_SERVER_ERROR).build(),
    }
}

fn bad_request(msg: &str) -> Response {
    let body = format!("{msg}\n").into_bytes();
    let mut builder = Response::builder(Status::BAD_REQUEST);
    if let (Ok(n), Ok(v)) = (
        HeaderName::new("content-type"),
        HeaderValue::new("text/plain; charset=utf-8"),
    ) {
        builder = builder.header(n, v);
    }
    builder.body(body).build()
}

/// Strip path components and dangerous characters from a filename.
fn sanitize_filename(name: &str) -> String {
    let base = name.rsplit(['/', '\\']).next().unwrap_or(name);
    base.chars().filter(|&c| c != '\0' && c != '/').collect()
}

/// Returns the index of the first occurrence of `needle` in `haystack`.
fn find_subsequence(haystack: &[u8], needle: &[u8]) -> Option<usize> {
    if needle.is_empty() {
        return Some(0);
    }
    haystack.windows(needle.len()).position(|w| w == needle)
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;
    use crate::application::ports::filesystem::fake::FakeFileSystem;
    use crate::domain::http::headers::{HeaderName, HeaderValue};
    use crate::domain::http::method::Method;
    use crate::domain::http::request::Request;
    use crate::domain::http::url::RequestPath;
    use crate::domain::http::version::HttpVersion;
    use std::net::SocketAddr;

    fn post_request(path: &str, body: Vec<u8>, content_type: &str) -> Request {
        Request::builder(
            Method::Post,
            RequestPath::parse(path).unwrap(),
            HttpVersion::Http11,
        )
        .header(
            HeaderName::new("content-type").unwrap(),
            HeaderValue::new(content_type).unwrap(),
        )
        .body(body)
        .build()
    }

    #[test]
    fn raw_upload_writes_file() {
        let fs = FakeFileSystem::new();
        fs.add_dir("/uploads");
        let req = post_request("/upload/hello.txt", b"hello world".to_vec(), "text/plain");
        let resp = handle(&req, Path::new("/uploads"), &fs);
        assert_eq!(resp.status(), Status::CREATED);
        assert_eq!(
            fs.read_file(Path::new("/uploads/hello.txt")).unwrap(),
            b"hello world"
        );
    }

    #[test]
    fn multipart_upload_writes_file() {
        let fs = FakeFileSystem::new();
        fs.add_dir("/uploads");
        let body = b"--boundary\r\nContent-Disposition: form-data; name=\"file\"; filename=\"test.txt\"\r\nContent-Type: text/plain\r\n\r\nfile content\r\n--boundary--".to_vec();
        let req = post_request("/upload", body, "multipart/form-data; boundary=boundary");
        let resp = handle(&req, Path::new("/uploads"), &fs);
        assert_eq!(resp.status(), Status::CREATED);
        assert_eq!(
            fs.read_file(Path::new("/uploads/test.txt")).unwrap(),
            b"file content"
        );
    }

    #[test]
    fn sanitize_strips_path_components() {
        assert_eq!(sanitize_filename("../evil.txt"), "evil.txt");
        assert_eq!(sanitize_filename("/abs/path/file.txt"), "file.txt");
        assert_eq!(sanitize_filename("normal.txt"), "normal.txt");
    }

    fn _dummy_addr() -> SocketAddr {
        "0.0.0.0:0".parse().unwrap()
    }
}
