//! `localhost` — single-threaded, non-blocking HTTP/1.1 server.
//!
//! Layered following clean architecture:
//! - `domain`         pure types, no I/O.
//! - `application`    use cases and orchestration; depends only on domain + ports.
//! - `infrastructure` real I/O (epoll, sockets, fs, processes). Implements ports.
//! - `interface`      boundary parsers/serializers (HTTP, config).
//!
//! The dependency rule points strictly inward: outer layers depend on inner,
//! never the reverse.

pub mod application;
pub mod domain;
pub mod infrastructure;
pub mod interface;
pub mod logger;

pub const VERSION: &str = env!("CARGO_PKG_VERSION");
