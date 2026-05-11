/**
 * Resolves or rejects as soon as the first promise in the array resolves or rejects.
 */
function race(promises) {
  return new Promise((resolve, reject) => {
    promises.forEach(promise => {
      // Promise.resolve normalizes both raw values and actual promises
      Promise.resolve(promise).then(resolve, reject);
    });
  });
}

/**
 * Resolves with an array of the first `count` resolved values.
 * If empty array or 0 count, resolves with undefined.
 */
function some(promises, count) {
  if (promises.length === 0 || count === 0) {
    return Promise.resolve([]);
  }

  return new Promise((resolve, reject) => {
    if (promises.length < count) {
      return reject(new Error("Not enough promises provided"));
    }

    const results = [];
    let resolvedCount = 0;
    let rejections = 0;

    // Notice we grab the 'index' from the forEach loop here
    promises.forEach((promise, index) => {
      Promise.resolve(promise).then(
        value => {
          if (resolvedCount < count) {
            // Store both the value and its original index
            results.push({ index, value });
            resolvedCount++;
            
            if (resolvedCount === count) {
              // Sort the results by their original index to preserve order
              results.sort((a, b) => a.index - b.index);
              // Map over the sorted array to return just the values
              resolve(results.map(r => r.value));
            }
          }
        },
        error => {
          rejections++;
          if (promises.length - rejections < count) {
            reject(new Error("Too many rejections to satisfy count"));
          }
        }
      );
    });
  });
}