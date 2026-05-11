function replica(target, ...sources) {
	for (const source of sources) {
		if (typeof source !== 'object' || source === null) continue;

		for (const key in source) {
			if (!source.hasOwnProperty(key)) continue;

			const sourceVal = source[key];
			const targetVal = target[key];

			if (
				typeof sourceVal === 'object' &&
				sourceVal !== null &&
				!Array.isArray(sourceVal) &&
				typeof targetVal === 'object' &&
				targetVal !== null &&
				!Array.isArray(targetVal)
			) {
				replica(targetVal, sourceVal);
			} else {
				target[key] = deepCopy(sourceVal);
			}
		}
	}

	return target;
}

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
