/**
 * Persists and renders high score data within localStorage-backed menus.
 */

export function loadHighScores() {
    try {
        const stored = localStorage.getItem('highScores');
        if (stored) return JSON.parse(stored);
    } catch (_) {
        // no-op: localStorage access failed
    }
    return [];
}

export function saveHighScores(scores) {
    try {
        localStorage.setItem('highScores', JSON.stringify(scores));
    } catch (_) {
        // no-op: localStorage access failed
    }
}

export function addHighScore(newScore) {
    const scores = loadHighScores();
    scores.push({ score: newScore, date: new Date().toISOString() });
    scores.sort((a, b) => b.score - a.score);
    const topScores = scores.slice(0, 5);
    saveHighScores(topScores);
    return topScores;
}

export function updateHighScoreListUI() {
    const listEl = document.getElementById('high-score-list');
    if (!listEl) return;

    const scores = loadHighScores();
    listEl.innerHTML = '';

    if (scores.length === 0) {
        const li = document.createElement('li');
        li.textContent = 'No scores yet';
        listEl.appendChild(li);
        return;
    }

    scores.forEach((entry, index) => {
        const li = document.createElement('li');
        li.className = 'score-item';
        li.innerHTML = `
            <span class="score-rank">${index + 1}.</span>
            <span class="score-value">${entry.score.toLocaleString()}</span>
        `;
        listEl.appendChild(li);
    });
}
