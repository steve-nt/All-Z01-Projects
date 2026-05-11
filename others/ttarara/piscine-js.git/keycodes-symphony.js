function compose() {
    const body = document.querySelector('body')
    addEventListener('keydown', (e) => {
        const div = document.createElement('div')
        if (e.key === 'Backspace') {
            body.removeChild(body.lastElementChild)
            return
        } else if (e.key === 'Escape') {
            document.querySelectorAll('.note').forEach(n => n.remove());
            return
        } else if (/^[a-zA-Z0-9]$/.test(e.key) || e.key === ' ') {
            div.textContent = e.key
            div.className = 'note'
            div.style.background = `#${e.keyCode.toString(16).padStart(6, '0')}`
            body.append(div)
        }
    })
}

export { compose }