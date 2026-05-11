use std::collections::{HashMap, HashSet};
use std::net::SocketAddr;
use std::path::PathBuf;

use crate::domain::config::route::RouteConfig;
use crate::domain::error::DomainError;

/// One virtual host listening on one or more sockets.
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct HostConfig {
    listeners: Vec<SocketAddr>,
    server_names: Vec<String>,
    client_max_body_size: u64,
    error_pages: HashMap<u16, PathBuf>,
    routes: Vec<RouteConfig>,
}

impl HostConfig {
    pub fn new(
        listeners: Vec<SocketAddr>,
        server_names: Vec<String>,
        client_max_body_size: u64,
        error_pages: HashMap<u16, PathBuf>,
        routes: Vec<RouteConfig>,
    ) -> Result<Self, DomainError> {
        if listeners.is_empty() {
            return Err(DomainError::InvalidConfig("host has no listeners".into()));
        }
        if routes.is_empty() {
            return Err(DomainError::InvalidConfig("host has no routes".into()));
        }
        if client_max_body_size == 0 {
            return Err(DomainError::InvalidConfig(
                "client_max_body_size must be > 0".into(),
            ));
        }

        // Reject duplicate route paths within a host.
        let mut seen = HashSet::with_capacity(routes.len());
        for r in &routes {
            if !seen.insert(r.path().to_owned()) {
                return Err(DomainError::DuplicateRoute(r.path().to_owned()));
            }
        }

        Ok(Self {
            listeners,
            server_names,
            client_max_body_size,
            error_pages,
            routes,
        })
    }

    pub fn listeners(&self) -> &[SocketAddr] {
        &self.listeners
    }

    pub fn server_names(&self) -> &[String] {
        &self.server_names
    }

    pub fn client_max_body_size(&self) -> u64 {
        self.client_max_body_size
    }

    pub fn error_page(&self, status_code: u16) -> Option<&std::path::Path> {
        self.error_pages.get(&status_code).map(PathBuf::as_path)
    }

    pub fn routes(&self) -> &[RouteConfig] {
        &self.routes
    }

    pub fn matches_name(&self, host_header: &str) -> bool {
        self.server_names
            .iter()
            .any(|n| n.eq_ignore_ascii_case(host_header))
    }
}

/// Top-level configuration: a collection of virtual hosts.
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct ServerConfig {
    hosts: Vec<HostConfig>,
}

impl ServerConfig {
    pub fn new(hosts: Vec<HostConfig>) -> Result<Self, DomainError> {
        if hosts.is_empty() {
            return Err(DomainError::InvalidConfig("no hosts configured".into()));
        }

        // (addr, server_name) must be unique. Two hosts may share an addr only
        // if they have distinct server_names; a single host with no name is
        // the default for that addr and there can be only one such default.
        let mut seen: HashMap<SocketAddr, HashSet<String>> = HashMap::new();
        let mut default_for: HashSet<SocketAddr> = HashSet::new();
        for h in &hosts {
            for addr in h.listeners() {
                let names = seen.entry(*addr).or_default();
                if h.server_names().is_empty() {
                    if !default_for.insert(*addr) {
                        return Err(DomainError::DuplicateListener(format!(
                            "{addr} has more than one default (unnamed) host"
                        )));
                    }
                } else {
                    for name in h.server_names() {
                        let key = name.to_ascii_lowercase();
                        if !names.insert(key) {
                            return Err(DomainError::DuplicateListener(format!(
                                "{addr} already has a host named {name:?}"
                            )));
                        }
                    }
                }
            }
        }

        Ok(Self { hosts })
    }

    pub fn hosts(&self) -> &[HostConfig] {
        &self.hosts
    }

    /// All distinct sockets to bind, deduplicated.
    pub fn distinct_listeners(&self) -> Vec<SocketAddr> {
        let mut out: Vec<SocketAddr> = Vec::new();
        for h in &self.hosts {
            for addr in h.listeners() {
                if !out.contains(addr) {
                    out.push(*addr);
                }
            }
        }
        out
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;
    use crate::domain::config::route::{RouteConfig, RouteKind};
    use crate::domain::http::method::Method;

    fn route(path: &str) -> RouteConfig {
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

    fn addr(s: &str) -> SocketAddr {
        s.parse().unwrap()
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

    #[test]
    fn host_requires_listener_and_route() {
        assert!(HostConfig::new(vec![], vec![], 1, HashMap::new(), vec![route("/")]).is_err());
        assert!(
            HostConfig::new(vec![addr("0.0.0.0:80")], vec![], 1, HashMap::new(), vec![]).is_err()
        );
    }

    #[test]
    fn host_rejects_zero_body_size() {
        assert!(
            HostConfig::new(
                vec![addr("0.0.0.0:80")],
                vec![],
                0,
                HashMap::new(),
                vec![route("/")]
            )
            .is_err()
        );
    }

    #[test]
    fn host_duplicate_routes_rejected() {
        let r1 = route("/api");
        let r2 = route("/api");
        assert!(matches!(
            HostConfig::new(
                vec![addr("0.0.0.0:80")],
                vec![],
                1,
                HashMap::new(),
                vec![r1, r2],
            ),
            Err(DomainError::DuplicateRoute(_))
        ));
    }

    #[test]
    fn server_rejects_two_unnamed_defaults_on_same_addr() {
        let h1 = host(vec![addr("0.0.0.0:80")], vec![], vec![route("/")]);
        let h2 = host(vec![addr("0.0.0.0:80")], vec![], vec![route("/")]);
        assert!(matches!(
            ServerConfig::new(vec![h1, h2]),
            Err(DomainError::DuplicateListener(_))
        ));
    }

    #[test]
    fn server_allows_same_addr_with_distinct_names() {
        let h1 = host(vec![addr("0.0.0.0:80")], vec!["a.test"], vec![route("/")]);
        let h2 = host(vec![addr("0.0.0.0:80")], vec!["b.test"], vec![route("/")]);
        assert!(ServerConfig::new(vec![h1, h2]).is_ok());
    }

    #[test]
    fn server_rejects_duplicate_names_on_same_addr() {
        let h1 = host(vec![addr("0.0.0.0:80")], vec!["a.test"], vec![route("/")]);
        let h2 = host(vec![addr("0.0.0.0:80")], vec!["A.TEST"], vec![route("/")]);
        assert!(matches!(
            ServerConfig::new(vec![h1, h2]),
            Err(DomainError::DuplicateListener(_))
        ));
    }

    #[test]
    fn distinct_listeners_dedup() {
        let h1 = host(vec![addr("0.0.0.0:80")], vec!["a"], vec![route("/")]);
        let h2 = host(
            vec![addr("0.0.0.0:80"), addr("0.0.0.0:81")],
            vec!["b"],
            vec![route("/")],
        );
        let cfg = ServerConfig::new(vec![h1, h2]).unwrap();
        let mut listeners = cfg.distinct_listeners();
        listeners.sort();
        assert_eq!(listeners, vec![addr("0.0.0.0:80"), addr("0.0.0.0:81")]);
    }

    #[test]
    fn matches_name_case_insensitive() {
        let h = host(
            vec![addr("0.0.0.0:80")],
            vec!["Example.COM"],
            vec![route("/")],
        );
        assert!(h.matches_name("example.com"));
        assert!(!h.matches_name("other.com"));
    }
}
