//! Config file parser.
//!
//! TOML is deserialized into thin `Raw*` shapes, then folded through the
//! validated domain constructors. Domain rules (path normalization, duplicate
//! listeners, redirect must be 3xx, etc.) live in `domain/`; this module is a
//! thin translator that surfaces those errors with file context.

use std::collections::HashMap;
use std::fs;
use std::net::SocketAddr;
use std::path::{Path, PathBuf};

use serde::Deserialize;

use crate::domain::config::cgi::CgiConfig;
use crate::domain::config::route::{RouteConfig, RouteKind};
use crate::domain::config::server::{HostConfig, ServerConfig};
use crate::domain::error::DomainError;
use crate::domain::http::method::Method;
use crate::domain::http::status::Status;

#[derive(Debug)]
pub enum ConfigError {
    Io(std::io::Error),
    Toml(toml::de::Error),
    Domain(DomainError),
    BadShape(String),
    BadAddress(String),
    BadStatusCode(u16),
}

impl std::fmt::Display for ConfigError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::Io(e) => write!(f, "config io error: {e}"),
            Self::Toml(e) => write!(f, "config parse error: {e}"),
            Self::Domain(e) => write!(f, "config validation error: {e}"),
            Self::BadShape(s) => write!(f, "config shape error: {s}"),
            Self::BadAddress(s) => write!(f, "invalid listen address: {s:?}"),
            Self::BadStatusCode(c) => write!(f, "invalid status code in error_pages: {c}"),
        }
    }
}

impl std::error::Error for ConfigError {}

impl From<DomainError> for ConfigError {
    fn from(e: DomainError) -> Self {
        Self::Domain(e)
    }
}

impl From<std::io::Error> for ConfigError {
    fn from(e: std::io::Error) -> Self {
        Self::Io(e)
    }
}

impl From<toml::de::Error> for ConfigError {
    fn from(e: toml::de::Error) -> Self {
        Self::Toml(e)
    }
}

// ---- Raw schema (deserialize-only) ----------------------------------------

#[derive(Debug, Deserialize)]
struct RawConfig {
    server: Vec<RawServer>,
}

#[derive(Debug, Deserialize)]
struct RawServer {
    listen: Vec<String>,
    #[serde(default)]
    server_names: Vec<String>,
    client_max_body_size: u64,
    #[serde(default)]
    error_pages: HashMap<String, PathBuf>,
    route: Vec<RawRoute>,
}

#[derive(Debug, Deserialize)]
struct RawRoute {
    path: String,
    methods: Vec<String>,
    // Static-mode fields
    root: Option<PathBuf>,
    index: Option<String>,
    #[serde(default)]
    autoindex: bool,
    upload_dir: Option<PathBuf>,
    // Redirect-mode fields
    redirect: Option<String>,
    redirect_status: Option<u16>,
    // Optional CGI
    cgi: Option<RawCgi>,
    // Built-in session counter demo
    #[serde(default)]
    session_counter: bool,
}

#[derive(Debug, Deserialize)]
struct RawCgi {
    extension: String,
    interpreter: PathBuf,
}

// ---- Public API -----------------------------------------------------------

/// Parse a TOML config string into a validated `ServerConfig`.
pub fn parse_str(s: &str) -> Result<ServerConfig, ConfigError> {
    let raw: RawConfig = toml::from_str(s)?;
    to_domain(raw)
}

/// Read and parse a TOML config file.
pub fn load_file(path: &Path) -> Result<ServerConfig, ConfigError> {
    let text = fs::read_to_string(path)?;
    parse_str(&text)
}

// ---- Raw -> Domain --------------------------------------------------------

fn to_domain(raw: RawConfig) -> Result<ServerConfig, ConfigError> {
    if raw.server.is_empty() {
        return Err(ConfigError::BadShape("no [[server]] entries".into()));
    }
    let mut hosts: Vec<HostConfig> = Vec::with_capacity(raw.server.len());
    for s in raw.server {
        hosts.push(host_from_raw(s)?);
    }
    Ok(ServerConfig::new(hosts)?)
}

fn host_from_raw(s: RawServer) -> Result<HostConfig, ConfigError> {
    let listeners = s
        .listen
        .iter()
        .map(|a| {
            a.parse::<SocketAddr>()
                .map_err(|_| ConfigError::BadAddress(a.clone()))
        })
        .collect::<Result<Vec<_>, _>>()?;

    let error_pages = s
        .error_pages
        .into_iter()
        .map(|(k, v)| {
            let code: u16 = k.parse().map_err(|_| ConfigError::BadStatusCode(0))?;
            let _ = Status::new(code).map_err(|_| ConfigError::BadStatusCode(code))?;
            Ok::<(u16, PathBuf), ConfigError>((code, v))
        })
        .collect::<Result<HashMap<_, _>, _>>()?;

    let routes = s
        .route
        .into_iter()
        .map(route_from_raw)
        .collect::<Result<Vec<_>, _>>()?;

    Ok(HostConfig::new(
        listeners,
        s.server_names,
        s.client_max_body_size,
        error_pages,
        routes,
    )?)
}

fn route_from_raw(r: RawRoute) -> Result<RouteConfig, ConfigError> {
    let methods = r
        .methods
        .iter()
        .map(|m| m.parse::<Method>())
        .collect::<Result<Vec<_>, _>>()?;

    let kind = match (r.session_counter, r.redirect, r.root) {
        (true, Some(_), _) | (true, _, Some(_)) => {
            return Err(ConfigError::BadShape(format!(
                "route {:?} sets session_counter alongside root or redirect",
                r.path
            )));
        }
        (true, None, None) => RouteKind::SessionCounter,
        (false, Some(_), Some(_)) => {
            return Err(ConfigError::BadShape(format!(
                "route {:?} sets both redirect and root",
                r.path
            )));
        }
        (false, Some(loc), None) => {
            let code = r.redirect_status.unwrap_or(302);
            let status = Status::new(code).map_err(|_| ConfigError::BadStatusCode(code))?;
            RouteKind::Redirect {
                location: loc,
                status,
            }
        }
        (false, None, Some(root)) => RouteKind::Static {
            root,
            index: r.index,
            autoindex: r.autoindex,
            upload_dir: r.upload_dir,
        },
        (false, None, None) => {
            return Err(ConfigError::BadShape(format!(
                "route {:?} has neither root, redirect, nor session_counter",
                r.path
            )));
        }
    };

    let cgi = r
        .cgi
        .map(|c| CgiConfig::new(c.extension, c.interpreter))
        .transpose()?;

    Ok(RouteConfig::new(r.path, methods, kind, cgi)?)
}

// ---- Tests ----------------------------------------------------------------

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;
    use crate::domain::config::route::RouteKind;

    fn ok(s: &str) -> ServerConfig {
        parse_str(s).unwrap()
    }

    fn err(s: &str) -> ConfigError {
        parse_str(s).unwrap_err()
    }

    const MIN: &str = r#"
[[server]]
listen = ["0.0.0.0:8080"]
client_max_body_size = 1048576
[[server.route]]
path = "/"
methods = ["GET"]
root = "/var/www"
"#;

    #[test]
    fn minimal_parses() {
        let cfg = ok(MIN);
        assert_eq!(cfg.hosts().len(), 1);
        let h = &cfg.hosts()[0];
        assert_eq!(h.listeners().len(), 1);
        assert_eq!(h.client_max_body_size(), 1_048_576);
        assert_eq!(h.routes().len(), 1);
    }

    #[test]
    fn empty_config_rejected() {
        // Missing `server` array surfaces as a TOML missing-field error.
        let e = err("");
        assert!(matches!(e, ConfigError::Toml(_)));
    }

    #[test]
    fn empty_server_array_rejected() {
        // An explicitly empty server list reaches our shape check.
        let e = err("server = []\n");
        assert!(matches!(e, ConfigError::BadShape(_)));
    }

    #[test]
    fn bad_listen_address_rejected() {
        let e = err(r#"
[[server]]
listen = ["not-an-address"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
root = "/var/www"
"#);
        assert!(matches!(e, ConfigError::BadAddress(_)));
    }

    #[test]
    fn bad_method_rejected() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GE T"]
root = "/var/www"
"#);
        assert!(matches!(
            e,
            ConfigError::Domain(DomainError::InvalidMethod(_))
        ));
    }

    #[test]
    fn route_must_be_absolute() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "no-slash"
methods = ["GET"]
root = "/var/www"
"#);
        assert!(matches!(
            e,
            ConfigError::Domain(DomainError::InvalidConfig(_))
        ));
    }

    #[test]
    fn route_needs_root_or_redirect() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
"#);
        assert!(matches!(e, ConfigError::BadShape(_)));
    }

    #[test]
    fn route_cant_have_root_and_redirect() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
root = "/var/www"
redirect = "/elsewhere"
"#);
        assert!(matches!(e, ConfigError::BadShape(_)));
    }

    #[test]
    fn redirect_default_status_is_302() {
        let cfg = ok(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/old"
methods = ["GET"]
redirect = "/new"
"#);
        let route = &cfg.hosts()[0].routes()[0];
        match route.kind() {
            RouteKind::Redirect { status, .. } => assert_eq!(*status, Status::FOUND),
            RouteKind::Static { .. } | RouteKind::SessionCounter => panic!("expected redirect"),
        }
    }

    #[test]
    fn redirect_non_3xx_rejected() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/old"
methods = ["GET"]
redirect = "/new"
redirect_status = 200
"#);
        assert!(matches!(
            e,
            ConfigError::Domain(DomainError::InvalidConfig(_))
        ));
    }

    #[test]
    fn session_counter_route_parses() {
        let cfg = ok(r#"
[[server]]
listen = ["0.0.0.0:8080"]
client_max_body_size = 1
[[server.route]]
path = "/session"
methods = ["GET"]
session_counter = true
"#);
        assert!(matches!(
            cfg.hosts()[0].routes()[0].kind(),
            RouteKind::SessionCounter
        ));
    }

    #[test]
    fn session_counter_with_root_rejected() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:8080"]
client_max_body_size = 1
[[server.route]]
path = "/session"
methods = ["GET"]
session_counter = true
root = "/www"
"#);
        assert!(matches!(e, ConfigError::BadShape(_)));
    }

    #[test]
    fn cgi_extension_must_start_with_dot() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/cgi"
methods = ["GET"]
root = "/var/www"
[server.route.cgi]
extension = "py"
interpreter = "/usr/bin/python3"
"#);
        assert!(matches!(
            e,
            ConfigError::Domain(DomainError::InvalidConfig(_))
        ));
    }

    #[test]
    fn duplicate_default_listener_rejected() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
root = "/a"

[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
root = "/b"
"#);
        assert!(matches!(
            e,
            ConfigError::Domain(DomainError::DuplicateListener(_))
        ));
    }

    #[test]
    fn same_addr_with_distinct_names_ok() {
        let cfg = ok(r#"
[[server]]
listen = ["0.0.0.0:80"]
server_names = ["a.test"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
root = "/a"

[[server]]
listen = ["0.0.0.0:80"]
server_names = ["b.test"]
client_max_body_size = 1
[[server.route]]
path = "/"
methods = ["GET"]
root = "/b"
"#);
        assert_eq!(cfg.hosts().len(), 2);
        assert_eq!(cfg.distinct_listeners().len(), 1);
    }

    #[test]
    fn duplicate_route_within_host_rejected() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/api"
methods = ["GET"]
root = "/a"
[[server.route]]
path = "/api"
methods = ["POST"]
root = "/b"
"#);
        assert!(matches!(
            e,
            ConfigError::Domain(DomainError::DuplicateRoute(_))
        ));
    }

    #[test]
    fn error_pages_loaded() {
        let cfg = ok(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1

[server.error_pages]
404 = "/srv/404.html"
500 = "/srv/500.html"

[[server.route]]
path = "/"
methods = ["GET"]
root = "/a"
"#);
        let h = &cfg.hosts()[0];
        assert!(h.error_page(404).is_some());
        assert!(h.error_page(500).is_some());
        assert!(h.error_page(403).is_none());
    }

    #[test]
    fn error_pages_invalid_code_rejected() {
        let e = err(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[server.error_pages]
99 = "/srv/x.html"
[[server.route]]
path = "/"
methods = ["GET"]
root = "/a"
"#);
        assert!(matches!(e, ConfigError::BadStatusCode(_)));
    }

    #[test]
    fn cgi_round_trip() {
        let cfg = ok(r#"
[[server]]
listen = ["0.0.0.0:80"]
client_max_body_size = 1
[[server.route]]
path = "/cgi"
methods = ["GET", "POST"]
root = "/var/www"
[server.route.cgi]
extension = ".py"
interpreter = "/usr/bin/python3"
"#);
        let route = &cfg.hosts()[0].routes()[0];
        let cgi = route.cgi().expect("cgi present");
        assert_eq!(cgi.extension(), ".py");
        assert!(route.allows(&Method::Get));
        assert!(route.allows(&Method::Post));
    }
}
