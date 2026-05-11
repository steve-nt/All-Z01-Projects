//! In-memory `SessionStore` with TTL eviction on access.

use std::collections::HashMap;

use crate::application::ports::session_store::SessionStore;
use crate::domain::session::session::Session;
use crate::domain::session::session_id::SessionId;

#[derive(Debug, Default)]
pub struct MemorySessionStore {
    sessions: HashMap<String, Session>,
}

impl MemorySessionStore {
    pub fn new() -> Self {
        Self::default()
    }
}

impl SessionStore for MemorySessionStore {
    fn get(&self, id: &SessionId) -> Option<Session> {
        self.sessions.get(id.as_str()).cloned()
    }

    fn put(&mut self, session: Session) {
        self.sessions
            .insert(session.id().as_str().to_owned(), session);
    }

    fn remove(&mut self, id: &SessionId) {
        self.sessions.remove(id.as_str());
    }

    fn evict_expired(&mut self, now_secs: u64, max_idle_secs: u64) {
        self.sessions
            .retain(|_, s| !s.is_expired(now_secs, max_idle_secs));
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn sid(s: &str) -> SessionId {
        SessionId::new(s).unwrap()
    }

    #[test]
    fn round_trip() {
        let mut store = MemorySessionStore::new();
        let id = sid("abcdefghijklmnop");
        let mut s = Session::new(id.clone(), 0);
        s.set("key", "val");
        store.put(s);
        assert_eq!(store.get(&id).unwrap().get("key"), Some("val"));
    }

    #[test]
    fn put_replaces_existing() {
        let mut store = MemorySessionStore::new();
        let id = sid("abcdefghijklmnop");
        let mut s1 = Session::new(id.clone(), 0);
        s1.set("k", "v1");
        store.put(s1);
        let mut s2 = Session::new(id.clone(), 1);
        s2.set("k", "v2");
        store.put(s2);
        assert_eq!(store.get(&id).unwrap().get("k"), Some("v2"));
    }

    #[test]
    fn evict_removes_idle_sessions() {
        let mut store = MemorySessionStore::new();
        let id = sid("abcdefghijklmnop");
        store.put(Session::new(id.clone(), 0));
        store.evict_expired(10_000, 5_000);
        assert!(store.get(&id).is_none());
    }

    #[test]
    fn evict_keeps_fresh_sessions() {
        let mut store = MemorySessionStore::new();
        let id = sid("abcdefghijklmnop");
        store.put(Session::new(id.clone(), 9_950));
        store.evict_expired(10_000, 5_000);
        assert!(store.get(&id).is_some());
    }
}
