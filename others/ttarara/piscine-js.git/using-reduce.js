function adder(numbers, initialValue = 0) {
    return numbers.reduce((acc, num) => acc + num, initialValue);
}

function sumOrMul(numbers, initialValue = 0) {
  return numbers.reduce((acc, num) =>
    (num % 2 === 0 ? acc * num : acc + num), initialValue);
}


function funcExec(funcs, initialValue = 0) {
    return funcs.reduce((acc, func) => func(acc), initialValue)
}