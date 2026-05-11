import {styles} from './pimp-my-style.data.js'

let i = 0

const pimp = () => {
    const button = document.querySelector('button')
    
    if (!button.classList.contains('unpimp')) {
        const styleAdd = styles[i]
        button.classList.add(styleAdd)
        i++
        
        if (i === styles.length) {
            button.classList.toggle('unpimp')
        }  
    } else {
        const styleRemove = styles[i-1]
        button.classList.remove(styleRemove)
        i--

        if (i === 0) {
            button.classList.toggle('unpimp')
        }
        return
    }
}

export {pimp}