export function filter(array, func) {
  const result = []
  for (let i = 0; i < array.length; i++) {
    if (func(array[i], i, array)) {
      result.push(array[i])
    }
  }
  return result
}

export function reject(array, func) {
  const result = []
  for (let i = 0; i < array.length; i++) {
    if (!func(array[i], i, array)) {
      result.push(array[i])
    }
  }
  return result
}

export function partition(array, func) {
  const passed = []
  const failed = []
  for (let i = 0; i < array.length; i++) {
    if (func(array[i], i, array)) {
      passed.push(array[i])
    } else {
      failed.push(array[i])
    }
  }
  return [passed, failed]
}
