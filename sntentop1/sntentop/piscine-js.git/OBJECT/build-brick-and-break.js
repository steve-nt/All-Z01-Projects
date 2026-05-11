// Build-Brick-and-Break: A DOM manipulation exercise demonstrating:
// - Dynamic element creation and addition to the DOM
// - Custom data attributes (dataset API)
// - SetInterval for timed operations
// - Event handling with custom DOM methods
// 
// Real-world example: Building a tower game where bricks are added at intervals,
// can be repaired (with different states for foundation vs. top bricks), and destroyed.
//
// Bonus challenges:
// - Add animations when bricks are placed
// - Implement gravity (bricks fall when support is removed)
// - Add score tracking for repairs

let brickCount = 0
let intervalId = null

export const build = (n) => {
  // Clear any existing interval
  if (intervalId) clearInterval(intervalId)
  
  // Get the container (assume it has id='container' or similar)
  const container = document.getElementById('brick-container') || document.body
  
  let count = 0
  
  intervalId = setInterval(() => {
    if (count >= n) {
      clearInterval(intervalId)
      return
    }
    
    count++
    brickCount++
    
    // Create a new brick div
    const brick = document.createElement('div')
    brick.id = `brick-${brickCount}`
    brick.className = 'brick'
    
    // Determine which column this brick belongs to (0, 1, or 2)
    const column = (brickCount - 1) % 3
    
    // If it's the middle column (column 1), set the foundation attribute
    if (column === 1) {
      brick.dataset.foundation = 'true'
    }
    
    // Add the brick to the container
    container.append(brick)
  }, 100)
}

export const repair = (...ids) => {
  ids.forEach(id => {
    const brick = document.getElementById(id)
    
    if (brick) {
      // Check if this brick is in the middle column (has foundation attribute)
      if (brick.dataset.foundation === 'true') {
        brick.dataset.repaired = 'in progress'
      } else {
        brick.dataset.repaired = 'true'
      }
    }
  })
}

export const destroy = () => {
  // Get all bricks
  const bricks = document.querySelectorAll('[id^="brick-"]')
  
  // Remove the last brick if it exists
  if (bricks.length > 0) {
    bricks[bricks.length - 1].remove()
  }
}
