import {colors} from './fifty-shades-of-cold.data.js'

function generateClasses() {
    const head = document.head
    const style = document.createElement('style')
    style.textContent = colors.map((color) => `.${color} { background: ${color}; }`).join('\n')
    head.append(style)
}

function generateColdShades() {
    const body = document.querySelector('body')
    const given = ['aqua', 'blue', 'turquoise', 'green', 'cyan', 'navy', 'purple']
    for (const color of colors) {
        if (given.some(c => color.includes(c))) {
            const div = document.createElement('div')
            div.textContent = color
            div.className = color
            body.append(div)
        }
    }
}

function choseShade(color) {
    const divs = [...document.querySelectorAll('div')]
    divs.forEach((div) => {
        div.classList.replace(div.classList[0], color)
    })
}

export {generateClasses, generateColdShades, choseShade}