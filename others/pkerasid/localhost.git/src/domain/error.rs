use std::fmt;

/// All validation failures originating from the domain layer.
///
/// Domain types use this in their constructors to refuse invalid inputs at the
/// boundary, so the rest of the system can assume well-formed values.
#[derive(Debug, Clone, Eq, PartialEq)]
pub enum DomainError {
    InvalidMethod(String),
    InvalidStatus(u16),
    InvalidVersion(String),
    InvalidHeaderName(String),
    InvalidHeaderValue(String),
    InvalidPath(String),
    PathTraversal,
    InvalidPercentEncoding,
    InvalidQuery(String),
    InvalidConfig(String),
    InvalidCookieName(String),
    InvalidCookieValue(String),
    InvalidSessionId,
    DuplicateRoute(String),
    DuplicateListener(String),
}

impl fmt::Display for DomainError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::InvalidMethod(s) => write!(f, "invalid HTTP method: {s:?}"),
            Self::InvalidStatus(c) => write!(f, "invalid HTTP status code: {c}"),
            Self::InvalidVersion(s) => write!(f, "invalid HTTP version: {s:?}"),
            Self::InvalidHeaderName(s) => write!(f, "invalid header name: {s:?}"),
            Self::InvalidHeaderValue(s) => write!(f, "invalid header value: {s:?}"),
            Self::InvalidPath(s) => write!(f, "invalid request path: {s:?}"),
            Self::PathTraversal => write!(f, "path escapes its root"),
            Self::InvalidPercentEncoding => write!(f, "malformed percent-encoded sequence"),
            Self::InvalidQuery(s) => write!(f, "invalid query string: {s:?}"),
            Self::InvalidConfig(s) => write!(f, "invalid configuration: {s}"),
            Self::InvalidCookieName(s) => write!(f, "invalid cookie name: {s:?}"),
            Self::InvalidCookieValue(s) => write!(f, "invalid cookie value: {s:?}"),
            Self::InvalidSessionId => write!(f, "invalid session id"),
            Self::DuplicateRoute(s) => write!(f, "duplicate route path: {s:?}"),
            Self::DuplicateListener(s) => write!(f, "duplicate listener: {s:?}"),
        }
    }
}

impl std::error::Error for DomainError {}
