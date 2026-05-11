use std::path::PathBuf;

use crate::domain::config::cgi::CgiConfig;
use crate::domain::error::DomainError;
use crate::domain::http::method::Method;
use crate::domain::http::status::Status;

/// What a route serves: static files, a redirect, session counter, or CGI (orthogonal).
#[derive(Debug, Clone, Eq, PartialEq)]
pub enum RouteKind {
    Static {
        root: PathBuf,
        index: Option<String>,
        autoindex: bool,
        upload_dir: Option<PathBuf>,
    },
    Redirect {
        location: String,
        status: Status,
    },
    /// Built-in demo handler: reads a per-session visit counter, increments it,
    /// and returns an HTML page — exercises the session middleware for the audit.
    SessionCounter,
}

#[derive(Debug, Clone, Eq, PartialEq)]
pub struct RouteConfig {
    /// Path prefix this route matches (must start with `/`).
    path: String,
    methods: Vec<Method>,
    kind: RouteKind,
    cgi: Option<CgiConfig>,
}

impl RouteConfig {
    pub fn new(
        path: impl Into<String>,
        methods: Vec<Method>,
        kind: RouteKind,
        cgi: Option<CgiConfig>,
    ) -> Result<Self, DomainError> {
        let path: String = path.into();
        if !path.starts_with('/') {
            return Err(DomainError::InvalidConfig(format!(
                "route path must start with '/': {path:?}"
            )));
        }
        if methods.is_empty() {
            return Err(DomainError::InvalidConfig(format!(
                "route {path:?} must allow at least one method"
            )));
        }
        if let RouteKind::Redirect { location, status } = &kind {
            if location.is_empty() {
                return Err(DomainError::InvalidConfig(format!(
                    "route {path:?} redirect target is empty"
                )));
            }
            if !status.is_redirection() {
                return Err(DomainError::InvalidConfig(format!(
                    "route {path:?} redirect uses non-3xx status {status}"
                )));
            }
        }
        Ok(Self {
            path,
            methods,
            kind,
            cgi,
        })
    }

    pub fn path(&self) -> &str {
        &self.path
    }

    pub fn methods(&self) -> &[Method] {
        &self.methods
    }

    pub fn kind(&self) -> &RouteKind {
        &self.kind
    }

    pub fn cgi(&self) -> Option<&CgiConfig> {
        self.cgi.as_ref()
    }

    pub fn allows(&self, m: &Method) -> bool {
        self.methods.contains(m)
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn static_kind() -> RouteKind {
        RouteKind::Static {
            root: PathBuf::from("/var/www"),
            index: Some("index.html".into()),
            autoindex: false,
            upload_dir: None,
        }
    }

    #[test]
    fn requires_leading_slash() {
        assert!(RouteConfig::new("public", vec![Method::Get], static_kind(), None).is_err());
        assert!(RouteConfig::new("/", vec![Method::Get], static_kind(), None).is_ok());
    }

    #[test]
    fn requires_at_least_one_method() {
        assert!(RouteConfig::new("/", vec![], static_kind(), None).is_err());
    }

    #[test]
    fn redirect_must_be_3xx() {
        let bad = RouteKind::Redirect {
            location: "/elsewhere".into(),
            status: Status::OK,
        };
        assert!(RouteConfig::new("/old", vec![Method::Get], bad, None).is_err());

        let good = RouteKind::Redirect {
            location: "/elsewhere".into(),
            status: Status::FOUND,
        };
        assert!(RouteConfig::new("/old", vec![Method::Get], good, None).is_ok());
    }

    #[test]
    fn redirect_target_non_empty() {
        let bad = RouteKind::Redirect {
            location: String::new(),
            status: Status::FOUND,
        };
        assert!(RouteConfig::new("/old", vec![Method::Get], bad, None).is_err());
    }

    #[test]
    fn allows_method() {
        let r =
            RouteConfig::new("/api", vec![Method::Get, Method::Post], static_kind(), None).unwrap();
        assert!(r.allows(&Method::Get));
        assert!(!r.allows(&Method::Delete));
    }
}
