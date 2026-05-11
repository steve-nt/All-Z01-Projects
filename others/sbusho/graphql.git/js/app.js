const app = {
    async init() {
        // Check authentication on load
        if (storage.isAuthenticated() && !storage.isCurrentTokenExpired()) {
            this.showProfile()
            await this.loadProfile()
        } else {
            this.showLogin()
        }

        // Setup event listeners
        document.getElementById('login-form').addEventListener('submit', (e) => {
            e.preventDefault()
            this.handleLogin()
        })

        document.getElementById('logout-button').addEventListener('click', () => {
            this.logout()
        })
    },

    showLogin() {
        ui.showLogin()
    },

    showProfile() {
        ui.showProfile()
    },

    async handleLogin() {
        const usernameInput = document.getElementById('username')
        const passwordInput = document.getElementById('password')
        const loginButton = document.getElementById('login-button')
        const errorDiv = document.getElementById('login-error')

        const username = usernameInput.value.trim()
        const password = passwordInput.value

        // Clear error
        errorDiv.style.display = 'none'
        errorDiv.textContent = ''

        // Validate
        if (!username) {
            errorDiv.textContent = 'Please enter your username or email'
            errorDiv.style.display = 'block'
            return
        }

        if (!password) {
            errorDiv.textContent = 'Please enter your password'
            errorDiv.style.display = 'block'
            return
        }

        // Disable form
        loginButton.disabled = true
        loginButton.textContent = 'Logging in...'
        usernameInput.disabled = true
        passwordInput.disabled = true

        try {
            const result = await auth.login(username, password)

            if (result.success) {
                this.showProfile()
                await this.loadProfile()
            } else {
                errorDiv.textContent = result.error || 'Login failed. Please try again.'
                errorDiv.style.display = 'block'
            }
        } catch (err) {
            errorDiv.textContent = 'An unexpected error occurred. Please try again.'
            errorDiv.style.display = 'block'
            console.error('Login error:', err)
        } finally {
            loginButton.disabled = false
            loginButton.textContent = 'Log In'
            usernameInput.disabled = false
            passwordInput.disabled = false
        }
    },

    async loadProfile() {
        ui.showLoading()

        try {
            const token = storage.getToken()
            if (!token) {
                throw new Error('No authentication token found')
            }

            if (storage.isCurrentTokenExpired()) {
                throw new Error('JWT_EXPIRED')
            }

            const data = await graphql.getAllUserData(token)
            ui.updateProfile(data)
            ui.showContent()
        } catch (err) {
            console.error('Failed to fetch user data:', err)

            const errorMessage = err.message || ''
            const isExpired = errorMessage.includes('JWTExpired') || 
                            errorMessage.includes('JWT_EXPIRED') || 
                            errorMessage === 'JWT_EXPIRED'

            if (isExpired) {
                ui.showError('Your session has expired. Please log in again to continue.', true)
            } else {
                ui.showError(err.message || 'Failed to load profile data', false)
            }
        }
    },

    logout() {
        auth.logout()
        this.showLogin()
        // Clear form
        document.getElementById('username').value = ''
        document.getElementById('password').value = ''
        document.getElementById('login-error').style.display = 'none'
    }
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => app.init())
} else {
    app.init()
}

