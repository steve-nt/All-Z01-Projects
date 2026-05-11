//! Session ID generator using /dev/urandom entropy.

use std::io::Read as _;

/// Generate a 64-character hex session ID (256 bits of entropy from /dev/urandom).
pub fn generate() -> String {
    let mut buf = [0u8; 32];
    if let Ok(mut f) = std::fs::File::open("/dev/urandom") {
        let _ = f.read_exact(&mut buf);
    }
    buf.iter().fold(String::with_capacity(64), |mut s, b| {
        use std::fmt::Write as _;
        let _ = write!(s, "{b:02x}");
        s
    })
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used)]

    use super::*;
    use crate::domain::session::session_id::SessionId;

    #[test]
    fn output_is_valid_session_id() {
        let id_str = generate();
        assert_eq!(id_str.len(), 64);
        assert!(SessionId::new(id_str).is_ok());
    }

    #[test]
    fn two_calls_are_different() {
        assert_ne!(generate(), generate());
    }
}
