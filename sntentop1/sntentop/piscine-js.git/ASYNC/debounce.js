// Basic debounce: Delays execution until 'wait' ms have passed since the last call.
const debounce = (func, wait) => {
  let timeout;
  
  return function (...args) {
    clearTimeout(timeout);
    
    // Use an arrow function so 'this' properly inherits from the returned function
    timeout = setTimeout(() => func.apply(this, args), wait);
  };
};

// opDebounce: Adds support for the { leading: true } option.
const opDebounce = (func, wait, options = {}) => {
  let timeout;
  
  return function (...args) {
    // If leading is true and no timeout is currently active, we trigger immediately
    const callNow = options.leading && !timeout;
    
    clearTimeout(timeout);
    
    timeout = setTimeout(() => {
      timeout = null;
      // If we aren't firing on the leading edge, fire on the trailing edge
      if (!options.leading) {
        func.apply(this, args);
      }
    }, wait);
    
    // Execute immediately if the conditions for a leading call are met
    if (callNow) {
      func.apply(this, args);
    }
  };
};