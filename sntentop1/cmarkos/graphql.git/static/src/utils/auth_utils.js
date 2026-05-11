const loginBtn = document.getElementById('loginBtn');
const errorMessage = document.getElementById('errorMessage');

export function showError(message) {
    errorMessage.textContent = message;
    errorMessage.classList.remove('hidden');
}

export function hideError() {
    errorMessage.classList.add('hidden');
}

export function shakeForm() {
    const loginForm = document.getElementById('loginForm');
    loginForm.classList.add('shake');
    setTimeout(() => {
        loginForm.classList.remove('shake');
    }, 500);
}