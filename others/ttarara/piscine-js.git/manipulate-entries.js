function filterEntries(obj, callback) {
	return Object.fromEntries(
		Object.entries(obj).filter(([key, value]) => callback([key, value]))
	);
}

function mapEntries(obj, callback) {
	return Object.fromEntries(
		Object.entries(obj).map(([key, value]) => callback([key, value]))
	);
}

function reduceEntries(obj, callback, initial) {
	return Object.entries(obj).reduce((acc, [key, value]) => {
		return callback(acc, [key, value]);
	}, initial);
}

function totalCalories(cart) {
	return reduceEntries(cart, (acc, [item, grams]) => {
		const itemData = nutritionDB[item];
		if (!itemData) return acc;
		return acc + (itemData.calories * grams) / 100;
	}, 0).toFixed(1) * 1;
}

function lowCarbs(cart) {
	return filterEntries(cart, ([item, grams]) => {
		const itemData = nutritionDB[item];
		if (!itemData) return false;
		const totalCarbs = (itemData.carbs * grams) / 100;
		return totalCarbs < 50;
	});
}

function cartTotal(cart) {
	return reduceEntries(cart, (acc, [item, grams]) => {
		const data = nutritionDB[item];
		if (!data) return acc;

		for (const [nutrient, value] of Object.entries(data)) {
			acc[item] ??= {};
			acc[item][nutrient] = +(value * grams / 100).toFixed(3);
		}

		return acc;
	}, {});
}
