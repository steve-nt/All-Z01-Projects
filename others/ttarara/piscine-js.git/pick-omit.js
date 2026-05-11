function pick(obj, keys) {
	const result = {};

	const keysToPick = Array.isArray(keys) ? keys : [keys];

	for (const key of keysToPick) {
		if (obj.hasOwnProperty(key)) {
			result[key] = obj[key];
		}
	}

	return result;
}

function omit(obj, keys) {
	if (obj == null || typeof obj !== 'object') {
		return {};
	}

	const result = {};
	const keysToOmit = new Set(Array.isArray(keys) ? keys : [keys]);

	const allKeys = Object.getOwnPropertyNames(obj);

	for (const key of allKeys) {
		if (!keysToOmit.has(key)) {
			result[key] = obj[key];
		}
	}

	return result;
}
