use std::fmt;

use crate::domain::error::DomainError;

/// Case-insensitive HTTP header name (stored canonicalized to lowercase).
#[derive(Debug, Clone, Eq, PartialEq, Hash)]
pub struct HeaderName(String);

impl HeaderName {
    /// Constructs a header name; rejects empty or non-token (RFC 7230) input.
    pub fn new(s: impl Into<String>) -> Result<Self, DomainError> {
        let owned: String = s.into();
        if !Self::is_token(&owned) {
            return Err(DomainError::InvalidHeaderName(owned));
        }
        Ok(Self(owned.to_ascii_lowercase()))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }

    fn is_token(s: &str) -> bool {
        !s.is_empty()
            && s.bytes().all(|b| {
                matches!(b,
                    b'!' | b'#' | b'$' | b'%' | b'&' | b'\'' | b'*' | b'+' |
                    b'-' | b'.' | b'^' | b'_' | b'`' | b'|' | b'~' |
                    b'0'..=b'9' | b'A'..=b'Z' | b'a'..=b'z')
            })
    }
}

impl fmt::Display for HeaderName {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.0)
    }
}

/// HTTP header value. Disallows CR/LF/NUL to prevent header injection.
#[derive(Debug, Clone, Eq, PartialEq, Hash)]
pub struct HeaderValue(String);

impl HeaderValue {
    pub fn new(s: impl Into<String>) -> Result<Self, DomainError> {
        let owned: String = s.into();
        if owned.bytes().any(|b| b == b'\r' || b == b'\n' || b == 0) {
            return Err(DomainError::InvalidHeaderValue(owned));
        }
        Ok(Self(owned.trim().to_owned()))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }
}

impl fmt::Display for HeaderValue {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.0)
    }
}

/// Ordered, case-insensitive header collection.
///
/// Preserves insertion order so the wire form is reproducible. Multiple values
/// for the same name are kept as separate entries (HTTP 1.1 allows it for
/// `Set-Cookie` etc.); set/replace semantics are also provided.
#[derive(Debug, Clone, Default, Eq, PartialEq)]
pub struct Headers {
    entries: Vec<(HeaderName, HeaderValue)>,
}

impl Headers {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn len(&self) -> usize {
        self.entries.len()
    }

    pub fn is_empty(&self) -> bool {
        self.entries.is_empty()
    }

    /// Append a header without removing existing entries with the same name.
    pub fn append(&mut self, name: HeaderName, value: HeaderValue) {
        self.entries.push((name, value));
    }

    /// Replace any existing entries for `name` with a single new value.
    pub fn set(&mut self, name: HeaderName, value: HeaderValue) {
        self.entries.retain(|(n, _)| n != &name);
        self.entries.push((name, value));
    }

    /// First value for the given name (case-insensitive lookup).
    pub fn get(&self, name: &str) -> Option<&str> {
        let needle = name.to_ascii_lowercase();
        self.entries
            .iter()
            .find(|(n, _)| n.as_str() == needle)
            .map(|(_, v)| v.as_str())
    }

    /// All values for the given name in insertion order.
    pub fn get_all<'a>(&'a self, name: &str) -> impl Iterator<Item = &'a str> {
        let needle = name.to_ascii_lowercase();
        self.entries
            .iter()
            .filter(move |(n, _)| n.as_str() == needle)
            .map(|(_, v)| v.as_str())
    }

    pub fn contains(&self, name: &str) -> bool {
        self.get(name).is_some()
    }

    pub fn remove(&mut self, name: &str) {
        let needle = name.to_ascii_lowercase();
        self.entries.retain(|(n, _)| n.as_str() != needle);
    }

    pub fn iter(&self) -> impl Iterator<Item = (&HeaderName, &HeaderValue)> {
        self.entries.iter().map(|(n, v)| (n, v))
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn n(s: &str) -> HeaderName {
        HeaderName::new(s).unwrap()
    }

    fn v(s: &str) -> HeaderValue {
        HeaderValue::new(s).unwrap()
    }

    #[test]
    fn name_validation() {
        assert!(HeaderName::new("Content-Type").is_ok());
        assert_eq!(
            HeaderName::new("Content-Type").unwrap().as_str(),
            "content-type"
        );
        assert!(HeaderName::new("").is_err());
        assert!(HeaderName::new("Bad Header").is_err());
        assert!(HeaderName::new("X-Header\n").is_err());
    }

    #[test]
    fn value_rejects_crlf_and_nul() {
        assert!(HeaderValue::new("ok value").is_ok());
        assert!(HeaderValue::new("bad\r\ninjection").is_err());
        assert!(HeaderValue::new("bad\nvalue").is_err());
        assert!(HeaderValue::new("bad\0value").is_err());
    }

    #[test]
    fn value_trims_whitespace() {
        assert_eq!(
            HeaderValue::new("  text/plain  ").unwrap().as_str(),
            "text/plain"
        );
    }

    #[test]
    fn case_insensitive_get() {
        let mut h = Headers::new();
        h.append(n("Content-Type"), v("text/html"));
        assert_eq!(h.get("content-type"), Some("text/html"));
        assert_eq!(h.get("CONTENT-TYPE"), Some("text/html"));
    }

    #[test]
    fn append_preserves_order_and_multiplicity() {
        let mut h = Headers::new();
        h.append(n("Set-Cookie"), v("a=1"));
        h.append(n("Set-Cookie"), v("b=2"));
        let all: Vec<&str> = h.get_all("set-cookie").collect();
        assert_eq!(all, vec!["a=1", "b=2"]);
    }

    #[test]
    fn set_replaces_all_existing() {
        let mut h = Headers::new();
        h.append(n("X"), v("1"));
        h.append(n("X"), v("2"));
        h.set(n("X"), v("3"));
        let all: Vec<&str> = h.get_all("x").collect();
        assert_eq!(all, vec!["3"]);
    }

    #[test]
    fn remove_clears_all() {
        let mut h = Headers::new();
        h.append(n("X"), v("1"));
        h.append(n("X"), v("2"));
        h.remove("X");
        assert!(h.get("x").is_none());
        assert!(h.is_empty());
    }
}
