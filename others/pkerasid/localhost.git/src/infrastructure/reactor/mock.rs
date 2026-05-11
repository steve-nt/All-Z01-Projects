//! `MockReactor` — in-memory reactor used by upper-layer tests.
//!
//! Records register/reregister/deregister calls and lets tests inject
//! readiness events. No syscalls. Works on any platform.

use std::collections::HashMap;
use std::io;
use std::os::fd::RawFd;

use crate::application::ports::reactor::{Event, Interest, Reactor, Token};

#[derive(Debug, Clone, Eq, PartialEq)]
pub struct Registration {
    pub token: Token,
    pub interest: Interest,
}

#[derive(Debug, Default)]
pub struct MockReactor {
    pub registered: HashMap<RawFd, Registration>,
    queued: Vec<Event>,
    /// Total number of `poll` calls. Lets tests assert "single poll per
    /// iteration" by snapshotting this around an iteration.
    pub poll_calls: u64,
}

impl MockReactor {
    pub fn new() -> Self {
        Self::default()
    }

    /// Push an event that the next `poll` will return.
    pub fn push_event(&mut self, ev: Event) {
        self.queued.push(ev);
    }
}

impl Reactor for MockReactor {
    fn register(&mut self, fd: RawFd, token: Token, interest: Interest) -> io::Result<()> {
        if self.registered.contains_key(&fd) {
            return Err(io::Error::new(
                io::ErrorKind::AlreadyExists,
                "fd already registered",
            ));
        }
        self.registered.insert(fd, Registration { token, interest });
        Ok(())
    }

    fn reregister(&mut self, fd: RawFd, token: Token, interest: Interest) -> io::Result<()> {
        match self.registered.get_mut(&fd) {
            Some(r) => {
                r.token = token;
                r.interest = interest;
                Ok(())
            }
            None => Err(io::Error::new(io::ErrorKind::NotFound, "fd not registered")),
        }
    }

    fn deregister(&mut self, fd: RawFd) -> io::Result<()> {
        self.registered.remove(&fd);
        Ok(())
    }

    fn poll(&mut self, out: &mut Vec<Event>, _timeout_ms: Option<i32>) -> io::Result<usize> {
        self.poll_calls += 1;
        let n = self.queued.len();
        out.append(&mut self.queued);
        Ok(n)
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn register_and_lookup() {
        let mut r = MockReactor::new();
        r.register(3, 1, Interest::READABLE).unwrap();
        assert_eq!(r.registered[&3].token, 1);
        assert!(r.registered[&3].interest.readable);
    }

    #[test]
    fn double_register_errors() {
        let mut r = MockReactor::new();
        r.register(3, 1, Interest::READABLE).unwrap();
        let e = r.register(3, 2, Interest::WRITABLE).unwrap_err();
        assert_eq!(e.kind(), io::ErrorKind::AlreadyExists);
    }

    #[test]
    fn reregister_updates() {
        let mut r = MockReactor::new();
        r.register(3, 1, Interest::READABLE).unwrap();
        r.reregister(3, 1, Interest::WRITABLE).unwrap();
        assert!(r.registered[&3].interest.writable);
        assert!(!r.registered[&3].interest.readable);
    }

    #[test]
    fn reregister_unknown_errors() {
        let mut r = MockReactor::new();
        let e = r.reregister(3, 1, Interest::READABLE).unwrap_err();
        assert_eq!(e.kind(), io::ErrorKind::NotFound);
    }

    #[test]
    fn deregister_idempotent() {
        let mut r = MockReactor::new();
        r.deregister(99).unwrap();
        r.register(3, 1, Interest::READABLE).unwrap();
        r.deregister(3).unwrap();
        assert!(!r.registered.contains_key(&3));
    }

    #[test]
    fn poll_drains_queued_events() {
        let mut r = MockReactor::new();
        r.push_event(Event {
            token: 5,
            readable: true,
            writable: false,
            error: false,
            hangup: false,
        });
        r.push_event(Event {
            token: 6,
            readable: false,
            writable: true,
            error: false,
            hangup: false,
        });
        let mut out = Vec::new();
        let n = r.poll(&mut out, Some(0)).unwrap();
        assert_eq!(n, 2);
        assert_eq!(out.len(), 2);
        assert_eq!(r.poll_calls, 1);
        // Second poll: queue is now empty.
        let n2 = r.poll(&mut out, Some(0)).unwrap();
        assert_eq!(n2, 0);
        assert_eq!(r.poll_calls, 2);
    }
}
