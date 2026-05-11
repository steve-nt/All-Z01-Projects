//! Minimal stderr logger. No async, no allocations beyond the formatted line.

use std::fmt;
use std::io::Write;
use std::sync::atomic::{AtomicU8, Ordering};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Copy, Clone, Debug, Eq, PartialEq, Ord, PartialOrd)]
#[repr(u8)]
pub enum Level {
    Error = 0,
    Warn = 1,
    Info = 2,
    Debug = 3,
    Trace = 4,
}

impl Level {
    fn tag(self) -> &'static str {
        match self {
            Self::Error => "ERROR",
            Self::Warn => "WARN ",
            Self::Info => "INFO ",
            Self::Debug => "DEBUG",
            Self::Trace => "TRACE",
        }
    }
}

static MAX_LEVEL: AtomicU8 = AtomicU8::new(Level::Info as u8);

pub fn set_max_level(level: Level) {
    MAX_LEVEL.store(level as u8, Ordering::Relaxed);
}

pub fn is_enabled(level: Level) -> bool {
    (level as u8) <= MAX_LEVEL.load(Ordering::Relaxed)
}

pub fn log(level: Level, args: fmt::Arguments<'_>) {
    if !is_enabled(level) {
        return;
    }
    let secs = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_secs())
        .unwrap_or(0);
    let mut stderr = std::io::stderr().lock();
    let _ = writeln!(stderr, "[{secs}] {} {args}", level.tag());
}

#[macro_export]
macro_rules! log_error {
    ($($arg:tt)*) => { $crate::logger::log($crate::logger::Level::Error, format_args!($($arg)*)) };
}

#[macro_export]
macro_rules! log_warn {
    ($($arg:tt)*) => { $crate::logger::log($crate::logger::Level::Warn, format_args!($($arg)*)) };
}

#[macro_export]
macro_rules! log_info {
    ($($arg:tt)*) => { $crate::logger::log($crate::logger::Level::Info, format_args!($($arg)*)) };
}

#[macro_export]
macro_rules! log_debug {
    ($($arg:tt)*) => { $crate::logger::log($crate::logger::Level::Debug, format_args!($($arg)*)) };
}

#[macro_export]
macro_rules! log_trace {
    ($($arg:tt)*) => { $crate::logger::log($crate::logger::Level::Trace, format_args!($($arg)*)) };
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn level_filtering() {
        set_max_level(Level::Warn);
        assert!(is_enabled(Level::Error));
        assert!(is_enabled(Level::Warn));
        assert!(!is_enabled(Level::Info));
        set_max_level(Level::Trace);
        assert!(is_enabled(Level::Trace));
    }
}
