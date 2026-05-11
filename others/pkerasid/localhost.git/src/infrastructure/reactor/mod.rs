//! Reactor implementations.
//!
//! Linux: `epoll` (audit-compliant). On other platforms only the `MockReactor`
//! is built, so unit tests still run and the project still compiles, but the
//! actual server binary is intended to run on Linux.

pub mod mock;

#[cfg(target_os = "linux")]
pub mod epoll;
