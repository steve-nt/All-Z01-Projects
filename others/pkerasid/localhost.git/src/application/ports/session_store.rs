//! SessionStore port: interface for session persistence.
//!
//! Implementations live in `infrastructure/session_store/`. A `FakeSessionStore`
//! is provided for unit tests in this module.

use std::collections::HashMap;

use crate::domain::session::session::Session;
use crate::domain::session::session_id::SessionId;

/// Persistence layer for HTTP sessions.
///
/// `get` returns a clone so callers may mutate the session and `put` it back
/// without holding a borrow on the store.
pub trait SessionStore {
    /// Retrieve a session by ID. Returns `None` if not found.
    fn get(&self, id: &SessionId) -> Option<Session>;

    /// Insert or replace a session.
    fn put(&mut self, session: Session);

    /// Remove a session. No-op if the ID is not found.
    fn remove(&mut self, id: &SessionId);

    /// Remove all sessions idle for longer than `max_idle_secs`.
    fn evict_expired(&mut self, now_secs: u64, max_idle_secs: u64);
}

/// In-test session store backed by a `HashMap`.
#[derive(Debug, Default)]
pub struct FakeSessionStore {
    sessions: HashMap<String, Session>,
}

impl FakeSessionStore {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn len(&self) -> usize {
        self.sessions.len()
    }

    pub fn is_empty(&self) -> bool {
        self.sessions.is_empty()
    }
}

impl SessionStore for FakeSessionStore {
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

    fn make(id: &str, last_seen: u64) -> Session {
        Session::new(sid(id), last_seen)
    }

    #[test]
    fn put_and_get_round_trip() {
        let mut store = FakeSessionStore::new();
        let id = sid("abcdefghijklmnop");
        store.put(Session::new(id.clone(), 100));
        let s = store.get(&id).unwrap();
        assert_eq!(s.id(), &id);
    }

    #[test]
    fn get_missing_returns_none() {
        let store = FakeSessionStore::new();
        assert!(store.get(&sid("abcdefghijklmnop")).is_none());
    }

    #[test]
    fn put_replaces_existing() {
        let mut store = FakeSessionStore::new();
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
    fn remove_works() {
        let mut store = FakeSessionStore::new();
        let id = sid("abcdefghijklmnop");
        store.put(Session::new(id.clone(), 0));
        store.remove(&id);
        assert!(store.get(&id).is_none());
    }

    #[test]
    fn evict_removes_idle_sessions() {
        let mut store = FakeSessionStore::new();
        // last_seen = 0; now = 10000; idle = 10000 > 5000 → evicted
        store.put(make("aaaaaaaaaaaaaaaa", 0));
        // last_seen = 9950; now = 10000; idle = 50 < 5000 → kept
        store.put(make("bbbbbbbbbbbbbbbb", 9_950));
        store.evict_expired(10_000, 5_000);
        assert!(store.get(&sid("aaaaaaaaaaaaaaaa")).is_none());
        assert!(store.get(&sid("bbbbbbbbbbbbbbbb")).is_some());
    }
}
