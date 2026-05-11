function defaultCurry(obj1) {
	return function(obj2) {
		return Object.assign({}, obj1, obj2);
	};
}

function mapCurry(fn) {
	return function(obj) {
		return Object.fromEntries(
			Object.entries(obj).map(fn)
		);
	}
}

function reduceCurry(red) {
	return function(obj, initial = 0) {
		return Object.entries(obj).reduce((acc, [key, value]) => red(acc, [key, value]), initial);
	}
}

function filterCurry(fn) {
	return function(obj) {
		return Object.fromEntries(Object.entries(obj).filter(([key, value]) => fn([key, value])));
	}
}

function reduceScore(personnel, initial = 0) {
	return reduceCurry((acc, [_, user]) => {
		return user.isForceUser
			? acc + user.pilotingScore + user.shootingScore
			: acc;
	})(personnel, initial);
}

function filterForce(personnel, threshold = 80) {
	return reduceCurry((acc, [key, user]) => {
		if (user.isForceUser && user.shootingScore >= threshold) {
			acc[key] = user;
		}
		return acc;
	})(personnel, {});
}

function mapAverage(personnel) {
	return mapCurry(([key, user]) => [
		key,
		{
			averageScore: (user.pilotingScore + user.shootingScore) / 2,
			...user
		}
	])(personnel);
}
