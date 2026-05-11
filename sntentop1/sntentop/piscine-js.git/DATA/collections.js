function arrToSet(arr) {
  return new Set(arr)
}

function arrToStr(arr) {
  return arr.join('')
}

function setToArr(set) {
  return Array.from(set)
}

function setToStr(set) {
  return Array.from(set).join('')
}

function strToArr(str) {
  return Array.from(str)
}

function strToSet(str) {
  return new Set(str)
}

function mapToObj(map) {
  return Object.fromEntries(map)
}

function objToArr(obj) {
  return Object.values(obj)
}

function objToMap(obj) {
  return new Map(Object.entries(obj))
}

function arrToObj(arr) {
  return { ...arr }
}

function strToObj(str) {
  return { ...str }
}

function superTypeOf(value) {
  if (value === null) return 'null'
  if (value === undefined) return 'undefined'
  
  const type = Object.prototype.toString.call(value).slice(8, -1)
  return type
}
