function generateLetters() {
    const letters = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z']
    const body = document.querySelector('body')
    const minSize = 11
    const maxSize = 130
    const count = 120
    for (let i = 0; i < count; i++) {
        const randomIndex = Math.floor(Math.random() * letters.length)
        const letter = letters[randomIndex]
        const div = document.createElement('div')
        div.textContent = letter
        body.append(div)
        const size = minSize + (maxSize - minSize) * (i / (count - 1));
        div.style.fontSize = size.toFixed(1) + 'px'
        if (i < count/3) {
            div.style.fontWeight = '300'
        } else if (i < count*2/3) {
            div.style.fontWeight = '400'
        } else {
            div.style.fontWeight = '600'
        }
    }
} 

export { generateLetters }