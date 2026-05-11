import { fetchScores } from '../api/client.js';

const TOP_FIVE_LIMIT = 5;

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

export function addHighScore(newScore, initials = '---') {
    const scores = loadHighScores();
    scores.push({ score: newScore, initials: initials.toUpperCase().slice(0, 3), date: new Date().toISOString() });
    scores.sort((a, b) => b.score - a.score);
    const topScores = scores.slice(0, TOP_FIVE_LIMIT);
    saveHighScores(topScores);
    return topScores;
}

export function isHighScore(score) {
    const scores = loadHighScores();
    // If less than 5 scores, always qualifies
    if (scores.length < TOP_FIVE_LIMIT) return true;
    // Check if score is higher than the lowest top 5 score
    const lowestTopScore = scores.length > 0 ? scores[scores.length - 1].score : 0;
    return score > lowestTopScore;
}

export async function updateHighScoreListUI() {
    const listEl = document.getElementById('high-score-list');
    if (!listEl) return;

    let entries = [];
    try {
        entries = await fetchTopFiveFromApi();
        saveHighScores(entries);
    } catch (error) {
        entries = loadHighScores();
    }

    listEl.innerHTML = '';

    if (!entries || entries.length === 0) {
        const li = document.createElement('li');
        li.className = 'arcade-score-item text-cyan-300 text-center';
        li.textContent = 'No scores yet';
        listEl.appendChild(li);
        return;
    }

    entries.forEach((entry, index) => {
        const li = document.createElement('li');
        li.className = 'arcade-score-item';
        const formattedScore = String(entry.score ?? 0).padStart(6, '0');
        const initials = formatInitials(entry.initials);
        li.innerHTML = `
            <span class="score-rank">${index + 1}.</span>
            <span class="score-value">${formattedScore}</span>
            <span class="score-initials">${initials}</span>
        `;
        listEl.appendChild(li);
    });
}

async function fetchTopFiveFromApi() {
    const response = await fetchScores({ page: 1, size: TOP_FIVE_LIMIT });
    const items = Array.isArray(response?.items) ? response.items : [];
    return items.slice(0, TOP_FIVE_LIMIT).map((item) => ({
        score: item.score ?? 0,
        initials: formatInitials(item.name),
        date: item.createdAt || new Date().toISOString(),
    }));
}

function formatInitials(value) {
    const safe = (typeof value === 'string' ? value : '').trim();
    if (!safe) return '---';
    return safe.toUpperCase().slice(0, 3).padEnd(3, '-');
}
