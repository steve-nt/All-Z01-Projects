function filterValues(obj, callback) {
	const result = {};
	for (const key in obj) {
		if (callback(obj[key], key, obj)) {
			result[key] = obj[key];
		}
	}

	return result;
}

function mapValues(obj, callback) {
	const result = {};
	for (const key in obj) {
		result[key] = callback(obj[key], key, obj);
	}

	return result;
}

function reduceValues(obj, callback, initial) {
	const keys = Object.keys(obj);
	let acc = initial !== undefined ? initial : obj[keys[0]];
	const startIndex = initial !== undefined ? 0 : 1;

	for (let i = startIndex; i < keys.length; i++) {
		acc = callback(acc, obj[keys[i]], keys[i], obj);
	}

	return acc;
}
