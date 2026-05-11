export function pick() {
  // Create SVG element for crosshairs
  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
  const axisX = document.createElementNS('http://www.w3.org/2000/svg', 'line')
  axisX.id = 'axisX'
  axisX.setAttribute('x1', 0)
  axisX.setAttribute('x2', 0)
  axisX.setAttribute('y1', 0)
  axisX.setAttribute('y2', '100vh')
  
  const axisY = document.createElementNS('http://www.w3.org/2000/svg', 'line')
  axisY.id = 'axisY'
  axisY.setAttribute('x1', 0)
  axisY.setAttribute('x2', '100vw')
  axisY.setAttribute('y1', 0)
  axisY.setAttribute('y2', 0)
  
  svg.append(axisX, axisY)
  document.body.append(svg)
  
  // Create HSL display div (center)
  const hslDiv = document.createElement('div')
  hslDiv.className = 'hsl text'
  hslDiv.style.left = '50%'
  hslDiv.style.top = '50%'
  hslDiv.style.transform = 'translate(-50%, -50%)'
  document.body.append(hslDiv)
  
  // Create hue display div (top right)
  const hueDiv = document.createElement('div')
  hueDiv.className = 'hue text'
  document.body.append(hueDiv)
  
  // Create luminosity display div (bottom left)
  const luminosityDiv = document.createElement('div')
  luminosityDiv.className = 'luminosity text'
  document.body.append(luminosityDiv)
  
  // Handle mouse move
  document.addEventListener('mousemove', (e) => {
    const x = e.clientX
    const y = e.clientY
    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight
    
    // Calculate hue (0-360) based on X axis
    const hue = Math.round((x / windowWidth) * 360)
    
    // Calculate luminosity (0-100) based on Y axis
    const luminosity = Math.round((y / windowHeight) * 100)
    
    // Saturation is fixed at 50%
    const saturation = 50
    
    // Create HSL string
    const hslValue = `hsl(${hue}, ${saturation}%, ${luminosity}%)`
    
    // Update background color
    document.body.style.background = hslValue
    
    // Update display divs
    hslDiv.textContent = hslValue
    hueDiv.textContent = `hue\n${hue}`
    luminosityDiv.textContent = `luminosity\n${luminosity}`
    
    // Update crosshair positions
    axisX.setAttribute('x1', x)
    axisX.setAttribute('x2', x)
    axisY.setAttribute('y1', y)
    axisY.setAttribute('y2', y)
  })
  
  // Handle click to copy to clipboard
  document.addEventListener('click', (e) => {
    const x = e.clientX
    const y = e.clientY
    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight
    
    const hue = Math.round((x / windowWidth) * 360)
    const luminosity = Math.round((y / windowHeight) * 100)
    const saturation = 50
    
    const hslValue = `hsl(${hue}, ${saturation}%, ${luminosity}%)`
    navigator.clipboard.writeText(hslValue)
  })
}
