function floor(num) {
	if (num >= 0) {
		let result = 0;
		let step = 1;
		while (result + step <= num) {
			result += step;
			step = step + step;
		}
		step = step / 2;
		while (step >= 1) {
			if (result + step <= num) {
				result += step;
			}
			step = step / 2;
		}
		return result;
	} else {
		let result = 0;
		let step = 1;
		while (result > num) {
			result -= step;
			step = step + step;
		}
		step = step / 2;
		while (step >= 1) {
			if (result + step <= num) {
				result += step;
			}
			step = step / 2;
		}
		return result;
	}
}

function ceil(num) {
	let f = floor(num);
	if (f === num) {
		return f;
	}
	return f + 1;
}

function round(num) {
	let f = floor(num);
	let frac = num - f;
	if (frac >= 0.5) {
		return f + 1;
	} else if (frac <= -0.5) {
		return f - 1;
	}
	return f;
}

function trunc(num) {
	if (num >= 0) {
		return floor(num);
	}
	let f = floor(num);
	if (f === num) {
		return f;
	}
	return f + 1;
}
