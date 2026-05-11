// Helper function to format XP numbers in kB
function formatXP(value) {
    if (!value || value === 0) return '0'
    const kB = value / 1000
    if (kB >= 100) {
        return `${Math.round(kB)} kB`
    } else {
        return `${kB.toFixed(1)} kB`
    }
}

// Helper function to format numbers in MB
function formatMB(value) {
    if (!value || value === 0) return '0'
    const MB = value / 1000000
    return `${MB.toFixed(2)} MB`
}

// Helper function to round audit ratio to one decimal place
function formatRatio(ratio) {
    if (!ratio || ratio === '∞' || ratio === 'Infinity') return ratio
    const numRatio = parseFloat(ratio)
    if (isNaN(numRatio)) return ratio
    return numRatio.toFixed(1)
}

const ui = {
    showLogin() {
        document.getElementById('login-page').style.display = 'flex'
        document.getElementById('profile-page').style.display = 'none'
    },

    showProfile() {
        document.getElementById('login-page').style.display = 'none'
        document.getElementById('profile-page').style.display = 'block'
    },

    showLoading() {
        document.getElementById('loading-state').style.display = 'block'
        document.getElementById('error-state').style.display = 'none'
        document.getElementById('profile-content').style.display = 'none'
    },

    showError(message, isExpired = false) {
        document.getElementById('loading-state').style.display = 'none'
        document.getElementById('error-state').style.display = 'block'
        document.getElementById('profile-content').style.display = 'none'
        document.getElementById('error-message').textContent = message
        const retryButton = document.getElementById('retry-button')
        retryButton.textContent = isExpired ? 'Go to Login' : 'Retry'
        retryButton.onclick = isExpired ? () => app.logout() : () => window.location.reload()
    },

    showContent() {
        document.getElementById('loading-state').style.display = 'none'
        document.getElementById('error-state').style.display = 'none'
        document.getElementById('profile-content').style.display = 'grid'
    },

    updateProfile(data) {
        const { user, xp, skills, audit, projects, latestProjects } = data

        // User info
        document.getElementById('user-name').textContent = user?.name || 'N/A'
        document.getElementById('user-login').textContent = user?.login || 'N/A'
        if (user?.email) {
            document.getElementById('user-email').textContent = user.email
            document.getElementById('user-email-item').style.display = 'flex'
        }

        // XP
        document.getElementById('xp-total').textContent = formatXP(xp?.total) || '0'

        // XP Projects SVG - Pie Chart (all projects)
        const xpProjectsContainer = document.getElementById('xp-projects-svg-container')
        if (xpProjectsContainer && projects && projects.length > 0) {
            xpProjectsContainer.innerHTML = ''
            try {
                // Use all projects for pie chart
                const xpChart = this.createXPByProjectPieChart(projects)
                if (xpChart) {
                    xpProjectsContainer.appendChild(xpChart)
                }
            } catch (error) {
                console.error('Error creating XP projects pie chart:', error)
            }
        }

        // Latest Projects SVG - Bar Chart (top 5 latest)
        const latestProjectsContainer = document.getElementById('xp-latest-projects-svg-container')
        if (latestProjectsContainer && latestProjects && latestProjects.length > 0) {
            latestProjectsContainer.innerHTML = ''
            try {
                const totalXP = xp?.total || 0
                const latestChart = this.createLatestProjectsChart(latestProjects, true, totalXP)
                if (latestChart) {
                    latestProjectsContainer.appendChild(latestChart)
                }
            } catch (error) {
                console.error('Error creating latest projects chart:', error)
            }
        }

        // Audit
        document.getElementById('audit-ratio').textContent = formatRatio(audit?.ratio) || '0'
        document.getElementById('audit-up').textContent = formatMB(audit?.up) || '0'
        document.getElementById('audit-down').textContent = formatMB(audit?.down) || '0'

        // Audit Ratio SVG
        const auditSvgContainer = document.getElementById('audit-ratio-svg-container')
        if (auditSvgContainer && audit) {
            auditSvgContainer.innerHTML = ''
            try {
                const gauge = this.createAuditRatioGauge(audit)
                if (gauge) {
                    auditSvgContainer.appendChild(gauge)
                }
            } catch (error) {
                console.error('Error creating audit ratio gauge:', error)
            }
        }

        // Skills
        const skillsList = document.getElementById('skills-list')
        const skillsError = document.getElementById('skills-error')
        skillsList.innerHTML = ''

        if (skills?.skills && skills.skills.length > 0) {
            skillsError.style.display = 'none'
            skills.skills.forEach(skill => {
                const skillItem = document.createElement('div')
                skillItem.className = 'skill-item'
                skillItem.innerHTML = `
                    <span class="skill-name">${skill.name}</span>
                    <span class="skill-amount">${skill.amount} %</span>
                `
                skillsList.appendChild(skillItem)
            })
        } else {
            skillsError.textContent = skills?.error || 'No skills data available'
            skillsError.style.display = 'block'
        }

    },

    renderStatistics(projectss) {
        const chartsGrid = document.getElementById('charts-grid')
        if (chartsGrid) {
            chartsGrid.innerHTML = ''
        }
    },

    // Create Latest Projects Chart (dot plot chart)
    createLatestProjectsChart(data, compact = false, totalXP = 0) {
        if (!data || data.length === 0) {
            return document.createElementNS('http://www.w3.org/2000/svg', 'svg')
        }

        const padding = { top: 40, right: 40, bottom: 60, left: 100 }
        const chartWidth = compact ? 300 : 450
        const chartHeight = Math.max(220, data.length * 40 + padding.top + padding.bottom)
        const plotWidth = chartWidth - padding.left - padding.right
        const plotHeight = chartHeight - padding.top - padding.bottom

        const maxXP = Math.max(...data.map(p => p.xp), 1)
        const dotRadius = 6
        const projectSpacing = plotHeight / data.length

        const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
        svg.setAttribute('viewBox', `0 0 ${chartWidth} ${chartHeight}`)
        svg.className = compact ? 'chart-svg compact-chart dot-plot-svg' : 'chart-svg dot-plot-svg'

        // Draw grid lines
        const gridLines = document.createElementNS('http://www.w3.org/2000/svg', 'g')
        gridLines.className = 'grid-lines'

        // Horizontal grid lines (one per project)
        for (let i = 0; i <= data.length; i++) {
            const line = document.createElementNS('http://www.w3.org/2000/svg', 'line')
            const y = padding.top + (i * projectSpacing)
            line.setAttribute('x1', padding.left)
            line.setAttribute('y1', y)
            line.setAttribute('x2', chartWidth - padding.right)
            line.setAttribute('y2', y)
            line.setAttribute('stroke', '#e0e0e0')
            line.setAttribute('stroke-width', '1')
            if (i === data.length) {
                line.setAttribute('stroke', '#999')
                line.setAttribute('stroke-width', '2')
            }
            gridLines.appendChild(line)
        }

        // Vertical grid lines (5 lines)
        const numGridLines = 5
        for (let i = 0; i <= numGridLines; i++) {
            const line = document.createElementNS('http://www.w3.org/2000/svg', 'line')
            const xpValue = (maxXP / numGridLines) * i
            const x = padding.left + (xpValue / maxXP) * plotWidth
            line.setAttribute('x1', x)
            line.setAttribute('y1', padding.top)
            line.setAttribute('x2', x)
            line.setAttribute('y2', chartHeight - padding.bottom)
            line.setAttribute('stroke', '#e0e0e0')
            line.setAttribute('stroke-width', '1')
            if (i === 0 || i === numGridLines) {
                line.setAttribute('stroke', '#999')
                line.setAttribute('stroke-width', '2')
            }
            gridLines.appendChild(line)

            // X-axis labels - only show for 0 and maxXP
            if (i === 0 || i === numGridLines) {
                const label = document.createElementNS('http://www.w3.org/2000/svg', 'text')
                label.setAttribute('x', x)
                label.setAttribute('y', chartHeight - padding.bottom + 20)
                label.setAttribute('text-anchor', 'middle')
                label.setAttribute('class', 'axis-label')
                label.textContent = formatXP(Math.round(xpValue))
                gridLines.appendChild(label)
            }
        }

        svg.appendChild(gridLines)

        // Draw dots and project labels
        data.forEach((project, index) => {
            const y = padding.top + (index * projectSpacing) + (projectSpacing / 2)
            const x = padding.left + (project.xp / maxXP) * plotWidth

            // Draw dot
            const circle = document.createElementNS('http://www.w3.org/2000/svg', 'circle')
            circle.setAttribute('cx', x)
            circle.setAttribute('cy', y)
            circle.setAttribute('r', '0')
            circle.setAttribute('fill', `hsl(${220 + (index * 360 / data.length)}, 70%, 50%)`)
            circle.className = 'dot-plot-dot'
            circle.setAttribute('data-name', project.name)
            circle.setAttribute('data-xp', project.xp)

            // Animate dot appearance
            const animate = document.createElementNS('http://www.w3.org/2000/svg', 'animate')
            animate.setAttribute('attributeName', 'r')
            animate.setAttribute('from', '0')
            animate.setAttribute('to', dotRadius)
            animate.setAttribute('dur', '0.5s')
            animate.setAttribute('fill', 'freeze')
            circle.appendChild(animate)

            // Hover effect
            circle.style.cursor = 'pointer'
            circle.style.transition = 'r 0.2s'
            circle.addEventListener('mouseenter', function () {
                this.setAttribute('r', dotRadius * 1.5)
            })
            circle.addEventListener('mouseleave', function () {
                this.setAttribute('r', dotRadius)
            })

            svg.appendChild(circle)

            // Project name label (Y-axis) - full name with text wrapping using foreignObject
            // Ensure it stays within the left padding area (red limits)
            const labelPadding = 13 // Small padding from left edge
            const maxLabelWidth = padding.left - labelPadding - 5 // Stay within left boundary
            const foreignObject = document.createElementNS('http://www.w3.org/2000/svg', 'foreignObject')
            foreignObject.setAttribute('x', labelPadding.toString())
            foreignObject.setAttribute('y', (y - 20).toString())
            foreignObject.setAttribute('width', maxLabelWidth.toString())
            foreignObject.setAttribute('height', '45')
            foreignObject.setAttribute('overflow', 'hidden')

            const div = document.createElement('div')
            div.className = 'project-label-wrapper'
            div.textContent = project.name
            foreignObject.appendChild(div)
            svg.appendChild(foreignObject)

            // Ratio text: "xp" next to the dot
            if (totalXP > 0) {
                const ratioText = document.createElementNS('http://www.w3.org/2000/svg', 'text')
                ratioText.setAttribute('x', x + dotRadius + 8)
                ratioText.setAttribute('y', y)
                ratioText.setAttribute('dominant-baseline', 'middle')
                ratioText.setAttribute('class', 'dot-ratio-text')
                ratioText.textContent = formatXP(project.xp)
                svg.appendChild(ratioText)
            }
        })

        // X-axis label
        const xAxisLabel = document.createElementNS('http://www.w3.org/2000/svg', 'text')
        xAxisLabel.setAttribute('x', chartWidth / 2)
        xAxisLabel.setAttribute('y', chartHeight - 10)
        xAxisLabel.setAttribute('text-anchor', 'middle')
        xAxisLabel.setAttribute('class', 'axis-title')
        xAxisLabel.textContent = 'total xp'
        svg.appendChild(xAxisLabel)

        // Y-axis label
        const yAxisLabel = document.createElementNS('http://www.w3.org/2000/svg', 'text')
        yAxisLabel.setAttribute('x', 15)
        yAxisLabel.setAttribute('y', chartHeight / 2)
        yAxisLabel.setAttribute('text-anchor', 'middle')
        yAxisLabel.setAttribute('transform', `rotate(-90, 15, ${chartHeight / 2})`)
        yAxisLabel.setAttribute('class', 'axis-title')
        yAxisLabel.textContent = 'project name'
        svg.appendChild(yAxisLabel)

        return svg
    },

    // Create XP by Project Pie Chart (all projects)
    createXPByProjectPieChart(data) {
        if (!data || data.length === 0) {
            return document.createElementNS('http://www.w3.org/2000/svg', 'svg')
        }

        const size = 180
        const center = size / 2
        const radius = 50
        const innerRadius = 0 // For donut effect (set to 0 for full pie)

        // Calculate total XP
        const totalXP = data.reduce((sum, project) => sum + project.xp, 0)
        if (totalXP === 0) {
            return document.createElementNS('http://www.w3.org/2000/svg', 'svg')
        }

        const container = document.createElement('div')
        container.className = 'pie-chart-container'

        const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
        svg.setAttribute('viewBox', `0 0 ${size} ${size}`)
        svg.className = 'chart-svg pie-chart-svg'

        // Create tooltip text (initially hidden) - positioned at top of pie
        const tooltipText = document.createElementNS('http://www.w3.org/2000/svg', 'text')
        tooltipText.setAttribute('class', 'pie-tooltip')
        tooltipText.setAttribute('fill', '#333') // Dark text for visibility
        tooltipText.setAttribute('font-size', '8')
        tooltipText.setAttribute('font-weight', '400')
        tooltipText.setAttribute('text-anchor', 'middle')
        tooltipText.setAttribute('dominant-baseline', 'middle')
        tooltipText.setAttribute('x', center)
        tooltipText.setAttribute('y', center - radius - 15) // Position above the pie
        tooltipText.style.opacity = '0'
        tooltipText.style.pointerEvents = 'none'
        tooltipText.style.transition = 'opacity 0.2s'
        // Add text shadow for better visibility (using filter)
        tooltipText.setAttribute('filter', 'drop-shadow(0 1px 2px rgba(255, 255, 255, 0.8))')

        let currentAngle = -90 // Start at top

        data.forEach((project, index) => {
            const percentage = (project.xp / totalXP) * 100
            const angle = (project.xp / totalXP) * 360

            // Create donut segment path
            const startAngle = currentAngle
            const endAngle = currentAngle + angle

            // Calculate points for outer and inner arcs
            const outerStart = this.polarToCartesian(center, center, radius, startAngle)
            const outerEnd = this.polarToCartesian(center, center, radius, endAngle)
            const innerStart = this.polarToCartesian(center, center, innerRadius, endAngle)
            const innerEnd = this.polarToCartesian(center, center, innerRadius, startAngle)

            const largeArcFlag = angle > 180 ? "1" : "0"

            // Create donut segment path: outer arc -> line to inner -> inner arc (reverse) -> close
            const fullPath = [
                "M", outerStart.x, outerStart.y,
                "A", radius, radius, 0, largeArcFlag, 1, outerEnd.x, outerEnd.y,
                "L", innerStart.x, innerStart.y,
                "A", innerRadius, innerRadius, 0, largeArcFlag, 0, innerEnd.x, innerEnd.y,
                "Z"
            ].join(" ")

            const pathEl = document.createElementNS('http://www.w3.org/2000/svg', 'path')
            pathEl.setAttribute('d', fullPath)
            pathEl.setAttribute('fill', `hsl(${160 + (index * 360 / data.length)}, 50%, 60%)`)
            pathEl.className = 'pie-segment'
            pathEl.setAttribute('data-name', project.name)
            pathEl.setAttribute('data-xp', project.xp)
            pathEl.setAttribute('data-percentage', percentage.toFixed(1))

            // Add hover effect with tooltip
            pathEl.style.cursor = 'pointer'
            pathEl.style.transition = 'opacity 0.2s'

            pathEl.addEventListener('mouseenter', function () {
                this.style.opacity = '0.8'

                // Show tooltip with project name and XP at top of pie
                tooltipText.textContent = `${project.name} - ${formatXP(project.xp)}`
                tooltipText.style.opacity = '1'
            })

            pathEl.addEventListener('mouseleave', function () {
                this.style.opacity = '1'
                tooltipText.style.opacity = '0'
            })

            // Animate the path
            const animate = document.createElementNS('http://www.w3.org/2000/svg', 'animate')
            animate.setAttribute('attributeName', 'd')
            animate.setAttribute('from', `M ${center} ${center} L ${center} ${center} Z`)
            animate.setAttribute('to', fullPath)
            animate.setAttribute('dur', '1s')
            animate.setAttribute('fill', 'freeze')
            pathEl.appendChild(animate)

            svg.appendChild(pathEl)

            currentAngle += angle
        })

        // Append tooltip text after all paths so it appears on top
        svg.appendChild(tooltipText)

        // Create legend
        const legend = document.createElement('div')
        legend.className = 'pie-legend'



        container.appendChild(svg)
        container.appendChild(legend)
        return container
    },

    // Helper function for arc paths
    describeArc(x, y, radius, startAngle, endAngle) {
        const start = this.polarToCartesian(x, y, radius, endAngle)
        const end = this.polarToCartesian(x, y, radius, startAngle)
        const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1"
        return [
            "M", start.x, start.y,
            "A", radius, radius, 0, largeArcFlag, 0, end.x, end.y
        ].join(" ")
    },

    polarToCartesian(centerX, centerY, radius, angleInDegrees) {
        const angleInRadians = (angleInDegrees - 90) * Math.PI / 180.0
        return {
            x: centerX + (radius * Math.cos(angleInRadians)),
            y: centerY + (radius * Math.sin(angleInRadians))
        }
    },

    // Create Audit Ratio Gauge Chart
    createAuditRatioGauge(auditData) {
        const { up, down, ratio } = auditData
        const size = 190
        const center = size / 2
        const radius = 70
        const strokeWidth = 20

        // Calculate the ratio value (handle infinity and edge cases)
        let ratioValue = 0
        let originalRatio = 0
        if (ratio === '∞' || ratio === 'Infinity') {
            ratioValue = 1 // Show as full gauge for infinity
            originalRatio = 999 // For color calculation
        } else {
            originalRatio = parseFloat(ratio) || 0
            // Normalize ratio: 0.0 = 0%, 1.0 = 50%, 2.0 = 100% of gauge
            // So we divide by 2 to get the fill percentage
            ratioValue = Math.min(originalRatio, 2.0) / 2.0
        }

        const container = document.createElement('div')
        container.className = 'audit-gauge-container'

        const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
        svg.setAttribute('viewBox', `0 0 ${size} ${size}`)
        svg.setAttribute('width', '180')
        svg.setAttribute('height', '180')
        svg.className = 'chart-svg audit-gauge-svg'

        // Background arc (full gauge) - always visible
        const bgPath = this.describeArc(center, center, radius, -90, 90)
        const bgPathEl = document.createElementNS('http://www.w3.org/2000/svg', 'path')
        bgPathEl.setAttribute('d', bgPath)
        bgPathEl.setAttribute('fill', 'none')
        bgPathEl.setAttribute('stroke', '#e0e0e0')
        bgPathEl.setAttribute('stroke-width', strokeWidth)
        bgPathEl.setAttribute('stroke-linecap', 'round')
        svg.appendChild(bgPathEl)

        // Filled arc (actual ratio)
        if (ratioValue > 0) {
            const fillAngle = -90 + (ratioValue * 180)
            const fillPath = this.describeArc(center, center, radius, -90, fillAngle)
            const fillPathEl = document.createElementNS('http://www.w3.org/2000/svg', 'path')
            fillPathEl.setAttribute('d', fillPath)
            fillPathEl.setAttribute('fill', 'none')

            // Color based on ratio: green for good (>1.0), yellow for ok (0.5-1.0), red for low (<0.5)
            let strokeColor = '#ef4444' // red
            if (originalRatio > 1.0) {
                strokeColor = '#10b981' // green
            } else if (originalRatio > 0.5) {
                strokeColor = '#f59e0b' // yellow
            }

            fillPathEl.setAttribute('stroke', strokeColor)
            fillPathEl.setAttribute('stroke-width', strokeWidth)
            fillPathEl.setAttribute('stroke-linecap', 'round')
            fillPathEl.className = 'audit-gauge-fill'

            // Calculate path length for animation (half circle = π * radius)
            const pathLength = Math.PI * radius

            // Set initial stroke-dasharray and animate
            fillPathEl.setAttribute('stroke-dasharray', pathLength)
            fillPathEl.setAttribute('stroke-dashoffset', pathLength)

            // Animate using stroke-dashoffset
            const animate = document.createElementNS('http://www.w3.org/2000/svg', 'animate')
            animate.setAttribute('attributeName', 'stroke-dashoffset')
            animate.setAttribute('from', pathLength)
            animate.setAttribute('to', pathLength * (1 - ratioValue))
            animate.setAttribute('dur', '1.5s')
            animate.setAttribute('fill', 'freeze')
            fillPathEl.appendChild(animate)

            svg.appendChild(fillPathEl)
        }


        container.appendChild(svg)
        return container
    }
}