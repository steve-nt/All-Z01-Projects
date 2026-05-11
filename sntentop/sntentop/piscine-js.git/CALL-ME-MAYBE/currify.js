export function currify(func) {
  return function curried(...args) {
    if (args.length >= func.length) {
      return func(...args)
    }
    return function(...nextArgs) {
      return curried(...args, ...nextArgs)
    }
  }
}
