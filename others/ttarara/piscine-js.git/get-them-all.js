function getArchitects() {
    const all = Array.from(document.body.children)
    const splitAt = all.findIndex(el => el.id === 'buttons')
    const peopleEls = splitAt >= 0 ? all.slice(0, splitAt) : all
    
    
    const architects = peopleEls.filter(el => el.tagName === 'A')
    const nonArchitects = peopleEls.filter(el => el.tagName !== 'A')

    return [architects, nonArchitects]
}

function getClassical() {
    const [architects] = getArchitects()
    const classical = architects.filter(el => el.classList.contains('classical'))
    const nonClassical = architects.filter(el => !el.classList.contains('classical'))

    return [classical, nonClassical]
}

function getActive() {
    const [classical] = getClassical()

    const active = classical.filter(el => el.classList.contains('active'))
    const nonActive = classical.filter(el => !el.classList.contains('active'))

    return [active, nonActive]
}

function getBonannoPisano() {
    const [active] = getActive()

    const bonannoPisano = active.filter(el => el.id === 'BonannoPisano')
    const found = bonannoPisano[0]
    const nonBonannoPisano = active.filter(el => el.id !== 'BonannoPisano')

    return [found, nonBonannoPisano]
}

export { getArchitects, getClassical, getActive, getBonannoPisano }