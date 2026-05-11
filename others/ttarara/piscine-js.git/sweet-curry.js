function mult2(number1) {
	return function(number2) {
		return number1 * number2;
	}
}

function add3(number1) {
	return function(number2) {
		return function(number3) {
			return number1 + number2 + number3;
		}
	}
}

function sub4(number1) {
	return function(number2) {
		return function(number3) {
			return function(number4) {
				return number1 - number2 - number3 - number4;
			}
		}
	}
}
