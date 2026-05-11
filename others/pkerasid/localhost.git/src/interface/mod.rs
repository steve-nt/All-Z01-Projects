//! Interface layer — boundary parsers and serializers.
//!
//! Translates external formats (TOML config, raw HTTP bytes) to/from the
//! validated domain types. No I/O here; that lives in `infrastructure`.

pub mod config_parser;
pub mod http_parser;
pub mod http_writer;
