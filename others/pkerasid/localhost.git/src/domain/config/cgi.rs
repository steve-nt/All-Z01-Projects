use std::path::PathBuf;

use crate::domain::error::DomainError;

/// CGI handler configuration: which extensions to dispatch to which interpreter.
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct CgiConfig {
    extension: String,
    interpreter: PathBuf,
}

impl CgiConfig {
    pub fn new(
        extension: impl Into<String>,
        interpreter: impl Into<PathBuf>,
    ) -> Result<Self, DomainError> {
        let ext: String = extension.into();
        let interp: PathBuf = interpreter.into();
        if ext.is_empty() {
            return Err(DomainError::InvalidConfig("cgi extension is empty".into()));
        }
        if !ext.starts_with('.') {
            return Err(DomainError::InvalidConfig(format!(
                "cgi extension must start with '.': {ext:?}"
            )));
        }
        if ext.len() < 2 {
            return Err(DomainError::InvalidConfig(format!(
                "cgi extension is missing characters after '.': {ext:?}"
            )));
        }
        if interp.as_os_str().is_empty() {
            return Err(DomainError::InvalidConfig(
                "cgi interpreter path is empty".into(),
            ));
        }
        Ok(Self {
            extension: ext,
            interpreter: interp,
        })
    }

    pub fn extension(&self) -> &str {
        &self.extension
    }

    pub fn interpreter(&self) -> &std::path::Path {
        &self.interpreter
    }

    /// Returns true if `path` matches this CGI's extension (case-sensitive).
    pub fn matches(&self, path: &str) -> bool {
        path.ends_with(&self.extension)
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn rejects_empty() {
        assert!(CgiConfig::new("", "/usr/bin/python3").is_err());
        assert!(CgiConfig::new(".py", "").is_err());
    }

    #[test]
    fn extension_must_start_with_dot() {
        assert!(CgiConfig::new("py", "/usr/bin/python3").is_err());
        assert!(CgiConfig::new(".py", "/usr/bin/python3").is_ok());
    }

    #[test]
    fn matches_works() {
        let c = CgiConfig::new(".py", "/usr/bin/python3").unwrap();
        assert!(c.matches("/scripts/hello.py"));
        assert!(!c.matches("/scripts/hello.txt"));
    }
}
