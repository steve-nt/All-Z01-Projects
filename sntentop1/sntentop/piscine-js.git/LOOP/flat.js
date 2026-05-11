function flat(arr, depth = 1) {
	let result = [];
	for (let i = 0; i < arr.length; i++) {
		if (Array.isArray(arr[i]) && depth > 0) {
			let flattened = flat(arr[i], depth - 1);
			for (let j = 0; j < flattened.length; j++) {
				result.push(flattened[j]);
			}
		} else {
			result.push(arr[i]);
		}
	}
	return result;
}
