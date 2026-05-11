export const  USER_INFO_QUERY = `
    query GetUserInfo {
        user {
            id
            login
            firstName
            lastName
            email
            auditRatio
            totalUp
            totalDown
        }
    }
`;

export const TRANSACTIONS_BY_TYPE_QUERY = `
    query GetTransactionsByType($userId: Int!, $type: String!) {
        transaction(
            where: {
                userId: { _eq: $userId },
                type: { _eq: $type }
            }
            order_by: { createdAt: asc }
        ) {
            id
            amount
            createdAt
            path
            object {
                id
                name
                type
            }
        }
    }
`;

export const USER_PROGRESS_QUERY = `
    query GetUserProgress($userId: Int!) {
        progress(
            where: {
                userId: { _eq: $userId }
                object: { type: { _eq: "project" } }
            }
            order_by: { updatedAt: desc }
        )
        {
            isDone
            group{
                captain{
                    login
                }
            }
            id
            updatedAt
            object {
                id
                name
                type
            }
        }
    }
`;

export const USER_XP_QUERY = `
    query GetUserXP($userId: Int!) {
        transaction_aggregate(
            where: {
                userId: { _eq: $userId },
                type: { _eq: "xp" },

                _or: [
                    {
                        object:{
                            type: {_nin: ["raid", "exercise"]}
                        }
                    },
                    {
                        _and: [
                            {
                                object: {
                                    type: { _eq: "exercise"}
                                }
                            },
                            {
                                path: { _ilike: "%checkpoint%"}
                            }
                        ]
                    }
                ]
            }
        ) 
        {
            aggregate {
                sum {
                    amount
                }
            }
        }
    }
`;