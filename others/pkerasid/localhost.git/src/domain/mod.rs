//! Domain layer — pure types, no I/O.
//!
//! Outer layers depend on this; this layer depends on nothing inside the crate.

pub mod config;
pub mod error;
pub mod http;
pub mod session;
