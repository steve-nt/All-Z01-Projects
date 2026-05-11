export function compose() {
  document.addEventListener('keydown', (e) => {
    const key = e.key
    
    // Check if it's a lowercase letter (a-z)
    if (key >= 'a' && key <= 'z') {
      const div = document.createElement('div')
      div.className = 'note'
      div.textContent = key
      
      // Generate unique color from key code
      const hue = e.keyCode * 5 % 360
      div.style.backgroundColor = `hsl(${hue}, 100%, 50%)`
      
      document.body.append(div)
    }
    
    // Backspace: delete last note
    if (key === 'Backspace') {
      const notes = document.querySelectorAll('.note')
      if (notes.length > 0) {
        notes[notes.length - 1].remove()
      }
    }
    
    // Escape: clear all notes
    if (key === 'Escape') {
      document.querySelectorAll('.note').forEach(note => note.remove())
    }
  })
}
