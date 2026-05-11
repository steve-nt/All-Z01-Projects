function triangle(str, n) {
	let result = '';
	for (let i = 1; i <= n; i++) {
		if (i > 1) {
			result += '\n';
		}
		for (let j = 0; j < i; j++) {
			result += str;
		}
	}
	return result;
}
