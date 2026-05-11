const deepCopy = (obj) => {
  // 1. Handle primitives and functions (returns them as-is)
  if (typeof obj !== 'object' || obj === null) {
    return obj
  }

  // 2. Handle specific object types like RegExp and Date
  if (obj instanceof RegExp) {
    return new RegExp(obj)
  }
  if (obj instanceof Date) {
    return new Date(obj)
  }

  // 3. Create an accumulator for standard Arrays and Objects
  const copy = Array.isArray(obj) ? [] : {}

  // 4. Recursively copy every property
  for (let key in obj) {
    copy[key] = deepCopy(obj[key])
  }

  return copy
}