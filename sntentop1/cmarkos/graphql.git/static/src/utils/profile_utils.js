export function showLoading() {
    const loading = document.getElementById('loadingState');
    const content = document.getElementById('profileContent');

    if (loading) loading.classList.remove('hidden');
    if (content) content.classList.add('hidden');
}

export function hideLoading() {
    const loading = document.getElementById('loadingState');
    const content = document.getElementById('profileContent');

    if (loading) loading.classList.add('hidden');
    if (content) content.classList.remove('hidden');
}

export function showError(message) {
    const errorMessage = document.getElementById('errorMessage');
    if (errorMessage) {
        errorMessage.textContent = message;
        errorMessage.classList.remove('hidden');
    }
    hideLoading();
}

export function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}