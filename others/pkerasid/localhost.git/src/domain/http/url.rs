use std::fmt;

use crate::domain::error::DomainError;

/// Percent-decode a URL segment per RFC 3986.
///
/// Leaves non-percent bytes untouched. Returns `InvalidPercentEncoding` for
/// truncated or non-hex sequences. Decoded output is validated as UTF-8.
pub fn percent_decode(input: &str) -> Result<String, DomainError> {
    let bytes = input.as_bytes();
    let mut out: Vec<u8> = Vec::with_capacity(bytes.len());
    let mut i = 0;
    while i < bytes.len() {
        let b = bytes[i];
        if b == b'%' {
            if i + 2 >= bytes.len() {
                return Err(DomainError::InvalidPercentEncoding);
            }
            let hi = hex_value(bytes[i + 1]).ok_or(DomainError::InvalidPercentEncoding)?;
            let lo = hex_value(bytes[i + 2]).ok_or(DomainError::InvalidPercentEncoding)?;
            out.push((hi << 4) | lo);
            i += 3;
        } else {
            out.push(b);
            i += 1;
        }
    }
    String::from_utf8(out).map_err(|_| DomainError::InvalidPercentEncoding)
}

fn hex_value(b: u8) -> Option<u8> {
    match b {
        b'0'..=b'9' => Some(b - b'0'),
        b'a'..=b'f' => Some(b - b'a' + 10),
        b'A'..=b'F' => Some(b - b'A' + 10),
        _ => None,
    }
}

/// Normalized request path.
///
/// - Must start with `/`.
/// - Percent-decoded, then split on `/`.
/// - `.` segments are dropped, `..` pops the parent.
/// - A `..` that would escape the root yields `PathTraversal`.
/// - NUL bytes are rejected.
#[derive(Debug, Clone, Eq, PartialEq, Hash)]
pub struct RequestPath {
    segments: Vec<String>,
    /// Original raw path (post percent-decode), with leading `/` and no `.`/`..`.
    normalized: String,
}

impl RequestPath {
    pub fn parse(raw: &str) -> Result<Self, DomainError> {
        if !raw.starts_with('/') {
            return Err(DomainError::InvalidPath(raw.to_owned()));
        }
        let decoded = percent_decode(raw)?;
        if decoded.as_bytes().contains(&0) {
            return Err(DomainError::InvalidPath(raw.to_owned()));
        }
        let mut segments: Vec<String> = Vec::new();
        for part in decoded.split('/') {
            match part {
                "" | "." => {}
                ".." => {
                    if segments.pop().is_none() {
                        return Err(DomainError::PathTraversal);
                    }
                }
                seg => segments.push(seg.to_owned()),
            }
        }
        let mut normalized = String::with_capacity(decoded.len());
        if segments.is_empty() {
            normalized.push('/');
        } else {
            for seg in &segments {
                normalized.push('/');
                normalized.push_str(seg);
            }
            // Preserve trailing slash (significant for directory routing).
            if decoded.ends_with('/') {
                normalized.push('/');
            }
        }
        Ok(Self {
            segments,
            normalized,
        })
    }

    pub fn segments(&self) -> &[String] {
        &self.segments
    }

    pub fn as_str(&self) -> &str {
        &self.normalized
    }

    pub fn ends_with_slash(&self) -> bool {
        self.normalized.ends_with('/')
    }
}

impl fmt::Display for RequestPath {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.normalized)
    }
}

/// URL query string (raw, validated). Decoding into key/value pairs is the
/// responsibility of upper layers.
#[derive(Debug, Clone, Eq, PartialEq, Hash)]
pub struct Query(String);

impl Query {
    pub fn new(s: impl Into<String>) -> Result<Self, DomainError> {
        let owned: String = s.into();
        if owned.bytes().any(|b| b == 0 || b == b'\r' || b == b'\n') {
            return Err(DomainError::InvalidQuery(owned));
        }
        Ok(Self(owned))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn percent_decode_basics() {
        assert_eq!(percent_decode("hello").unwrap(), "hello");
        assert_eq!(percent_decode("a%20b").unwrap(), "a b");
        assert_eq!(percent_decode("%2Fpath").unwrap(), "/path");
        assert_eq!(percent_decode("%C3%A9").unwrap(), "é");
    }

    #[test]
    fn percent_decode_errors() {
        assert!(percent_decode("%2").is_err());
        assert!(percent_decode("%ZZ").is_err());
        assert!(percent_decode("%FF").is_err()); // invalid utf-8
    }

    #[test]
    fn path_must_be_absolute() {
        assert!(RequestPath::parse("relative").is_err());
        assert!(RequestPath::parse("/").is_ok());
    }

    #[test]
    fn path_normalizes_dots() {
        assert_eq!(RequestPath::parse("/a/./b").unwrap().as_str(), "/a/b");
        assert_eq!(RequestPath::parse("/a/b/../c").unwrap().as_str(), "/a/c");
        assert_eq!(RequestPath::parse("/a//b").unwrap().as_str(), "/a/b");
    }

    #[test]
    fn path_traversal_blocked() {
        assert!(matches!(
            RequestPath::parse("/.."),
            Err(DomainError::PathTraversal)
        ));
        assert!(matches!(
            RequestPath::parse("/a/../.."),
            Err(DomainError::PathTraversal)
        ));
    }

    #[test]
    fn path_preserves_trailing_slash() {
        assert!(RequestPath::parse("/dir/").unwrap().ends_with_slash());
        assert!(!RequestPath::parse("/dir/file").unwrap().ends_with_slash());
    }

    #[test]
    fn path_decodes_percent() {
        assert_eq!(RequestPath::parse("/a%20b").unwrap().as_str(), "/a b");
    }

    #[test]
    fn path_rejects_nul() {
        assert!(RequestPath::parse("/a%00b").is_err());
    }

    #[test]
    fn query_validation() {
        assert!(Query::new("a=1&b=2").is_ok());
        assert!(Query::new("bad\nquery").is_err());
    }
}
