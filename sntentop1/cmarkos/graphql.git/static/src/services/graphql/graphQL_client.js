import { getGraphqlUrl } from "../../config/config.js";
import { AuthService } from "../auth_service.js";

export const authService = new AuthService();

export class GraphQLClient {
    constructor() {
        this.endpoint = getGraphqlUrl();
    }

    async executeQuery(query, variables = {}) {
        if (!authService.isAuthenticated()) {
            throw new Error('User is not authenticated. Please log in to continue.');
        }

        console.log('Executing GraphQL query:', { query, variables });

        try {
            const response = await fetch(this.endpoint, {
                method: 'POST',
                headers: authService.getAuthHeaders(),
                body: JSON.stringify({ query, variables })
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const result = await response.json();

            if (result.errors) {
                console.error('GraphQL errors:', result.errors);
                throw new Error(result.errors.map(e => e.message).join(', '));
            }

            console.log('GraphQL query result:', result.data);
            return result.data;
        } catch (error) {
            console.error('GraphQL query error:', error);
            throw error;
        }
    }
}