function all(promisesObj = {}) {
  return new Promise((resolve, reject) => {
    const keys = Object.keys(promisesObj);
    
    // If the object is empty, resolve immediately with an empty object
    if (keys.length === 0) {
      return resolve({});
    }

    const result = {};
    let completedCount = 0;

    keys.forEach(key => {
      // Wrap the value in Promise.resolve() to handle plain values safely
      Promise.resolve(promisesObj[key])
        .then(value => {
          result[key] = value;
          completedCount++;
          
          // Once all keys have resolved, resolve the main promise
          if (completedCount === keys.length) {
            resolve(result);
          }
        })
        .catch(reject); // If any promise fails, reject the main promise immediately
    });
  });
}