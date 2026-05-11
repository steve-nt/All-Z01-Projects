import { AuthService } from "../services/auth_service.js";
import { Graphs } from "../services/graphs/graphs.js";
import GraphQL from "../services/graphql/queries/query_functions.js";
import { showLoading, showError, hideLoading, escapeHtml } from "../utils/profile_utils.js";
import { format } from "../utils/format.js";

const authService = new AuthService();
const graphs = new Graphs();

class ProfileController {
    constructor() {
        this.userData = null;
        this.userId = null;
    }

    async init() {
        if (!authService.isAuthenticated()) {
            window.location.href = 'index.html';
            return;
        }

        try {
            showLoading();

            await this.loadUserData();

            this.renderUserInfo();
            this.renderStatistics();

            await this.loadAndRenderProjects();
            await this.loadAndRenderGraphs();

            hideLoading();
        } catch (error) {
            showError(error.message);
        }
    }

    async loadUserData() {
        try {
            this.userData = await GraphQL.getUserInfo();
            this.userId = this.userData.id;

            authService.setUserData(this.userData);
        } catch (error) {
            console.error('Error loading user data:', error);
            throw new Error('Failed to load user data. Please try again later.');
        }
    }

    renderUserInfo() {
        const name = document.getElementById('userName');
        const email = document.getElementById('userEmail');
        const avatar = document.getElementById('userAvatar');

        if (!name || !email || !avatar || !this.userData) {
            return;
        }

        if (this.userData.firstName && this.userData.lastName) {
            name.innerHTML = `${this.userData.firstName} ${this.userData.lastName}
                <span class="user-login"> (${this.userData.login})</span>`;
        } else {
            name.textContent = this.userData.login;
        }       
        
        avatar.textContent = name.textContent.charAt(0).toUpperCase();

        email.textContent = this.userData.email || this.userData.login;

    }

    async renderStatistics() {
        try {
            const totalXP = await GraphQL.getUserXP(this.userId);
            const totalXPValue = totalXP.transaction_aggregate.aggregate.sum.amount || 0;
            document.getElementById('totalXP').textContent = format.formatNumber(totalXPValue);

            const auditRatio = this.userData.auditRatio || 0;
            const auditRatioElement = document.getElementById('auditRatio');
            auditRatioElement.textContent = auditRatio.toFixed(2);

            if (auditRatio >= 1) {
                auditRatioElement.classList.add('audit-ratio-good');
            } else if (auditRatio >= 0.7) {
                auditRatioElement.classList.add('audit-ratio-warning');
            } else {
                auditRatioElement.classList.add('audit-ratio-bad');
            }

            const progressData = await GraphQL.getUserProgress(this.userId);
            const completedProjects = progressData.progress
                .filter(p => p.isDone === true)
                .length;
            document.getElementById('projectsCompleted').textContent = completedProjects;
            
        } catch (error) {
            console.error('Error loading statistics:', error);
        }
    }

    async loadAndRenderProjects() {
        try {
            const progress = await GraphQL.getUserProgress(this.userId);
            const projectsList = document.getElementById('projectsList');

            if (!progress.progress || progress.progress.length === 0) {
                projectsList.innerHTML = '<p>There are no completed projects yet.</p>';
                return;
            }

            const topProjects = progress
                .progress
                .filter(p => p.object && p.object.name)
                .slice(0, 10);

            projectsList.innerHTML = topProjects.map(project => {
                const captainLogin = project.group && project.group.captain ? project.group.captain.login : 'Unknown';
                const captainClass = captainLogin === this.userData.login ? 'captain-self' : 'captain-other';
                const isDoneClass = project.isDone ? 'done' : 'in-progress';

                return `
                    <div class="project-item">
                        <div class="project-info">
                            <div class="project-name">${escapeHtml(project.object.name)}</div>
                            <div class="text-muted" style="font-size: 0.85rem;">
                                ${new Date(project.updatedAt).toLocaleDateString()}
                            </div>
                        </div>
                        <div class="project-captain ${captainClass}">
                            Group Captain: ${escapeHtml(captainLogin)}
                        </div>
                        <div class="project-status ${isDoneClass}">
                            ${project.isDone ? 'Completed' : 'In Progress'}
                        </div>
                    </div>
                `;
            }).join('');

        } catch (error) {
            console.error('Error loading projects:', error);
        }
    }

    async loadAndRenderGraphs() {
        try {
            const xpData = await GraphQL.getTransactionsByType(this.userId, 'xp');
            const xpTransactions = xpData.transaction
                .filter(t => {
                    const type = t.object.type
                    const path = t.path

                    return (
                        type !== 'raid' &&
                        type !== 'exercise' ||
                        (type === 'exercise' && path.toLowerCase().includes('checkpoint'))
                    )
                });

            if (xpTransactions && xpTransactions.length > 0) {
                graphs.renderXPProgressionChart('xpProgressionChart', xpTransactions);
                graphs.renderXPByProjectChart('projectXPChart', xpTransactions, 10);
            }

            if (this.userData) {
                graphs.renderAuditRatioChart(
                    'auditRatioChart',
                    this.userData.auditRatio || 0,
                    this.userData.totalUp || 0,
                    this.userData.totalDown || 0
                );
            }
        }
        catch (error) {
            console.error('Error loading graphs:', error);
        }
    }

    logout() {
        if (confirm('Are you sure you want to log out?')) {
            authService.logout();
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.profileController = new ProfileController();
    window.profileController.init();
});