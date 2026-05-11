use std::fmt;

use crate::domain::error::DomainError;

/// `SameSite` attribute for `Set-Cookie`.
#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash)]
pub enum SameSite {
    Strict,
    Lax,
    None,
}

impl SameSite {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Strict => "Strict",
            Self::Lax => "Lax",
            Self::None => "None",
        }
    }
}

/// A single HTTP cookie (request side: name+value; response side adds attributes).
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct Cookie {
    name: String,
    value: String,
    path: Option<String>,
    domain: Option<String>,
    max_age: Option<i64>,
    http_only: bool,
    secure: bool,
    same_site: Option<SameSite>,
}

impl Cookie {
    /// Build a name=value cookie. Names follow the RFC 6265 token rules; values
    /// allow most printable ASCII excluding control characters, comma, semicolon,
    /// and double-quote (which would need DQUOTE-wrapping not handled here).
    pub fn new(name: impl Into<String>, value: impl Into<String>) -> Result<Self, DomainError> {
        let name: String = name.into();
        let value: String = value.into();
        if !is_token(&name) {
            return Err(DomainError::InvalidCookieName(name));
        }
        if !is_cookie_value(&value) {
            return Err(DomainError::InvalidCookieValue(value));
        }
        Ok(Self {
            name,
            value,
            path: None,
            domain: None,
            max_age: None,
            http_only: false,
            secure: false,
            same_site: None,
        })
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn value(&self) -> &str {
        &self.value
    }

    #[must_use]
    pub fn with_path(mut self, p: impl Into<String>) -> Self {
        self.path = Some(p.into());
        self
    }

    #[must_use]
    pub fn with_domain(mut self, d: impl Into<String>) -> Self {
        self.domain = Some(d.into());
        self
    }

    #[must_use]
    pub fn with_max_age(mut self, secs: i64) -> Self {
        self.max_age = Some(secs);
        self
    }

    #[must_use]
    pub fn http_only(mut self) -> Self {
        self.http_only = true;
        self
    }

    #[must_use]
    pub fn secure(mut self) -> Self {
        self.secure = true;
        self
    }

    #[must_use]
    pub fn with_same_site(mut self, ss: SameSite) -> Self {
        self.same_site = Some(ss);
        self
    }

    /// Render as a `Set-Cookie` value (does not include the header name).
    pub fn render_set_cookie(&self) -> String {
        let mut s = format!("{}={}", self.name, self.value);
        if let Some(p) = &self.path {
            s.push_str("; Path=");
            s.push_str(p);
        }
        if let Some(d) = &self.domain {
            s.push_str("; Domain=");
            s.push_str(d);
        }
        if let Some(a) = self.max_age {
            s.push_str("; Max-Age=");
            s.push_str(&a.to_string());
        }
        if self.http_only {
            s.push_str("; HttpOnly");
        }
        if self.secure {
            s.push_str("; Secure");
        }
        if let Some(ss) = self.same_site {
            s.push_str("; SameSite=");
            s.push_str(ss.as_str());
        }
        s
    }
}

impl fmt::Display for Cookie {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}={}", self.name, self.value)
    }
}

/// Parse the value of an inbound `Cookie:` header into name/value pairs.
///
/// Robust against optional whitespace and missing semicolons; ignores malformed
/// pairs rather than failing the whole request.
pub fn parse_cookie_header(value: &str) -> Vec<(String, String)> {
    value
        .split(';')
        .filter_map(|part| {
            let trimmed = part.trim();
            if trimmed.is_empty() {
                return None;
            }
            let (n, v) = trimmed.split_once('=')?;
            let n = n.trim();
            let v = v.trim();
            if !is_token(n) || !is_cookie_value(v) {
                return None;
            }
            Some((n.to_owned(), v.to_owned()))
        })
        .collect()
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

fn is_cookie_value(s: &str) -> bool {
    s.bytes()
        .all(|b| matches!(b, 0x21 | 0x23..=0x2B | 0x2D..=0x3A | 0x3C..=0x5B | 0x5D..=0x7E))
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn rejects_bad_name() {
        assert!(Cookie::new("", "v").is_err());
        assert!(Cookie::new("bad name", "v").is_err());
    }

    #[test]
    fn rejects_bad_value() {
        assert!(Cookie::new("k", "with;semi").is_err());
        assert!(Cookie::new("k", "with,comma").is_err());
        assert!(Cookie::new("k", "with\"quote").is_err());
    }

    #[test]
    fn render_simple() {
        let c = Cookie::new("SID", "abc123").unwrap();
        assert_eq!(c.render_set_cookie(), "SID=abc123");
    }

    #[test]
    fn render_with_attrs() {
        let c = Cookie::new("SID", "abc123")
            .unwrap()
            .with_path("/")
            .with_max_age(3600)
            .http_only()
            .secure()
            .with_same_site(SameSite::Lax);
        assert_eq!(
            c.render_set_cookie(),
            "SID=abc123; Path=/; Max-Age=3600; HttpOnly; Secure; SameSite=Lax"
        );
    }

    #[test]
    fn parse_request_cookie_header() {
        let pairs = parse_cookie_header("SID=abc123; theme=dark; bad");
        assert_eq!(
            pairs,
            vec![
                ("SID".into(), "abc123".into()),
                ("theme".into(), "dark".into()),
            ]
        );
    }

    #[test]
    fn parse_handles_extra_whitespace() {
        let pairs = parse_cookie_header("  SID = abc123 ;theme=dark  ");
        assert_eq!(
            pairs,
            vec![
                ("SID".into(), "abc123".into()),
                ("theme".into(), "dark".into()),
            ]
        );
    }
}
