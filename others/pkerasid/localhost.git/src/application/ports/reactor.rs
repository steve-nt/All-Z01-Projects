//! Reactor port: the abstraction the event loop uses to wait for I/O.
//!
//! Implementations live in `infrastructure/reactor/` (Linux: epoll). A
//! `MockReactor` is provided for tests so application logic stays
//! OS-independent.
//!
//! The audit requires that "the server uses only one select (or equivalent) to
//! read the client requests and write answers." This trait encodes that
//! contract: every reachable I/O readiness event must come from a single call
//! to `poll`.

use std::io;
use std::os::fd::RawFd;

/// Caller-chosen identifier returned with every event for the registered fd.
pub type Token = usize;

/// Readiness interest flags. The reactor delivers events when the underlying
/// fd transitions to a state matching one of these.
#[derive(Debug, Copy, Clone, Eq, PartialEq, Default)]
pub struct Interest {
    pub readable: bool,
    pub writable: bool,
}

impl Interest {
    pub const READABLE: Self = Self {
        readable: true,
        writable: false,
    };
    pub const WRITABLE: Self = Self {
        readable: false,
        writable: true,
    };
    pub const READ_WRITE: Self = Self {
        readable: true,
        writable: true,
    };

    pub fn is_empty(self) -> bool {
        !self.readable && !self.writable
    }
}

/// A single readiness event delivered by the reactor.
///
/// Each bool maps directly to an independent epoll/kqueue flag bit; modelling
/// them as four fields keeps the wire-to-domain translation 1:1 and avoids a
/// bespoke bitflag type for this small surface.
#[derive(Debug, Copy, Clone, Eq, PartialEq)]
#[allow(clippy::struct_excessive_bools)]
pub struct Event {
    pub token: Token,
    pub readable: bool,
    pub writable: bool,
    /// `EPOLLERR` or kqueue `EV_ERROR` was set.
    pub error: bool,
    /// `EPOLLHUP` / `EPOLLRDHUP` — peer closed its half.
    pub hangup: bool,
}

impl Event {
    pub fn is_terminal(&self) -> bool {
        self.error || self.hangup
    }
}

/// I/O readiness multiplexer.
///
/// All implementations must satisfy: a single `poll` call returns readiness for
/// all currently registered fds, and the loop performs no other I/O wait
/// syscalls (audit invariant).
pub trait Reactor {
    /// Add `fd` to the interest set under `token`.
    fn register(&mut self, fd: RawFd, token: Token, interest: Interest) -> io::Result<()>;

    /// Replace the interest for an already-registered `fd`.
    fn reregister(&mut self, fd: RawFd, token: Token, interest: Interest) -> io::Result<()>;

    /// Remove `fd` from the interest set. Idempotent on best effort: callers
    /// may re-deregister an fd whose registration is unknown without error.
    fn deregister(&mut self, fd: RawFd) -> io::Result<()>;

    /// Block until at least one event is ready, or the timeout elapses.
    ///
    /// `timeout_ms = None` blocks indefinitely; `Some(0)` returns immediately.
    /// Events are appended into `out`. Returns the number of events written.
    fn poll(&mut self, out: &mut Vec<Event>, timeout_ms: Option<i32>) -> io::Result<usize>;
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn interest_constants() {
        const _: () = {
            assert!(Interest::READABLE.readable);
            assert!(!Interest::READABLE.writable);
            assert!(Interest::WRITABLE.writable);
            assert!(!Interest::WRITABLE.readable);
            assert!(Interest::READ_WRITE.readable && Interest::READ_WRITE.writable);
        };
        assert!(Interest::default().is_empty());
    }

    #[test]
    fn event_terminality() {
        let normal = Event {
            token: 1,
            readable: true,
            writable: false,
            error: false,
            hangup: false,
        };
        assert!(!normal.is_terminal());

        let err = Event {
            error: true,
            ..normal
        };
        assert!(err.is_terminal());

        let hup = Event {
            hangup: true,
            ..normal
        };
        assert!(hup.is_terminal());
    }
}
