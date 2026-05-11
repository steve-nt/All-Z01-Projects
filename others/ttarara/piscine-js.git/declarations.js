
const escapeStr = "`\\/\"'";


const arr = Object.freeze([4, '2']);

const nested = Object.freeze({
  arr: Object.freeze([4, undefined, '2']),
  obj: Object.freeze({
    str: "Nested String",
    num: 456,
    bool: false
  })
});

const obj = Object.freeze({
  str: "Hello",
  num: 123,
  bool: true,
  undef: undefined,
  nested
});
