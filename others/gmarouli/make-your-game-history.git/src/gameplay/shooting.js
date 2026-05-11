/**
 * @file shooting.js
 * @module shooting
 * @description
 * Handles all player shooting interactions:
 * - Detects hits on fish entities
 * - Updates score, combos, and lives
 * - Displays feedback (popups, combo messages)
 * - Applies penalties for missed shots
 *
 * Dependencies:
 * - `feedback.js` for visual feedback and combo updates
 * - `lifeFish.js` and `turtle.js` for special entity handling
 * - `gameLoop.js` for live entity list
 */

export function createShootingSystem({
    gameEl,
    entitiesLayer,
    lifeFishManager,
    turtleManager,
    removeEntity,
    rand,
    constants,
    state,
    feedback,
}) {
    // 🔊 Preload SFX once (fail-safe on browsers blocking autoplay)
    const hitSfx = new Audio('sounds/Fish-hit.mp3');
    const lifeFishSfx = new Audio('sounds/ExtraLife.mp3');
    hitSfx.volume = 0.5;       // Optional: hook to a global SFX volume slider
    lifeFishSfx.volume = 0.6;  // Optional: as above

    /**
     * Handles the "octopus capture" animation that pulls a caught fish away.
     * This visual effect triggers slightly after a fish dies.
     */
    function catchWithOctopus(entity) {
        if (!entity || !entity.dead) return;
        entity.capturing = true;

        const startX = entity.x;
        const startY = entity.y;
        const startScale = entity.scale || 1;
        const captureDuration = 260; // ms
        const t0 = performance.now();

        function captureStep(now) {
            if (!entity.capturing) return;
            const t = Math.min(1, (now - t0) / captureDuration);
            const k = 1 - Math.pow(1 - t, 2); // smooth ease-in
            const scale = startScale * (1 - 0.75 * k);
            const rotation = t * (entity.flip >= 0 ? 15 : -15);
            const scaleX = entity.flip || 1;

            entity.el.style.transform = `
                translate3d(${entity.x - entity.width / 2}px, ${entity.y - entity.height / 2}px, 0)
                scale(${scale}) scaleX(${scaleX}) rotate(${rotation}deg)
            `;

            if (t < 1) {
                requestAnimationFrame(captureStep);
            } else {
                // End of capture
                entity.alive = false;
                entity.dead = false;
                entity.capturing = false;
                removeEntity(entity);
            }
        }

        requestAnimationFrame(captureStep);
    }

    /**
     * Displays a floating red “miss” popup and deducts points.
     * @param {number} x - Screen X coordinate.
     * @param {number} y - Screen Y coordinate.
     * @param {number} penalty - Points deducted.
     */
    function showMissPopup(x, y, penalty) {
        if (!entitiesLayer) return;
        const popup = document.createElement('div');
        popup.className = 'miss-feedback';
        popup.textContent = `-${penalty}`;
        popup.style.left = `${x}px`;
        popup.style.top = `${y}px`;
        entitiesLayer.appendChild(popup);
        setTimeout(() => popup.parentNode?.removeChild(popup), 1500);
    }

    /**
     * Attempts a single shot based on mouse click position.
     * Handles hit detection, scoring, combos, and misses.
     */
    function attemptShot(clientX, clientY) {
        if (!gameEl) return;

        const now = performance.now();
        if (now - state.getLastShotTimestamp() < constants.shotRateLimitMs) return; // Rate limit
        state.setLastShotTimestamp(now);
        state.incrementTotalShots();

        // Translate click position to game coordinates
        const rect = gameEl.getBoundingClientRect();
        const x = clientX - rect.left;
        const y = clientY - rect.top;

        const hitPadding = 10; // leniency
        let hitSomething = false;
        const entities = state.getEntities();

        // Iterate in reverse to prioritize top-most entities
        for (let i = entities.length - 1; i >= 0; i--) {
            const entity = entities[i];
            if (!entity.alive) continue;

            const ex = entity.x - entity.width / 2;
            const ey = entity.y - entity.height / 2;

            // Simple AABB hit check
            if (
                x >= ex - hitPadding &&
                x <= ex + entity.width + hitPadding &&
                y >= ey - hitPadding &&
                y <= ey + entity.height + hitPadding
            ) {
                hitSomething = true;

                // 🔊 Play suitable SFX (safe try/catch for autoplay policies)
                try {
                    if (entity.isLifeFish) {
                        lifeFishSfx.currentTime = 0;
                        lifeFishSfx.play();
                    } else {
                        hitSfx.currentTime = 0;
                        hitSfx.play();
                    }
                } catch (_) {}

                // 🐢 TURTLE — hazard, deduct 30 points but DON'T kill it
                if (entity.isTurtle || entity.variant === 'turtle') {
                    const turtlePenalty = 30;
                    const currentScore = state.getScore();
                    const newScore = Math.max(0, currentScore - turtlePenalty);
                    state.setScore(newScore);

                    // Visual feedback - flash the turtle
                    entity.el.classList.add('hit');
                    setTimeout(() => {
                        if (entity.el) entity.el.classList.remove('hit');
                    }, 150);

                    // Show penalty popup only if score was above 0 before penalty
                    if (currentScore > 0) {
                        const actualPenalty = currentScore - newScore; // Actual penalty applied (may be less than turtlePenalty)
                        showMissPopup(entity.x, entity.y, actualPenalty);
                    }

                    // Reset combo & update
                    state.setCombo(0);
                    state.setComboTimer(0);
                    feedback.updateComboDisplay();

                    // Don't mark as dead, don't capture - turtle keeps swimming
                    return;
                }

                if (entity.dead) return;
                entity.dead = true;

                // Tracking stats
                state.incrementTotalHits();
                state.incrementFishCaught();

                // Life Fish — reward extra life, reset combo
                if (entity.isLifeFish) {
                    const newLives = lifeFishManager.handleLifeFishCaught(
                        state.getLives(),
                        constants.maxLives,
                        feedback.showCenterNotification
                    );
                    state.setLives(newLives);
                    state.setCombo(0);
                    state.setComboTimer(0);
                    feedback.updateComboDisplay();
                }
                // (Legacy path) Regular turtle handler if flagged separately
                else if (entity.isTurtle) {
                    state.setCombo(0);
                    state.setComboTimer(0);
                    feedback.updateComboDisplay();
                    turtleManager.handleTurtleHit(() => {});
                }
                // Regular fish — increment combo and score
                else {
                    const combo = state.getCombo() + 1;
                    state.setCombo(combo);
                    if (combo > state.getMaxCombo()) state.setMaxCombo(combo);
                    state.setComboTimer(constants.comboWindow);
                    feedback.updateComboDisplay();

                    // Base score scaling by variant and level
                    const basePoints = {
                        fish1: 10,
                        fish2: 15,
                        fish3: 20,
                        fish4: 25,
                        fish5: 30,
                    }[entity.variant] || 10;

                    const points =
                        basePoints *
                        Math.max(1, state.getLevel()) *
                        Math.max(1, combo);

                    state.setScore(state.getScore() + points);
                    // Use click position for popup, or entity center as fallback
                    const popupX = x; // Click position in game coordinates
                    const popupY = y; // Click position in game coordinates
                    feedback.showScorePopup(popupX, popupY, points, combo);
                }

                // Trigger death animation (sink and capture)
                entity.el.classList.add('dead');
                const startX = entity.x;
                const startY = entity.y;
                const rotation = (rand() < 0.5 ? -1 : 1) * (10 + rand() * 20);
                const sinkDuration = 500 + rand() * 350;
                const t0 = performance.now();

                function sinkStep(nowTs) {
                    if (!entity.dead) return;
                    const t = Math.min(1, (nowTs - t0) / sinkDuration);
                    const sx = entity.flip || 1;
                    const scale = 1 - t * 0.2;
                    const rot = t * rotation;
                    entity.y = startY + t * 90;
                    entity.el.style.transform = `
                        translate3d(${startX - entity.width / 2}px, ${entity.y - entity.height / 2}px, 0)
                        scale(${scale}) scaleX(${sx}) rotate(${rot}deg)
                    `;
                    if (t < 1) requestAnimationFrame(sinkStep);
                }

                requestAnimationFrame(sinkStep);
                setTimeout(() => catchWithOctopus(entity), 220);
                return;
            }
        }

        // Missed all entities → penalty
        if (!hitSomething) {
            state.incrementMissedShots();
            const missPenalty = 10;
            const currentScore = state.getScore();
            const newScore = Math.max(0, currentScore - missPenalty);
            state.setScore(newScore);
            
            // Show penalty popup only if score was above 0 before penalty
            if (currentScore > 0) {
                const actualPenalty = currentScore - newScore; // Actual penalty applied (may be less than missPenalty)
                showMissPopup(x, y, actualPenalty);
            }
            
            state.setCombo(0);
            state.setComboTimer(0);
            feedback.updateComboDisplay();
        }
    }

    return { attemptShot };
}
