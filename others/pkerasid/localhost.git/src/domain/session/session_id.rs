use std::fmt;

use crate::domain::error::DomainError;

/// Opaque session identifier.
///
/// Restricted to URL-safe characters so it can be carried in a Cookie value
/// without quoting. Length range protects against trivial enumeration and
/// nonsense input.
#[derive(Debug, Clone, Eq, PartialEq, Hash)]
pub struct SessionId(String);

impl SessionId {
    pub const MIN_LEN: usize = 16;
    pub const MAX_LEN: usize = 128;

    pub fn new(s: impl Into<String>) -> Result<Self, DomainError> {
        let owned: String = s.into();
        if !(Self::MIN_LEN..=Self::MAX_LEN).contains(&owned.len()) {
            return Err(DomainError::InvalidSessionId);
        }
        if !owned.bytes().all(Self::is_valid_byte) {
            return Err(DomainError::InvalidSessionId);
        }
        Ok(Self(owned))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }

    fn is_valid_byte(b: u8) -> bool {
        matches!(b, b'A'..=b'Z' | b'a'..=b'z' | b'0'..=b'9' | b'-' | b'_')
    }
}

impl fmt::Display for SessionId {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.0)
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn accepts_url_safe() {
        assert!(SessionId::new("abcDEF123456-_xy").is_ok());
    }

    #[test]
    fn rejects_too_short() {
        assert!(SessionId::new("short").is_err());
    }

    #[test]
    fn rejects_too_long() {
        let s: String = (0..200).map(|_| 'a').collect();
        assert!(SessionId::new(s).is_err());
    }

    #[test]
    fn rejects_bad_chars() {
        assert!(SessionId::new("abcdefghijklmnop!").is_err());
        assert!(SessionId::new("abcdefghijklmnop ").is_err());
        assert!(SessionId::new("abcdefghijklmnop/").is_err());
    }
}
