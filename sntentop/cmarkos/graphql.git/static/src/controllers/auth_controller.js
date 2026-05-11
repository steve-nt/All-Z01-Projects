import { AuthService } from "../services/auth_service.js";
import { showError, hideError, shakeForm } from "../utils/auth_utils.js";

const authService = new AuthService();
const loginForm = document.getElementById('loginForm');

if (authService.isAuthenticated()) {
    window.location.href = 'profile.html';
}

loginForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const identifier = document.getElementById('identifier').value.trim();
    const password = document.getElementById('password').value;

    if (!identifier || !password) {
        showError('Please enter both username/email and password.');
        return;
    }
    
    hideError();

    try {
        const result = await authService.login(identifier, password);

        if (result.success) {
            window.location.href = 'profile.html';
        } else {
            showError(result.error || 'Login failed. Please try again.');
            shakeForm();
        }
    } catch (error) {
        showError('An unexpected error occurred. Please try again later.');
        console.error('Login error:', error);
    }
});