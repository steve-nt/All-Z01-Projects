const BASE_URL = 'http://localhost:8090';
const DEFAULT_TIMEOUT = 4000;

async function request(path, { timeoutMs = DEFAULT_TIMEOUT, ...fetchOptions } = {}) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), timeoutMs);

    try {
        const response = await fetch(`${BASE_URL}${path}`, {
            ...fetchOptions,
            signal: controller.signal,
        });

        if (!response.ok) {
            const errorBody = await safeReadBody(response);
            throw new Error(`Scoreboard API error (${response.status}): ${errorBody}`);
        }

        return response.json();
    } catch (error) {
        if (error.name === 'AbortError') {
            throw new Error('Scoreboard request timed out');
        }
        throw error;
    } finally {
        clearTimeout(timeout);
    }
}

async function safeReadBody(response) {
    try {
        return await response.text();
    } catch (err) {
        return 'unreadable response';
    }
}

export async function postScore({ name, score, timeSeconds }, { timeoutMs = DEFAULT_TIMEOUT } = {}) {
    return request('/scores', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, score, timeSeconds }),
        timeoutMs,
    });
}

export async function fetchScores({ page = 1, size = 5 } = {}, { timeoutMs = DEFAULT_TIMEOUT } = {}) {
    const params = new URLSearchParams({ page: String(page), size: String(size) });
    return request(`/scores?${params.toString()}`, { timeoutMs });
}

export { BASE_URL };
