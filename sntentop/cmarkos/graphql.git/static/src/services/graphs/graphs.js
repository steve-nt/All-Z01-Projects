import { format } from '../../utils/format.js';

export class Graphs {
    renderXPProgressionChart(containerId, transactions) {
        const container = document.getElementById(containerId);
        if (!container || !transactions || transactions.length === 0) {
            container.innerHTML = '<p class="no-data">No XP data available.</p>';
            return;
        }

        let cumulativeXP = 0;
        let name = '';
        const data = transactions.map(transaction => {
            cumulativeXP += transaction.amount;
            name = transaction.object.name
            
            return {
                date: new Date(transaction.createdAt),
                xp: cumulativeXP,
                name: name
            };
        });

        const width = 800;
        const height = 400;
        const padding = 60;
        const graphWidth = width - 2 * padding;
        const graphHeight = height - 2 * padding;

        const maxXP = Math.max(...data.map(d => d.xp));
        const minDate = data[0].date;
        const maxDate = data[data.length - 1].date;
        const dateRange = maxDate - minDate;

        let svg = `<svg viewBox="0 0 ${width} ${height}" xmlns="http://www.w3.org/2000/svg" class="xp-chart">`;

        svg += `<rect width="${width}" height="${height}" fill="#0c0518" rx="8"/>`;
        
        const gridLines = 5;
        for (let i = 0; i <= gridLines; i++) {
            const y = padding + (graphHeight / gridLines) * i;
            svg += `<line x1="${padding}" y1="${y}" x2="${width - padding}" y2="${y}" 
                    stroke="#1e0f2e" stroke-width="1" stroke-dasharray="5, 5"/>`;
            
            const value = Math.round(maxXP * (1 - i / gridLines));
            svg += `<text x="${padding - 10}" y="${y + 5}" fill="#5a7a6c" 
                    text-anchor="end" font-size="12">${format.formatNumber(value)}</text>`;
        };

        const dateLabels = 5;
        for (let i = 0; i <= dateLabels; i++) {
            const date = new Date(minDate.getTime() + (dateRange / dateLabels) * i);
            const x = padding + (graphWidth / dateLabels) * i;
            svg += `<text x="${x}" y="${height - padding + 30}" fill="#5a7a6c"
                    text-anchor="middle" font-size="12">${date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric'})}</text>`;
        };

        let pathData = 'M ';
        data.forEach((point) => {
            const x = padding + ((point.date - minDate) / dateRange) * graphWidth;
            const y = padding + graphHeight - (point.xp / maxXP) * graphHeight;
            pathData += `${x},${y} `;
        });

        svg += `<defs>
            <linearGradient id="xpGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#00fff7;stop-opacity:0.8" />
                <stop offset="100%" style="stop-color:#00fff7;stop-opacity:0.1" />
            </linearGradient>
        </defs>`;

        let areaPath = pathData + `L ${width - padding},${height - padding} L ${padding},${height - padding} Z`;
        svg += `<path d="${areaPath}" fill="url(#xpGradient)" stroke="none" opacity="0.3"/>`;

        svg += `<path d="${pathData}" fill="none" stroke="#00fff7" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>`;

        data.forEach((point) => {
            const x = padding + ((point.date - minDate) / dateRange) * graphWidth;
            const y = padding + graphHeight - (point.xp / maxXP) * graphHeight;
            const name = point.name;

            svg += `<circle cx="${x}" cy="${y}" r="4" fill="#ff00ff" stroke="#000000" stroke-width="2">
                        <title>${name}\n${point.date.toLocaleDateString()}: ${format.formatNumber(point.xp)} XP</title>
                    </circle>`;
        });

        svg += `<text x="${width / 2}" y="30" fill="#e0ffe8" text-anchor="middle" 
        font-size="20" font-weight="bold">XP Progression Over Time</text>`;

        svg += '</svg>';
        container.innerHTML = svg;
    }

    renderXPByProjectChart(containerId, transactions, limit) {
        const container = document.getElementById(containerId);
        if (!container || !transactions || transactions.length === 0) {
            container.innerHTML = '<p class="no-data">No XP data available.</p>';
            return;
        }

        const xpByProject = {};
        transactions.forEach(transaction => {
            const projectName = transaction.object?.name || 'Unknown Project';
            if (!xpByProject[projectName]) {
                xpByProject[projectName] = 0;
            }
            xpByProject[projectName] += transaction.amount;
        });

        const projectArray = Object.entries(xpByProject)
            .map(([name, xp]) => ({ name, xp }))
            .sort((a, b) => b.xp - a.xp);
        const topData = projectArray.slice(0, limit);

        const width = 800;
        const height = 400;
        const padding = 60;
        const barPadding = 10;
        const graphWidth = width - 2 * padding;
        const graphHeight = height - 2 * padding;

        const barWidth = (graphWidth / topData.length) - barPadding;
        const maxXP = Math.max(...topData.map(d => d.xp));

        let svg = `<svg viewBox="0 0 ${width} ${height}" xmlns="http://www.w3.org/2000/svg" class="project-chart">`;

        svg += `<rect width="${width}" height="${height}" fill="#0c0518" rx="8"/>`;

        const gridLines = 5;
        for (let i = 0; i <= gridLines; i++) {
            const y = padding + (graphHeight / gridLines) * i;
            svg += `<line x1="${padding}" y1="${y}" x2="${width - padding}" y2="${y}" 
                    stroke="#1e0f2e" stroke-width="1" stroke-dasharray="5, 5"/>`;
            
            const value = Math.round(maxXP * (1 - i / gridLines));
            svg += `<text x="${padding - 10}" y="${y + 5}" fill="#5a7a6c" 
                    text-anchor="end" font-size="12">${format.formatNumber(value)}</text>`;
        };

        const colors = ['#00fff7', '#ffee00', '#b4ff39', '#ff2eea', '#00ff88', '#ff6b4a', '#7afcff', '#ffa500', '#9d4edd', '#f72585'];
        topData.forEach((project, index) => {
            const barHeight = (project.xp / maxXP) * graphHeight;
            const x = padding + index * (barWidth + barPadding) + barPadding / 2;
            const y = padding + graphHeight - barHeight;
            const color = colors[index % colors.length];

            svg += `<defs>
                <linearGradient id="barGradient${index}" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" style="stop-color:${color};stop-opacity:1" />
                    <stop offset="100%" style="stop-color:${color};stop-opacity:0.6" />
                </linearGradient>
            </defs>`;

            svg += `<rect x="${x}" y="${y}" width="${barWidth}" height="${barHeight}" fill="url(#barGradient${index})" rx="4">
                        <title>${project.name}: ${format.formatNumber(project.xp)} XP</title>
                    </rect>`;

            svg += `<text x="${x + barWidth / 2}" y="${y - 10}" fill="#e0ffe8" text-anchor="middle" font-size="12">${format.formatNumber(project.xp)}</text>`;

            const textY = height - padding + 35;
            const textX = x + barWidth / 2;
            const projectName = project.name.length > 15 ? project.name.slice(0, 12) + '...' : project.name;
            svg += `<text x="${textX}" y="${textY}" fill="#5a7a6c" 
                    text-anchor="middle" font-size="11"
                    transform="rotate(-35, ${textX}, ${textY})">${projectName}</text>`;
        });

        svg += `<text x="${width / 2}" y="30" fill="#e0ffe8" text-anchor="middle" 
        font-size="20" font-weight="bold">Top ${limit} Projects by XP</text>`;

        svg += '</svg>';
        container.innerHTML = svg;
    }

    renderAuditRatioChart(containerId, auditRatio, totalUp, totalDown) {
        const container = document.getElementById(containerId);
        if (!container) {
            return;
        }

        const width = 400;
        const height = 400;
        const centerX = width / 2;
        const centerY = height / 2 + 40;
        const radius = 120;

        let svg = `<svg viewBox="0 0 ${width} ${height}" xmlns="http://www.w3.org/2000/svg" class="audit-chart">`;

        svg += `<rect width="${width}" height="${height}" fill="#0c0518" rx="8"/>`;

        const total = totalUp + totalDown;
        const percentageUp = totalUp > 0 ? totalUp / total : 0.5;

        const colorUp = '#b4ff39';
        const colorDown = '#ff2eea';

        svg += `<defs>
            <linearGradient id="colorUp" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:${colorUp};stop-opacity:1" />
                <stop offset="100%" style="stop-color:${colorUp};stop-opacity:0.5" />
            </linearGradient>
            <linearGradient id="colorDown" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:${colorDown};stop-opacity:1" />
                <stop offset="100%" style="stop-color:${colorDown};stop-opacity:0.5" />
            </linearGradient>
        </defs>`;

        const titleUp = `Done: ${format.formatBytes(totalUp)}`;
        const titleDown = `Received: ${format.formatBytes(totalDown)}`;

        svg += this.createPieSlice(centerX, centerY, radius, 0, percentageUp * 360, "url(#colorUp)", titleUp);
        svg += this.createPieSlice(centerX, centerY, radius, percentageUp * 360, 360, "url(#colorDown)", titleDown);

        svg += `<circle cx="${centerX}" cy="${centerY}" r="${radius * 0.6}" fill="#0c0518"/>`;
        svg += `<text x="${centerX}" y="${centerY - 10}" fill="#e0ffe8" text-anchor="middle" 
                font-size="36" font-weight="bold">${auditRatio.toFixed(1)}</text>`;
        svg += `<text x="${centerX}" y="${centerY + 20}" fill="#5a7a6c" text-anchor="middle" 
                font-size="16">Audit Ratio</text>`;

        svg += `<circle cx="80" cy="50" r="8" fill="${colorUp}"/>`;
        svg += `<text x="95" y="55" fill="#e0ffe8" font-size="14">Done: ${format.formatBytes(totalUp)}</text>`;
        
        svg += `<circle cx="80" cy="80" r="8" fill="${colorDown}"/>`;
        svg += `<text x="95" y="85" fill="#e0ffe8" font-size="14">Received: ${format.formatBytes(totalDown)}</text>`;

        svg += `<text x="${centerX}" y="20" fill="#e0ffe8" text-anchor="middle" 
                font-size="20" font-weight="bold">Audit Statistics</text>`;

        svg += '</svg>';
        container.innerHTML = svg;
    }

    createPieSlice(cx, cy, radius, startAngle, endAngle, color, title) {
        const startRad = (startAngle - 90) * Math.PI / 180;
        const endRad = (endAngle - 90) * Math.PI / 180;
        
        const x1 = cx + radius * Math.cos(startRad);
        const y1 = cy + radius * Math.sin(startRad);
        const x2 = cx + radius * Math.cos(endRad);
        const y2 = cy + radius * Math.sin(endRad);
        
        const largeArc = endAngle - startAngle > 180 ? 1 : 0;
        
        const path = `M ${cx} ${cy} L ${x1} ${y1} A ${radius} ${radius} 0 ${largeArc} 1 ${x2} ${y2} Z`;
        
        return `<path d="${path}" fill="${color}" opacity="0.9" stroke="#0c0518" stroke-width="2">
                    <title>${title}</title>
                </path>`;
    }
};