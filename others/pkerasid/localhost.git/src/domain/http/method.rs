use std::fmt;
use std::str::FromStr;

use crate::domain::error::DomainError;

/// HTTP method.
///
/// The audit only requires GET, POST, DELETE. Other tokens are still parsed
/// (kept as-is) so the request can be answered with `405 Method Not Allowed`
/// instead of `400 Bad Request`.
#[derive(Debug, Clone, Eq, PartialEq, Hash)]
pub enum Method {
    Get,
    Post,
    Delete,
    Other(String),
}

impl Method {
    pub fn as_str(&self) -> &str {
        match self {
            Self::Get => "GET",
            Self::Post => "POST",
            Self::Delete => "DELETE",
            Self::Other(s) => s,
        }
    }

    pub fn is_supported(&self) -> bool {
        matches!(self, Self::Get | Self::Post | Self::Delete)
    }

    /// RFC 7230 token: 1*tchar.
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

impl FromStr for Method {
    type Err = DomainError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        if !Self::is_token(s) {
            return Err(DomainError::InvalidMethod(s.to_owned()));
        }
        Ok(match s {
            "GET" => Self::Get,
            "POST" => Self::Post,
            "DELETE" => Self::Delete,
            _ => Self::Other(s.to_owned()),
        })
    }
}

impl fmt::Display for Method {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn parses_known() {
        assert_eq!("GET".parse::<Method>().ok(), Some(Method::Get));
        assert_eq!("POST".parse::<Method>().ok(), Some(Method::Post));
        assert_eq!("DELETE".parse::<Method>().ok(), Some(Method::Delete));
    }

    #[test]
    fn parses_unknown_as_other() {
        assert_eq!(
            "PATCH".parse::<Method>().ok(),
            Some(Method::Other("PATCH".into()))
        );
    }

    #[test]
    fn rejects_empty_and_invalid() {
        assert!("".parse::<Method>().is_err());
        assert!("GE T".parse::<Method>().is_err());
        assert!("GET\n".parse::<Method>().is_err());
    }

    #[test]
    fn supported_check() {
        assert!(Method::Get.is_supported());
        assert!(!Method::Other("PATCH".into()).is_supported());
    }

    #[test]
    fn case_sensitive() {
        // RFC 7230: methods are case-sensitive; "get" != "GET".
        assert_eq!(
            "get".parse::<Method>().ok(),
            Some(Method::Other("get".into()))
        );
    }
}
