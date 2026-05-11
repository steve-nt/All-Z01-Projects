export function mult2(a) {
  return function(b) {
    return a * b
  }
}

export function add3(a) {
  return function(b) {
    return function(c) {
      return a + b + c
    }
  }
}

export function sub4(a) {
  return function(b) {
    return function(c) {
      return function(d) {
        return a - b - c - d
      }
    }
  }
}
