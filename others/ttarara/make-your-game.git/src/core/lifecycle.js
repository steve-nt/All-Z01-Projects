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
        hidePauseMenu();
        hideGameOver();
        state.setPaused(false);
        state.setRunning(true);
        requestFrame(frame);
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
        hideStartMenu();
        start();
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
