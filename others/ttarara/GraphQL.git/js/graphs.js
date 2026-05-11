// SVG Graph generation utilities

// Theme helpers for dark backgrounds
const THEME = {
    fg: 'rgba(255,255,255,0.92)',
    fgMuted: 'rgba(255,255,255,0.70)',
    fgFaint: 'rgba(255,255,255,0.45)',
    grid: 'rgba(255,255,255,0.12)',
    axis: 'rgba(255,255,255,0.55)',
    indigo: '#6366f1',
    cyan: '#06b6d4',
    pink: '#ff00d6',
    green: '#10b981',
    red: '#ef4444',
    amber: '#f59e0b',
};

function formatXP(value) {
    const n = Number(value) || 0;
    if (n === 0) return '0';
    if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(2)} MB`;
    if (n >= 1000) return `${(n / 1000).toFixed(1)} kB`;
    return String(Math.round(n));
}

function formatRatio(ratio) {
    if (ratio === '∞' || ratio === 'Infinity') return '∞';
    const x = Number(ratio);
    if (Number.isNaN(x)) return String(ratio ?? '0.00');
    return x.toFixed(2);
}

// Pie chart: XP by project (expects array [{name, xp}])
function createXPByProjectPieGraph(projects) {
    if (!Array.isArray(projects) || projects.length === 0) {
        return '<p class="graph-placeholder">No project XP data</p>';
    }

    // "Color wheel" donut with labels around (like the reference image).
    // Keep top N and group the rest as "Other" to avoid clutter.
    const sorted = [...projects].sort((a, b) => (Number(b.xp) || 0) - (Number(a.xp) || 0));
    const total = sorted.reduce((sum, p) => sum + (Number(p.xp) || 0), 0);
    if (total <= 0) return '<p class="graph-placeholder">No project XP data</p>';

    const topN = 7;
    const top = sorted.slice(0, topN);
    const rest = sorted.slice(topN);
    const otherXP = rest.reduce((sum, p) => sum + (Number(p.xp) || 0), 0);
    const series = otherXP > 0 ? [...top, { name: 'Other', xp: otherXP, isOther: true }] : top;

    const escapeHtml = (s) =>
        String(s ?? '')
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');

    const size = 640; // a bit smaller so labels fit comfortably
    const cx = size / 2;
    const cy = size / 2;
    const outerR = 150;
    const innerR = 96;
    const labelR = 214;
    const elbowR = 182;
    let currentAngle = -90;

    const polar = (r, angleDeg) => {
        const rad = (angleDeg * Math.PI) / 180;
        return { x: cx + r * Math.cos(rad), y: cy + r * Math.sin(rad) };
    };

    const donutPath = (startAngle, endAngle) => {
        const a0 = polar(outerR, startAngle);
        const a1 = polar(outerR, endAngle);
        const b0 = polar(innerR, endAngle);
        const b1 = polar(innerR, startAngle);
        const large = endAngle - startAngle > 180 ? 1 : 0;
        return [
            `M ${a0.x} ${a0.y}`,
            `A ${outerR} ${outerR} 0 ${large} 1 ${a1.x} ${a1.y}`,
            `L ${b0.x} ${b0.y}`,
            `A ${innerR} ${innerR} 0 ${large} 0 ${b1.x} ${b1.y}`,
            'Z'
        ].join(' ');
    };

    const colorForIndex = (i) => {
        // Bright wheel palette on dark bg
        const hue = (200 + i * 34) % 360;
        return `hsl(${hue}, 90%, 62%)`;
    };

    let slices = '';
    let labels = '';

    series.forEach((p, i) => {
        const value = Number(p.xp) || 0;
        if (value <= 0) return;
        const angle = (value / total) * 360;
        const start = currentAngle;
        const end = currentAngle + angle;
        const mid = (start + end) / 2;

        const color = p.isOther ? 'rgba(255,255,255,0.30)' : colorForIndex(i);
        const name = (p.name || 'Project').toString();
        const pct = ((value / total) * 100).toFixed(1);

        slices += `<path d="${donutPath(start, end)}" fill="${color}" stroke="rgba(255,255,255,0.85)" stroke-width="2">
            <title>${escapeHtml(name)}: ${formatXP(value)} (${pct}%)</title>
        </path>`;

        // label position + leader line
        const pOuter = polar(outerR, mid);
        const pElbow = polar(elbowR, mid);
        const pText = polar(labelR, mid);

        const isRight = pText.x >= cx;
        const anchor = isRight ? 'start' : 'end';
        const textX = pText.x + (isRight ? 14 : -14);
        const textY = pText.y;

        const labelName = name.length > 16 ? name.slice(0, 16) + '…' : name;

        labels += `
            <path d="M ${pOuter.x} ${pOuter.y} L ${pElbow.x} ${pElbow.y} L ${textX} ${textY}"
                  stroke="rgba(255,255,255,0.35)" stroke-width="2" fill="none" />
            <circle cx="${pOuter.x}" cy="${pOuter.y}" r="3" fill="rgba(255,255,255,0.55)" />
            <text x="${textX}" y="${textY - 5}" text-anchor="${anchor}" font-size="13" font-weight="900" fill="${THEME.fg}">
                ${escapeHtml(labelName)}
            </text>
            <text x="${textX}" y="${textY + 14}" text-anchor="${anchor}" font-size="11" font-weight="900" fill="${THEME.fgMuted}">
                ${pct}% • ${formatXP(value)}
            </text>
        `;

        currentAngle = end;
    });

    return `
        <svg viewBox="0 0 ${size} ${size}" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid meet">
            <defs>
                <radialGradient id="donutGlow" cx="50%" cy="50%" r="60%">
                    <stop offset="0%" stop-color="rgba(0,0,0,0.25)" />
                    <stop offset="100%" stop-color="rgba(0,0,0,0.00)" />
                </radialGradient>
            </defs>
            <circle cx="${cx}" cy="${cy}" r="${outerR + 34}" fill="url(#donutGlow)" />
            ${slices}
            <circle cx="${cx}" cy="${cy}" r="${innerR - 6}" fill="rgba(0,0,0,0.20)" stroke="rgba(255,255,255,0.18)" stroke-width="1"></circle>
            <text x="${cx}" y="${cy - 6}" text-anchor="middle" font-size="20" font-weight="900" fill="${THEME.fg}">${formatXP(total)}</text>
            <text x="${cx}" y="${cy + 16}" text-anchor="middle" font-size="11" font-weight="900" fill="${THEME.fgMuted}">total xp</text>
            ${labels}
        </svg>
    `;
}

// Dot plot: latest projects (expects array [{name, xp}])
function createLatestProjectsDotPlot(projects, totalXP = 0) {
    if (!Array.isArray(projects) || projects.length === 0) {
        return '<p class="graph-placeholder">No latest projects</p>';
    }

    // Larger + more readable for dashboard cards
    // Balanced sizing (consistent with other cards)
    const width = 880;
    const rowH = 72;
    const padding = { top: 34, right: 38, bottom: 52, left: 360 };
    const height = padding.top + padding.bottom + rowH * projects.length;
    const plotW = width - padding.left - padding.right;
    const maxXP = Math.max(...projects.map(p => Number(p.xp) || 0), 1);

    const fmtLabel = (s) => {
        const str = String(s ?? '');
        return str.length > 22 ? str.slice(0, 22) + '…' : str;
    };

    let svg = `<svg viewBox="0 0 ${width} ${height}" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid meet">`;
    svg += `<text x="${padding.left}" y="26" font-size="20" font-weight="900" fill="${THEME.fgMuted}">project</text>`;
    svg += `<text x="${width - padding.right}" y="26" text-anchor="end" font-size="20" font-weight="900" fill="${THEME.fgMuted}">xp</text>`;

    projects.forEach((p, idx) => {
        const y = padding.top + idx * rowH + rowH / 2;
        const xp = Number(p.xp) || 0;
        const x = padding.left + (xp / maxXP) * plotW;
        const color = `hsl(${200 + (idx * 35)}, 85%, 62%)`;
        svg += `<line x1="${padding.left}" y1="${y}" x2="${width - padding.right}" y2="${y}" stroke="${THEME.grid}" stroke-width="1" />`;
        svg += `<text x="${padding.left - 18}" y="${y + 7}" text-anchor="end" font-size="18" font-weight="900" fill="${THEME.fgMuted}">${fmtLabel(p.name || 'Project')}</text>`;
        svg += `<circle cx="${x}" cy="${y}" r="10" fill="${color}" stroke="rgba(255,255,255,0.55)" stroke-width="1.8"><title>${p.name || 'Project'}: ${formatXP(xp)}</title></circle>`;
        svg += `<text x="${width - padding.right}" y="${y + 7}" text-anchor="end" font-size="18" font-weight="900" fill="${THEME.fgMuted}">${formatXP(xp)}</text>`;
    });

    if (totalXP > 0) {
        svg += `<text x="${width - padding.right}" y="${height - 16}" text-anchor="end" font-size="14" font-weight="900" fill="${THEME.fgFaint}">total: ${formatXP(totalXP)}</text>`;
    }
    svg += `</svg>`;
    return svg;
}

// Audit gauge (expects {up, down, ratio})
function createAuditGaugeGraph(audit) {
    if (!audit) return '<p class="graph-placeholder">No audit data</p>';
    const up = Number(audit.up) || 0;
    const down = Number(audit.down) || 0;
    const ratioText = formatRatio(audit.ratio);

    const size = 200;
    const cx = size / 2;
    const cy = size / 2;
    const r = 70;
    const stroke = 18;

    // normalize ratio to 0..2 for gauge fill (like the working example)
    let ratioVal = 0;
    if (ratioText === '∞') ratioVal = 1;
    else ratioVal = Math.min(Number(ratioText) || 0, 2) / 2;

    const path = (startAngle, endAngle) => {
        const polar = (angle) => {
            const rad = ((angle - 90) * Math.PI) / 180;
            return { x: cx + r * Math.cos(rad), y: cy + r * Math.sin(rad) };
        };
        const start = polar(endAngle);
        const end = polar(startAngle);
        const large = endAngle - startAngle <= 180 ? 0 : 1;
        return `M ${start.x} ${start.y} A ${r} ${r} 0 ${large} 0 ${end.x} ${end.y}`;
    };

    const bg = path(-90, 90);
    const fillAngle = -90 + ratioVal * 180;
    const fg = ratioVal > 0 ? path(-90, fillAngle) : '';

    let color = THEME.red;
    const ratioNum = ratioText === '∞' ? 999 : (Number(ratioText) || 0);
    if (ratioNum > 1) color = THEME.green;
    else if (ratioNum > 0.5) color = THEME.amber;

    const halfCirc = Math.PI * r;
    const dashFrom = halfCirc;
    const dashTo = halfCirc * (1 - ratioVal);

    return `
        <svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}" xmlns="http://www.w3.org/2000/svg">
            <path d="${bg}" fill="none" stroke="${THEME.grid}" stroke-width="${stroke}" stroke-linecap="round" />
            ${fg ? `
            <path d="${fg}" fill="none" stroke="${color}" stroke-width="${stroke}" stroke-linecap="round"
                  stroke-dasharray="${halfCirc}" stroke-dashoffset="${dashFrom}">
                <animate attributeName="stroke-dashoffset" from="${dashFrom}" to="${dashTo}" dur="1.2s" fill="freeze" />
            </path>` : ''}
            <text x="${cx}" y="${cy + 10}" text-anchor="middle" font-size="22" font-weight="900" fill="${THEME.fg}">${ratioText}</text>
            <text x="${cx}" y="${cy + 30}" text-anchor="middle" font-size="11" font-weight="800" fill="${THEME.fgMuted}">ratio</text>
            <title>Audit up: ${formatXP(up)} | down: ${formatXP(down)} | ratio: ${ratioText}</title>
        </svg>
    `;
}

// Create XP Over Time Line Chart
function createXPOverTimeGraph(data) {
    if (!data || !data.transaction || data.transaction.length === 0) {
        return '<p class="graph-placeholder">No data available for XP over time</p>';
    }

    const transactions = data.transaction;
    // Smaller footprint so it doesn't dominate the dashboard
    const width = 760;
    const height = 320;
    const padding = { top: 34, right: 34, bottom: 52, left: 74 };
    const graphWidth = width - padding.left - padding.right;
    const graphHeight = height - padding.top - padding.bottom;

    // Calculate cumulative XP
    let cumulativeXP = 0;
    const points = transactions.map((t, index) => {
        cumulativeXP += t.amount;
        return {
            x: index,
            y: cumulativeXP,
            date: new Date(t.createdAt),
            amount: t.amount
        };
    });

    // Find min/max values
    const maxXP = Math.max(...points.map(p => p.y));
    const minXP = 0;
    const xScale = graphWidth / (points.length - 1 || 1);
    const yScale = graphHeight / (maxXP - minXP || 1);

    // Create SVG
    let svg = `<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">`;
    
    // Draw grid lines
    svg += `<defs>
        <linearGradient id="lineGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" style="stop-color:${THEME.cyan};stop-opacity:0.95" />
            <stop offset="55%" style="stop-color:${THEME.indigo};stop-opacity:0.65" />
            <stop offset="100%" style="stop-color:${THEME.pink};stop-opacity:0.15" />
        </linearGradient>
    </defs>`;

    // Y-axis grid lines
    for (let i = 0; i <= 5; i++) {
        const y = padding.top + (graphHeight / 5) * i;
        const value = maxXP - (maxXP / 5) * i;
        svg += `<line x1="${padding.left}" y1="${y}" x2="${width - padding.right}" y2="${y}" stroke="${THEME.grid}" stroke-width="1" stroke-dasharray="2,2"/>`;
        svg += `<text x="${padding.left - 10}" y="${y + 5}" text-anchor="end" font-size="12" fill="${THEME.fgMuted}">${Math.round(value).toLocaleString()}</text>`;
    }

    // X-axis labels (sample some dates)
    const labelInterval = Math.max(1, Math.floor(points.length / 6));
    for (let i = 0; i < points.length; i += labelInterval) {
        const x = padding.left + i * xScale;
        const date = points[i].date;
        const dateStr = date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
        svg += `<text x="${x}" y="${height - padding.bottom + 20}" text-anchor="middle" font-size="11" fill="${THEME.fgMuted}">${dateStr}</text>`;
    }

    // Draw area under line
    let areaPath = `M ${padding.left} ${padding.top + graphHeight}`;
    points.forEach((point, index) => {
        const x = padding.left + index * xScale;
        const y = padding.top + graphHeight - (point.y - minXP) * yScale;
        areaPath += ` L ${x} ${y}`;
    });
    areaPath += ` L ${padding.left + (points.length - 1) * xScale} ${padding.top + graphHeight} Z`;
    svg += `<path d="${areaPath}" fill="url(#lineGradient)" opacity="0.55"/>`;

    // Draw line
    let linePath = '';
    points.forEach((point, index) => {
        const x = padding.left + index * xScale;
        const y = padding.top + graphHeight - (point.y - minXP) * yScale;
        if (index === 0) {
            linePath = `M ${x} ${y}`;
        } else {
            linePath += ` L ${x} ${y}`;
        }
    });
    svg += `<path d="${linePath}" fill="none" stroke="${THEME.cyan}" stroke-width="3"/>`;

    // Draw points
    points.forEach((point, index) => {
        const x = padding.left + index * xScale;
        const y = padding.top + graphHeight - (point.y - minXP) * yScale;
        svg += `<circle cx="${x}" cy="${y}" r="4" fill="${THEME.pink}" stroke="rgba(255,255,255,0.75)" stroke-width="2"/>`;
    });

    // Axis lines
    svg += `<line x1="${padding.left}" y1="${padding.top}" x2="${padding.left}" y2="${padding.top + graphHeight}" stroke="${THEME.axis}" stroke-width="2"/>`;
    svg += `<line x1="${padding.left}" y1="${padding.top + graphHeight}" x2="${width - padding.right}" y2="${padding.top + graphHeight}" stroke="${THEME.axis}" stroke-width="2"/>`;

    // Title
    svg += `<text x="${width / 2}" y="24" text-anchor="middle" font-size="16" font-weight="900" fill="${THEME.fg}">XP Progress Over Time</text>`;

    svg += `</svg>`;
    return svg;
}

// Create XP by Project Bar Chart
function createXPByProjectGraph(data) {
    if (!data || !data.transaction || data.transaction.length === 0) {
        return '<p class="graph-placeholder">No data available for XP by project</p>';
    }

    const transactions = data.transaction;
    
    // Group by project path
    const projectMap = new Map();
    transactions.forEach(t => {
        const path = t.path || 'Unknown';
        const projectName = path.split('/').pop() || 'Unknown';
        if (!projectMap.has(projectName)) {
            projectMap.set(projectName, 0);
        }
        projectMap.set(projectName, projectMap.get(projectName) + t.amount);
    });

    // Sort by XP and take top 10
    const projects = Array.from(projectMap.entries())
        .sort((a, b) => b[1] - a[1])
        .slice(0, 10);

    if (projects.length === 0) {
        return '<p class="graph-placeholder">No project data available</p>';
    }

    const width = 780;
    const height = 340;
    // extra bottom padding so rotated labels are readable
    const padding = { top: 38, right: 34, bottom: 130, left: 84 };
    const graphWidth = width - padding.left - padding.right;
    const graphHeight = height - padding.top - padding.bottom;

    const maxXP = Math.max(...projects.map(p => p[1]));
    const barWidth = graphWidth / projects.length;
    const barSpacing = barWidth * 0.1;
    const actualBarWidth = barWidth - barSpacing;

    let svg = `<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">`;
    svg += `<defs>
        <filter id="txtShadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="1" stdDeviation="1.2" flood-color="rgba(0,0,0,0.45)"/>
        </filter>
    </defs>`;

    // Draw bars
    projects.forEach((project, index) => {
        const [name, xp] = project;
        const barHeight = (xp / maxXP) * graphHeight;
        const x = padding.left + index * barWidth + barSpacing / 2;
        const y = padding.top + graphHeight - barHeight;

        // Bar with gradient
        const gradientId = `gradient${index}`;
        svg += `<defs>
            <linearGradient id="${gradientId}" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:${THEME.cyan};stop-opacity:0.95" />
                <stop offset="55%" style="stop-color:${THEME.indigo};stop-opacity:0.9" />
                <stop offset="100%" style="stop-color:${THEME.pink};stop-opacity:0.75" />
            </linearGradient>
        </defs>`;

        svg += `<rect x="${x}" y="${y}" width="${actualBarWidth}" height="${barHeight}" fill="url(#${gradientId})" rx="4"/>`;
        
        // Value label on top of bar
        svg += `<text x="${x + actualBarWidth / 2}" y="${y - 7}" text-anchor="middle"
            font-size="12" font-weight="900" fill="${THEME.fg}"
            filter="url(#txtShadow)" style="paint-order:stroke;stroke:rgba(0,0,0,0.35);stroke-width:3px;">
            ${xp.toLocaleString()}
        </text>`;
        
        // Project name (rotated)
        const textX = x + actualBarWidth / 2;
        const textY = height - padding.bottom + 28;
        const short = name.length > 18 ? name.substring(0, 18) + '…' : name;
        svg += `<text x="${textX}" y="${textY}" text-anchor="end"
            font-size="12" font-weight="800" fill="${THEME.fgMuted}"
            filter="url(#txtShadow)"
            style="paint-order:stroke;stroke:rgba(0,0,0,0.40);stroke-width:3px;"
            transform="rotate(-35 ${textX} ${textY})">
            <title>${name}</title>${short}
        </text>`;
    });

    // Y-axis
    for (let i = 0; i <= 5; i++) {
        const y = padding.top + (graphHeight / 5) * i;
        const value = maxXP - (maxXP / 5) * i;
        svg += `<line x1="${padding.left}" y1="${y}" x2="${width - padding.right}" y2="${y}" stroke="${THEME.grid}" stroke-width="1" stroke-dasharray="2,2"/>`;
        svg += `<text x="${padding.left - 10}" y="${y + 5}" text-anchor="end"
            font-size="12" font-weight="800" fill="${THEME.fgMuted}"
            filter="url(#txtShadow)">${Math.round(value).toLocaleString()}</text>`;
    }

    // Axis lines
    svg += `<line x1="${padding.left}" y1="${padding.top}" x2="${padding.left}" y2="${padding.top + graphHeight}" stroke="${THEME.axis}" stroke-width="2"/>`;
    svg += `<line x1="${padding.left}" y1="${padding.top + graphHeight}" x2="${width - padding.right}" y2="${padding.top + graphHeight}" stroke="${THEME.axis}" stroke-width="2"/>`;

    // Title
    svg += `<text x="${width / 2}" y="24" text-anchor="middle" font-size="16" font-weight="900" fill="${THEME.fg}">XP by Project (Top 10)</text>`;

    svg += `</svg>`;
    return svg;
}

// Create Pass/Fail Ratio Pie Chart
function createPassFailRatioGraph(progressData, resultData) {
    let passed = 0;
    let failed = 0;

    // Count from progress
    if (progressData && progressData.progress) {
        progressData.progress.forEach(p => {
            if (p.grade === 1) passed++;
            else if (p.grade === 0) failed++;
        });
    }

    // Count from results
    if (resultData && resultData.result) {
        resultData.result.forEach(r => {
            if (r.grade === 1) passed++;
            else if (r.grade === 0) failed++;
        });
    }

    if (passed === 0 && failed === 0) {
        return '<p class="graph-placeholder">No pass/fail data available</p>';
    }

    const total = passed + failed;
    const passedPercent = (passed / total) * 100;
    const failedPercent = (failed / total) * 100;

    const width = 600;
    const height = 400;
    const centerX = width / 2;
    const centerY = height / 2;
    const radius = 120;

    let svg = `<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">`;

    // Calculate angles
    const passedAngle = (passedPercent / 100) * 360;
    const failedAngle = (failedPercent / 100) * 360;

    // Draw passed slice
    const passedStartAngle = -90;
    const passedEndAngle = passedStartAngle + passedAngle;
    const passedPath = createArcPath(centerX, centerY, radius, passedStartAngle, passedEndAngle);
    svg += `<path d="${passedPath}" fill="${THEME.green}" stroke="rgba(255,255,255,0.85)" stroke-width="3"/>`;

    // Draw failed slice
    const failedStartAngle = passedEndAngle;
    const failedEndAngle = failedStartAngle + failedAngle;
    const failedPath = createArcPath(centerX, centerY, radius, failedStartAngle, failedEndAngle);
    svg += `<path d="${failedPath}" fill="${THEME.pink}" stroke="rgba(255,255,255,0.85)" stroke-width="3"/>`;

    // Legend and labels
    const legendX = centerX + radius + 40;
    const legendY = centerY - 30;

    // Passed legend
    svg += `<rect x="${legendX}" y="${legendY}" width="20" height="20" fill="${THEME.green}" rx="4" stroke="rgba(255,255,255,0.25)" stroke-width="1"/>`;
    svg += `<text x="${legendX + 30}" y="${legendY + 15}" font-size="14" font-weight="800" fill="${THEME.fg}">Passed: ${passed} (${passedPercent.toFixed(1)}%)</text>`;

    // Failed legend
    svg += `<rect x="${legendX}" y="${legendY + 35}" width="20" height="20" fill="${THEME.pink}" rx="4" stroke="rgba(255,255,255,0.25)" stroke-width="1"/>`;
    svg += `<text x="${legendX + 30}" y="${legendY + 50}" font-size="14" font-weight="800" fill="${THEME.fg}">Failed: ${failed} (${failedPercent.toFixed(1)}%)</text>`;

    // Center text
    svg += `<text x="${centerX}" y="${centerY - 10}" text-anchor="middle" font-size="24" font-weight="900" fill="${THEME.fg}">${total}</text>`;
    svg += `<text x="${centerX}" y="${centerY + 15}" text-anchor="middle" font-size="14" fill="${THEME.fgMuted}">Total</text>`;

    // Title
    svg += `<text x="${width / 2}" y="30" text-anchor="middle" font-size="18" font-weight="900" fill="${THEME.fg}">Pass/Fail Ratio</text>`;

    svg += `</svg>`;
    return svg;
}

// Create Audit Ratio Graph
function createAuditRatioGraph(data) {
    // Backward compatible: if we receive aggregated audit object, render gauge
    if (data && typeof data === 'object' && ('up' in data || 'down' in data || 'ratio' in data)) {
        return createAuditGaugeGraph(data);
    }
    return '<p class="graph-placeholder">No audit data available</p>';
}

// Helper function to create arc path for pie chart
function createArcPath(centerX, centerY, radius, startAngle, endAngle) {
    const start = polarToCartesian(centerX, centerY, radius, endAngle);
    const end = polarToCartesian(centerX, centerY, radius, startAngle);
    const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";

    return [
        "M", centerX, centerY,
        "L", start.x, start.y,
        "A", radius, radius, 0, largeArcFlag, 0, end.x, end.y,
        "Z"
    ].join(" ");
}

function polarToCartesian(centerX, centerY, radius, angleInDegrees) {
    const angleInRadians = (angleInDegrees - 90) * Math.PI / 180.0;
    return {
        x: centerX + (radius * Math.cos(angleInRadians)),
        y: centerY + (radius * Math.sin(angleInRadians))
    };
}
