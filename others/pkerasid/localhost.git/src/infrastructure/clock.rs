//! `SystemClock` — monotonic clock backed by `std::time::Instant`.

use std::time::Instant;

use crate::application::ports::clock::Clock;

/// Process-monotonic clock anchored at construction time.
#[derive(Debug)]
pub struct SystemClock {
    origin: Instant,
}

impl SystemClock {
    pub fn new() -> Self {
        Self {
            origin: Instant::now(),
        }
    }
}

impl Default for SystemClock {
    fn default() -> Self {
        Self::new()
    }
}

impl Clock for SystemClock {
    fn now_millis(&self) -> u64 {
        // Instant::elapsed is guaranteed monotonic, so as_millis truncation
        // direction is well-defined and never goes backwards.
        u64::try_from(self.origin.elapsed().as_millis()).unwrap_or(u64::MAX)
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use super::*;

    #[test]
    fn monotonic_non_decreasing() {
        let c = SystemClock::new();
        let a = c.now_millis();
        // Tight loop — both calls may return the same value but never
        // a smaller one.
        for _ in 0..100 {
            let b = c.now_millis();
            assert!(b >= a, "clock went backwards: {a} -> {b}");
        }
    }
}
