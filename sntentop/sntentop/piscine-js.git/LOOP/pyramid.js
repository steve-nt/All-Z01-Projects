function pyramid(str, n) {
	let result = '';
	for (let i = 1; i <= n; i++) {
		if (i > 1) {
			result += '\n';
		}
		for (let j = 0; j < (n - i) * str.length; j++) {
			result += ' ';
		}
		for (let j = 0; j < 2 * i - 1; j++) {
			result += str;
		}
	}
	return result;
}
