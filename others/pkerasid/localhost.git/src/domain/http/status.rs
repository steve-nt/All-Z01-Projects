use std::fmt;

use crate::domain::error::DomainError;

/// HTTP status code with canonical reason phrase.
#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash)]
pub struct Status(u16);

impl Status {
    pub const OK: Self = Self(200);
    pub const CREATED: Self = Self(201);
    pub const NO_CONTENT: Self = Self(204);
    pub const PARTIAL_CONTENT: Self = Self(206);
    pub const MOVED_PERMANENTLY: Self = Self(301);
    pub const FOUND: Self = Self(302);
    pub const SEE_OTHER: Self = Self(303);
    pub const NOT_MODIFIED: Self = Self(304);
    pub const TEMPORARY_REDIRECT: Self = Self(307);
    pub const PERMANENT_REDIRECT: Self = Self(308);
    pub const BAD_REQUEST: Self = Self(400);
    pub const UNAUTHORIZED: Self = Self(401);
    pub const FORBIDDEN: Self = Self(403);
    pub const NOT_FOUND: Self = Self(404);
    pub const METHOD_NOT_ALLOWED: Self = Self(405);
    pub const REQUEST_TIMEOUT: Self = Self(408);
    pub const LENGTH_REQUIRED: Self = Self(411);
    pub const PAYLOAD_TOO_LARGE: Self = Self(413);
    pub const URI_TOO_LONG: Self = Self(414);
    pub const UNSUPPORTED_MEDIA_TYPE: Self = Self(415);
    pub const EXPECTATION_FAILED: Self = Self(417);
    pub const INTERNAL_SERVER_ERROR: Self = Self(500);
    pub const NOT_IMPLEMENTED: Self = Self(501);
    pub const BAD_GATEWAY: Self = Self(502);
    pub const SERVICE_UNAVAILABLE: Self = Self(503);
    pub const GATEWAY_TIMEOUT: Self = Self(504);
    pub const HTTP_VERSION_NOT_SUPPORTED: Self = Self(505);

    /// Constructs a status from a numeric code; rejects anything outside 100..=599.
    pub fn new(code: u16) -> Result<Self, DomainError> {
        if (100..=599).contains(&code) {
            Ok(Self(code))
        } else {
            Err(DomainError::InvalidStatus(code))
        }
    }

    pub fn code(self) -> u16 {
        self.0
    }

    pub fn is_informational(self) -> bool {
        (100..200).contains(&self.0)
    }

    pub fn is_success(self) -> bool {
        (200..300).contains(&self.0)
    }

    pub fn is_redirection(self) -> bool {
        (300..400).contains(&self.0)
    }

    pub fn is_client_error(self) -> bool {
        (400..500).contains(&self.0)
    }

    pub fn is_server_error(self) -> bool {
        (500..600).contains(&self.0)
    }

    pub fn is_error(self) -> bool {
        self.is_client_error() || self.is_server_error()
    }

    pub fn reason(self) -> &'static str {
        match self.0 {
            200 => "OK",
            201 => "Created",
            204 => "No Content",
            206 => "Partial Content",
            301 => "Moved Permanently",
            302 => "Found",
            303 => "See Other",
            304 => "Not Modified",
            307 => "Temporary Redirect",
            308 => "Permanent Redirect",
            400 => "Bad Request",
            401 => "Unauthorized",
            403 => "Forbidden",
            404 => "Not Found",
            405 => "Method Not Allowed",
            408 => "Request Timeout",
            411 => "Length Required",
            413 => "Payload Too Large",
            414 => "URI Too Long",
            415 => "Unsupported Media Type",
            417 => "Expectation Failed",
            500 => "Internal Server Error",
            501 => "Not Implemented",
            502 => "Bad Gateway",
            503 => "Service Unavailable",
            504 => "Gateway Timeout",
            505 => "HTTP Version Not Supported",
            _ => "",
        }
    }
}

impl fmt::Display for Status {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let r = self.reason();
        if r.is_empty() {
            write!(f, "{}", self.0)
        } else {
            write!(f, "{} {r}", self.0)
        }
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn constants_match_codes() {
        assert_eq!(Status::OK.code(), 200);
        assert_eq!(Status::NOT_FOUND.code(), 404);
        assert_eq!(Status::PAYLOAD_TOO_LARGE.code(), 413);
        assert_eq!(Status::INTERNAL_SERVER_ERROR.code(), 500);
    }

    #[test]
    fn reason_phrases() {
        assert_eq!(Status::OK.reason(), "OK");
        assert_eq!(Status::METHOD_NOT_ALLOWED.reason(), "Method Not Allowed");
        assert_eq!(Status::PAYLOAD_TOO_LARGE.reason(), "Payload Too Large");
    }

    #[test]
    fn classification() {
        assert!(Status::OK.is_success());
        assert!(Status::FOUND.is_redirection());
        assert!(Status::NOT_FOUND.is_client_error());
        assert!(Status::INTERNAL_SERVER_ERROR.is_server_error());
        assert!(Status::FORBIDDEN.is_error());
        assert!(!Status::OK.is_error());
    }

    #[test]
    fn rejects_out_of_range() {
        assert!(Status::new(99).is_err());
        assert!(Status::new(600).is_err());
        assert!(Status::new(200).is_ok());
        assert!(Status::new(599).is_ok());
    }

    #[test]
    fn unknown_code_displays_number_only() {
        let s = Status::new(299).expect("valid range");
        assert_eq!(format!("{s}"), "299");
    }
}
