// Authentication – uses storage.js (load storage first)
const AUTH_ENDPOINT = 'https://platform.zone01.gr/api/auth/signin';

function normalizeToken(raw) {
    if (!raw) return '';
    let t = String(raw).trim();
    t = t.replace(/^Bearer\s+/i, '').trim();
    if ((t.startsWith('"') && t.endsWith('"')) || (t.startsWith("'") && t.endsWith("'"))) {
        t = t.slice(1, -1).trim();
    }
    return t;
}

function isJWTLike(token) {
    return /^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$/.test(token);
}

const auth = {
    async login(username, password) {
        try {
            const credentials = btoa(`${username}:${password}`);

            const response = await fetch(AUTH_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Authorization': `Basic ${credentials}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Login failed:', response.status, errorText);
                if (response.status === 401) {
                    return { success: false, error: 'Invalid username/email or password' };
                }
                return { success: false, error: errorText || 'Login failed. Please try again.' };
            }

            // 1) Try headers first (may fail due to CORS expose-headers)
            let token =
                response.headers.get('Authorization')?.replace(/^Bearer\s+/i, '').trim() ||
                response.headers.get('X-Token')?.trim() ||
                response.headers.get('token')?.trim() ||
                null;

            // 2) Always fall back to body parsing (most reliable)
            if (!token) {
                const raw = (await response.text() || '').trim();

                // Case A: raw token as text
                const rawNormalized = normalizeToken(raw);
                if (rawNormalized.includes('.') && rawNormalized.split('.').length === 3) {
                    token = rawNormalized;
                } else {
                    // Case B: JSON string token:  "eyJ..."
                    // Case C: JSON object token: { "token": "eyJ..." }
                    try {
                        const parsed = JSON.parse(raw);
                        if (typeof parsed === 'string') {
                            token = normalizeToken(parsed);
                        } else if (parsed && typeof parsed === 'object') {
                            token =
                                (parsed.token || parsed.access_token || parsed.accessToken || parsed.jwt || parsed.jwt_token || null);
                            if (typeof token === 'string') token = normalizeToken(token);
                        }
                    } catch (_) {
                        // Not JSON; keep token null
                    }
                }
            }

            token = normalizeToken(token);

            if (!token || token.length === 0) {
                console.error('Signin succeeded but token was not found in headers/body.');
                return { success: false, error: 'No token received from server' };
            }

            if (!isJWTLike(token)) {
                console.error('Signin succeeded but token is not a valid JWT format:', token);
                return { success: false, error: 'Signin succeeded but received an invalid token. Please try again.' };
            }

            storage.saveToken(token);
            return { success: true, token: token };
        } catch (error) {
            console.error('Login error:', error);
            return { success: false, error: error.message || 'Network error. Please check your connection.' };
        }
    },

    logout() {
        storage.removeToken();
        window.location.href = 'index.html';
    },

    getToken() {
        return storage.getToken();
    },

    isAuthenticated() {
        return storage.isAuthenticated();
    },

    getUserIdFromToken() {
        const payload = storage.decodeJWT(storage.getToken());
        return payload ? (payload.id || payload.userId) : null;
    },

    requireAuth() {
        if (!storage.isAuthenticated()) {
            window.location.href = 'index.html';
        }
    },

    redirectIfAuthenticated() {
        if (storage.isAuthenticated()) {
            window.location.href = 'profile.html';
        }
    }
};
