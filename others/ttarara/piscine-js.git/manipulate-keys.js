function filterKeys(obj, callback) {
	const result = {};
	for (const key in obj) {
		if (callback(key, obj[key], obj)) {
			result[key] = obj[key];
		}
	}
	return result;
}

function mapKeys(obj, callback) {
	const result = {};
	for (const key in obj) {
		const newKey = callback(key, obj[key], obj);
		result[newKey] = obj[key];
	}
	return result;
}

function reduceKeys(obj, callback, initial) {
	const keys = Object.keys(obj);
	let acc = initial !== undefined ? initial : keys[0];
	const startIndex = initial !== undefined ? 0 : 1;

	for (let i = startIndex; i < keys.length; i++) {
		acc = callback(acc, keys[i], obj);
	}

	return acc;
}
