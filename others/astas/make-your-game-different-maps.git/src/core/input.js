/**
 * @file input.js
 * @module input
 * @description
 * Sets up all input handling for player control, including mouse and keyboard events.
 * Controls the crosshair movement, shooting behavior, and directional navigation.
 *
 * Dependencies:
 * - `shooting.js` for the `attemptShot` callback
 * - Shared `state` object for crosshair position and pressed keys
 */

export function createInputSystem({
    gameEl,
    crosshair,
    attemptShot,
    documentRef,
    windowRef,
    constants,
    state,
}) {
    /** Move crosshair to follow mouse position */
    function handleMouseMove(event) {
        const x = event.clientX;
        const y = event.clientY;
        state.setCrosshairX(x);
        state.setCrosshairY(y);
        crosshair.style.transform = `translate3d(${x}px, ${y}px, 0)`;
    }

    /** Fade crosshair when cursor leaves the game window */
    function handleMouseLeave() {
        crosshair.style.opacity = '0.3';
    }

    /** Restore crosshair visibility on re-entry */
    function handleMouseEnter() {
        crosshair.style.opacity = '1';
    }

    /** Trigger shooting on left click if allowed */
    function handleMouseDown(event) {
        if (!state.getRunning() || state.getPaused()) return;
        
        // Don't trigger shooting if clicking on UI elements (buttons, HUD, etc.)
        const target = event.target;
        if (target.closest('button') || target.closest('#hud') || target.closest('[role="dialog"]')) {
            event.stopPropagation();
            event.stopImmediatePropagation();
            return;
        }
        
        // Check if clicking on a fish or game entity - stop all propagation
        if (target.closest('.fish') || target.closest('.sprite') || target.closest('#entities-layer')) {
            event.stopPropagation();
            event.stopImmediatePropagation();
        }
        
        if (event.button === 0) {
            // Stop all event propagation to prevent any other handlers from firing
            event.stopPropagation();
            event.stopImmediatePropagation();
            attemptShot(event.clientX, event.clientY);
        }
    }

    /** Track keydown state and allow spacebar to fire */
    function handleKeyDown(event) {
        if (state.getPaused()) return;
        const keys = state.getKeys();
        keys[event.key] = true;
        if ((event.key === ' ' || event.key === 'Spacebar') && state.getRunning()) {
            attemptShot(state.getCrosshairX(), state.getCrosshairY());
        }
    }

    /** Clear key flag on keyup */
    function handleKeyUp(event) {
        const keys = state.getKeys();
        keys[event.key] = false;
    }

    /**
     * Continuously updates crosshair position based on arrow key input.
     * Runs independently from the main game loop.
     */
    function crosshairLoop(timestamp) {
        const lastTs = state.getLastCrossTs();
        const dt = Math.min(0.05, (timestamp - lastTs) / 1000);
        state.setLastCrossTs(timestamp);

        const keys = state.getKeys();
        let dx = 0;
        let dy = 0;

        // Movement directions
        if (keys['ArrowUp']) dy -= 1;
        if (keys['ArrowDown']) dy += 1;
        if (keys['ArrowLeft']) dx -= 1;
        if (keys['ArrowRight']) dx += 1;

        if (dx !== 0 || dy !== 0) {
            // Normalize diagonal movement
            const length = Math.hypot(dx, dy);
            dx /= length;
            dy /= length;

            const distance = constants.speedPxPerSec * dt;
            let x = state.getCrosshairX() + dx * distance;
            let y = state.getCrosshairY() + dy * distance;

            // Clamp to game boundaries
            const rect = gameEl.getBoundingClientRect();
            const minX = rect.left + constants.crosshairHalf;
            const maxX = rect.right - constants.crosshairHalf;
            const minY = rect.top + constants.crosshairHalf;
            const maxY = rect.bottom - constants.crosshairHalf;

            x = Math.max(minX, Math.min(maxX, x));
            y = Math.max(minY, Math.min(maxY, y));

            state.setCrosshairX(x);
            state.setCrosshairY(y);
        }

        // Render updated crosshair position
        crosshair.style.transform = `translate3d(${state.getCrosshairX()}px, ${state.getCrosshairY()}px, 0)`;
        windowRef.requestAnimationFrame(crosshairLoop);
    }

    /** Initialize all event listeners and start crosshair motion loop */
    function init() {
        state.setLastCrossTs(windowRef.performance.now());
        gameEl.addEventListener('mousemove', handleMouseMove);
        gameEl.addEventListener('mouseleave', handleMouseLeave);
        gameEl.addEventListener('mouseenter', handleMouseEnter);
        gameEl.addEventListener('mousedown', handleMouseDown);
        documentRef.addEventListener('keydown', handleKeyDown);
        documentRef.addEventListener('keyup', handleKeyUp);
        windowRef.requestAnimationFrame(crosshairLoop);
    }

    return { init };
}
