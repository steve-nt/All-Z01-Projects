function get(src, path) {
	let keys = path.split('.');
	let current = src;
	for (let i = 0; i < keys.length; i++) {
		if (current === undefined || current === null) {
			return undefined;
		}
		current = current[keys[i]];
	}	
	return current;
}
