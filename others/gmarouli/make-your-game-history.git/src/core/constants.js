/**
 * @file constants.js
 * @module constants
 * @description
 * Defines core game configuration constants used across multiple modules.
 * These values control global game pacing, difficulty progression, and
 * entity behavior. Shared between logic modules to maintain consistency.
 *
 * Exported constants are immutable and referenced by modules such as:
 * - `gameLoop.js` (for spawn logic and level timing)
 * - `shooting.js` (for combo and rate limit timing)
 * - `lifecycle.js` (for resets and initialization)
 */

export const GAME_DURATION_S = 30;        // Seconds per level duration
export const MAX_LEVEL = 10;              // Total number of playable levels
export const WORLD = { width: 0, height: 0 }; // Runtime-populated world bounds
export const PARROT_BASE_SPEED = 90;      // Base movement speed of fish entities
export const BASE_SPAWN_MIN = 1.0;        // Minimum spawn interval (seconds)
export const BASE_SPAWN_MAX = 2.0;        // Maximum spawn interval (seconds)
export const BASE_MAX_FISH = 5;           // Base maximum number of concurrent fish
export const LEVEL_DURATION = 30;         // Level time limit (seconds)
export const MAX_LIVES = 3;               // Player starting lives
export const SHOT_RATE_LIMIT_MS = 90;     // Minimum delay between shots (ms)
export const COMBO_WINDOW = 2.5;          // Combo reset window (seconds)
