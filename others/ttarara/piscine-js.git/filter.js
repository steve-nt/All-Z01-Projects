function filter(array, callback) {
	const result = [];

	for (let i = 0; i < array.length; i++) {
		if (callback(array[i], i, array)) {
			result.push(array[i]);
		}
	}

	return result;
}

function reject(array, callback) {
	const result = [];

	for (let i = 0; i < array.length; i++) {
		if (!callback(array[i], i, array)) {
			result.push(array[i]);
		}
	}

	return result;
}

function partition(array, callback) {
	const group1 = filter(array, callback);
	const group2 = reject(array, callback);

	const result = [group1, group2];

	return result;
}