function filterEntries(obj, callback) {
  return Object.fromEntries(Object.entries(obj).filter(callback));
}

function mapEntries(obj, callback) {
  return Object.fromEntries(Object.entries(obj).map(callback));
}

function reduceEntries(obj, callback, ...args) {
  // ...args elegantly catches the initialValue if it exists, or stays empty if it doesn't
  return Object.entries(obj).reduce(callback, ...args);
}

function totalCalories(cart) {
  return Number(
    reduceEntries(
      cart,
      (acc, [item, grams]) => acc + (nutritionDB[item].calories * grams) / 100,
      0
    ).toFixed(3)
  );
}

function lowCarbs(cart) {
  return filterEntries(cart, ([item, grams]) => {
    return (nutritionDB[item].carbs * grams) / 100 < 50;
  });
}

function cartTotal(cart) {
  return mapEntries(cart, ([item, grams]) => {
    const scaled = {};
    for (const key in nutritionDB[item]) {
      // .toFixed(3) caps decimals, Number() converts the string back into a pure number
      scaled[key] = Number(((nutritionDB[item][key] * grams) / 100).toFixed(3));
    }
    return [item, scaled];
  });
}

export { filterEntries, mapEntries, reduceEntries, totalCalories, lowCarbs, cartTotal };