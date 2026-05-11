export const CONFIG = {
    BASE_URL: 'https://platform.zone01.gr',

    AUTH_ENDPOINT: '/api/auth/signin',

    GRAPHQL_ENDPOINT: '/api/graphql-engine/v1/graphql',
}

export const getApiUrl = (path) => {
    return `${CONFIG.BASE_URL}${path}`;
}

export const getAuthUrl = () => {
    return getApiUrl(CONFIG.AUTH_ENDPOINT);
}

export const getGraphqlUrl = () => {
    return getApiUrl(CONFIG.GRAPHQL_ENDPOINT);
}

export default CONFIG;