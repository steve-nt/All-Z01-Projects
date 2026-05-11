import { fetchScores } from '../api/client.js';

const PAGE_SIZE = 5;
let teardownActiveOverlay = null;

export function showScoreboard({ recentSubmission } = {}) {
    if (typeof document === 'undefined') return;

    if (typeof teardownActiveOverlay === 'function') {
        teardownActiveOverlay();
    }

    const overlay = document.createElement('div');
    overlay.dataset.scoreboardOverlay = 'true';
    overlay.style.cssText = `
        position: fixed;
        inset: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        background: rgba(15, 23, 42, 0.75);
        backdrop-filter: blur(4px);
        z-index: 5000;
        padding: 1rem;
    `;

    const modal = document.createElement('div');
    modal.style.cssText = `
        background: rgba(2, 6, 23, 0.92);
        color: #f8fafc;
        padding: 1.5rem;
        border-radius: 1rem;
        width: min(520px, 90vw);
        font-family: 'Inter', 'Segoe UI', sans-serif;
        box-shadow: 0 20px 45px rgba(2, 6, 23, 0.7);
        border: 1px solid rgba(148, 163, 184, 0.2);
    `;

    const title = document.createElement('h2');
    title.textContent = 'Global Scoreboard';
    title.style.cssText = `
        margin: 0 0 0.5rem 0;
        font-size: 1.5rem;
        letter-spacing: 0.05em;
    `;

    const banner = document.createElement('p');
    banner.style.cssText = `
        margin: 0 0 1rem 0;
        font-size: 1rem;
        color: #cbd5f5;
    `;

    const table = document.createElement('table');
    table.style.cssText = `
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 0.75rem;
        font-size: 0.95rem;
    `;

    const thead = document.createElement('thead');
    thead.innerHTML = `
        <tr>
            <th style="text-align:left;padding:0.35rem 0;">Rank</th>
            <th style="text-align:left;padding:0.35rem 0;">Name</th>
            <th style="text-align:right;padding:0.35rem 0;">Score</th>
            <th style="text-align:right;padding:0.35rem 0;">Time</th>
        </tr>
    `;
    const tbody = document.createElement('tbody');

    const pagination = document.createElement('div');
    pagination.style.cssText = `
        display:flex;
        align-items:center;
        justify-content:space-between;
        margin-top:0.75rem;
        font-size:0.9rem;
        color:#f1f5f9;
    `;

    const prevButton = document.createElement('button');
    prevButton.textContent = '←';
    styleButton(prevButton);

    const pageLabel = document.createElement('span');
    pageLabel.style.cssText = `
        flex:1;
        text-align:center;
        font-weight:600;
        letter-spacing:0.05em;
    `;

    const nextButton = document.createElement('button');
    nextButton.textContent = '→';
    styleButton(nextButton);

    const status = document.createElement('div');
    status.style.cssText = `
        min-height: 1.25rem;
        font-size: 0.85rem;
        color: #fda4af;
        text-align: center;
        margin-top: 0.5rem;
    `;

    const closeHint = document.createElement('p');
    closeHint.textContent = 'Press Esc or tap outside to close';
    closeHint.style.cssText = `
        margin: 1rem 0 0 0;
        text-align: center;
        font-size: 0.8rem;
        color: #94a3b8;
    `;

    overlay.appendChild(modal);
    modal.appendChild(title);
    modal.appendChild(banner);
    modal.appendChild(table);
    table.appendChild(thead);
    table.appendChild(tbody);
    modal.appendChild(pagination);
    pagination.appendChild(prevButton);
    pagination.appendChild(pageLabel);
    pagination.appendChild(nextButton);
    modal.appendChild(status);
    modal.appendChild(closeHint);

    const initialPage = 1;
    const recentPage = getInitialPage(recentSubmission);
    const shouldOfferJump = Boolean(recentSubmission) && recentPage !== initialPage;
    if (shouldOfferJump) {
        const jumpRow = document.createElement('div');
        jumpRow.style.cssText = `
            display:flex;
            justify-content:flex-end;
            margin-top:0.5rem;
        `;
        const jumpButton = document.createElement('button');
        jumpButton.textContent = 'Jump to my rank';
        styleButton(jumpButton);
        jumpButton.style.fontSize = '0.8rem';
        jumpButton.style.padding = '0.35rem 0.6rem';
        jumpButton.addEventListener('click', () => {
            loadPage(recentPage);
        });
        jumpRow.appendChild(jumpButton);
        modal.appendChild(jumpRow);
    }

    document.body.appendChild(overlay);

    let currentPage = initialPage;
    let totalPages = 1;
    let isClosing = false;

    const onKeyDown = (event) => {
        if (event.key === 'Escape') {
            closeOverlay();
        }
    };

    function closeOverlay() {
        if (isClosing) return;
        isClosing = true;
        document.removeEventListener('keydown', onKeyDown);
        if (overlay.parentNode) {
            overlay.parentNode.removeChild(overlay);
        }
        if (teardownActiveOverlay === closeOverlay) {
            teardownActiveOverlay = null;
        }
    }

    teardownActiveOverlay = closeOverlay;
    document.addEventListener('keydown', onKeyDown);
    overlay.addEventListener('click', (event) => {
        if (event.target === overlay) {
            closeOverlay();
        }
    });

    prevButton.addEventListener('click', () => {
        if (currentPage > 1) {
            loadPage(currentPage - 1);
        }
    });

    nextButton.addEventListener('click', () => {
        if (currentPage < totalPages) {
            loadPage(currentPage + 1);
        }
    });

    updateBanner(banner, recentSubmission);
    loadPage(initialPage);

    async function loadPage(page) {
        status.textContent = 'Loading scores...';
        prevButton.disabled = true;
        nextButton.disabled = true;
        try {
            const data = await fetchScores({ page, size: PAGE_SIZE });
            const resolvedPage = Math.max(1, data.page || page);
            currentPage = resolvedPage;
            totalPages = Math.max(1, data.totalPages || 1);
            renderRows(tbody, data.items || [], recentSubmission);
            updatePaginationControls();
            status.textContent = data.items && data.items.length ? '' : 'No scores yet. Be the first!';
        } catch (error) {
            status.textContent = 'Unable to load scores right now.';
        } finally {
            prevButton.disabled = currentPage <= 1;
            nextButton.disabled = currentPage >= totalPages;
        }
    }

    function updatePaginationControls() {
        pageLabel.textContent = `Page ${currentPage}/${totalPages}`;
    }

    return closeOverlay;
}

function styleButton(button) {
    button.type = 'button';
    button.style.cssText = `
        background: rgba(248, 250, 252, 0.08);
        border: 1px solid rgba(148, 163, 184, 0.3);
        color: #f8fafc;
        border-radius: 0.5rem;
        padding: 0.4rem 0.75rem;
        cursor: pointer;
        transition: background 0.2s ease;
    `;
    button.addEventListener('mouseenter', () => {
        button.style.background = 'rgba(248, 250, 252, 0.15)';
    });
    button.addEventListener('mouseleave', () => {
        button.style.background = 'rgba(248, 250, 252, 0.08)';
    });
}

function updateBanner(node, recent) {
    if (!node) return;
    const name = recent?.name?.trim() || 'Anon';
    const percentile = Number.isFinite(recent?.percentile) ? recent.percentile : 100;
    const rank = Number.isFinite(recent?.rank) ? recent.rank : 1;
    const ordinalRank = formatOrdinal(rank);
    node.textContent = `Congrats ${name}, you're in the top ${percentile}%, on the ${ordinalRank} position.`;
}

function renderRows(tbody, rows, recentSubmission) {
    if (!tbody) return;
    tbody.innerHTML = '';
    rows.forEach((row) => {
        const tr = document.createElement('tr');
        tr.style.background = row.id === recentSubmission?.id ? 'rgba(147, 197, 253, 0.1)' : 'transparent';
        tr.innerHTML = `
            <td style="padding:0.35rem 0;">${row.rank}</td>
            <td style="padding:0.35rem 0;">${escapeHTML(row.name)}</td>
            <td style="padding:0.35rem 0;text-align:right;">${row.score}</td>
            <td style="padding:0.35rem 0;text-align:right;">${formatTime(row.timeSeconds)}</td>
        `;
        tbody.appendChild(tr);
    });
}

function formatTime(totalSeconds) {
    const safeSeconds = Number.isFinite(totalSeconds) ? totalSeconds : 0;
    const minutes = Math.floor(safeSeconds / 60)
        .toString()
        .padStart(2, '0');
    const seconds = Math.floor(safeSeconds % 60)
        .toString()
        .padStart(2, '0');
    return `${minutes}:${seconds}`;
}

function escapeHTML(value) {
    const div = document.createElement('div');
    div.textContent = value ?? '';
    return div.innerHTML;
}

function getInitialPage(recent) {
    if (!recent?.rank) return 1;
    return Math.max(1, Math.ceil(recent.rank / PAGE_SIZE));
}

function formatOrdinal(value) {
    const n = Math.abs(Number(value)) || 0;
    const remainder100 = n % 100;
    if (remainder100 >= 11 && remainder100 <= 13) {
        return `${n}th`;
    }
    switch (n % 10) {
        case 1:
            return `${n}st`;
        case 2:
            return `${n}nd`;
        case 3:
            return `${n}rd`;
        default:
            return `${n}th`;
    }
}
