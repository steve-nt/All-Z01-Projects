function reverse(value) {
	if (typeof value === 'string') {
		let result = '';
		for (let i = value.length - 1; i >= 0; i--) {
			result += value[i];
		}
		return result;
	} else {
		for (let i = 0; i < value.length / 2; i++) {
			let temp = value[i];
			value[i] = value[value.length - 1 - i];
			value[value.length - 1 - i] = temp;
		}
		return value;
	}
}
