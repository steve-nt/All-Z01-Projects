export function fold(array, func, accumulator) {
  let acc = accumulator
  for (let i = 0; i < array.length; i++) {
    acc = func(acc, array[i])
  }
  return acc
}

export function foldRight(array, func, accumulator) {
  let acc = accumulator
  for (let i = array.length - 1; i >= 0; i--) {
    acc = func(acc, array[i])
  }
  return acc
}

export function reduce(array, func) {
  if (array.length < 1) {
    throw new Error('Reduce of empty array with no initial value')
  }
  let acc = array[0]
  for (let i = 1; i < array.length; i++) {
    acc = func(acc, array[i])
  }
  return acc
}

export function reduceRight(array, func) {
  if (array.length < 1) {
    throw new Error('Reduce of empty array with no initial value')
  }
  let acc = array[array.length - 1]
  for (let i = array.length - 2; i >= 0; i--) {
    acc = func(acc, array[i])
  }
  return acc
}
