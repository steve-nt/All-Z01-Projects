use std::collections::HashMap;

use crate::domain::session::session_id::SessionId;

/// In-memory session record. Time-based concerns (creation, expiry) are
/// modelled as plain integers so the domain stays free of any clock dependency.
#[derive(Debug, Clone, Eq, PartialEq)]
pub struct Session {
    id: SessionId,
    created_at_secs: u64,
    last_seen_secs: u64,
    data: HashMap<String, String>,
}

impl Session {
    pub fn new(id: SessionId, now_secs: u64) -> Self {
        Self {
            id,
            created_at_secs: now_secs,
            last_seen_secs: now_secs,
            data: HashMap::new(),
        }
    }

    pub fn id(&self) -> &SessionId {
        &self.id
    }

    pub fn created_at_secs(&self) -> u64 {
        self.created_at_secs
    }

    pub fn last_seen_secs(&self) -> u64 {
        self.last_seen_secs
    }

    pub fn touch(&mut self, now_secs: u64) {
        self.last_seen_secs = now_secs;
    }

    pub fn get(&self, key: &str) -> Option<&str> {
        self.data.get(key).map(String::as_str)
    }

    pub fn set(&mut self, key: impl Into<String>, value: impl Into<String>) {
        self.data.insert(key.into(), value.into());
    }

    pub fn remove(&mut self, key: &str) -> Option<String> {
        self.data.remove(key)
    }

    /// True if this session has been idle for more than `max_idle_secs`.
    pub fn is_expired(&self, now_secs: u64, max_idle_secs: u64) -> bool {
        now_secs.saturating_sub(self.last_seen_secs) > max_idle_secs
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn sid() -> SessionId {
        SessionId::new("abcdefghijklmnop").unwrap()
    }

    #[test]
    fn touch_updates_last_seen() {
        let mut s = Session::new(sid(), 100);
        assert_eq!(s.last_seen_secs(), 100);
        s.touch(150);
        assert_eq!(s.last_seen_secs(), 150);
    }

    #[test]
    fn data_round_trip() {
        let mut s = Session::new(sid(), 0);
        s.set("user", "alice");
        assert_eq!(s.get("user"), Some("alice"));
        assert_eq!(s.remove("user").as_deref(), Some("alice"));
        assert_eq!(s.get("user"), None);
    }

    #[test]
    fn expiry_check() {
        let s = Session::new(sid(), 1000);
        assert!(!s.is_expired(1100, 200));
        assert!(s.is_expired(1300, 200));
        // Clock-skew safe: now < last_seen returns not expired.
        assert!(!s.is_expired(500, 200));
    }
}
