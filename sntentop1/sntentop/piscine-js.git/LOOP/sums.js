function sums(n) {
	let result = [];
	
	function helper(remaining, min, current) {
		if (remaining === 0) {
			if (current.length > 1) {
				result.push([...current]);
			}
			return;
		}
		
		for (let i = min; i <= remaining; i++) {
			current.push(i);
			helper(remaining - i, i, current);
			current.pop();
		}
	}
	
	helper(n, 1, []);
	return result;
}
