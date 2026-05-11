#![allow(clippy::cast_precision_loss)]

use crate::config::{INTER_LEFT, INTER_TOP, LANE_W, WINDOW_H, WINDOW_W};
use crate::domain::lane::{Direction, Route};

const OFFSCREEN_MARGIN: f32 = 40.0;

/// Minimal 2D vector used by sampled path geometry.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Vec2 {
    pub x: f32,
    pub y: f32,
}

impl Vec2 {
    #[must_use]
    pub const fn new(x: f32, y: f32) -> Self {
        Self { x, y }
    }

    #[must_use]
    pub fn distance_to(self, other: Self) -> f32 {
        let dx = other.x - self.x;
        let dy = other.y - self.y;
        dx.hypot(dy)
    }

    #[must_use]
    pub fn lerp(self, other: Self, t: f32) -> Self {
        Self::new(
            self.x + (other.x - self.x) * t,
            self.y + (other.y - self.y) * t,
        )
    }
}

/// Position plus tangent heading sampled from a route path.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct PathSample {
    pub position: Vec2,
    pub heading_deg: f32,
}

/// Sampled polyline route used by vehicles to move through the map.
///
/// Each route is an L-shaped (or straight) polyline with exactly 2 or 3 waypoints:
///
/// - Straight: `[entry_offscreen, exit_offscreen]`
/// - Turn:     `[entry_offscreen, turn_waypoint, exit_offscreen]`
///
/// The turn waypoint is the center of the single tile where the 90° change of
/// direction happens. This matches the conflict-map tile grid exactly.
#[derive(Debug, Clone, PartialEq)]
pub struct RoutePath {
    points: Vec<Vec2>,
    segment_lengths: Vec<f32>,
    pub total_length: f32,
}

impl RoutePath {
    /// Build a path from polyline waypoints.
    ///
    /// # Panics
    ///
    /// Panics if fewer than two points are provided.
    #[must_use]
    pub fn from_points(points: Vec<Vec2>) -> Self {
        assert!(points.len() >= 2, "route path needs at least two points");

        let segment_lengths: Vec<f32> = points.windows(2).map(|s| s[0].distance_to(s[1])).collect();
        let total_length = segment_lengths.iter().sum();

        Self {
            points,
            segment_lengths,
            total_length,
        }
    }

    /// Build the canonical path for one of the 12 lane routes.
    ///
    /// ### Coordinate system
    ///
    /// The 6-lane road is divided into two halves per direction.
    /// Columns (x) and rows (y) are numbered 0-5 from left/top of the
    /// intersection box.  Helper functions `col_x(n)` and `row_y(n)` return the
    /// pixel center of lane n.
    ///
    /// ```text
    /// cols:  0   1   2   3   4   5
    ///       NR  NS  NL  SL  SS  SR   ← North/South lanes
    ///
    /// rows:  0   1   2   3   4   5
    ///       ER  ES  EL  WL  WS  WR   ← East/West lanes
    /// ```
    ///
    /// Inner 4×4 conflict-map tile (c, r) maps to outer lane col c+1 / row r+1.
    ///
    /// ### Turn waypoints (must match the conflict map exactly)
    ///
    /// | Route      | Turn tile (inner) | Waypoint px |
    /// |------------|-------------------|-------------|
    /// | North-Left | (1,2) = col2,row3 | (620, 380)  |
    /// | South-Left | (2,1) = col3,row2 | (660, 340)  |
    /// | East-Left  | (1,1) = col2,row2 | (620, 340)  |
    /// | West-Left  | (2,2) = col3,row3 | (660, 380)  |
    ///
    /// Right turns use the outer corner tiles (excluded from the conflict zone):
    ///
    /// | Route       | Corner px  |
    /// |-------------|------------|
    /// | North-Right / East-Right  | (540, 260) — top-left  |
    /// | South-Right / West-Right  | (740, 460) — bot-right |
    #[must_use]
    pub fn for_lane(origin: Direction, route: Route) -> Self {
        use Direction::{East, North, South, West};
        use Route::{Left, Right, Straight};

        let w = WINDOW_W as f32;
        let h = WINDOW_H as f32;
        let m = OFFSCREEN_MARGIN;

        match (origin, route) {
            // ── North (enters top, travels south) ──────────────────────────
            (North, Right) => Self::from_points(vec![
                Vec2::new(col_x(0), -m),
                Vec2::new(col_x(0), row_y(0)), // corner tile top-left
                Vec2::new(-m, row_y(0)),
            ]),
            (North, Straight) => {
                Self::from_points(vec![Vec2::new(col_x(1), -m), Vec2::new(col_x(1), h + m)])
            }
            (North, Left) => Self::from_points(vec![
                Vec2::new(col_x(2), -m),
                Vec2::new(col_x(2), row_y(3)), // inner (1,2): col2, row3
                Vec2::new(w + m, row_y(3)),
            ]),

            // ── South (enters bottom, travels north) ───────────────────────
            (South, Right) => Self::from_points(vec![
                Vec2::new(col_x(5), h + m),
                Vec2::new(col_x(5), row_y(5)), // corner tile bot-right
                Vec2::new(w + m, row_y(5)),
            ]),
            (South, Straight) => {
                Self::from_points(vec![Vec2::new(col_x(4), h + m), Vec2::new(col_x(4), -m)])
            }
            (South, Left) => Self::from_points(vec![
                Vec2::new(col_x(3), h + m),
                Vec2::new(col_x(3), row_y(2)), // inner (2,1): col3, row2
                Vec2::new(-m, row_y(2)),
            ]),

            // ── East (enters right, travels west) ──────────────────────────
            (East, Right) => Self::from_points(vec![
                Vec2::new(w + m, row_y(0)),
                Vec2::new(col_x(5), row_y(0)), // corner tile top-right (first tile entering from East)
                Vec2::new(col_x(5), -m),
            ]),
            (East, Straight) => {
                Self::from_points(vec![Vec2::new(w + m, row_y(1)), Vec2::new(-m, row_y(1))])
            }
            (East, Left) => Self::from_points(vec![
                Vec2::new(w + m, row_y(2)),
                Vec2::new(col_x(2), row_y(2)), // inner (1,1): col2, row2
                Vec2::new(col_x(2), h + m),
            ]),

            // ── West (enters left, travels east) ───────────────────────────
            (West, Right) => Self::from_points(vec![
                Vec2::new(-m, row_y(5)),
                Vec2::new(col_x(0), row_y(5)), // corner tile bot-left (first tile entering from West)
                Vec2::new(col_x(0), h + m),
            ]),
            (West, Straight) => {
                Self::from_points(vec![Vec2::new(-m, row_y(4)), Vec2::new(w + m, row_y(4))])
            }
            (West, Left) => Self::from_points(vec![
                Vec2::new(-m, row_y(3)),
                Vec2::new(col_x(3), row_y(3)), // inner (2,2): col3, row3
                Vec2::new(col_x(3), -m),
            ]),
        }
    }

    #[must_use]
    pub fn sample(&self, s: f32) -> PathSample {
        let clamped = s.clamp(0.0, self.total_length);

        let mut walked = 0.0;
        for (idx, length) in self.segment_lengths.iter().copied().enumerate() {
            if walked + length >= clamped {
                let local_t = if length <= f32::EPSILON {
                    0.0
                } else {
                    (clamped - walked) / length
                };
                let start = self.points[idx];
                let end = self.points[idx + 1];
                let position = start.lerp(end, local_t);
                let heading_deg = (end.y - start.y).atan2(end.x - start.x).to_degrees();

                return PathSample {
                    position,
                    heading_deg,
                };
            }
            walked += length;
        }

        let start = self.points[self.points.len() - 2];
        let end = self.points[self.points.len() - 1];
        PathSample {
            position: end,
            heading_deg: (end.y - start.y).atan2(end.x - start.x).to_degrees(),
        }
    }

    #[must_use]
    pub fn point_count(&self) -> usize {
        self.points.len()
    }
}

/// Tracks how far a vehicle has traveled along its route, in world units.
/// `s` advances by `speed * dt` each tick. Path is complete when `s >= total_length`.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct PathProgress {
    pub s: f32,
    pub total_length: f32,
}

impl PathProgress {
    #[must_use]
    pub fn new(total_length: f32) -> Self {
        Self {
            s: 0.0,
            total_length,
        }
    }

    pub fn advance(&mut self, speed: f32, dt: f32) {
        self.s = (self.s + speed * dt).min(self.total_length);
    }

    #[must_use]
    pub fn is_complete(self) -> bool {
        self.s >= self.total_length
    }
}

// ── Lane center helpers ───────────────────────────────────────────────────────

/// Pixel x-coordinate of the center of lane column `col` (0 = leftmost).
///
/// Columns 0-2 are North-entry / South-exit lanes; columns 3-5 are South-entry / North-exit.
fn col_x(col: u32) -> f32 {
    INTER_LEFT as f32 + col as f32 * LANE_W as f32 + LANE_W as f32 / 2.0
}

/// Pixel y-coordinate of the center of lane row `row` (0 = topmost).
///
/// Rows 0-2 are East-entry / West-exit lanes; rows 3-5 are West-entry / East-exit.
fn row_y(row: u32) -> f32 {
    INTER_TOP as f32 + row as f32 * LANE_W as f32 + LANE_W as f32 / 2.0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn advance_moves_progress_forward() {
        let mut p = PathProgress::new(100.0);
        p.advance(50.0, 1.0);
        assert!((p.s - 50.0).abs() < f32::EPSILON);
    }

    #[test]
    fn advance_clamps_at_total_length() {
        let mut p = PathProgress::new(10.0);
        p.advance(100.0, 1.0);
        assert!((p.s - 10.0).abs() < f32::EPSILON);
    }

    #[test]
    fn is_complete_when_at_end() {
        let mut p = PathProgress::new(10.0);
        p.advance(10.0, 1.0);
        assert!(p.is_complete());
    }

    #[test]
    fn is_not_complete_at_start() {
        assert!(!PathProgress::new(10.0).is_complete());
    }

    #[test]
    fn all_twelve_lane_paths_exist() {
        let paths = [
            RoutePath::for_lane(Direction::North, Route::Right),
            RoutePath::for_lane(Direction::North, Route::Straight),
            RoutePath::for_lane(Direction::North, Route::Left),
            RoutePath::for_lane(Direction::South, Route::Right),
            RoutePath::for_lane(Direction::South, Route::Straight),
            RoutePath::for_lane(Direction::South, Route::Left),
            RoutePath::for_lane(Direction::East, Route::Right),
            RoutePath::for_lane(Direction::East, Route::Straight),
            RoutePath::for_lane(Direction::East, Route::Left),
            RoutePath::for_lane(Direction::West, Route::Right),
            RoutePath::for_lane(Direction::West, Route::Straight),
            RoutePath::for_lane(Direction::West, Route::Left),
        ];
        assert!(paths.iter().all(|p| p.total_length > 0.0));
    }

    #[test]
    fn straight_path_keeps_constant_heading() {
        let path = RoutePath::for_lane(Direction::West, Route::Straight);
        let start = path.sample(0.0);
        let mid = path.sample(path.total_length * 0.5);
        let end = path.sample(path.total_length);

        // West-Straight travels due east → heading 0°.
        assert!((start.heading_deg - 0.0).abs() < 0.01);
        assert!((mid.heading_deg - 0.0).abs() < 0.01);
        assert!((end.heading_deg - 0.0).abs() < 0.01);
    }

    #[test]
    fn left_turn_has_two_distinct_headings() {
        // North-Left: approach heading ≈90° (south), exit heading ≈0° (east).
        // The turn is a single sharp waypoint, not a smooth curve.
        let path = RoutePath::for_lane(Direction::North, Route::Left);
        let start = path.sample(0.0);
        let end = path.sample(path.total_length);

        assert!(start.heading_deg > 60.0 && start.heading_deg < 120.0);
        assert!(end.heading_deg.abs() < 10.0);
    }

    #[test]
    fn turning_path_has_exactly_three_waypoints() {
        // Turn paths are entry → turn-tile → exit: exactly 3 points.
        let path = RoutePath::for_lane(Direction::East, Route::Right);
        assert_eq!(path.point_count(), 3);
    }

    #[test]
    fn straight_path_has_exactly_two_waypoints() {
        let path = RoutePath::for_lane(Direction::East, Route::Straight);
        assert_eq!(path.point_count(), 2);
    }

    #[test]
    fn each_right_turn_uses_its_own_entry_corner() {
        // Each right-turn route turns at the FIRST tile it encounters when
        // entering the intersection, not the last tile before exiting.
        //
        // Corner tile centres (outer col/row at the entry edge):
        //   North-Right → top-left  (col_x(0)=540, row_y(0)=260)
        //   East-Right  → top-right (col_x(5)=740, row_y(0)=260)
        //   South-Right → bot-right (col_x(5)=740, row_y(5)=460)
        //   West-Right  → bot-left  (col_x(0)=540, row_y(5)=460)

        // North-Right: (540,-40)→(540,260)→(−40,260). Approach len = 260+40 = 300.
        let nr = RoutePath::for_lane(Direction::North, Route::Right);
        let p = nr.sample(300.0);
        assert!((p.position.x - 540.0).abs() < 1.0, "NR x");
        assert!((p.position.y - 260.0).abs() < 1.0, "NR y");

        // East-Right: (1320,260)→(740,260)→(740,−40). Approach len = 1320−740 = 580.
        let er = RoutePath::for_lane(Direction::East, Route::Right);
        let p = er.sample(580.0);
        assert!((p.position.x - 740.0).abs() < 1.0, "ER x");
        assert!((p.position.y - 260.0).abs() < 1.0, "ER y");

        // South-Right: (740,760)→(740,460)→(1320,460). Approach len = 760−460 = 300.
        let sr = RoutePath::for_lane(Direction::South, Route::Right);
        let p = sr.sample(300.0);
        assert!((p.position.x - 740.0).abs() < 1.0, "SR x");
        assert!((p.position.y - 460.0).abs() < 1.0, "SR y");

        // West-Right: (−40,460)→(540,460)→(540,760). Approach len = 540+40 = 580.
        let wr = RoutePath::for_lane(Direction::West, Route::Right);
        let p = wr.sample(580.0);
        assert!((p.position.x - 540.0).abs() < 1.0, "WR x");
        assert!((p.position.y - 460.0).abs() < 1.0, "WR y");
    }

    #[test]
    fn turn_waypoints_match_conflict_map_tiles() {
        // Verify that each left-turn path passes through its expected tile center.
        // Inner tile (c,r) → pixel (col_x(c+1), row_y(r+1)).
        //
        // NE (1,2) → (col_x(2), row_y(3)) = (620, 380)
        let nl = RoutePath::for_lane(Direction::North, Route::Left);
        let nl_turn = nl.sample(380.0 + OFFSCREEN_MARGIN); // approach length = 380+40=420
        assert!((nl_turn.position.x - 620.0).abs() < 1.0);
        assert!((nl_turn.position.y - 380.0).abs() < 1.0);

        // SW (2,1) → (col_x(3), row_y(2)) = (660, 340)
        let sl = RoutePath::for_lane(Direction::South, Route::Left);
        // South-Left: (660, 760) → (660, 340) → (−40, 340)
        // Approach length = (H + MARGIN) − row_y(2) = 760 − 340 = 420.
        let sl_turn = sl.sample(760.0 - 340.0);
        assert!((sl_turn.position.x - 660.0).abs() < 1.0);
        assert!((sl_turn.position.y - 340.0).abs() < 1.0);
    }
}
