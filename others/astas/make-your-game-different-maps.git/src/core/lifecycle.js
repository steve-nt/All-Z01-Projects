/**
 * @file lifecycle.js
 * @module lifecycle
 * @description
 * Manages high-level game lifecycle states: start, pause, resume, restart, and returning
 * to the main menu. Coordinates entity resets, overlay visibility, and state resets.
 *
 * Works closely with:
 * - `hud.js` for UI updates
 * - `feedback.js` for endgame handling
 * - `gameLoop.js` for restarting the loop
 */

export function createLifecycleSystem({
    lifeFishManager,
    turtleManager,
    clearEntities,
    clearCelebration,
    updateComboDisplay,
    clearCenterNotification,
    showCenterNotification,
    updateHud,
    updateMenu,
    updateHighScoreListUI,
    hideGameOver,
    createPauseBubbles,
    pauseOverlay,
    startOverlay,
    requestFrame,
    frame,
    constants,
    state,
    storyManager,
}) {
    /** Show pause overlay with live stats */
    function showPauseMenu() {
        updateMenu();
        createPauseBubbles();
        pauseOverlay.classList.remove('hidden');
    }

    /** Hide pause overlay */
    function hidePauseMenu() {
        pauseOverlay.classList.add('hidden');
    }

    /** Display start menu (main menu) */
    function showStartMenu() {
        if (startOverlay) startOverlay.classList.remove('hidden');
        state.setPaused(true);
    }

    /** Hide main start menu */
    function hideStartMenu() {
        if (startOverlay) startOverlay.classList.add('hidden');
    }

    /** Begin/resume frame updates after initialization */
    function start() {
        state.setRunning(true);
        state.setPaused(false);
        requestFrame(frame);
    }

    /** Resets all game-relevant runtime state to initial defaults */
    function resetCoreState() {
        clearCelebration();
        clearEntities(state.getEntities());
        state.setTimeLeft(constants.gameDuration);
        state.setScore(0);
        state.setLives(constants.maxLives);
        state.setSpawnTimer(0.2);
        // Ensure ambient bubbles start spawning soon after reset
        if (typeof state.setBubbleTimer === 'function') {
            state.setBubbleTimer(0.4);
        }
        // Reset entity id counter for new spawns
        if (typeof state.setNextEntityId === 'function') {
            state.setNextEntityId(1);
        }
        state.setRngSeed((Math.random() * 1e9) | 0);
        state.setLastFrameTs(0);
        state.setLevel(1);
        state.setLevelTimer(constants.levelDuration);
        state.setCombo(0);
        state.setMaxCombo(0);
        state.setComboTimer(0);
        state.setTotalShots(0);
        state.setTotalHits(0);
        state.setFishCaught(0);
        state.setMissedShots(0);
        updateComboDisplay();
    }

    /** Pause gameplay; optionally show pause menu */
    function pauseGame(forceMenu = false) {
        state.setPaused(true);
        if (forceMenu) showPauseMenu();
    }

    /** Resume gameplay and close pause menu */
    function resumeGame() {
        state.setPaused(false);
        hidePauseMenu();
    }

    /** Restart current session while keeping the game running */
    function restartGame() {
        lifeFishManager.reset();
        turtleManager.reset();
        resetCoreState();
        if (clearCenterNotification) clearCenterNotification();
        hidePauseMenu();
        hideGameOver();
        state.setPaused(false);
        state.setRunning(true);
        // Use requestAnimationFrame to ensure overlays are hidden before starting the game
        // This ensures notifications can display properly
        requestAnimationFrame(() => {
            requestFrame(frame);
            // Show LEVEL 1 notification after the game starts
            if (showCenterNotification) {
                // Small delay to ensure overlay is fully hidden
                setTimeout(() => {
                    showCenterNotification('LEVEL 1', 'level-up', 2000);
                }, 100);
            }
        });
    }

    /** Fully return to main menu (clear entities, stop loop) */
    function returnToMainMenu() {
        state.setRunning(false);
        state.setPaused(true);
        clearEntities(state.getEntities());
        hidePauseMenu();
        hideGameOver();
        updateHighScoreListUI();
        showStartMenu();
        state.setTimeLeft(constants.gameDuration);
        state.setScore(0);
        state.setLives(constants.maxLives);
        state.setLevel(1);
        state.setLevelTimer(constants.levelDuration);
        state.setCombo(0);
        state.setMaxCombo(0);
        state.setComboTimer(0);
        state.setTotalShots(0);
        state.setTotalHits(0);
        state.setFishCaught(0);
        state.setMissedShots(0);
        updateHud();
        updateComboDisplay();
    }

    /** Start a completely new game run */
    function startNewRun() {
        lifeFishManager.reset();
        turtleManager.reset();
        resetCoreState();
        if (clearCenterNotification) clearCenterNotification();
        if (storyManager) {
            storyManager.reset();
        }
        hideStartMenu();
        hidePauseMenu();
        hideGameOver();
        start();
        // Show LEVEL 1 notification after the game starts
        // Use multiple requestAnimationFrame calls to ensure overlays are fully hidden
        if (showCenterNotification) {
            requestAnimationFrame(() => {
                requestAnimationFrame(() => {
                    setTimeout(() => {
                        showCenterNotification('LEVEL 1', 'level-up', 2000);
                    }, 300);
                });
            });
        }
    }

    return {
        pauseGame,
        resumeGame,
        restartGame,
        returnToMainMenu,
        showPauseMenu,
        hidePauseMenu,
        showStartMenu,
        hideStartMenu,
        start,
        startNewRun,
    };
}
