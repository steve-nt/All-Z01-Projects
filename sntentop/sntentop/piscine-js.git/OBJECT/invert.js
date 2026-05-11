function invert(obj) {
  const result = {}
//  for (const key in obj) {
  for (const key of Object.keys(obj)) { 
    result[obj[key]] = key
  }
  return result
}

//module.exports = invert
export default invert
