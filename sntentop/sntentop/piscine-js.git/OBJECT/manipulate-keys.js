function filterKeys(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).filter(([key, value]) => callback(key))
  )
}

function mapKeys(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [callback(key), value])
  )
}

function reduceKeys(obj, callback, initialValue) {
  const keys = Object.keys(obj)
  return initialValue !== undefined
    ? keys.reduce(callback, initialValue)
    : keys.reduce(callback)
}

export { filterKeys, mapKeys, reduceKeys }
