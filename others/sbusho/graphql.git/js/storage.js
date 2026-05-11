const TOKEN_KEY = 'graphql_jwt_token'

const storage = {
    saveToken(token) {
        localStorage.setItem(TOKEN_KEY, token)
    },

    getToken() {
        return localStorage.getItem(TOKEN_KEY)
    },

    removeToken() {
        localStorage.removeItem(TOKEN_KEY)
    },

    isAuthenticated() {
        return this.getToken() !== null
    },

    decodeJWT(token) {
        try {
            const parts = token.split('.')
            if (parts.length !== 3) return null
            const payload = JSON.parse(atob(parts[1]))
            return payload
        } catch (error) {
            console.error('Failed to decode JWT:', error)
            return null
        }
    },

    isTokenExpired(token) {
        if (!token) return true
        try {
            const payload = this.decodeJWT(token)
            if (!payload || !payload.exp) return false
            const expirationTime = payload.exp * 1000
            const currentTime = Date.now()
            return currentTime >= expirationTime
        } catch (error) {
            console.error('Failed to check token expiration:', error)
            return true
        }
    },

    isCurrentTokenExpired() {
        const token = this.getToken()
        return this.isTokenExpired(token)
    }
}
