// GraphQL – uses storage.js and auth (load after storage and auth)
const GRAPHQL_ENDPOINT = 'https://platform.zone01.gr/api/graphql-engine/v1/graphql';

function normalizeToken(raw) {
    if (!raw) return '';
    let t = String(raw).trim();
    t = t.replace(/^Bearer\s+/i, '').trim();
    if ((t.startsWith('"') && t.endsWith('"')) || (t.startsWith("'") && t.endsWith("'"))) {
        t = t.slice(1, -1).trim();
    }
    return t;
}

const graphql = {
    async request(query, variables = {}) {
        const token = normalizeToken(storage.getToken());
        if (!token) throw new Error('Not authenticated');

        try {
            const response = await fetch(GRAPHQL_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ query, variables })
            });

            if (!response.ok) {
                if (response.status === 401) {
                    auth.logout();
                    throw new Error('Authentication failed. Please login again.');
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();
            if (result.errors) {
                throw new Error(result.errors[0]?.message || 'GraphQL error');
            }
            return result.data;
        } catch (error) {
            console.error('GraphQL request error:', error);
            throw error;
        }
    },

    async getUserInfo() {
        const query = `query { user { id login email firstName lastName } }`;
        return this.request(query);
    },

    async getXPTransactions() {
        const query = `
            query {
                transaction(where: { type: { _eq: "xp" } }, order_by: { createdAt: desc }) {
                    id amount createdAt path objectId
                    user { id login }
                }
            }
        `;
        return this.request(query);
    },

    async getTotalXP() {
        const query = `
            query {
                transaction_aggregate(
                    where: {
                        type: { _eq: "xp" }
                        event: { path: { _ilike: "%/div-01%" } }
                    }
                ) {
                    aggregate { sum { amount } count }
                }
            }
        `;
        return this.request(query);
    },

    async getProgress() {
        const query = `
            query {
                progress(order_by: { createdAt: desc }) {
                    id grade createdAt updatedAt path objectId
                    user { id login }
                }
            }
        `;
        return this.request(query);
    },

    async getResults() {
        const query = `
            query {
                result(order_by: { createdAt: desc }) {
                    id grade type createdAt updatedAt path objectId
                    user { id login }
                }
            }
        `;
        return this.request(query);
    },

    async getObjectById(objectId) {
        const query = `query GetObject($id: Int!) { object(where: { id: { _eq: $id } }) { id name type attrs } }`;
        return this.request(query, { id: objectId });
    },

    async getXPByProject() {
        const query = `
            query {
                transaction(
                    where: { 
                        type: { _eq: "xp" }
                        object: { type: { _eq: "project" } }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                    order_by: { createdAt: desc }
                ) {
                    id
                    amount
                    objectId
                    createdAt
                    object {
                        id
                        name
                        type
                    }
                }
            }
        `;
        const data = await this.request(query);
        const transactions = data?.transaction || [];
        
        // Group by project
        const projectMap = new Map();
        transactions.forEach(t => {
            const projectId = t.objectId;
            const projectName = t.object?.name || `Project ${projectId}`;
            
            if (!projectMap.has(projectId)) {
                projectMap.set(projectId, {
                    id: projectId,
                    name: projectName,
                    xp: 0,
                    latestDate: t.createdAt || null
                });
            }
            
            const project = projectMap.get(projectId);
            project.xp += t.amount || 0;
            if (t.createdAt && (!project.latestDate || t.createdAt > project.latestDate)) {
                project.latestDate = t.createdAt;
            }
        });
        
        return Array.from(projectMap.values()).sort((a, b) => b.xp - a.xp);
    },

    async getLatestProjects() {
        const query = `
            query {
                transaction(
                    where: { 
                        type: { _eq: "xp" }
                        object: { type: { _eq: "project" } }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                    order_by: { createdAt: desc }
                ) {
                    id
                    amount
                    objectId
                    createdAt
                    object {
                        id
                        name
                        type
                    }
                }
            }
        `;
        const data = await this.request(query);
        const transactions = data?.transaction || [];
        
        // Group by project
        const projectMap = new Map();
        transactions.forEach(t => {
            const projectId = t.objectId;
            const projectName = t.object?.name || `Project ${projectId}`;
            
            if (!projectMap.has(projectId)) {
                projectMap.set(projectId, {
                    id: projectId,
                    name: projectName,
                    xp: 0,
                    latestDate: t.createdAt || null
                });
            }
            
            const project = projectMap.get(projectId);
            project.xp += t.amount || 0;
            if (t.createdAt && (!project.latestDate || t.createdAt > project.latestDate)) {
                project.latestDate = t.createdAt;
            }
        });
        
        // Sort by latest date and take top 5
        return Array.from(projectMap.values())
            .sort((a, b) => (b.latestDate || '').localeCompare(a.latestDate || ''))
            .slice(0, 5);
    },

    async getSkills() {
        try {
            let query = `
                query {
                    transaction(
                        where: { 
                            eventId: { _eq: 200 }
                            _and: [{ type: { _like: "skill%" } }]
                        }
                        distinct_on: [type]
                        order_by: { type: asc, amount: desc }
                    ) {
                        id
                        type
                        amount
                    }
                }
            `;
            
            let result;
            try {
                result = await this.request(query);
            } catch (e) {
                // Fallback query without _and
                query = `
                    query {
                        transaction(
                            where: { 
                                eventId: { _eq: 200 }
                                type: { _ilike: "skill%" }
                            }
                            distinct_on: [type]
                            order_by: { type: asc, amount: desc }
                        ) {
                            id
                            type
                            amount
                        }
                    }
                `;
                result = await this.request(query);
            }
            
            const transactions = result?.transaction || [];
            const skills = transactions
                .map((t, i) => ({
                    id: t.id || i,
                    name: (t.type || '').replace(/^skill_/, '').replace(/-/g, ' ') || `Skill ${i + 1}`,
                    amount: t.amount || 0
                }))
                .sort((a, b) => b.amount - a.amount)
                .slice(0, 5);
            
            return { skills };
        } catch (error) {
            console.warn('Failed to fetch skills:', error.message);
            return { skills: [], error: 'Could not fetch skills from available data' };
        }
    },

    async getAuditRatio() {
        const query = `
            query {
                audit_down: transaction_aggregate(
                    where: { 
                        type: { _eq: "down" }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                ) {
                    aggregate {
                        sum { amount }
                    }
                }
                audit_up: transaction_aggregate(
                    where: { 
                        type: { _eq: "up" }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                ) {
                    aggregate {
                        sum { amount }
                    }
                }
            }
        `;
        
        const result = await this.request(query);
        const auditDown = result?.audit_down?.aggregate?.sum?.amount || 0;
        const auditUp = result?.audit_up?.aggregate?.sum?.amount || 0;
        const ratio = auditDown > 0 
            ? (auditUp / auditDown).toFixed(2) 
            : auditUp > 0 
                ? '∞' 
                : '0.00';
        
        return { up: auditUp, down: auditDown, ratio };
    },

    async getXPOverTime() {
        const query = `
            query {
                transaction(where: { type: { _eq: "xp" } }, order_by: { createdAt: asc }) {
                    amount createdAt
                }
            }
        `;
        return this.request(query);
    }
};
