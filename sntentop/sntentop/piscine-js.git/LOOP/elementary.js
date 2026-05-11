function multiply(a, b) {
	let result = 0;
	if (b < 0) {
		for (let i = 0; i > b; i--) {
			result -= a;
		}
	} else {
		for (let i = 0; i < b; i++) {
			result += a;
		}
	}
	return result;
}

function divide(a, b) {
	let result = 0;
	let isNegative = (a < 0) !== (b < 0);
	a = a < 0 ? -a : a;
	b = b < 0 ? -b : b;
	while (a >= b) {
		a -= b;
		result++;
	}
	return isNegative ? -result : result;
}

function modulo(a, b) {
	let isNegative = a < 0;
	a = a < 0 ? -a : a;
	b = b < 0 ? -b : b;
	while (a >= b) {
		a -= b;
	}
	return isNegative ? -a : a;
}
