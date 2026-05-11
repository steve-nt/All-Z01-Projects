use std::net::SocketAddr;
use std::path::{Path, PathBuf};
use std::time::Duration;

use crate::application::ports::process_runner::{CgiRunSpec, ProcessRunner};
use crate::domain::config::cgi::CgiConfig;
use crate::domain::config::server::HostConfig;
use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::request::Request;
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

const CGI_TIMEOUT: Duration = Duration::from_secs(5);

pub fn handle(
    req: &Request,
    host: &HostConfig,
    local_addr: SocketAddr,
    script_filename: &Path,
    cgi: &CgiConfig,
    runner: &dyn ProcessRunner,
) -> Response {
    let spec = CgiRunSpec {
        interpreter: cgi.interpreter().to_path_buf(),
        script_filename: script_filename.to_path_buf(),
        cwd: script_filename
            .parent()
            .map_or_else(|| PathBuf::from("."), Path::to_path_buf),
        env: cgi_env(req, host, local_addr, script_filename),
        stdin: req.body().to_vec(),
    };

    match runner.run_cgi(&spec, CGI_TIMEOUT) {
        Ok(out) => {
            if !matches!(out.exit_code, Some(0)) {
                return Response::builder(Status::BAD_GATEWAY)
                    .body(out.stderr)
                    .build();
            }
            parse_cgi_output(&out.stdout)
        }
        Err(_) => Response::builder(Status::GATEWAY_TIMEOUT)
            .body(b"Gateway Timeout\n".to_vec())
            .build(),
    }
}

fn cgi_env(
    req: &Request,
    host: &HostConfig,
    local_addr: SocketAddr,
    script_filename: &Path,
) -> Vec<(String, String)> {
    let mut env = vec![
        ("REQUEST_METHOD".into(), req.method().as_str().into()),
        ("SERVER_PROTOCOL".into(), req.version().as_str().into()),
        (
            "SCRIPT_FILENAME".into(),
            script_filename.to_string_lossy().into_owned(),
        ),
        ("SCRIPT_NAME".into(), req.path().as_str().into()),
        ("PATH_INFO".into(), req.path().as_str().into()),
        (
            "QUERY_STRING".into(),
            req.query()
                .map_or_else(String::new, |query| query.as_str().to_owned()),
        ),
        ("SERVER_PORT".into(), local_addr.port().to_string()),
    ];

    if let Some(name) = req
        .host()
        .or_else(|| host.server_names().first().map(String::as_str))
    {
        env.push(("SERVER_NAME".into(), name.to_owned()));
    }
    if let Some(content_type) = req.headers().get("content-type") {
        env.push(("CONTENT_TYPE".into(), content_type.to_owned()));
    }
    if !req.body().is_empty() {
        env.push(("CONTENT_LENGTH".into(), req.body().len().to_string()));
    }

    for (name, value) in req.headers().iter() {
        let key = format!(
            "HTTP_{}",
            name.as_str().replace('-', "_").to_ascii_uppercase()
        );
        env.push((key, value.as_str().to_owned()));
    }
    env
}

fn parse_cgi_output(stdout: &[u8]) -> Response {
    let Some((raw_headers, body)) = split_headers_body(stdout) else {
        return Response::builder(Status::BAD_GATEWAY)
            .body(b"Bad Gateway\n".to_vec())
            .build();
    };

    let mut status = Status::OK;
    let mut headers = Vec::new();
    for line in raw_headers.lines() {
        let line = line.trim_end_matches('\r');
        if line.is_empty() {
            continue;
        }
        if let Some(rest) = line.strip_prefix("Status:") {
            let code = rest
                .trim()
                .split_once(' ')
                .map_or(rest.trim(), |(c, _)| c)
                .parse::<u16>();
            if let Ok(code) = code
                && let Ok(parsed) = Status::new(code)
            {
                status = parsed;
            }
            continue;
        }
        if let Some((name, value)) = line.split_once(':')
            && let (Ok(name), Ok(value)) = (HeaderName::new(name.trim()), HeaderValue::new(value))
        {
            headers.push((name, value));
        }
    }

    let mut builder = Response::builder(status);
    for (name, value) in headers {
        builder = builder.header(name, value);
    }
    builder.body(body.to_vec()).build()
}

fn split_headers_body(bytes: &[u8]) -> Option<(&str, &[u8])> {
    if let Some(idx) = bytes.windows(4).position(|w| w == b"\r\n\r\n") {
        let head = std::str::from_utf8(&bytes[..idx]).ok()?;
        let body = &bytes[idx + 4..];
        return Some((head, body));
    }
    if let Some(idx) = bytes.windows(2).position(|w| w == b"\n\n") {
        let head = std::str::from_utf8(&bytes[..idx]).ok()?;
        let body = &bytes[idx + 2..];
        return Some((head, body));
    }
    None
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::collections::HashMap;
    use std::fs;

    use super::*;
    use crate::application::ports::process_runner::ProcessOutput;
    use crate::domain::config::route::{RouteConfig, RouteKind};
    use crate::domain::http::method::Method;
    use crate::domain::http::url::{Query, RequestPath};
    use crate::domain::http::version::HttpVersion;
    use crate::infrastructure::cgi::OsProcessRunner;

    #[derive(Debug)]
    struct FakeRunner {
        out: Result<ProcessOutput, std::io::ErrorKind>,
    }

    impl ProcessRunner for FakeRunner {
        fn run_cgi(
            &self,
            _spec: &CgiRunSpec,
            _timeout: Duration,
        ) -> std::io::Result<ProcessOutput> {
            self.out.clone().map_err(std::io::Error::from)
        }
    }

    fn host() -> HostConfig {
        let route = RouteConfig::new(
            "/",
            vec![Method::Get],
            RouteKind::Static {
                root: "/tmp".into(),
                index: None,
                autoindex: false,
                upload_dir: None,
            },
            None,
        )
        .unwrap();
        HostConfig::new(
            vec!["127.0.0.1:8080".parse().unwrap()],
            vec!["example.test".into()],
            1024 * 1024,
            HashMap::new(),
            vec![route],
        )
        .unwrap()
    }

    #[test]
    fn parses_status_and_headers() {
        let req = Request::builder(
            Method::Get,
            RequestPath::parse("/cgi/hello.py").unwrap(),
            HttpVersion::Http11,
        )
        .query(Query::new("x=1").unwrap())
        .build();
        let cgi = CgiConfig::new(".py", "/usr/bin/python3").unwrap();
        let runner = FakeRunner {
            out: Ok(ProcessOutput {
                stdout: b"Status: 201 Created\r\nContent-Type: text/plain\r\n\r\nok".to_vec(),
                stderr: Vec::new(),
                exit_code: Some(0),
            }),
        };
        let resp = handle(
            &req,
            &host(),
            "127.0.0.1:8080".parse().unwrap(),
            Path::new("/www/cgi/hello.py"),
            &cgi,
            &runner,
        );
        assert_eq!(resp.status(), Status::CREATED);
        assert_eq!(resp.headers().get("content-type"), Some("text/plain"));
        assert_eq!(resp.body(), b"ok");
    }

    #[test]
    fn timeout_maps_to_504() {
        let req = Request::builder(
            Method::Get,
            RequestPath::parse("/cgi/hello.py").unwrap(),
            HttpVersion::Http11,
        )
        .build();
        let cgi = CgiConfig::new(".py", "/usr/bin/python3").unwrap();
        let runner = FakeRunner {
            out: Err(std::io::ErrorKind::TimedOut),
        };
        let resp = handle(
            &req,
            &host(),
            "127.0.0.1:8080".parse().unwrap(),
            Path::new("/www/cgi/hello.py"),
            &cgi,
            &runner,
        );
        assert_eq!(resp.status(), Status::GATEWAY_TIMEOUT);
    }

    #[cfg(target_os = "linux")]
    #[test]
    fn os_runner_executes_python_script() {
        let Some(interpreter) = ["/usr/bin/python3", "/usr/bin/python"]
            .into_iter()
            .find(|p| Path::new(p).exists())
        else {
            return;
        };
        let temp = tempfile::tempdir().unwrap();
        let script = temp.path().join("echo.py");
        let source = r#"#!/usr/bin/env python3
import sys
body = sys.stdin.read()
print("Status: 201 Created")
print("Content-Type: text/plain")
print()
print("cgi:" + body)
"#;
        fs::write(&script, source).unwrap();
        let req = Request::builder(
            Method::Post,
            RequestPath::parse("/cgi/echo.py").unwrap(),
            HttpVersion::Http11,
        )
        .body(b"hello".to_vec())
        .build();
        let cgi = CgiConfig::new(".py", interpreter).unwrap();
        let runner = OsProcessRunner::new();
        let resp = handle(
            &req,
            &host(),
            "127.0.0.1:8080".parse().unwrap(),
            &script,
            &cgi,
            &runner,
        );
        assert_eq!(resp.status(), Status::CREATED);
        assert_eq!(resp.headers().get("content-type"), Some("text/plain"));
        assert_eq!(resp.body(), b"cgi:hello\n");
    }

    #[test]
    fn nonzero_exit_maps_to_502_with_stderr() {
        let req = Request::builder(
            Method::Get,
            RequestPath::parse("/cgi/bad.py").unwrap(),
            HttpVersion::Http11,
        )
        .build();
        let cgi = CgiConfig::new(".py", "/usr/bin/python3").unwrap();
        let runner = FakeRunner {
            out: Ok(ProcessOutput {
                stdout: Vec::new(),
                stderr: b"boom".to_vec(),
                exit_code: Some(1),
            }),
        };
        let resp = handle(
            &req,
            &host(),
            "127.0.0.1:8080".parse().unwrap(),
            Path::new("/www/cgi/bad.py"),
            &cgi,
            &runner,
        );
        assert_eq!(resp.status(), Status::BAD_GATEWAY);
        assert!(!resp.body().is_empty());
    }
}
