const DOMAIN = 'platform.zone01.gr'
const SIGNIN_URL = `https://${DOMAIN}/api/auth/signin`
const GRAPHQL_URL = `https://${DOMAIN}/api/graphql-engine/v1/graphql`

const auth = {
    async login(usernameOrEmail, password) {
        try {
            const credentials = btoa(`${usernameOrEmail}:${password}`)
            const response = await fetch(SIGNIN_URL, {
                method: 'POST',
                headers: {
                    'Authorization': `Basic ${credentials}`,
                    'Content-Type': 'application/json',
                },
            })

            if (!response.ok) {
                let errorMessage = 'Invalid credentials. Please check your username/email and password.'

                try {
                    // Clone response so we can try multiple parsing methods
                    const clonedResponse = response.clone()

                    // Try to parse as JSON first (server typically returns JSON error)
                    try {
                        const errorData = await response.json()
                        // Extract server error message, or use custom message
                        const serverError = errorData.error || errorData.message
                        errorMessage = serverError || errorMessage
                        // Optionally customize the message here:
                        if (serverError && serverError.includes('does not exist')) {
                            errorMessage = 'Invalid credentials. Please check your username/email or password.'
                        }
                    } catch (jsonError) {
                        // If JSON parsing fails, try to get text from cloned response
                        try {
                            const errorText = await clonedResponse.text()
                            // Try parsing the text as JSON (in case content-type was wrong)
                            try {
                                const parsed = JSON.parse(errorText)
                                errorMessage = parsed.error || parsed.message || errorText || errorMessage
                            } catch {
                                // Not JSON, use text as-is
                                errorMessage = errorText || errorMessage
                            }
                        } catch (textError) {
                            // Could not read text, use default message
                            console.warn('Could not read error response:', textError)
                        }
                    }
                } catch (error) {
                    // If anything fails, use default message
                    console.warn('Error handling failed response:', error)
                }

                return {
                    success: false,
                    error: errorMessage,
                }
            }

            let token = null
            token = response.headers.get('Authorization')?.replace('Bearer ', '') ||
                response.headers.get('X-Token')

            if (!token) {
                const contentType = response.headers.get('content-type')
                if (contentType && contentType.includes('application/json')) {
                    const data = await response.json()
                    token = data.token || data.access_token || data.jwt || (typeof data === 'string' ? data : null)
                } else {
                    token = await response.text()
                    token = token.trim()
                }
            }

            if (!token || token.length === 0) {
                return {
                    success: false,
                    error: 'No token received from server',
                }
            }

            storage.saveToken(token)
            return {
                success: true,
                token: token,
            }
        } catch (error) {
            return {
                success: false,
                error: error.message || 'Network error. Please check your connection.',
            }
        }
    },

    logout() {
        storage.removeToken()
    }
}

const toggleBtn = document.getElementById('togglePassword');
const passInput = document.getElementById('password');
const icon = toggleBtn.querySelector('.eye-icon');

// Click toggle
toggleBtn.addEventListener('click', () => {
    const isHidden = passInput.type === 'password';
    passInput.type = isHidden ? 'text' : 'password';

    // Toggle between eye open and eye closed SVG
    if (isHidden) {
        // Show closed eye (password visible)
        icon.innerHTML = `
            <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
            <line x1="1" y1="1" x2="23" y2="23"/>
        `;
    } else {
        // Show open eye (password hidden)
        icon.innerHTML = `
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
            <circle cx="12" cy="12" r="3"/>
        `;
    }
    toggleBtn.setAttribute('aria-label', isHidden ? 'Hide password' : 'Show password');
});