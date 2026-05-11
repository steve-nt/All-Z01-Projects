// Accurately checks if an item is a plain JSON-like object
const isPlainObject = (item) => {
  return Object.prototype.toString.call(item) === '[object Object]';
};

const replica = (target, ...sources) => {
  sources.forEach((source) => {
    Object.keys(source).forEach((key) => {
      if (isPlainObject(source[key])) {
        // If the property doesn't exist or isn't a plain object, initialize it
        if (!isPlainObject(target[key])) {
          target[key] = {};
        }
        // Deep dive into the nested object
        replica(target[key], source[key]);
      } else {
        // Direct assignment for primitives, arrays, functions, regex, dates, etc.
        target[key] = source[key];
      }
    });
  });
  return target;
};