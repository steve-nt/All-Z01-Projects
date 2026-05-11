/// Minimum gap between the front of the follower and the rear of the leader.
/// `gap = leader_progress - follower_progress - vehicle_length`
#[must_use]
pub const fn min_safe_gap(_vehicle_length: f32, buffer: f32) -> f32 {
    buffer
}

/// Minimum center-to-center distance between two vehicles when spawning.
#[must_use]
pub const fn min_spawn_center_distance(vehicle_length: f32, buffer: f32) -> f32 {
    vehicle_length + buffer
}

/// Returns true when the follower is closer than the safe gap to the vehicle ahead.
/// `gap` is the current distance between front of follower and rear of leader.
#[must_use]
pub fn is_too_close(gap: f32, vehicle_length: f32, buffer: f32) -> bool {
    gap < min_safe_gap(vehicle_length, buffer)
}

/// Compute the appropriate target speed for a follower vehicle given the current gap to its leader.
///
/// Two zones based on the minimum safe gap:
/// - gap < safe   → match leader speed so the gap stops shrinking
/// - gap ≥ safe   → free to drive at phase speed
///
/// This lets queued vehicles close up to the same target empty gap instead of
/// stopping early at uneven distances.
#[must_use]
pub fn adjusted_follower_speed(
    gap: f32,
    natural_speed: f32,
    leader_speed: f32,
    vehicle_length: f32,
    buffer: f32,
    _low_speed: f32,
    _medium_speed: f32,
) -> f32 {
    let safe = min_safe_gap(vehicle_length, buffer);
    if gap < safe {
        leader_speed.min(natural_speed)
    } else {
        natural_speed
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const LEN: f32 = 20.0;
    const BUF: f32 = 10.0;

    #[test]
    fn safe_gap_is_buffer() {
        assert!((min_safe_gap(LEN, BUF) - BUF).abs() < f32::EPSILON);
    }

    #[test]
    fn spawn_center_distance_is_length_plus_buffer() {
        assert!((min_spawn_center_distance(LEN, BUF) - 30.0).abs() < f32::EPSILON);
    }

    #[test]
    fn gap_below_minimum_is_too_close() {
        assert!(is_too_close(9.9, LEN, BUF));
    }

    #[test]
    fn gap_at_minimum_is_not_too_close() {
        assert!(!is_too_close(10.0, LEN, BUF));
    }

    #[test]
    fn gap_above_minimum_is_safe() {
        assert!(!is_too_close(50.0, LEN, BUF));
    }

    #[test]
    fn follower_keeps_closing_until_safe_gap() {
        let speed = adjusted_follower_speed(15.0, 50.0, 0.0, LEN, BUF, 20.0, 40.0);
        assert!((speed - 50.0).abs() < f32::EPSILON);
    }
}
