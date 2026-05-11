const graphql = {
    async request(query, token) {
        try {
            const response = await fetch(GRAPHQL_URL, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query }),
            })

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }

            const result = await response.json()

            if (result.errors) {
                console.error('GraphQL errors:', result.errors)
                throw new Error(result.errors[0]?.message || 'GraphQL query failed')
            }

            return {
                data: result.data,
                errors: result.errors || [],
            }
        } catch (error) {
            console.error('GraphQL request failed:', error)
            throw error
        }
    },

    async getUserInfo(token) {
        const query = `
            {
                user {
                    login
                    email
                    firstName
                    lastName
                }
            }
        `

        const result = await this.request(query, token)
        const user = result.data?.user?.[0] || null

        if (!user) {
            throw new Error('User information not found')
        }

        const fullName = user.firstName && user.lastName
            ? `${user.firstName} ${user.lastName}`.trim()
            : user.firstName || user.lastName || null

        return {
            login: user.login,
            email: user.email || null,
            name: fullName,
            firstName: user.firstName || null,
            lastName: user.lastName || null,
        }
    },

    async getXPByProject(token) {
        const query = `
            {
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
        `

        const result = await this.request(query, token)
        const transactions = result.data?.transaction || []

        const projectMap = new Map()
        transactions.forEach(transaction => {
            const projectId = transaction.objectId
            const projectName = transaction.object?.name || `Project ${projectId}`

            if (!projectMap.has(projectId)) {
                projectMap.set(projectId, {
                    id: projectId,
                    name: projectName,
                    xp: 0,
                    latestDate: transaction.createdAt || null,
                })
            }

            const project = projectMap.get(projectId)
            project.xp += transaction.amount || 0
            // Update latest date if this transaction is newer
            if (transaction.createdAt && (!project.latestDate || transaction.createdAt > project.latestDate)) {
                project.latestDate = transaction.createdAt
            }
        })

        const projects = Array.from(projectMap.values())
            .sort((a, b) => b.xp - a.xp)

        return projects
    },

    async getLatestProjects(token) {
        const query = `
            {
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
        `

        const result = await this.request(query, token)
        const transactions = result.data?.transaction || []

        // Group by project and track latest transaction
        const projectMap = new Map()
        transactions.forEach(transaction => {
            const projectId = transaction.objectId
            const projectName = transaction.object?.name || `Project ${projectId}`

            if (!projectMap.has(projectId)) {
                projectMap.set(projectId, {
                    id: projectId,
                    name: projectName,
                    xp: 0,
                    latestDate: transaction.createdAt || null,
                })
            }

            const project = projectMap.get(projectId)
            project.xp += transaction.amount || 0
            // Update latest date if this transaction is newer
            if (transaction.createdAt && (!project.latestDate || transaction.createdAt > project.latestDate)) {
                project.latestDate = transaction.createdAt
            }
        })

        // Sort by latest date (most recent first) and take top 5
        const latestProjects = Array.from(projectMap.values())
            .sort((a, b) => {
                const dateA = a.latestDate || ''
                const dateB = b.latestDate || ''
                return dateB.localeCompare(dateA) // Most recent first
            })
            .slice(0, 5)

        return latestProjects
    },

    async getTotalXP(token) {
        const query = `
            {
                transaction_aggregate(
                    where: { 
                        type: { _eq: "xp" }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                ) {
                    aggregate {
                        sum {
                            amount
                        }
                    }
                }
            }
        `

        const result = await this.request(query, token)
        const totalXP = result.data?.transaction_aggregate?.aggregate?.sum?.amount || 0

        return {
            total: totalXP,
            transactions: [],
        }
    },

    async getSkills(token) {
        try {
            let query = `
                {
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
            `

            let result
            try {
                result = await this.request(query, token)
            } catch (error1) {
                query = `
                    {
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
                `
                result = await this.request(query, token)
            }

            const transactions = result.data?.transaction || []

            const skills = transactions
                .map((transaction, index) => ({
                    id: transaction.id || index,
                    name: transaction.type?.replace('skill_', '').replace(/-/g, ' ') || `Skill ${index + 1}`,
                    amount: transaction.amount || 0,
                }))
                .sort((a, b) => b.amount - a.amount)
                .slice(0, 5)

            return {
                skills: skills,
            }
        } catch (error) {
            console.warn('Failed to fetch skills from transactions:', error.message)
            return {
                skills: [],
                error: 'Could not fetch skills from available data',
            }
        }
    },

    async getAuditRatio(token) {
        const query = `
            {
                audit_down: transaction_aggregate(
                    where: { 
                        type: { _eq: "down" }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                ) {
                    aggregate {
                        sum {
                            amount
                        }
                    }
                }
                audit_up: transaction_aggregate(
                    where: { 
                        type: { _eq: "up" }
                        event: { path: { _ilike: "%/div-01" } }
                    }
                ) {
                    aggregate {
                        sum {
                            amount
                        }
                    }
                }
            }
        `

        const result = await this.request(query, token)

        const auditDown = result.data?.audit_down?.aggregate?.sum?.amount || 0
        const auditUp = result.data?.audit_up?.aggregate?.sum?.amount || 0

        const ratio = auditDown > 0 ? (auditUp / auditDown).toFixed(2) : auditUp > 0 ? '∞' : '0.00'

        return {
            up: auditUp,
            down: auditDown,
            ratio: ratio,
        }
    },



    async getAllUserData(token) {
        try {
            const [userInfo, projectsData, latestProjectsData, xpData, auditData] = await Promise.all([
                this.getUserInfo(token),
                this.getXPByProject(token).catch(err => {
                    console.warn('Could not fetch projects:', err.message)
                    return []
                }),
                this.getLatestProjects(token).catch(err => {
                    console.warn('Could not fetch latest projects:', err.message)
                    return []
                }),
                this.getTotalXP(token),
                this.getAuditRatio(token).catch(err => {
                    console.warn('Could not fetch audit ratio:', err.message)
                    return { up: 0, down: 0, ratio: '0.00' }
                }),
            ])

            let skillsData = { skills: [], error: 'Not available' }
            try {
                skillsData = await this.getSkills(token)
            } catch (error) {
                console.warn('Could not fetch skills:', error.message)
            }

            return {
                user: userInfo,
                xp: xpData,
                skills: skillsData,
                audit: auditData,
                projects: projectsData,
                latestProjects: latestProjectsData,
            }
        } catch (error) {
            console.error('Failed to fetch all user data:', error)
            throw error
        }
    }
}
