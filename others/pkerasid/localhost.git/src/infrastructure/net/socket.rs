//! Socket helpers: `O_NONBLOCK`, `SO_REUSEADDR`, accept-non-block.
//!
//! All I/O on these sockets must be driven from the reactor; helpers here
//! configure them but do not block.

use std::io;
use std::net::{SocketAddr, TcpListener, TcpStream};
use std::os::fd::{AsRawFd, RawFd};

/// Bind a non-blocking, `SO_REUSEADDR`-enabled TCP listener.
pub fn bind_listener(addr: SocketAddr) -> io::Result<TcpListener> {
    let listener = TcpListener::bind(addr)?;
    listener.set_nonblocking(true)?;
    set_reuseaddr(listener.as_raw_fd(), true)?;
    Ok(listener)
}

/// Mark an accepted stream as non-blocking. (`accept(2)` on Linux inherits
/// `O_NONBLOCK` from the listener with `accept4`, but `std`'s `accept` does
/// not, so we set it explicitly.)
pub fn make_nonblocking(stream: &TcpStream) -> io::Result<()> {
    stream.set_nonblocking(true)
}

/// Toggle `SO_REUSEADDR` on a raw fd.
pub fn set_reuseaddr(fd: RawFd, on: bool) -> io::Result<()> {
    let flag: libc::c_int = i32::from(on);
    let optlen = libc::socklen_t::try_from(std::mem::size_of_val(&flag))
        .map_err(|_| io::Error::other("optlen overflow"))?;
    // SAFETY: setsockopt with SOL_SOCKET / SO_REUSEADDR is well-defined for
    // any valid socket fd; we pass a correctly-sized int and the matching
    // `optlen`. Errors are checked below.
    let rc = unsafe {
        libc::setsockopt(
            fd,
            libc::SOL_SOCKET,
            libc::SO_REUSEADDR,
            std::ptr::from_ref(&flag).cast::<libc::c_void>(),
            optlen,
        )
    };
    if rc == 0 {
        Ok(())
    } else {
        Err(io::Error::last_os_error())
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    fn loopback() -> SocketAddr {
        // Port 0 lets the OS pick a free port.
        "127.0.0.1:0".parse().unwrap()
    }

    #[test]
    fn bind_listener_is_nonblocking() {
        let l = bind_listener(loopback()).unwrap();
        // accept on an idle listener returns WouldBlock when non-blocking.
        match l.accept() {
            Err(e) if e.kind() == io::ErrorKind::WouldBlock => {}
            other => panic!("expected WouldBlock, got {other:?}"),
        }
    }

    #[test]
    fn make_nonblocking_succeeds_on_connected_pair() {
        let server = bind_listener(loopback()).unwrap();
        let local = server.local_addr().unwrap();
        let client = TcpStream::connect(local).unwrap();
        make_nonblocking(&client).unwrap();
    }
}
