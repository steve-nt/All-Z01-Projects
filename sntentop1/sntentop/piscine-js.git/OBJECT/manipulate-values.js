function filterValues(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).filter(([key, value]) => callback(value))
  )
}

function mapValues(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [key, callback(value)])
  )
}

function reduceValues(obj, callback, initialValue) {
  const values = Object.values(obj)
  return initialValue !== undefined 
    ? values.reduce(callback, initialValue)
    : values.reduce(callback)
}

export { filterValues, mapValues, reduceValues }
