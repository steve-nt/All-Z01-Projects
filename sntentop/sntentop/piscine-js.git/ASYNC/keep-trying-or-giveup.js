// Retries the callback up to `count` times if it fails
export const retry = (count, callback) => async (...args) => {
  for (let i = 0; i <= count; i++) {
    try {
      // If it succeeds, immediately return the result
      return await callback(...args);
    } catch (error) {
      // If we've reached the maximum number of retries, throw the last error
      if (i === count) throw error;
    }
  }
};

// Races the callback against a timeout limit
export const timeout = (delay, callback) => async (...args) => {
  // Create a promise that rejects after the given delay
  const timeLimit = new Promise((_, reject) => {
    setTimeout(() => reject(new Error('timeout')), delay);
  });

  // Return whichever finishes first: the callback or the timeLimit
  return Promise.race([callback(...args), timeLimit]);
};