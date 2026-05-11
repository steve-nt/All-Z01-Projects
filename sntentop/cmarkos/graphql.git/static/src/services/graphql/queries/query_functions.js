import { 
    USER_INFO_QUERY, 
    TRANSACTIONS_BY_TYPE_QUERY, 
    USER_PROGRESS_QUERY,
    USER_XP_QUERY
} from './queries.js';

import { GraphQLClient } from '../graphQL_client.js';

const client = new GraphQLClient();

const GraphQL = {
    async getUserInfo() {
        const query = USER_INFO_QUERY;
        const data = await client.executeQuery(query);
        return data.user[0];
    },

    async getTransactionsByType(userId, type) {
        const query = TRANSACTIONS_BY_TYPE_QUERY;
        const variables = { userId, type };
        return await client.executeQuery(query, variables);
    },

    async getUserProgress(userId) {
        const query = USER_PROGRESS_QUERY;
        const variables = { userId };
        return await client.executeQuery(query, variables);
    },

    async getUserXP(userId) {
        const query = USER_XP_QUERY;
        const variables = { userId };
        return await client.executeQuery(query, variables);
    }
};

export default GraphQL;