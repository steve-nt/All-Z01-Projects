const escapeStr = '`\\/"\''

const arr = [4, '2']

const obj = {
  str: 'string value',
  num: 123,
  bool: false,
  undef: undefined
}

const nested = Object.freeze({
  arr: Object.freeze([4, undefined, '2']),
  obj: Object.freeze({
    str: 'nested string',
    num: 456,
    bool: true
  })
})
