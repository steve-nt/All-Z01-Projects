//! Clock port: monotonic time for timeout bookkeeping.
//!
//! Decoupled from `std::time` so tests can inject deterministic time.

/// Monotonic clock in milliseconds.
///
/// Implementations must be monotonic — values returned from successive calls
/// never decrease. Used for connection idle/timeout tracking.
pub trait Clock {
    fn now_millis(&self) -> u64;
}

/// In-test clock that returns a programmable value.
#[derive(Debug, Default, Clone)]
pub struct FixedClock {
    now: u64,
}

impl FixedClock {
    pub fn new(now: u64) -> Self {
        Self { now }
    }

    pub fn advance(&mut self, delta_ms: u64) {
        self.now = self.now.saturating_add(delta_ms);
    }

    pub fn set(&mut self, now: u64) {
        self.now = now;
    }
}

impl Clock for FixedClock {
    fn now_millis(&self) -> u64 {
        self.now
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn fixed_clock_advances() {
        let mut c = FixedClock::new(100);
        assert_eq!(c.now_millis(), 100);
        c.advance(50);
        assert_eq!(c.now_millis(), 150);
        c.set(0);
        assert_eq!(c.now_millis(), 0);
    }

    #[test]
    fn advance_saturates() {
        let mut c = FixedClock::new(u64::MAX);
        c.advance(1);
        assert_eq!(c.now_millis(), u64::MAX);
    }
}
