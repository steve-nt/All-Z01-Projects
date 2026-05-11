export function map(array, func) {
  const result = []
  for (let i = 0; i < array.length; i++) {
    result.push(func(array[i], i, array))
  }
  return result
}

export function flatMap(array, func) {
  const result = []
  for (let i = 0; i < array.length; i++) {
    const mapped = func(array[i], i, array)
    if (Array.isArray(mapped)) {
      for (let j = 0; j < mapped.length; j++) {
        result.push(mapped[j])
      }
    } else {
      result.push(mapped)
    }
  }
  return result
}
