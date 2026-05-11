function fusion(obj1, obj2) {
	const result = {};

	const keys = new Set([...Object.keys(obj1), ...Object.keys(obj2)]);

	for (const key of keys) {
		const has1 = key in obj1;
		const has2 = key in obj2;

		const val1 = obj1[key];
		const val2 = obj2[key];

		if (has1 && has2) {
			if (Array.isArray(val1) && Array.isArray(val2)) {
				result[key] = val1.concat(val2);
			} else if (typeof val1 === "string" && typeof val2 === "string") {
				result[key] = val1 + " " + val2;
			} else if (typeof val1 === "number" && typeof val2 === "number") {
				result[key] = val1 + val2;
			} else if (isObject(val1) && isObject(val2)) {
				result[key] = fusion(val1, val2);
			} else {
				result[key] = val2;
			}
		} else {
			result[key] = has1 ? val1 : val2;
		}
	}

	return result;
}

function isObject(value) {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
