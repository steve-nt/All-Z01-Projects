import { styles } from './pimp-my-style.data.js'

let currentIndex = 0
let isAdding = true

export function pimp() {
  const button = document.querySelector('.button')
  
  if (isAdding) {
    // Add classes in order
    if (currentIndex < styles.length) {
      button.classList.add(styles[currentIndex])
      button.classList.remove('unpimp')
      currentIndex++
      
      // When all classes are added, switch to removing mode
      if (currentIndex === styles.length) {
        isAdding = false
        button.classList.add('unpimp')
      }
    }
  } else {
    // Remove classes in reverse order
    if (currentIndex > 0) {
      currentIndex--
      button.classList.remove(styles[currentIndex])
      
      // If we've removed all classes, switch back to adding mode
      if (currentIndex === 0) {
        isAdding = true
        button.classList.remove('unpimp')
      } else {
        button.classList.add('unpimp')
      }
    }
  }
}
