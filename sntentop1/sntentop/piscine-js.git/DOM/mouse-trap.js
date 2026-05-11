let lastCircle = null
let box = null
let trappedCircles = new Set()

export function setBox() {
  box = document.createElement('div')
  box.className = 'box'
  document.body.append(box)
}

export function createCircle() {
  document.addEventListener('click', (e) => {
    const circle = document.createElement('div')
    circle.className = 'circle'
    circle.style.background = 'white'
    circle.style.left = (e.clientX - 25) + 'px'
    circle.style.top = (e.clientY - 25) + 'px'
    document.body.append(circle)
    lastCircle = circle
  })
}

export function moveCircle() {
  document.addEventListener('mousemove', (e) => {
    if (!lastCircle) return
    
    // If circle is trapped, keep it inside the box
    if (trappedCircles.has(lastCircle)) {
      const boxRect = box.getBoundingClientRect()
      const minX = boxRect.left + 1
      const maxX = boxRect.right - 51
      const minY = boxRect.top + 1
      const maxY = boxRect.bottom - 51
      
      const x = Math.max(minX, Math.min(e.clientX - 25, maxX))
      const y = Math.max(minY, Math.min(e.clientY - 25, maxY))
      
      lastCircle.style.left = x + 'px'
      lastCircle.style.top = y + 'px'
    } else {
      // Move circle freely
      lastCircle.style.left = (e.clientX - 25) + 'px'
      lastCircle.style.top = (e.clientY - 25) + 'px'
      
      // Check if circle is entirely inside the box
      const circleRect = lastCircle.getBoundingClientRect()
      const boxRect = box.getBoundingClientRect()
      
      if (
        circleRect.left >= boxRect.left + 1 &&
        circleRect.right <= boxRect.right - 1 &&
        circleRect.top >= boxRect.top + 1 &&
        circleRect.bottom <= boxRect.bottom - 1
      ) {
        // Circle is trapped
        trappedCircles.add(lastCircle)
        lastCircle.style.background = 'var(--purple)'
      }
    }
  })
}
