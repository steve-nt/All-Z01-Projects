// Main app – uses storage, auth, graphql, graphs (load last)

document.addEventListener('DOMContentLoaded', () => {
    const page = (window.location.pathname || '').split('/').pop() || '';
    const isProfile = page === 'profile.html';

    if (isProfile) {
        initProfilePage();
    } else {
        initLoginPage();
    }
});

function initLoginPage() {
    document.body.classList.add('page-login');
    document.body.classList.remove('page-profile');
    auth.redirectIfAuthenticated();

    const form = document.getElementById('loginForm');
    const err = document.getElementById('errorMessage');
    const btn = document.getElementById('loginBtn');
    const btnText = btn?.querySelector('.btn-text');
    const btnLoader = btn?.querySelector('.btn-loader');

    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = (document.getElementById('username')?.value || '').trim();
        const password = document.getElementById('password')?.value || '';

        if (!username || !password) {
            showError('Please enter both username/email and password');
            return;
        }

        if (err) { err.style.display = 'none'; err.textContent = ''; }
        if (btn) btn.disabled = true;
        if (btnText) btnText.style.display = 'none';
        if (btnLoader) btnLoader.style.display = 'inline-block';

        const result = await auth.login(username, password);

        if (result.success) {
            window.location.href = 'profile.html';
            return;
        }

        showError(result.error || 'Login failed. Please try again.');
        if (btn) btn.disabled = false;
        if (btnText) btnText.style.display = 'inline-block';
        if (btnLoader) btnLoader.style.display = 'none';
    });
}

function showError(msg) {
    const el = document.getElementById('errorMessage');
    if (el) {
        el.textContent = msg;
        el.style.display = 'block';
    }
}

function initProfilePage() {
    document.body.classList.add('page-profile');
    document.body.classList.remove('page-login');
    auth.requireAuth();

    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            if (confirm('Are you sure you want to logout?')) auth.logout();
        });
    }

    loadProfileData();
    setupGraphControls();
}

async function loadProfileData() {
    try {
        const [
            userInfo,
            xpData,
            xpTransactions,
            progressData,
            resultData,
            audit,
            skills,
            projects,
            latestProjects
        ] = await Promise.all([
            graphql.getUserInfo(),
            graphql.getTotalXP(),
            graphql.getXPTransactions(),
            graphql.getProgress(),
            graphql.getResults(),
            graphql.getAuditRatio(),
            graphql.getSkills(),
            graphql.getXPByProject(),
            graphql.getLatestProjects()
        ]);

        window.profileData = { userInfo, xpData, xpTransactions, progressData, resultData, audit, skills, projects, latestProjects };

        displayUserInfo(userInfo, projects);
        displayXPInfo(xpData);
        displayProgressInfo(progressData, resultData);
        displayAuditInfo(audit);
        displaySkills(skills);
        renderSideCharts(projects, latestProjects, xpData);
    } catch (error) {
        console.error('Error loading profile data:', error);
        alert('Error loading profile data: ' + (error.message || 'Unknown error'));
        if (String(error.message || '').toLowerCase().includes('auth')) {
            auth.logout();
        }
    }
}

function displayUserInfo(data, projects = []) {
    const u = data?.user?.[0];
    const idBadge = document.getElementById('user-id-badge');
    const nameEl = document.getElementById('user-name');
    const loginEl = document.getElementById('user-login');
    const emailRow = document.getElementById('user-email-row');
    const emailEl = document.getElementById('user-email');

    if (!u) {
        if (idBadge) idBadge.textContent = 'ID: -';
        if (nameEl) nameEl.textContent = 'N/A';
        if (loginEl) loginEl.textContent = '-';
        if (emailRow) emailRow.style.display = 'none';
        return;
    }

    const fullName = [u.firstName, u.lastName].filter(Boolean).join(' ').trim();
    if (idBadge) idBadge.textContent = `ID: ${u.id ?? '-'}`;
    if (nameEl) nameEl.textContent = fullName || u.login || 'N/A';
    if (loginEl) loginEl.textContent = u.login ?? '-';

    if (u.email) {
        if (emailRow) emailRow.style.display = '';
        if (emailEl) emailEl.textContent = u.email;
    } else {
        if (emailRow) emailRow.style.display = 'none';
    }

    // Use an arguments-based query (required by subject): fetch object by ID and display it.
    // We'll pick the top XP project id if available.
    const featuredRow = document.getElementById('featured-object-row');
    const featuredEl = document.getElementById('featured-object');
    const candidateId = Array.isArray(projects) && projects.length ? projects[0].id : null;
    if (!candidateId) {
        if (featuredRow) featuredRow.style.display = 'none';
        return;
    }

    (async () => {
        try {
            const obj = await graphql.getObjectById(candidateId);
            const o = obj?.object?.[0];
            if (o?.name) {
                if (featuredRow) featuredRow.style.display = '';
                if (featuredEl) featuredEl.textContent = o.name;
            } else {
                if (featuredRow) featuredRow.style.display = 'none';
            }
        } catch {
            if (featuredRow) featuredRow.style.display = 'none';
        }
    })();
}

function displayXPInfo(xpData) {
    const total = xpData?.transaction_aggregate?.aggregate?.sum?.amount ?? 0;
    const count = xpData?.transaction_aggregate?.aggregate?.count ?? 0;
    const totalEl = document.getElementById('totalXP');
    const countEl = document.getElementById('xpTxCount');
    if (totalEl) totalEl.textContent = String(total).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
    if (countEl) countEl.textContent = `tx: ${String(count).replace(/\B(?=(\d{3})+(?!\d))/g, ',')}`;
}

function displayProgressInfo(progressData, resultData) {
    let total = 0, passed = 0, failed = 0;

    if (progressData?.progress?.length) {
        progressData.progress.forEach(p => {
            total++;
            if (p.grade === 1) passed++; else if (p.grade === 0) failed++;
        });
    }
    if (resultData?.result?.length) {
        resultData.result.forEach(r => {
            total++;
            if (r.grade === 1) passed++; else if (r.grade === 0) failed++;
        });
    }

    const fmt = n => String(n).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
    const totalEl = document.getElementById('totalProgress');
    const passEl = document.getElementById('passedProjects');
    const failEl = document.getElementById('failedProjects');
    if (totalEl) totalEl.textContent = fmt(total);
    if (passEl) passEl.textContent = fmt(passed);
    if (failEl) failEl.textContent = fmt(failed);
}

function displayAuditInfo(audit) {
    const ratioEl = document.getElementById('auditRatioValue');
    const upEl = document.getElementById('auditUpValue');
    const downEl = document.getElementById('auditDownValue');
    const gaugeEl = document.getElementById('auditGauge');

    const ratio = audit?.ratio ?? '0.00';
    const up = audit?.up ?? 0;
    const down = audit?.down ?? 0;

    if (ratioEl) ratioEl.textContent = String(ratio);
    if (upEl) upEl.textContent = String(up);
    if (downEl) downEl.textContent = String(down);
    if (gaugeEl) gaugeEl.innerHTML = createAuditGaugeGraph(audit);
}

function displaySkills(skillsData) {
    const list = document.getElementById('skillsList');
    const err = document.getElementById('skillsError');
    if (!list) return;
    list.innerHTML = '';

    const skills = skillsData?.skills || [];
    if (skills.length === 0) {
        if (err) {
            err.textContent = skillsData?.error || 'No skills data available';
            err.style.display = '';
        }
        return;
    }

    if (err) err.style.display = 'none';
    skills.forEach(s => {
        const row = document.createElement('div');
        row.className = 'skill-row';
        row.innerHTML = `<span class="skill-name">${String(s.name || '').replace(/_/g, ' ')}</span><span class="skill-amt">${s.amount}%</span>`;
        list.appendChild(row);
    });
}

function renderSideCharts(projects, latestProjects, xpData) {
    const total = xpData?.transaction_aggregate?.aggregate?.sum?.amount ?? 0;
    const pie = document.getElementById('xpByProjectPie');
    const dot = document.getElementById('xpLatestDotPlot');
    if (pie) pie.innerHTML = createXPByProjectPieGraph(projects);
    if (dot) dot.innerHTML = createLatestProjectsDotPlot(latestProjects, total);
}

function setupGraphControls() {
    const buttons = document.querySelectorAll('.btn-graph');
    const display = document.getElementById('graphDisplay');
    if (!display) return;

    buttons.forEach(btn => {
        btn.addEventListener('click', async () => {
            buttons.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            display.innerHTML = '<p class="graph-placeholder">Loading graph...</p>';
            await displayGraph(btn.getAttribute('data-graph'));
        });
    });
}

async function displayGraph(type) {
    const display = document.getElementById('graphDisplay');
    if (!display) return;

    try {
        let html = '';

        switch (type) {
            case 'xpOverTime':
                html = createXPOverTimeGraph(
                    window.profileData?.xpTransactions || await graphql.getXPOverTime()
                );
                break;
            case 'xpByProjectBar':
                html = createXPByProjectGraph(
                    window.profileData?.xpTransactions || await graphql.getXPTransactions()
                );
                break;
            case 'passFailRatio':
                html = createPassFailRatioGraph(
                    window.profileData?.progressData || await graphql.getProgress(),
                    window.profileData?.resultData || await graphql.getResults()
                );
                break;
            case 'auditGauge':
                html = createAuditGaugeGraph(window.profileData?.audit || await graphql.getAuditRatio());
                break;
            default:
                html = '<p class="graph-placeholder">Unknown graph type</p>';
        }

        display.innerHTML = html || '<p class="graph-placeholder">No data</p>';
    } catch (e) {
        console.error('Graph error:', e);
        display.innerHTML = '<p class="graph-placeholder">Error: ' + (e.message || 'Unknown') + '</p>';
    }
}
