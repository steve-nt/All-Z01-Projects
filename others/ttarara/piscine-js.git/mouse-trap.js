let lastCircle = null

function createCircle() {
    addEventListener('click', (event) => {
        const circle = document.createElement('div')
        circle.classList.add('circle')
        circle.style.background = 'white' 
        document.body.append(circle)

        const radius = 25;
        circle.style.left = `${event.clientX - radius}px`;
        circle.style.top  = `${event.clientY - radius}px`;
        circle.isTrapped = false

        lastCircle = circle
    })
}

function moveCircle() {
    addEventListener('mousemove', (event) => {
        if (!lastCircle) return

        const radius = 25
        const box = document.querySelector('.box')
        const boxRect = box.getBoundingClientRect()

        const boxLeft = boxRect.left + 1
        const boxRight = boxRect.right - 1
        const boxTop = boxRect.top + 1
        const boxBottom = boxRect.bottom - 1

        const mx = event.clientX
        const my = event.clientY
        
        let newX = mx - radius
        let newY = my - radius

        const circleRight = newX + 2 * radius
        const circleBottom = newY + 2 * radius
 
        if (!lastCircle.isTrapped) {
            if (newX >= boxLeft && circleRight <= boxRight && newY >= boxTop && circleBottom <= boxBottom) {
                lastCircle.isTrapped = true
            }
        }

        if (lastCircle.isTrapped) {
            const minCenterX = boxLeft + radius
            const maxCenterX = boxRight - radius
            const minCenterY = boxTop + radius
            const maxCenterY = boxBottom - radius

            const centerX = Math.min(Math.max(mx, minCenterX), maxCenterX)
            const centerY = Math.min(Math.max(my, minCenterY), maxCenterY)

            lastCircle.style.left = `${centerX - radius}px`
            lastCircle.style.top = `${centerY - radius}px`
            lastCircle.style.background = 'var(--purple)'
        } else {
            lastCircle.style.left = `${newX}px`
            lastCircle.style.top = `${newY}px`
            lastCircle.style.background = 'white'
        }
    })
}

function setBox() {
    const box = document.createElement('div')
    box.classList.add('box')
    document.body.append(box)

    box.style.position = 'absolute'
    box.style.left = '50%'
    box.style.top = '50%'
    box.style.transform = 'translate(-50%, -50%)'
}

export {createCircle, moveCircle, setBox}