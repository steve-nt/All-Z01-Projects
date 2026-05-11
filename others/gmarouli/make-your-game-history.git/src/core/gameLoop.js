/**
 * @file gameLoop.js
 * @module gameLoop
 * @description
 * Provides the main update and render loop that drives all in-game activity.
 * Handles entity updates, spawning, timing, and transitions between levels.
 *
 * Dependencies:
 * - Entity managers (`lifeFishManager`, `turtleManager`)
 * - Gameplay spawners (`spawnBubble`, `spawnParrot`)
 * - UI feedback (`hud`, `feedback`)
 * - Shared state getters/setters from main.js
 *
 * Exports:
 * - `createGameLoop()` → returns { update, render, frame }
 */

export function createGameLoop({
    lifeFishManager,
    turtleManager,
    spawnBubble,
    spawnParrot,
    positionElement,
    removeEntity,
    entitiesLayer,
    bubblesLayer,
    world,
    constants,
    hud,
    feedback,
    state,
    rand,
    storyManager,
}) {
    /**
     * Main update routine.
     * @param {number} dt - Time delta since last frame (seconds)
     */
    function update(dt) {
        // Check for score milestones and show development story scenes
        // Only check if game is running and not already paused
        if (storyManager && state.getRunning() && !state.getPaused()) {
            const currentScore = state.getScore();
            const showedScene = storyManager.showDevelopment(currentScore, () => {
                // Resume gameplay after development scene
                state.setPaused(false);
            });
            if (showedScene) {
                // Pause game while showing story
                state.setPaused(true);
                return; // Don't update game state while story is showing
            }
        }

        // Update special entity spawn timers
        const lifeReady = lifeFishManager.updateSpawnTimer(dt, state.getLives(), constants.maxLives, rand);
        turtleManager.updateSpawnTimer(dt, rand);

        // If life fish timer elapsed, spawn it immediately
        if (lifeReady && !lifeFishManager.isActive()) {
            const nextId = state.getNextEntityId();
            const entity = lifeFishManager.spawnLifeFish(
                entitiesLayer,
                state.getEntities(),
                nextId,
                world,
                constants.parrotBaseSpeed,
                rand,
                feedback.showCenterNotification
            );
            positionElement(entity);
            lifeFishManager.markSpawnedThisLevel();
            state.setNextEntityId(nextId + 1);
        }

        // Handle combo countdown and reset when expired
        if (state.getComboTimer() > 0) {
            const timer = state.getComboTimer() - dt;
            state.setComboTimer(timer);
            if (timer <= 0) {
                state.setCombo(0);
                feedback.updateComboDisplay();
            }
        }

        // Decrease time and level timers
        let timeLeft = state.getTimeLeft() - dt;
        state.setTimeLeft(timeLeft);
        let levelTimer = state.getLevelTimer() - dt;
        state.setLevelTimer(levelTimer);

        // Level transition logic
        if (timeLeft <= 0) {
            if (state.getLevel() < constants.maxLevel) {
                // Advance to next level
                state.setLevel(state.getLevel() + 1);
                state.setTimeLeft(constants.levelDuration);
                state.setLevelTimer(constants.levelDuration);
                lifeFishManager.resetForNewLevel();
                turtleManager.resetForNewLevel();
                feedback.showCenterNotification(`LEVEL ${state.getLevel()}`, 'level-up', 2000);
            } else {
                // Game ends after final level
                state.setTimeLeft(0);
                feedback.showGameOver();
                return;
            }
        }

        // Handle generic fish spawning
        let spawnTimer = state.getSpawnTimer() - dt;
        state.setSpawnTimer(spawnTimer);
        if (spawnTimer <= 0) {
            const spawnResult = spawnParrot({
                lifeFishManager,
                lives: state.getLives(),
                maxLives: constants.maxLives,
                rand,
                entitiesLayer,
                entities: state.getEntities(),
                nextEntityId: state.getNextEntityId(),
                world,
                parrotBaseSpeed: constants.parrotBaseSpeed,
                baseMaxFish: constants.baseMaxFish,
                level: state.getLevel(),
                showCenterNotification: feedback.showCenterNotification,
            });
            if (spawnResult && typeof spawnResult.nextEntityId === 'number') {
                state.setNextEntityId(spawnResult.nextEntityId);
            }

            // Spawn interval scales down slightly with level
            const levelSpawnMin = Math.max(0.4, constants.baseSpawnMin * Math.pow(0.94, state.getLevel() - 1));
            const levelSpawnMax = Math.max(0.9, constants.baseSpawnMax * Math.pow(0.94, state.getLevel() - 1));
            state.setSpawnTimer(levelSpawnMin + rand() * (levelSpawnMax - levelSpawnMin));
        }

        // Spawn turtle hazard when conditions met
        if (!turtleManager.turtleSpawnedThisLevel && !turtleManager.turtleActive && turtleManager.readyToSpawn) {
            let nextId = state.getNextEntityId();
            const turtle = turtleManager.spawnTurtle(
                entitiesLayer,
                state.getEntities(),
                nextId,
                world,
                constants.parrotBaseSpeed,
                rand
            );
            positionElement(turtle);
            turtleManager.markSpawnedThisLevel();
            turtleManager.readyToSpawn = false;
            state.setNextEntityId(nextId + 1);
        }

        // Bubble ambience spawner
        let bubbleTimer = state.getBubbleTimer() - dt;
        state.setBubbleTimer(bubbleTimer);
        if (bubbleTimer <= 0) {
            const count = rand() < 0.15 ? 2 : 1;
            for (let i = 0; i < count; i++) {
                // Route bubbles into the entities layer so they render above decor
                spawnBubble({ bubblesLayer: entitiesLayer, world, rand });
            }
            state.setBubbleTimer(0.5 + rand() * 0.9);
        }

        // Update entity positions and check bounds
        const entities = state.getEntities();
        for (let i = entities.length - 1; i >= 0; i--) {
            const entity = entities[i];
            if (!entity.alive) continue;

            // Movement update (special case for life fish)
            if (!entity.dead) {
                if (entity.isLifeFish) {
                    lifeFishManager.updateLifeFishMovement(entity, dt);
                }
                entity.x += entity.vx * dt;
                entity.y += entity.vy * dt;

                // Small idle wave for normal fish
                if (!entity.isLifeFish) {
                    entity.y += Math.sin(performance.now() / 450 + entity.id) * 0.4;
                }
            }

            // Out-of-bounds cleanup and penalty handling
            if (
                entity.x < -80 ||
                entity.x > world.width + 80 ||
                entity.y < -80 ||
                entity.y > world.height + 80
            ) {
                entity.alive = false;
                removeEntity(entity);

                if (entity.isLifeFish) {
                    lifeFishManager.handleLifeFishEscaped();
                } else if (entity.isTurtle) {
                    turtleManager.handleTurtleEscaped();
                } else if (!entity.dead) {
                    // Normal fish escape: lose one life
                    const newLives = Math.max(0, state.getLives() - 1);
                    state.setLives(newLives);
                    if (newLives > 0) {
                        feedback.showCenterNotification('💔', 'life-lost', 1500);
                    }
                }

                if (state.getLives() === 0) {
                    feedback.showGameOver();
                    return;
                }
            }
        }
    }

    /** Renders entities to DOM positions. */
    function render() {
        const entities = state.getEntities();
        for (let i = 0; i < entities.length; i++) {
            const entity = entities[i];
            if (!entity.alive) continue;
            positionElement(entity);
        }
    }

    /** Core animation frame callback for `requestAnimationFrame`. */
    function frame(timestamp) {
        if (!state.getRunning()) return;

        // Initialize timestamp baseline on first frame
        if (state.getLastFrameTs() === 0) state.setLastFrameTs(timestamp);
        const dt = Math.min(0.05, (timestamp - state.getLastFrameTs()) / 1000);
        state.setLastFrameTs(timestamp);

        // Approximate FPS tracking over short intervals
        // --- Accurate per-frame FPS like Chrome DevTools ---
        const deltaMs = timestamp - state.getLastFrameTs();
        const instantaneousFps = 1000 / deltaMs;
        // Avoid NaN/infinity on first frame
        if (Number.isFinite(instantaneousFps)) {
            state.setFps(instantaneousFps);
        }

        // Game logic update only if not paused
        if (!state.getPaused()) {
            update(dt);
            render();
            hud.updateHud();
        }

        requestAnimationFrame(frame);
    }

    return { update, render, frame };
}
