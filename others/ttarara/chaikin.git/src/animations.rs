//! **Person 3 — Animation + rendering pipeline (logic layer).**
//!
//! Task mapping:
//! - **Step-by-step Chaikin:** precompute all subdivision levels from the control polygon supplied
//!   by the input layer (`start` + `build_steps`), then show one level at a time.
//! - **Seven iterations + loop:** advance `current_step` through the original polygon plus seven
//!   Chaikin steps; after the last step, wrap back to step 0 (initial points).
//! - **Animation timing:** advance the visible step every `STEP_DURATION_SEC` seconds (wall-clock),
//!   driven from `update(now)` each frame.
//! - **Special cases (1 point / 2 points):** not started here when fewer than three control points;
//!   the app keeps showing raw control geometry. `curve_points` still supplies whatever the UI
//!   should draw for the current mode.

use crate::algorithm::build_steps;
use macroquad::prelude::Vec2;

/// Number of Chaikin iterations to precompute (seven), plus step 0 = original polygon → eight frames per cycle.
pub const MAX_ANIMATION_STEPS: usize = 7;

/// How long each visible subdivision level stays on screen before moving to the next (animation speed).
pub const STEP_DURATION_SEC: f64 = 0.6;

/// Holds animation state: cached Chaikin levels, whether playback is running, and timing for step changes.
pub struct ChaikinAnimator {
    /// Precomputed polylines: index `0` = input polygon, then each Chaikin iteration up to `MAX_ANIMATION_STEPS`.
    steps: Vec<Vec<Vec2>>,
    /// True while cycling through `steps` after a successful `start`.
    is_animating: bool,
    /// Index into `steps` for the polyline currently shown (0 … `MAX_ANIMATION_STEPS`).
    current_step: usize,
    /// Last wall time (`get_time()`) at which `current_step` was advanced — used with `STEP_DURATION_SEC`.
    last_step_time: f64,
}

impl ChaikinAnimator {
    /// Initial idle state: no cached steps, not animating.
    pub fn new() -> Self {
        Self {
            steps: Vec::new(),
            is_animating: false,
            current_step: 0,
            last_step_time: 0.0,
        }
    }

    /// Whether the 7-step loop is actively advancing (for UI / control-point overlay).
    pub fn is_animating(&self) -> bool {
        self.is_animating
    }

    /// Current Chaikin level index shown (0 = original control polygon when animating).
    pub fn current_step(&self) -> usize {
        self.current_step
    }

    /// Stop animation and discard precomputed steps (e.g. clear board or Enter with 1–2 points).
    pub fn clear(&mut self) {
        self.steps.clear();
        self.is_animating = false;
        self.current_step = 0;
    }

    /// Stop playback and drop cached geometry (e.g. user clicked to add a point while animating).
    pub fn stop_and_reset_step(&mut self) {
        self.is_animating = false;
        self.current_step = 0;
        self.steps.clear();
    }

    /// Begin the repeating animation: build all Chaikin levels from `control_points`, show step 0, reset timer.
    /// Caller should only invoke this with ≥3 points; otherwise the input layer keeps static 1-point / 2-point views.
    pub fn start(&mut self, control_points: &[Vec2], now: f64) {
        self.steps = build_steps(control_points, MAX_ANIMATION_STEPS);
        self.is_animating = true;
        self.current_step = 0;
        self.last_step_time = now;
    }

    /// Advance to the next subdivision level when enough real time has passed; wrap to step 0 after the last level.
    pub fn update(&mut self, now: f64) {
        if !self.is_animating || self.steps.is_empty() {
            return;
        }

        if now - self.last_step_time >= STEP_DURATION_SEC {
            self.current_step += 1;
            if self.current_step > MAX_ANIMATION_STEPS {
                self.current_step = 0;
            }
            self.last_step_time = now;
        }
    }

    /// Polyline to draw for this frame: during animation, the cached level `steps[current_step]`; otherwise the
    /// live control points from the input layer (used for editing and for 0–2 point special rendering in `App::draw`).
    pub fn curve_points<'a>(&'a self, control_points: &'a [Vec2]) -> &'a [Vec2] {
        if self.steps.is_empty() {
            control_points
        } else {
            &self.steps[self.current_step]
        }
    }
}
