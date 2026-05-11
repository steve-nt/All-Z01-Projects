function build(bricks) {
    const body = document.querySelector('body')
    let interval
    let i = 0

    interval = setInterval(() => {
        const brick = document.createElement('div')
        brick.id = `brick-${i+1}`
        body.append(brick)
        if (i % 3 === 1) {
            brick.dataset.foundation = 'true'
        }
        i++
        if (i === bricks) {
            clearInterval(interval)
        }
    }, 100)
}

function repair(...ids) {
    ids.forEach(id => {
        const brick = document.getElementById(id)
        const index = parseInt(id.split('-')[1])
        if (!brick) return
        brick.dataset.repaired = (index % 3 === 1) ? 'in progress' : 'true'
    })
}

function destroy() {
    const bricks = document.querySelectorAll('div')
    bricks.forEach(brick => {
        const id = brick.id
        const index = parseInt(id.split('-')[1])
        if (index === Math.max(...Array.from(bricks).map(brick => parseInt(brick.id.split('-')[1])))) {
            brick.remove()
        }
    })
}

export { build, repair, destroy }