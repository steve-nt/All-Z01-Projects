function map(array, callback) {
    let result = []
    for (let i = 0; i < array.length; i++) {
        result.push(callback(array[i], i, array))
    }
    return result
}

function flatMap(array, callback) {
    let result = []
    for (let i = 0; i < array.length; i++) {
        result = result.concat(callback(array[i], i, array))
    }
    return result
}