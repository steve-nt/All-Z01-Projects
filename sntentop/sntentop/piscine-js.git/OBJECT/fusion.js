function fusion(obj1, obj2) {
  const result = { ...obj1 }
  
  for (const key in obj2) {
    if (key in result) {
      const val1 = result[key]
      const val2 = obj2[key]
      
      if (Array.isArray(val1) && Array.isArray(val2)) {
        result[key] = [...val1, ...val2]
      } else if (
        typeof val1 === 'object' && val1 !== null && !Array.isArray(val1) &&
        typeof val2 === 'object' && val2 !== null && !Array.isArray(val2)
      ) {
        result[key] = fusion(val1, val2)
      } else if (typeof val1 === 'string' && typeof val2 === 'string') {
        result[key] = val1 + ' ' + val2
      } else if (typeof val1 === 'number' && typeof val2 === 'number') {
        result[key] = val1 + val2
      } else {
        result[key] = val2
      }
    } else {
      result[key] = obj2[key]
    }
  }
  
  return result
}

export default fusion
