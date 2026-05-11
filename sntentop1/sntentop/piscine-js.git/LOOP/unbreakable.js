function split(str, separator) {
	if (separator === '') {
		let result = [];
		for (let i = 0; i < str.length; i++) {
			result.push(str[i]);
		}
		return result;
	}
	
	let result = [];
	let current = '';
	for (let i = 0; i < str.length; i++) {
		let match = true;
		if (i + separator.length <= str.length) {
			for (let j = 0; j < separator.length; j++) {
				if (str[i + j] !== separator[j]) {
					match = false;
					break;
				}
			}
		} else {
			match = false;
		}
		
		if (match) {
			result.push(current);
			current = '';
			i += separator.length - 1;
		} else {
			current += str[i];
		}
	}
	result.push(current);
	return result;
}

function join(arr, separator) {
	let result = '';
	for (let i = 0; i < arr.length; i++) {
		if (i > 0) {
			result += separator;
		}
		result += arr[i];
	}
	return result;
}
