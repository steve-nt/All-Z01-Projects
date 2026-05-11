import { colors } from './fifty-shades-of-cold.data.js'

export function generateClasses() {
  const style = document.createElement('style')
  
  let cssText = ''
  colors.forEach(color => {
    cssText += `.${color} { background: ${color}; }\n`
  })
  
  style.textContent = cssText
  document.head.append(style)
}

export function generateColdShades() {
  const coldKeywords = ['aqua', 'blue', 'turquoise', 'green', 'cyan', 'navy', 'purple']
  
  colors.forEach(color => {
    const isColdfShade = coldKeywords.some(keyword => color.includes(keyword))
    
    if (isColdfShade) {
      const div = document.createElement('div')
      div.className = color
      div.textContent = color
      document.body.append(div)
    }
  })
}

export function choseShade(shade) {
  const divs = document.querySelectorAll('div')
  
  divs.forEach(div => {
    div.className = shade
  })
}
