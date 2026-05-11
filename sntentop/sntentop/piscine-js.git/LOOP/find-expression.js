function findExpression(target) {
	if (target < 1) {
		return undefined;
	}
	
	if (target === 1) {
		return '1';
	}
	
	let queue = [[1, '1']];
	let visited = new Set([1]);
	
	while (queue.length > 0) {
		let current = queue.shift();
		let num = current[0];
		let path = current[1];
		
		// Try multiply by 2
		let next = num * 2;
		if (next === target) {
			return path + ' ' + mul2;
		}
		if (next < target && !visited.has(next)) {
			visited.add(next);
			queue.push([next, path + ' ' + mul2]);
		}
		
		// Try add 4
		next = num + 4;
		if (next === target) {
			return path + ' ' + add4;
		}
		if (next < target && !visited.has(next)) {
			visited.add(next);
			queue.push([next, path + ' ' + add4]);
		}
	}
	
	return undefined;
}
