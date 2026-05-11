import { getAuthUrl } from "../config/config.js";
import { STORAGE_KEYS } from "../config/constants.js";

export class AuthService {

    async login(identifier, password) {
        try {
            const credentials = btoa(`${identifier}:${password}`);

            const response = await fetch(getAuthUrl(), {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${credentials}`
                }
            });

            if (!response.ok) {
                const errorMessage = await this.handleErrorResponse(response);
                throw new Error(errorMessage);
            }

            const data = await response.json();
            const token = data;

            if (!token) {
                throw new Error('Authentication failed: No token received');
            }

            this.setToken(token);
            
            return { success: true, token };
        } catch (error) {
            console.error('Login error:', error);
            return { success: false, error: error.message };
        }
    }

    logout() {
        localStorage.removeItem(STORAGE_KEYS.JWT_TOKEN_KEY);
        window.location.href = 'index.html';
    }

    isAuthenticated() {
        const token = this.getToken();
        if (!token) return false;

        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            if (payload.exp && payload.exp * 1000 < Date.now()) {
                localStorage.removeItem(STORAGE_KEYS.JWT_TOKEN_KEY);
                return false;
            }
        } catch {
            return false;
        }

        return true;
    }

    getAuthHeaders() {
        const token = this.getToken();
        if (!token) {
            throw new Error('User is not authenticated. Please log in to continue.');
        }
        return {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        };
    }

    async handleErrorResponse(response) {
        const statusCode = response.status;

        if (statusCode === 401) {
        return 'Invalid credentials. Please check your username/email and password.';
        }

        if (statusCode === 429) {
        return 'Too many login attempts. Please try again later.';
        }

        try {
        const data = await response.json();
        return data.message || data.error || response.statusText;
        } catch {
        return response.statusText || 'Authentication failed';
        }
    }

    setToken(token) {
        localStorage.setItem(STORAGE_KEYS.JWT_TOKEN_KEY, token);
    }

    getToken() {
        return localStorage.getItem(STORAGE_KEYS.JWT_TOKEN_KEY);
    }

    setUserData(userData) {
        localStorage.setItem(STORAGE_KEYS.USER_DATA_KEY, JSON.stringify(userData));
    }
}