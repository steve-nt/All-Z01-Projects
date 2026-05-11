function deepCopy(value) {
	if (value === null || typeof value !== 'object') {
		return value;
	}

	if (value instanceof Date) {
		return new Date(value);
	}

	if (value instanceof RegExp) {
		return new RegExp(value);
	}

	if (typeof value === 'function') {
		return value;
	}

	if (Array.isArray(value)) {
		return value.map(deepCopy);
	}

	const result = {};
	for (const key in value) {
		if (value.hasOwnProperty(key)) {
			result[key] = deepCopy(value[key]);
		}
	}

	return result;
}
