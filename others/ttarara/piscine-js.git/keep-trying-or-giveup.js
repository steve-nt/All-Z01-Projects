function retry(count, callback) {
	return async function(...args) {
		let attempts = 0;
		let lastError;
		while (attempts <= count) {
			try {
				return await callback(...args);
			} catch (err) {
				lastError = err;
				attempts++;
			}
		}
		throw lastError;
	};
}

function timeout(delay, callback) {
	return async function(...args) {
		return Promise.race([
			callback(...args),
			new Promise((_, reject) =>
				setTimeout(() => reject(new Error('timeout')), delay)
			),
		]);
	};
}
