/// Window dimensions in pixels.
pub const WINDOW_W: u32 = 1280;
pub const WINDOW_H: u32 = 720;

/// Pixels per lane.
pub const LANE_W: u32 = 40;

/// Number of lanes per direction (inbound or outbound).
pub const LANES_PER_DIR: u32 = 3;

/// Total road width: 3 inbound + 3 outbound lanes.
pub const ROAD_W: u32 = LANE_W * LANES_PER_DIR * 2; // 240 px

/// Left pixel column where the vertical (N/S) road starts.
pub const INTER_LEFT: u32 = (WINDOW_W - ROAD_W) / 2; // 520

/// Top pixel row where the horizontal (E/W) road starts.
pub const INTER_TOP: u32 = (WINDOW_H - ROAD_W) / 2; // 240

/// Right pixel column where the vertical road ends.
pub const INTER_RIGHT: u32 = INTER_LEFT + ROAD_W; // 760

/// Bottom pixel row where the horizontal road ends.
pub const INTER_BOTTOM: u32 = INTER_TOP + ROAD_W; // 480

/// Fixed simulation step (seconds).
pub const FIXED_DT: f32 = 1.0 / 60.0;

/// Target frame rate used to compute the fixed step.
pub const TARGET_FPS: u64 = 60;

/// Physical length of a vehicle sprite in world units.
/// Used by the safety distance system and spawn clearance checks.
pub const VEHICLE_LENGTH: f32 = 36.0;

/// Minimum empty gap between the front of a follower and the rear of the leader.
///
/// Spawn checks add this to `VEHICLE_LENGTH` because they compare center-to-center
/// distance. Lane following uses this value directly because it already subtracts
/// vehicle length from the path-progress distance.
pub const SAFETY_BUFFER: f32 = 30.0;
