//! Application layer — use cases, orchestration, and ports.
//!
//! Depends only on `domain`; the infrastructure layer adapts to the ports
//! defined here. The dependency rule points strictly inward.

pub mod connection;
pub mod error_pages;
pub mod event_loop;
pub mod handlers;
pub mod ports;
pub mod request_pipeline;
pub mod router;
