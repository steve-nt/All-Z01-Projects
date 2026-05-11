//! Linux `epoll` reactor.
//!
//! A thin, safe wrapper over `epoll_create1`, `epoll_ctl`, and `epoll_wait`.
//! Edge-triggered mode (`EPOLLET`) plus `EPOLLRDHUP` so we observe peer
//! half-closes promptly. The application drains until `EAGAIN` after every
//! readiness event, satisfying ET semantics.
//!
//! Audit invariant: the event loop calls `poll` (which calls `epoll_wait`)
//! exactly once per iteration; all reads/writes happen through readiness
//! events from that call.

use std::io;
use std::os::fd::{AsRawFd, FromRawFd, OwnedFd, RawFd};

use crate::application::ports::reactor::{Event, Interest, Reactor, Token};

/// Linux epoll-based reactor.
#[derive(Debug)]
pub struct EpollReactor {
    epfd: OwnedFd,
}

impl EpollReactor {
    pub fn new() -> io::Result<Self> {
        // SAFETY: epoll_create1 with EPOLL_CLOEXEC returns a new fd or -1.
        // We check the return value before constructing OwnedFd.
        let raw = unsafe { libc::epoll_create1(libc::EPOLL_CLOEXEC) };
        if raw < 0 {
            return Err(io::Error::last_os_error());
        }
        // SAFETY: raw is a fresh, owned epoll fd we just created.
        let epfd = unsafe { OwnedFd::from_raw_fd(raw) };
        Ok(Self { epfd })
    }

    fn ctl(&self, op: libc::c_int, fd: RawFd, token: Token, interest: Interest) -> io::Result<()> {
        let mut ev = libc::epoll_event {
            events: events_mask(interest),
            u64: token as u64,
        };
        // SAFETY: epoll_ctl with a valid epoll fd, valid op, valid fd, and a
        // properly initialized epoll_event is well-defined. For DEL the event
        // pointer is unused but accepting non-null is safe per kernel docs.
        let rc =
            unsafe { libc::epoll_ctl(self.epfd.as_raw_fd(), op, fd, std::ptr::from_mut(&mut ev)) };
        if rc == 0 {
            Ok(())
        } else {
            Err(io::Error::last_os_error())
        }
    }
}

#[allow(clippy::cast_sign_loss)]
fn events_mask(interest: Interest) -> u32 {
    let mut m: u32 = 0;
    if interest.readable {
        m |= libc::EPOLLIN as u32;
    }
    if interest.writable {
        m |= libc::EPOLLOUT as u32;
    }
    m |= libc::EPOLLRDHUP as u32;
    m |= libc::EPOLLET as u32;
    m
}

impl Reactor for EpollReactor {
    fn register(&mut self, fd: RawFd, token: Token, interest: Interest) -> io::Result<()> {
        self.ctl(libc::EPOLL_CTL_ADD, fd, token, interest)
    }

    fn reregister(&mut self, fd: RawFd, token: Token, interest: Interest) -> io::Result<()> {
        self.ctl(libc::EPOLL_CTL_MOD, fd, token, interest)
    }

    fn deregister(&mut self, fd: RawFd) -> io::Result<()> {
        let rc = unsafe {
            libc::epoll_ctl(
                self.epfd.as_raw_fd(),
                libc::EPOLL_CTL_DEL,
                fd,
                std::ptr::null_mut(),
            )
        };
        if rc == 0 {
            return Ok(());
        }
        // ENOENT is acceptable — caller may have already closed the fd.
        let err = io::Error::last_os_error();
        if err.raw_os_error() == Some(libc::ENOENT) {
            Ok(())
        } else {
            Err(err)
        }
    }

    fn poll(&mut self, out: &mut Vec<Event>, timeout_ms: Option<i32>) -> io::Result<usize> {
        const MAX_EVENTS: usize = 1024;
        let mut buf: [libc::epoll_event; MAX_EVENTS] = unsafe { std::mem::zeroed() };
        let timeout = timeout_ms.unwrap_or(-1);

        // SAFETY: buf is a stack array of MAX_EVENTS valid epoll_event slots.
        // The kernel writes at most MAX_EVENTS entries.
        #[allow(clippy::cast_possible_truncation, clippy::cast_possible_wrap)]
        let n = unsafe {
            libc::epoll_wait(
                self.epfd.as_raw_fd(),
                buf.as_mut_ptr(),
                MAX_EVENTS as libc::c_int,
                timeout,
            )
        };
        if n < 0 {
            let err = io::Error::last_os_error();
            // EINTR is benign — caller polls again.
            if err.raw_os_error() == Some(libc::EINTR) {
                return Ok(0);
            }
            return Err(err);
        }
        #[allow(clippy::cast_sign_loss)]
        let n_usize = n as usize;
        out.reserve(n_usize);
        for ev in &buf[..n_usize] {
            out.push(decode(ev));
        }
        Ok(n_usize)
    }
}

#[allow(clippy::cast_sign_loss, clippy::cast_possible_truncation)]
fn decode(ev: &libc::epoll_event) -> Event {
    let bits = ev.events;
    let read = bits & (libc::EPOLLIN as u32) != 0;
    let write = bits & (libc::EPOLLOUT as u32) != 0;
    let err = bits & (libc::EPOLLERR as u32) != 0;
    let hup = bits & ((libc::EPOLLHUP as u32) | (libc::EPOLLRDHUP as u32)) != 0;
    Event {
        token: ev.u64 as Token,
        readable: read,
        writable: write,
        error: err,
        hangup: hup,
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::io::{Read, Write};
    use std::net::{TcpListener, TcpStream};

    use super::*;
    use crate::infrastructure::net::socket::{bind_listener, make_nonblocking};

    fn loopback() -> std::net::SocketAddr {
        "127.0.0.1:0".parse().unwrap()
    }

    #[test]
    fn create_reactor() {
        let _r = EpollReactor::new().unwrap();
    }

    #[test]
    fn detects_listener_readable() {
        let listener: TcpListener = bind_listener(loopback()).unwrap();
        let local = listener.local_addr().unwrap();
        let mut reactor = EpollReactor::new().unwrap();

        reactor
            .register(listener.as_raw_fd(), 42, Interest::READABLE)
            .unwrap();

        // Connect from a peer to make the listener readable.
        let _client = TcpStream::connect(local).unwrap();

        let mut events: Vec<Event> = Vec::new();
        let n = reactor.poll(&mut events, Some(1000)).unwrap();
        assert!(n >= 1, "expected at least one event");
        let ev = events.iter().find(|e| e.token == 42).expect("token match");
        assert!(ev.readable);
    }

    #[test]
    fn detects_stream_readable_after_write() {
        let listener: TcpListener = bind_listener(loopback()).unwrap();
        let local = listener.local_addr().unwrap();
        let client_writer = TcpStream::connect(local).unwrap();

        let (server_stream, _) = {
            // The listener is non-blocking; loop briefly until accept succeeds.
            let mut acc = listener.accept();
            let start = std::time::Instant::now();
            while matches!(&acc, Err(e) if e.kind() == io::ErrorKind::WouldBlock) {
                assert!(
                    start.elapsed() <= std::time::Duration::from_secs(2),
                    "accept timed out"
                );
                std::thread::sleep(std::time::Duration::from_millis(5));
                acc = listener.accept();
            }
            acc.unwrap()
        };

        make_nonblocking(&server_stream).unwrap();
        let mut reactor = EpollReactor::new().unwrap();
        reactor
            .register(server_stream.as_raw_fd(), 7, Interest::READABLE)
            .unwrap();

        // Push some bytes from the client side.
        (&client_writer).write_all(b"hello").unwrap();

        let mut events = Vec::new();
        let n = reactor.poll(&mut events, Some(1000)).unwrap();
        assert!(n >= 1);
        let ev = events.iter().find(|e| e.token == 7).expect("token match");
        assert!(ev.readable);

        // Drain to keep ET happy.
        let mut buf = [0u8; 16];
        let read = (&server_stream).read(&mut buf).unwrap();
        assert_eq!(&buf[..read], b"hello");
    }

    #[test]
    fn deregister_unknown_fd_is_ok() {
        let mut reactor = EpollReactor::new().unwrap();
        // Some fd that's open but never registered — use stderr.
        reactor.deregister(libc::STDERR_FILENO).unwrap();
    }

    #[test]
    fn poll_zero_timeout_returns_quickly() {
        let mut reactor = EpollReactor::new().unwrap();
        let mut events = Vec::new();
        // No registrations: poll returns 0 immediately.
        let n = reactor.poll(&mut events, Some(0)).unwrap();
        assert_eq!(n, 0);
        assert!(events.is_empty());
    }
}
