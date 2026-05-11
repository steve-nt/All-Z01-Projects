/**
 * Bundles UI feedback helpers for scores, combos, notifications, and Game Over.
 */

export function createFeedbackSystem({
    entitiesLayer,
    comboDisplay,
    centerNotification,
    gameoverOverlay,
    finalScore,
    finalAccuracy,
    finalCombo,
    finalLevel,
    finalCaught,
    celebrateHighScore,
    updateHighScoreListUI,
    addHighScore,
    state,
}) {
    function showScorePopup(x, y, points, comboCount) {
        if (!entitiesLayer) return;
        const popup = document.createElement('div');
        popup.className = 'score-popup';
        const comboText = comboCount > 1 ? ` (${comboCount}x)` : '';
        popup.textContent = `+${points}${comboText}`;
        popup.style.left = `${x}px`;
        popup.style.top = `${y}px`;
        entitiesLayer.appendChild(popup);
        setTimeout(() => {
            if (popup.parentNode) popup.parentNode.removeChild(popup);
        }, 1000);
    }

    function updateComboDisplay() {
        if (!comboDisplay) return;
        if (state.getCombo() >= 2) {
            comboDisplay.textContent = `${state.getCombo()}x COMBO!`;
            comboDisplay.classList.add('active');
        } else {
            comboDisplay.classList.remove('active');
        }
    }

    function showCenterNotification(message, type = '', duration = 2000) {
        if (!centerNotification) return;

        centerNotification.textContent = message;
        centerNotification.className = `show ${type}`;

        setTimeout(() => {
            centerNotification.classList.remove('show');
        }, duration);
    }

    function showGameOver() {
        state.setRunning(false);
        state.setPaused(true);

        showCenterNotification('GAME OVER', 'life-lost', 1500);
        if (!gameoverOverlay) return;

        const totalShots = state.getTotalShots();
        const totalHits = state.getTotalHits();
        const accuracy = totalShots > 0 ? Math.round((totalHits / totalShots) * 100) : 0;

        if (finalScore) finalScore.textContent = String(state.getScore());
        if (state.getScore() > state.getHighScore()) {
            state.setHighScore(state.getScore());
            addHighScore(state.getScore());
            showCenterNotification('🏆 NEW HIGH SCORE!', 'level-up', 2500);
            celebrateHighScore();
        }

        if (finalAccuracy) finalAccuracy.textContent = `${accuracy}%`;
        if (finalCombo) finalCombo.textContent = `${state.getMaxCombo()}x`;
        if (finalLevel) finalLevel.textContent = String(state.getLevel());
        if (finalCaught) finalCaught.textContent = String(state.getFishCaught());

        setTimeout(() => {
            gameoverOverlay.classList.remove('hidden');
        }, 1600);

        updateHighScoreListUI();
    }

    function hideGameOver() {
        if (gameoverOverlay) gameoverOverlay.classList.add('hidden');
    }

    return {
        showScorePopup,
        updateComboDisplay,
        showCenterNotification,
        showGameOver,
        hideGameOver,
    };
}
