//! Infrastructure layer — real I/O adapters.
//!
//! Implements the ports defined in `application::ports` against the OS
//! (sockets, epoll, filesystem, processes). Outer layers may depend inward on
//! `application` and `domain`; nothing inward depends on this module.

pub mod cgi;
pub mod clock;
pub mod fs;
pub mod net;
pub mod reactor;
pub mod session_store;
