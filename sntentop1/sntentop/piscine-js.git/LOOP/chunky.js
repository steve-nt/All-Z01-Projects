function chunk(arr, size) {
	let result = [];
	for (let i = 0; i < arr.length; i += size) {
		let subArray = [];
		for (let j = 0; j < size && i + j < arr.length; j++) {
			subArray.push(arr[i + j]);
		}
		result.push(subArray);
	}
	return result;
}
