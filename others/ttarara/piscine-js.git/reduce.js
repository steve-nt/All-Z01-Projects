function fold(array, callback, accum) {
	for (let i = 0; i < array.length; i++) {
		accum = callback(accum, array[i], i, array);
	}

	return accum;
}

function foldRight(array, callback, accum) {
	for (let i = array.length - 1; i >= 0; i--) {
		accum = callback(accum, array[i], i, array);
	}

	return accum;
}

function reduce(array, reducer) {
	if (array.length < 1) {
		throw new TypeError("Reduce of empty arrayay with no initial value");
	}

	let accum = array[0];

	for (let i = 1; i < array.length; i++) {
		accum = reducer(accum, array[i], i, array);
	}

	return accum;
}

function reduceRight(array, reducer) {
	if (array.length < 1) {
		throw new TypeError("Reduce of empty arrayay with no initial value");
	}

	let accum = array[array.length - 1];

	for (let i = array.length - 2; i >= 0; i--) {
		accum = reducer(accum, array[i], i, array);
	}

	return accum;
}