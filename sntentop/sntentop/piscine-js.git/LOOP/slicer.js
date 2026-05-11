function slice(value, start, end) {
	let result = '';
	let length = value.length;
	
	// Normalize start index
	if (start < 0) {
		start = length + start;
	}
	start = start < 0 ? 0 : start;
	
	// Normalize end index
	if (end === undefined) {
		end = length;
	} else if (end < 0) {
		end = length + end;
	}
	end = end > length ? length : end;
	
	// Build result
	if (typeof value === 'string') {
		for (let i = start; i < end; i++) {
			result += value[i];
		}
		return result;
	} else {
		result = [];
		for (let i = start; i < end; i++) {
			result.push(value[i]);
		}
		return result;
	}
}
