function longWords(array) {
    return array.every(word => word.length >= 5 && typeof word === 'string');
}

function oneLongWord(array) {
    return array.some(word => word.length >= 10 && typeof word === 'string');
}

function noLongWords(array) {
    return array.every(word => word.length < 7 || typeof word !== 'string');
}