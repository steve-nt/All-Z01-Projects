//! Request pipeline: route → validate → dispatch → `Response`.
//!
//! `PipelineContext` holds the shared config and filesystem adapter needed by
//! all handlers. `Connection` owns an `Rc<PipelineContext>` so the event loop
//! creates it once and hands it to every accepted connection.
//!
//! Session middleware: when `session_store` is `Some`, the pipeline reads the
//! `SID` cookie from every inbound request, loads or creates a session, passes
//! it into the dispatch function, and flushes the (possibly mutated) session
//! back to the store. A `Set-Cookie: SID=…` header is added to the response
//! whenever a new session is created.

use std::cell::RefCell;
use std::net::SocketAddr;
use std::path::PathBuf;
use std::rc::Rc;
use std::time::{SystemTime, UNIX_EPOCH};

use crate::application::error_pages;
use crate::application::handlers::{cgi, delete, redirect, session_counter, static_file, upload};
use crate::application::ports::filesystem::FileSystem;
use crate::application::ports::process_runner::ProcessRunner;
use crate::application::ports::session_store::SessionStore;
use crate::application::router;
use crate::domain::config::route::RouteKind;
use crate::domain::config::server::ServerConfig;
use crate::domain::http::headers::{HeaderName, HeaderValue};
use crate::domain::http::method::Method;
use crate::domain::http::request::Request;
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;
use crate::domain::session::cookie::parse_cookie_header;
use crate::domain::session::session::Session;
use crate::domain::session::session_id::SessionId;

/// Session idle timeout: 30 minutes.
const SESSION_MAX_IDLE_SECS: u64 = 30 * 60;

/// Shared state passed to every connection.
pub struct PipelineContext {
    pub config: Rc<ServerConfig>,
    pub fs: Rc<dyn FileSystem>,
    pub process_runner: Rc<dyn ProcessRunner>,
    /// When `Some`, sessions are enabled for every request.
    pub session_store: Option<Rc<RefCell<dyn SessionStore>>>,
}

impl std::fmt::Debug for PipelineContext {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("PipelineContext").finish_non_exhaustive()
    }
}

/// Produce a `Response` for `req` arriving on `local_addr`.
pub fn handle(req: &Request, local_addr: SocketAddr, ctx: &PipelineContext) -> Response {
    let config = ctx.config.as_ref();
    let fs = ctx.fs.as_ref();

    // Resolve host + route.
    let Some(matched) = router::route(config, local_addr, req.host(), req.path().as_str()) else {
        return error_pages::error_response(Status::NOT_FOUND, None, fs);
    };
    let host = matched.host;
    let route = matched.route;

    // Method allowed?
    if !route.allows(req.method()) {
        let allow: String = route
            .methods()
            .iter()
            .map(Method::as_str)
            .collect::<Vec<_>>()
            .join(", ");
        let body = format!("{}\n", Status::METHOD_NOT_ALLOWED.reason()).into_bytes();
        let mut builder = Response::builder(Status::METHOD_NOT_ALLOWED);
        if let (Ok(n), Ok(v)) = (HeaderName::new("allow"), HeaderValue::new(allow)) {
            builder = builder.header(n, v);
        }
        return builder.body(body).build();
    }

    // Session middleware: load or create session before dispatch.
    let (mut session_opt, is_new_session) = load_or_create_session(req, ctx);

    // Dispatch.
    let mut resp = match route.kind() {
        RouteKind::Redirect { location, status } => redirect::handle(location, *status),

        RouteKind::SessionCounter => {
            let visits = session_opt.as_mut().map_or(0, |s| {
                let v = s
                    .get("visits")
                    .and_then(|n| n.parse::<u64>().ok())
                    .unwrap_or(0)
                    .saturating_add(1);
                s.set("visits", v.to_string());
                v
            });
            session_counter::handle(visits)
        }

        RouteKind::Static {
            root,
            index,
            autoindex,
            upload_dir,
        } => {
            let file_path = resolve_path(route.path(), root, req.path().as_str());
            if let Some(cgi_cfg) = route.cgi()
                && cgi_cfg.matches(req.path().as_str())
            {
                cgi::handle(req, host, local_addr, &file_path, cgi_cfg, fs_runner(ctx))
            } else {
                match req.method() {
                    Method::Get => {
                        static_file::handle(&file_path, index.as_deref(), *autoindex, fs)
                    }
                    Method::Post => {
                        let dir = upload_dir.as_deref().unwrap_or(root.as_path());
                        upload::handle(req, dir, fs)
                    }
                    Method::Delete => delete::handle(&file_path, fs),
                    Method::Other(_) => {
                        error_pages::error_response(Status::NOT_IMPLEMENTED, Some(host), fs)
                    }
                }
            }
        }
    };

    // Session middleware: flush session and inject Set-Cookie if new.
    if let Some(session) = session_opt {
        if is_new_session {
            inject_set_cookie(&mut resp, session.id());
        }
        if let Some(store_rc) = &ctx.session_store {
            store_rc.borrow_mut().put(session);
        }
    }

    resp
}

/// Extract the `SID` cookie from the request and look up the session in the
/// store, or create a fresh one. Returns `(session, is_new)`.
///
/// Returns `(None, false)` when no store is configured.
fn load_or_create_session(req: &Request, ctx: &PipelineContext) -> (Option<Session>, bool) {
    let Some(store_rc) = &ctx.session_store else {
        return (None, false);
    };

    let now_secs = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map_or(0, |d| d.as_secs());

    // Evict expired sessions opportunistically on each request.
    store_rc
        .borrow_mut()
        .evict_expired(now_secs, SESSION_MAX_IDLE_SECS);

    // Look for an existing valid session in the Cookie header.
    let existing = req.headers().get("cookie").and_then(|raw| {
        parse_cookie_header(raw)
            .into_iter()
            .find(|(name, _)| name == "SID")
            .and_then(|(_, val)| SessionId::new(val).ok())
            .and_then(|id| {
                let store = store_rc.borrow();
                let s = store.get(&id)?;
                if s.is_expired(now_secs, SESSION_MAX_IDLE_SECS) {
                    None
                } else {
                    Some(s)
                }
            })
    });

    if let Some(mut s) = existing {
        s.touch(now_secs);
        (Some(s), false)
    } else {
        let id_str = crate::infrastructure::session_store::id_gen::generate();
        if let Ok(id) = SessionId::new(id_str) {
            (Some(Session::new(id, now_secs)), true)
        } else {
            (None, false)
        }
    }
}

/// Append a `Set-Cookie: SID=<id>; Path=/; HttpOnly; SameSite=Lax` header.
fn inject_set_cookie(resp: &mut Response, id: &SessionId) {
    use crate::domain::session::cookie::{Cookie, SameSite};
    if let Ok(cookie) = Cookie::new("SID", id.as_str()) {
        let value = cookie
            .with_path("/")
            .http_only()
            .with_same_site(SameSite::Lax)
            .render_set_cookie();
        if let (Ok(n), Ok(v)) = (HeaderName::new("set-cookie"), HeaderValue::new(value)) {
            resp.add_header(n, v);
        }
    }
}

/// Strip the route prefix from the request path, then join with root.
///
/// Example: route `/static`, root `/var/www/html`, path `/static/img/logo.png`
/// → `/var/www/html/img/logo.png`.
fn resolve_path(route_prefix: &str, root: &std::path::Path, req_path: &str) -> PathBuf {
    let relative = if route_prefix == "/" {
        req_path
    } else {
        req_path.strip_prefix(route_prefix).unwrap_or(req_path)
    };
    root.join(relative.trim_start_matches('/'))
}

fn fs_runner(ctx: &PipelineContext) -> &dyn ProcessRunner {
    ctx.process_runner.as_ref()
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::collections::HashMap;
    use std::net::SocketAddr;
    use std::path::{Path, PathBuf};
    use std::rc::Rc;

    use super::*;
    use crate::application::ports::filesystem::fake::FakeFileSystem;
    use crate::application::ports::process_runner::ProcessOutput;
    use crate::application::ports::process_runner::fake::FakeProcessRunner;
    use crate::domain::config::cgi::CgiConfig;
    use crate::domain::config::route::{RouteConfig, RouteKind};
    use crate::domain::config::server::{HostConfig, ServerConfig};
    use crate::domain::http::method::Method;
    use crate::domain::http::request::Request;
    use crate::domain::http::url::RequestPath;
    use crate::domain::http::version::HttpVersion;

    fn addr(s: &str) -> SocketAddr {
        s.parse().unwrap()
    }

    fn get(path: &str) -> Request {
        Request::builder(
            Method::Get,
            RequestPath::parse(path).unwrap(),
            HttpVersion::Http11,
        )
        .build()
    }

    fn make_config(local: SocketAddr, routes: Vec<RouteConfig>) -> Rc<ServerConfig> {
        let host =
            HostConfig::new(vec![local], vec![], 1024 * 1024, HashMap::new(), routes).unwrap();
        Rc::new(ServerConfig::new(vec![host]).unwrap())
    }

    fn runner() -> Rc<FakeProcessRunner> {
        Rc::new(FakeProcessRunner::new())
    }

    fn static_route(path: &str, root: &str, methods: Vec<Method>) -> RouteConfig {
        RouteConfig::new(
            path,
            methods,
            RouteKind::Static {
                root: PathBuf::from(root),
                index: Some("index.html".into()),
                autoindex: false,
                upload_dir: None,
            },
            None,
        )
        .unwrap()
    }

    #[test]
    fn serves_static_file() {
        let local = addr("127.0.0.1:8080");
        let fs = FakeFileSystem::new();
        fs.add_dir("/www");
        fs.add_file("/www/index.html", b"<h1>Hi</h1>");
        let cfg = make_config(local, vec![static_route("/", "/www", vec![Method::Get])]);
        let ctx = PipelineContext {
            config: cfg,
            fs: Rc::new(fs),
            process_runner: runner(),
            session_store: None,
        };
        let resp = handle(&get("/"), local, &ctx);
        assert_eq!(resp.status(), Status::OK);
        assert_eq!(resp.body(), b"<h1>Hi</h1>");
    }

    #[test]
    fn method_not_allowed_has_allow_header() {
        let local = addr("127.0.0.1:8080");
        let cfg = make_config(local, vec![static_route("/", "/www", vec![Method::Get])]);
        let ctx = PipelineContext {
            config: cfg,
            fs: Rc::new(FakeFileSystem::new()),
            process_runner: runner(),
            session_store: None,
        };
        let req = Request::builder(
            Method::Delete,
            RequestPath::parse("/file.txt").unwrap(),
            HttpVersion::Http11,
        )
        .build();
        let resp = handle(&req, local, &ctx);
        assert_eq!(resp.status(), Status::METHOD_NOT_ALLOWED);
        assert!(resp.headers().get("allow").is_some());
    }

    #[test]
    fn redirect_route() {
        let local = addr("127.0.0.1:8080");
        let route = RouteConfig::new(
            "/old",
            vec![Method::Get],
            RouteKind::Redirect {
                location: "/new".into(),
                status: Status::MOVED_PERMANENTLY,
            },
            None,
        )
        .unwrap();
        let cfg = make_config(local, vec![route]);
        let ctx = PipelineContext {
            config: cfg,
            fs: Rc::new(FakeFileSystem::new()),
            process_runner: runner(),
            session_store: None,
        };
        let resp = handle(&get("/old"), local, &ctx);
        assert_eq!(resp.status(), Status::MOVED_PERMANENTLY);
        assert_eq!(resp.headers().get("location"), Some("/new"));
    }

    #[test]
    fn no_matching_host_returns_404() {
        let local = addr("127.0.0.1:8080");
        let cfg = make_config(local, vec![static_route("/", "/www", vec![Method::Get])]);
        let ctx = PipelineContext {
            config: cfg,
            fs: Rc::new(FakeFileSystem::new()),
            process_runner: runner(),
            session_store: None,
        };
        let resp = handle(&get("/"), addr("127.0.0.1:9999"), &ctx);
        assert_eq!(resp.status(), Status::NOT_FOUND);
    }

    #[test]
    fn cgi_route_executes_runner() {
        let local = addr("127.0.0.1:8080");
        let cgi_route = RouteConfig::new(
            "/cgi",
            vec![Method::Get],
            RouteKind::Static {
                root: PathBuf::from("/www"),
                index: None,
                autoindex: false,
                upload_dir: None,
            },
            Some(CgiConfig::new(".py", "/usr/bin/python3").unwrap()),
        )
        .unwrap();
        let cfg = make_config(local, vec![cgi_route]);
        let fake_runner = runner();
        fake_runner.push_output(Ok(ProcessOutput {
            stdout: b"Status: 200 OK\r\nContent-Type: text/plain\r\n\r\ncgi-ok".to_vec(),
            stderr: Vec::new(),
            exit_code: Some(0),
        }));
        let ctx = PipelineContext {
            config: cfg,
            fs: Rc::new(FakeFileSystem::new()),
            process_runner: fake_runner,
            session_store: None,
        };

        let resp = handle(&get("/cgi/hello.py"), local, &ctx);
        assert_eq!(resp.status(), Status::OK);
        assert_eq!(resp.body(), b"cgi-ok");
        assert_eq!(resp.headers().get("content-type"), Some("text/plain"));
    }

    #[test]
    fn resolve_path_strips_prefix() {
        let p = resolve_path("/static", Path::new("/var/www"), "/static/img/logo.png");
        assert_eq!(p, PathBuf::from("/var/www/img/logo.png"));
    }

    #[test]
    fn resolve_path_root_prefix() {
        let p = resolve_path("/", Path::new("/var/www"), "/index.html");
        assert_eq!(p, PathBuf::from("/var/www/index.html"));
    }

    fn session_counter_route(path: &str) -> RouteConfig {
        RouteConfig::new(path, vec![Method::Get], crate::domain::config::route::RouteKind::SessionCounter, None).unwrap()
    }

    fn make_ctx_with_store(local: SocketAddr, store: impl SessionStore + 'static) -> PipelineContext {
        let cfg = make_config(local, vec![session_counter_route("/session")]);
        PipelineContext {
            config: cfg,
            fs: Rc::new(FakeFileSystem::new()),
            process_runner: runner(),
            session_store: Some(Rc::new(RefCell::new(store))),
        }
    }

    #[test]
    fn session_counter_increments_on_each_visit() {
        use crate::application::ports::session_store::FakeSessionStore;
        use std::cell::RefCell;

        let local = addr("127.0.0.1:8080");
        let store = FakeSessionStore::new();
        let ctx = make_ctx_with_store(local, store);

        let resp1 = handle(&get("/session"), local, &ctx);
        assert_eq!(resp1.status(), Status::OK);
        let body1 = String::from_utf8_lossy(resp1.body());
        assert!(body1.contains("<strong>1</strong>"), "first visit should be 1: {body1}");
        assert!(resp1.headers().get("set-cookie").is_some(), "new session needs Set-Cookie");

        // Second visit: re-use the SID from the first response.
        let sid_header = resp1.headers().get("set-cookie").unwrap();
        let sid_val = sid_header.split('=').nth(1).unwrap().split(';').next().unwrap();
        let req2 = Request::builder(
            Method::Get,
            RequestPath::parse("/session").unwrap(),
            HttpVersion::Http11,
        )
        .header(
            HeaderName::new("cookie").unwrap(),
            HeaderValue::new(format!("SID={sid_val}")).unwrap(),
        )
        .build();
        let resp2 = handle(&req2, local, &ctx);
        let body2 = String::from_utf8_lossy(resp2.body());
        assert!(body2.contains("<strong>2</strong>"), "second visit should be 2: {body2}");
        assert!(resp2.headers().get("set-cookie").is_none(), "existing session must not re-set cookie");
    }

    #[test]
    fn no_set_cookie_when_no_store() {
        let local = addr("127.0.0.1:8080");
        let cfg = make_config(local, vec![session_counter_route("/session")]);
        let ctx = PipelineContext {
            config: cfg,
            fs: Rc::new(FakeFileSystem::new()),
            process_runner: runner(),
            session_store: None,
        };
        let resp = handle(&get("/session"), local, &ctx);
        assert_eq!(resp.status(), Status::OK);
        assert!(resp.headers().get("set-cookie").is_none());
    }
}
