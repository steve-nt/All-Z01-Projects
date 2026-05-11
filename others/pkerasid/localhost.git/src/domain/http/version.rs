use std::fmt;
use std::str::FromStr;

use crate::domain::error::DomainError;

#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash)]
pub enum HttpVersion {
    Http10,
    Http11,
}

impl HttpVersion {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Http10 => "HTTP/1.0",
            Self::Http11 => "HTTP/1.1",
        }
    }
}

impl FromStr for HttpVersion {
    type Err = DomainError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "HTTP/1.0" => Ok(Self::Http10),
            "HTTP/1.1" => Ok(Self::Http11),
            other => Err(DomainError::InvalidVersion(other.to_owned())),
        }
    }
}

impl fmt::Display for HttpVersion {
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
        assert_eq!(
            "HTTP/1.0".parse::<HttpVersion>().ok(),
            Some(HttpVersion::Http10)
        );
        assert_eq!(
            "HTTP/1.1".parse::<HttpVersion>().ok(),
            Some(HttpVersion::Http11)
        );
    }

    #[test]
    fn rejects_others() {
        assert!("HTTP/2.0".parse::<HttpVersion>().is_err());
        assert!("http/1.1".parse::<HttpVersion>().is_err());
        assert!("".parse::<HttpVersion>().is_err());
    }
}
