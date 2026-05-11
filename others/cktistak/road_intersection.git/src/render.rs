use macroquad::{color::Color, shapes::draw_rectangle};

pub const W: u32 = 800;
pub const H: u32 = 800;

/// Screen center (same as intersection center).
pub const CX: f32 = W as f32 / 2.;
pub const CY: f32 = H as f32 / 2.;

/// Half-width of one *road* (both directions / both lanes together).
/// Full road width on screen = 2 * ROAD_HALF.
pub const ROAD_HALF: f32 = 36.;

pub fn window_size() -> (u32, u32) {
    (W, H)
}

/// Two filled rectangles: vertical strip + horizontal strip → cross shape.
pub fn draw_roads() {
    let color = Color::from_rgba(48, 52, 58, 255);
    draw_rectangle(CX - ROAD_HALF, 0., 2. * ROAD_HALF, H as f32, color);
    draw_rectangle(0., CY - ROAD_HALF, W as f32, 2. * ROAD_HALF, color);
}

/// One lane each way: dashed line down the middle of each road band.
/// Vertical road → split at x = CX. Horizontal road → split at y = CY.
pub fn draw_lane_dividers() {
    let color = Color::from_rgba(220, 220, 210, 255);

    let dash: f32 = 14.;
    let gap: f32 = 10.;
    let margin = 8.;

    // Middle of north–south road (two lanes: e.g. up vs down).
    let mid_x = CX;
    let mut y = margin;
    while y < H as f32 - margin {
        draw_rectangle(mid_x - 1., y, 3., dash, color);
        y += dash + gap;
    }

    // Middle of east–west road (two lanes: e.g. left vs right).
    let mid_y: f32 = CY as f32;
    let mut x = margin;
    while x < W as f32 - margin {
        draw_rectangle(x, mid_y - 1., dash, 3., color);
        x += dash + gap;
    }
}

/// Per approach: `true` = green, `false` = red (no yellow phase in this milestone).
#[derive(Clone, Copy, Debug)]
pub struct TrafficLightState {
    pub north: bool,
    pub east: bool,
    pub south: bool,
    pub west: bool,
}

impl Default for TrafficLightState {
    /// N/S green, E/W red — placeholder until Person C drives timing.
    fn default() -> Self {
        Self {
            north: true,
            east: false,
            south: true,
            west: false,
        }
    }
}

/// Thick white bars at the four edges of the intersection box (one per approach).
pub fn draw_stop_lines() {
    let color = Color::from_rgba(245, 245, 238, 255);
    let t: f32 = 4.;
    let h = t as f32 / 2.;

    draw_rectangle(CX - ROAD_HALF, CY - ROAD_HALF - h, 2. * ROAD_HALF, t, color);
    draw_rectangle(CX - ROAD_HALF, CY + ROAD_HALF - h, 2. * ROAD_HALF, t, color);
    draw_rectangle(CX - ROAD_HALF - h, CY - ROAD_HALF, t, 2. * ROAD_HALF, color);
    draw_rectangle(CX + ROAD_HALF - h, CY - ROAD_HALF, t, 2. * ROAD_HALF, color);
}

// One small square per approach: green or red (no yellow).
pub fn draw_traffic_lights(state: TrafficLightState) {
    let d: f32 = 10.;
    let r = d / 2.;
    let o = ROAD_HALF + 14.;
    let c = |g: bool| {
        if !g {
            Color::from_rgba(50, 210, 80, 255)
        } else {
            Color::from_rgba(235, 45, 45, 255)
        }
    };
    draw_rectangle(CX + o - r, CY - ROAD_HALF - 14. - r, d, d, c(state.north));
    draw_rectangle(CX - o - r, CY + ROAD_HALF + 14. - r, d, d, c(state.south));
    draw_rectangle(CX - ROAD_HALF - 14. - r, CY - o - r, d, d, c(state.west));
    draw_rectangle(CX + ROAD_HALF + 14. - r, CY + o - r, d, d, c(state.east));
}
