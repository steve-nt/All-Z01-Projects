function every(array, callback) {
	for (let i = 0; i < array.length; i++) {
		if (!callback(array[i], i, array)) {
			return false;
		}
	}
	return true;
}

function some(array, callback) {
	for (let i = 0; i < array.length; i++) {
		if (callback(array[i], i, array)) {
			return true;
		}
	}
	return false;

}

function none(array, callback) {
	return every(array, (element, index, arrayay) => !callback(element, index, arrayay));
}
