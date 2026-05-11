//! Route resolution: (`local_addr`, `host_header`, `path`) → (`HostConfig`, `RouteConfig`).
//!
//! Host selection: prefer a host whose `server_name` matches the `Host` header;
//! fall back to the unnamed default for that address, then the first host.
//!
//! Route selection: longest-prefix match on the normalized path.

use std::net::SocketAddr;

use crate::domain::config::route::RouteConfig;
use crate::domain::config::server::{HostConfig, ServerConfig};

#[derive(Debug)]
pub struct RouteMatch<'a> {
    pub host: &'a HostConfig,
    pub route: &'a RouteConfig,
}

/// Resolve the best (host, route) pair for an incoming request.
///
/// Returns `None` when no host listens on `local_addr` (should not happen in
/// normal operation) or when no route prefix matches `path`.
pub fn route<'a>(
    config: &'a ServerConfig,
    local_addr: SocketAddr,
    host_header: Option<&str>,
    path: &str,
) -> Option<RouteMatch<'a>> {
    let host = resolve_host(config, local_addr, host_header)?;
    let route = longest_prefix_match(host.routes(), path)?;
    Some(RouteMatch { host, route })
}

fn resolve_host<'a>(
    config: &'a ServerConfig,
    local_addr: SocketAddr,
    host_header: Option<&str>,
) -> Option<&'a HostConfig> {
    let candidates: Vec<&HostConfig> = config
        .hosts()
        .iter()
        .filter(|h| h.listeners().contains(&local_addr))
        .collect();

    if candidates.is_empty() {
        return None;
    }

    // Named match on Host header.
    if let Some(hdr) = host_header {
        let hdr_host = hdr.split(':').next().unwrap_or(hdr);
        if let Some(h) = candidates.iter().find(|h| h.matches_name(hdr_host)) {
            return Some(h);
        }
    }

    // Default (unnamed) host for this addr, then first.
    candidates
        .iter()
        .find(|h| h.server_names().is_empty())
        .or_else(|| candidates.first())
        .copied()
}

fn longest_prefix_match<'a>(routes: &'a [RouteConfig], path: &str) -> Option<&'a RouteConfig> {
    routes
        .iter()
        .filter(|r| path_has_prefix(path, r.path()))
        .max_by_key(|r| r.path().len())
}

/// True when `path` starts with the route `prefix`.
///
/// The root prefix `/` matches everything. For other prefixes the match
/// requires that either `path == prefix` or the next char is `/`, preventing
/// `/api` from matching `/apiv2`.
fn path_has_prefix(path: &str, prefix: &str) -> bool {
    if prefix == "/" {
        return true;
    }
    if path == prefix {
        return true;
    }
    path.starts_with(prefix) && path.as_bytes().get(prefix.len()).copied() == Some(b'/')
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::collections::HashMap;
    use std::path::PathBuf;

    use super::*;
    use crate::domain::config::route::{RouteConfig, RouteKind};
    use crate::domain::config::server::{HostConfig, ServerConfig};
    use crate::domain::http::method::Method;

    fn addr(s: &str) -> SocketAddr {
        s.parse().unwrap()
    }

    fn static_route(path: &str) -> RouteConfig {
        RouteConfig::new(
            path,
            vec![Method::Get],
            RouteKind::Static {
                root: PathBuf::from("/var/www"),
                index: None,
                autoindex: false,
                upload_dir: None,
            },
            None,
        )
        .unwrap()
    }

    fn host(addrs: Vec<SocketAddr>, names: Vec<&str>, routes: Vec<RouteConfig>) -> HostConfig {
        HostConfig::new(
            addrs,
            names.into_iter().map(String::from).collect(),
            1024 * 1024,
            HashMap::new(),
            routes,
        )
        .unwrap()
    }

    fn config_two_vhosts() -> ServerConfig {
        let a = host(
            vec![addr("127.0.0.1:8080")],
            vec!["a.test"],
            vec![static_route("/"), static_route("/api")],
        );
        let b = host(
            vec![addr("127.0.0.1:8080")],
            vec!["b.test"],
            vec![static_route("/")],
        );
        ServerConfig::new(vec![a, b]).unwrap()
    }

    #[test]
    fn resolves_by_server_name() {
        let cfg = config_two_vhosts();
        let m = route(&cfg, addr("127.0.0.1:8080"), Some("a.test"), "/index.html").unwrap();
        assert_eq!(m.host.server_names(), &["a.test"]);
        assert_eq!(m.route.path(), "/");

        let m = route(&cfg, addr("127.0.0.1:8080"), Some("b.test"), "/").unwrap();
        assert_eq!(m.host.server_names(), &["b.test"]);
    }

    #[test]
    fn longest_prefix_wins() {
        let cfg = config_two_vhosts();
        let m = route(&cfg, addr("127.0.0.1:8080"), Some("a.test"), "/api/users").unwrap();
        assert_eq!(m.route.path(), "/api");
    }

    #[test]
    fn no_false_prefix_match() {
        let cfg = config_two_vhosts();
        // /apiv2 must NOT match the /api route.
        let m = route(&cfg, addr("127.0.0.1:8080"), Some("a.test"), "/apiv2/foo").unwrap();
        assert_eq!(m.route.path(), "/");
    }

    #[test]
    fn unknown_host_header_falls_back_to_default() {
        let default_host = host(
            vec![addr("127.0.0.1:9090")],
            vec![],
            vec![static_route("/")],
        );
        let named = host(
            vec![addr("127.0.0.1:9090")],
            vec!["named.test"],
            vec![static_route("/")],
        );
        let cfg = ServerConfig::new(vec![default_host, named]).unwrap();
        let m = route(&cfg, addr("127.0.0.1:9090"), Some("unknown.host"), "/").unwrap();
        assert!(
            m.host.server_names().is_empty(),
            "should pick the default host"
        );
    }

    #[test]
    fn no_match_on_wrong_addr() {
        let cfg = config_two_vhosts();
        assert!(route(&cfg, addr("127.0.0.1:9999"), Some("a.test"), "/").is_none());
    }

    #[test]
    fn root_prefix_matches_all_paths() {
        assert!(path_has_prefix("/anything/goes", "/"));
        assert!(path_has_prefix("/", "/"));
    }

    #[test]
    fn prefix_boundary_check() {
        assert!(path_has_prefix("/api/v1", "/api"));
        assert!(path_has_prefix("/api", "/api"));
        assert!(!path_has_prefix("/apiv2", "/api"));
        assert!(!path_has_prefix("/ap", "/api"));
    }
}
