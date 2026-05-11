/**
 * Produces HUD and pause menu render helpers that reflect current game state.
 */

export function createHudSystem({
    hudTime,
    hudScore,
    hudLives,
    hudFps,
    hudLevel,
    menuTime,
    menuScore,
    menuLives,
    menuFps,
    menuLevel,
    documentRef,
    getTimeLeft,
    getScore,
    getFps,
    getLevel,
    getLives,
    getHighScore,
}) {
    const ACTIVE_HEART = '❤️';
    const LOST_HEART = '🖤';

    function updateHearts(container) {
        if (!container) return;

        const hearts = container.querySelectorAll('.heart');
        hearts.forEach((heart, index) => {
            const isActive = index < getLives();
            heart.textContent = isActive ? ACTIVE_HEART : LOST_HEART;

            if (isActive) {
                heart.classList.add('active');
                heart.classList.remove('lost');
            } else {
                heart.classList.remove('active');
                heart.classList.add('lost');
            }
        });
    }

    function updateHud() {
        if (hudTime) hudTime.textContent = getTimeLeft().toFixed(1);
        if (hudScore) hudScore.textContent = String(getScore());
        if (hudFps) hudFps.textContent = getFps().toFixed(1);
        if (hudLevel) hudLevel.textContent = String(getLevel());

        if (hudTime) {
            if (getTimeLeft() <= 10 && !hudTime.classList.contains('low-time')) {
                hudTime.classList.add('low-time');
            } else if (getTimeLeft() > 10 && hudTime.classList.contains('low-time')) {
                hudTime.classList.remove('low-time');
            }
        }

        updateHearts(hudLives);

        const hudHighScore = documentRef.getElementById('hud-highscore');
        if (hudHighScore) hudHighScore.textContent = String(getHighScore());
    }

    function updateMenu() {
        if (menuTime) menuTime.textContent = getTimeLeft().toFixed(1);
        if (menuScore) menuScore.textContent = String(getScore());
        if (menuFps) menuFps.textContent = getFps().toFixed(1);
        if (menuLevel) menuLevel.textContent = String(getLevel());

        updateHearts(menuLives);
    }

    return { updateHud, updateMenu };
}
