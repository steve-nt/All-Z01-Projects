function pick(obj, keys) {
  const keyArray = typeof keys === 'string' ? [keys] : keys
  const result = {}
  for (const key of Object.keys(obj)) {
    if (keyArray.includes(key)) {
      result[key] = obj[key]
    }
  }
  return result
}

function omit(obj, keys) {
  const keyArray = typeof keys === 'string' ? [keys] : keys
  const result = {}
  for (const key of Object.keys(obj)) {
    if (!keyArray.includes(key)) {
      result[key] = obj[key]
    }
  }
  return result
}

export { pick, omit }
