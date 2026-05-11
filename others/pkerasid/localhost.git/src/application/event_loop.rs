//! Single-threaded event loop.
//!
//! Audit invariant: each iteration calls the reactor's `poll` exactly once;
//! all reads/writes happen via readiness events from that call. Listeners are
//! registered for `READABLE` only; on accept, the new connection is added to
//! the same reactor with `READABLE`, then re-armed to `WRITABLE` once it has
//! a queued response.
//!
//! Tokens are allocated monotonically; listeners and connections share the
//! token space and are distinguished by which map they live in. Idle
//! connections (no progress for `timeout_ms`) are evicted at the end of every
//! iteration.

use std::collections::HashMap;
use std::io;
use std::net::{SocketAddr, TcpListener};
use std::os::fd::AsRawFd;
use std::rc::Rc;

use crate::application::connection::{Connection, ConnectionAction};
use crate::application::ports::clock::Clock;
use crate::application::ports::reactor::{Event, Interest, Reactor, Token};
use crate::application::request_pipeline::PipelineContext;
use crate::infrastructure::net::socket::{bind_listener, make_nonblocking};

/// Default idle-connection timeout: 30 seconds (PLAN Phase 3b).
pub const DEFAULT_TIMEOUT_MS: u64 = 30_000;

/// Default reactor `poll` timeout. Bounded so idle-eviction runs even when no
/// I/O is happening.
pub const DEFAULT_POLL_TIMEOUT_MS: i32 = 1_000;

#[derive(Debug)]
pub struct EventLoop<R: Reactor, C: Clock> {
    reactor: R,
    clock: C,
    next_token: Token,
    listeners: HashMap<Token, TcpListener>,
    /// Local address each listener is bound to, needed for virtual-host routing.
    listener_addrs: HashMap<Token, SocketAddr>,
    connections: HashMap<Token, Connection>,
    timeout_ms: u64,
    poll_timeout_ms: i32,
    /// When `Some`, new connections are dispatched through the real pipeline.
    pipeline: Option<Rc<PipelineContext>>,
}

impl<R: Reactor, C: Clock> EventLoop<R, C> {
    pub fn new(reactor: R, clock: C) -> Self {
        Self {
            reactor,
            clock,
            next_token: 1,
            listeners: HashMap::new(),
            listener_addrs: HashMap::new(),
            connections: HashMap::new(),
            timeout_ms: DEFAULT_TIMEOUT_MS,
            poll_timeout_ms: DEFAULT_POLL_TIMEOUT_MS,
            pipeline: None,
        }
    }

    /// Attach a real request pipeline; must be called before `run` / `tick`.
    #[must_use]
    pub fn with_pipeline(mut self, pipeline: Rc<PipelineContext>) -> Self {
        self.pipeline = Some(pipeline);
        self
    }

    #[must_use]
    pub fn with_timeout_ms(mut self, timeout_ms: u64) -> Self {
        self.timeout_ms = timeout_ms;
        self
    }

    #[must_use]
    pub fn with_poll_timeout_ms(mut self, poll_timeout_ms: i32) -> Self {
        self.poll_timeout_ms = poll_timeout_ms;
        self
    }

    /// Bind, set non-blocking + `SO_REUSEADDR`, register on the reactor.
    pub fn bind(&mut self, addr: SocketAddr) -> io::Result<Token> {
        let listener = bind_listener(addr)?;
        let local_addr = listener.local_addr()?;
        let token = self.alloc_token();
        self.reactor
            .register(listener.as_raw_fd(), token, Interest::READABLE)?;
        self.listeners.insert(token, listener);
        self.listener_addrs.insert(token, local_addr);
        Ok(token)
    }

    pub fn listener_count(&self) -> usize {
        self.listeners.len()
    }

    pub fn connection_count(&self) -> usize {
        self.connections.len()
    }

    pub fn connection_tokens(&self) -> Vec<Token> {
        self.connections.keys().copied().collect()
    }

    fn alloc_token(&mut self) -> Token {
        let t = self.next_token;
        self.next_token = self.next_token.checked_add(1).unwrap_or(1);
        t
    }

    /// Run a single iteration: one `poll`, dispatch events, evict idle
    /// connections. Tests drive the loop one tick at a time.
    pub fn tick(&mut self) -> io::Result<()> {
        let mut events: Vec<Event> = Vec::with_capacity(64);
        self.reactor.poll(&mut events, Some(self.poll_timeout_ms))?;
        for ev in events {
            self.dispatch(ev);
        }
        self.evict_idle();
        Ok(())
    }

    /// Run forever (until an unrecoverable reactor error).
    pub fn run(&mut self) -> io::Result<()> {
        loop {
            self.tick()?;
        }
    }

    fn dispatch(&mut self, ev: Event) {
        if self.listeners.contains_key(&ev.token) {
            self.on_listener_event(ev);
        } else if self.connections.contains_key(&ev.token) {
            self.on_connection_event(ev);
        }
        // Unknown token: stale event for an already-closed connection. Drop.
    }

    fn on_listener_event(&mut self, ev: Event) {
        if !ev.readable && !ev.is_terminal() {
            return;
        }
        // Drain all pending accepts up front so we don't hold an immutable
        // borrow of `self.listeners` across the alloc_token / register calls.
        let mut accepted: Vec<std::net::TcpStream> = Vec::new();
        if let Some(listener) = self.listeners.get(&ev.token) {
            loop {
                match listener.accept() {
                    Ok((stream, _addr)) => {
                        if make_nonblocking(&stream).is_err() {
                            continue;
                        }
                        accepted.push(stream);
                    }
                    Err(e) if e.kind() == io::ErrorKind::WouldBlock => break,
                    Err(e) if e.kind() == io::ErrorKind::Interrupted => {}
                    Err(_) => break,
                }
            }
        }
        let local_addr = self.listener_addrs.get(&ev.token).copied();
        for stream in accepted {
            let fd = stream.as_raw_fd();
            let token = self.alloc_token();
            if self
                .reactor
                .register(fd, token, Interest::READABLE)
                .is_err()
            {
                continue;
            }
            let now = self.clock.now_millis();
            let conn = if let (Some(addr), Some(pipeline)) = (local_addr, self.pipeline.as_ref()) {
                Connection::with_pipeline(token, stream, now, addr, Rc::clone(pipeline))
            } else {
                Connection::new(token, stream, now)
            };
            self.connections.insert(token, conn);
        }
    }

    fn on_connection_event(&mut self, ev: Event) {
        let now = self.clock.now_millis();
        let mut close = ev.is_terminal();
        let mut next_interest: Option<Interest> = None;

        if let Some(conn) = self.connections.get_mut(&ev.token) {
            if !close && ev.readable {
                match conn.on_readable(now) {
                    ConnectionAction::Close => close = true,
                    ConnectionAction::Rearm(i) => next_interest = Some(i),
                }
            }
            if !close && ev.writable {
                match conn.on_writable(now) {
                    ConnectionAction::Close => close = true,
                    ConnectionAction::Rearm(i) => next_interest = Some(i),
                }
            }
        }

        if close {
            self.close_connection(ev.token);
        } else if let Some(interest) = next_interest
            && let Some(conn) = self.connections.get(&ev.token)
            && self
                .reactor
                .reregister(conn.fd(), ev.token, interest)
                .is_err()
        {
            self.close_connection(ev.token);
        }
    }

    fn close_connection(&mut self, token: Token) {
        if let Some(conn) = self.connections.remove(&token) {
            // Best-effort deregister; the kernel auto-cleans on close anyway.
            let _ = self.reactor.deregister(conn.fd());
            // Drop closes the socket.
        }
    }

    fn evict_idle(&mut self) {
        let now = self.clock.now_millis();
        let timeout = self.timeout_ms;
        let stale: Vec<Token> = self
            .connections
            .iter()
            .filter_map(|(t, c)| {
                if c.is_idle(now, timeout) {
                    Some(*t)
                } else {
                    None
                }
            })
            .collect();
        for token in stale {
            // Try to send a 408 response before closing; if the connection is
            // already past the response phase, close immediately.
            let action = self
                .connections
                .get_mut(&token)
                .map(|c| c.start_timeout(now));
            match action {
                Some(ConnectionAction::Rearm(interest)) => {
                    if let Some(conn) = self.connections.get(&token) {
                        if self.reactor.reregister(conn.fd(), token, interest).is_err() {
                            self.close_connection(token);
                        }
                    }
                }
                _ => self.close_connection(token),
            }
        }
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::io::{Read as _, Write as _};
    use std::net::TcpStream;
    use std::time::{Duration, Instant};

    use super::*;
    use crate::application::ports::clock::FixedClock;
    use crate::application::ports::reactor::Event;
    use crate::infrastructure::reactor::mock::MockReactor;

    fn loopback() -> SocketAddr {
        "127.0.0.1:0".parse().unwrap()
    }

    fn make_event(token: Token, readable: bool, writable: bool) -> Event {
        Event {
            token,
            readable,
            writable,
            error: false,
            hangup: false,
        }
    }

    fn read_with_timeout(stream: &mut TcpStream, want_prefix: &[u8]) -> Vec<u8> {
        let start = Instant::now();
        let mut got = Vec::new();
        let mut buf = [0u8; 4096];
        while got.len() < want_prefix.len() {
            match stream.read(&mut buf) {
                Ok(0) => break,
                Ok(n) => got.extend_from_slice(&buf[..n]),
                Err(e) if e.kind() == io::ErrorKind::WouldBlock => {
                    assert!(
                        start.elapsed() <= Duration::from_secs(2),
                        "timeout, got: {got:?}"
                    );
                    std::thread::sleep(Duration::from_millis(5));
                }
                Err(e) => panic!("read error: {e}"),
            }
        }
        got
    }

    #[test]
    fn bind_registers_listener_as_readable() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0));
        let token = el.bind(loopback()).unwrap();
        assert_eq!(el.listener_count(), 1);
        let registration = &el.reactor.registered;
        let listener_reg = registration
            .values()
            .find(|r| r.token == token)
            .expect("registration");
        assert!(listener_reg.interest.readable);
        assert!(!listener_reg.interest.writable);
    }

    #[test]
    fn tick_calls_poll_exactly_once() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0));
        let _ = el.bind(loopback()).unwrap();
        let before = el.reactor.poll_calls;
        el.tick().unwrap();
        assert_eq!(el.reactor.poll_calls, before + 1);
        el.tick().unwrap();
        assert_eq!(el.reactor.poll_calls, before + 2);
    }

    #[test]
    fn listener_event_accepts_and_registers_connection() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0));
        let listener_token = el.bind(loopback()).unwrap();
        let listener_addr = el
            .listeners
            .get(&listener_token)
            .unwrap()
            .local_addr()
            .unwrap();

        let _client = TcpStream::connect(listener_addr).unwrap();
        // Settle the SYN/ACK.
        std::thread::sleep(Duration::from_millis(20));

        el.reactor
            .push_event(make_event(listener_token, true, false));
        el.tick().unwrap();

        assert_eq!(el.connection_count(), 1);
        // The new connection's fd must be registered with READABLE only.
        let conn_token = el.connection_tokens()[0];
        let reg = el
            .reactor
            .registered
            .values()
            .find(|r| r.token == conn_token)
            .expect("connection registration");
        assert!(reg.interest.readable);
        assert!(!reg.interest.writable);
    }

    #[test]
    fn full_round_trip_via_skeleton_response() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0));
        let listener_token = el.bind(loopback()).unwrap();
        let addr = el
            .listeners
            .get(&listener_token)
            .unwrap()
            .local_addr()
            .unwrap();

        let mut client = TcpStream::connect(addr).unwrap();
        std::thread::sleep(Duration::from_millis(20));

        // Accept.
        el.reactor
            .push_event(make_event(listener_token, true, false));
        el.tick().unwrap();
        let conn_token = el.connection_tokens()[0];

        // Client sends a complete request and asks the server to close so
        // the assertions are deterministic regardless of keep-alive default.
        client
            .write_all(b"GET / HTTP/1.1\r\nHost: x\r\nConnection: close\r\n\r\n")
            .unwrap();
        std::thread::sleep(Duration::from_millis(20));

        // Drive the read.
        el.reactor.push_event(make_event(conn_token, true, false));
        el.tick().unwrap();
        // Should now be re-armed for WRITABLE.
        let reg = el
            .reactor
            .registered
            .values()
            .find(|r| r.token == conn_token)
            .expect("registration after read");
        assert!(reg.interest.writable);

        // Drive the write.
        el.reactor.push_event(make_event(conn_token, false, true));
        client.set_nonblocking(true).unwrap();
        el.tick().unwrap();

        // `Connection: close` means we close after writing.
        assert_eq!(el.connection_count(), 0);
        let got = read_with_timeout(&mut client, b"HTTP/1.1 200 OK\r\n");
        assert!(got.starts_with(b"HTTP/1.1 200 OK\r\n"), "got: {got:?}");
    }

    #[test]
    fn terminal_event_closes_connection() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0));
        let listener_token = el.bind(loopback()).unwrap();
        let addr = el
            .listeners
            .get(&listener_token)
            .unwrap()
            .local_addr()
            .unwrap();

        let _client = TcpStream::connect(addr).unwrap();
        std::thread::sleep(Duration::from_millis(20));
        el.reactor
            .push_event(make_event(listener_token, true, false));
        el.tick().unwrap();
        assert_eq!(el.connection_count(), 1);

        let conn_token = el.connection_tokens()[0];
        let mut hup = make_event(conn_token, false, false);
        hup.hangup = true;
        el.reactor.push_event(hup);
        el.tick().unwrap();
        assert_eq!(el.connection_count(), 0);
    }

    #[test]
    fn idle_connections_are_evicted() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0)).with_timeout_ms(100);
        let listener_token = el.bind(loopback()).unwrap();
        let addr = el
            .listeners
            .get(&listener_token)
            .unwrap()
            .local_addr()
            .unwrap();

        let _client = TcpStream::connect(addr).unwrap();
        std::thread::sleep(Duration::from_millis(20));
        el.reactor
            .push_event(make_event(listener_token, true, false));
        el.tick().unwrap();
        assert_eq!(el.connection_count(), 1);

        // Advance the clock past the timeout.
        // evict_idle queues a 408 and re-arms the connection for WRITABLE.
        el.clock.set(10_000);
        el.tick().unwrap();
        // Connection is still alive — waiting to flush the 408 response.
        assert_eq!(el.connection_count(), 1);
        let conn_token = el.connection_tokens()[0];

        // Fire the WRITABLE event: 408 is written, connection closes.
        el.reactor.push_event(make_event(conn_token, false, true));
        el.tick().unwrap();
        assert_eq!(el.connection_count(), 0);
    }

    #[test]
    fn unknown_token_event_is_dropped_silently() {
        let mut el = EventLoop::new(MockReactor::new(), FixedClock::new(0));
        el.reactor.push_event(make_event(9999, true, false));
        // No panic, no insert.
        el.tick().unwrap();
        assert_eq!(el.connection_count(), 0);
        assert_eq!(el.listener_count(), 0);
    }
}
