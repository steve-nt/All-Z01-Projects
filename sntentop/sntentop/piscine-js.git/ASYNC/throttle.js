// 1. Basic Throttle (No options - Drops calls during wait period)
function throttle(func, wait) {
  let lastTime = 0;

  return function (...args) {
    const now = Date.now();
    
    // If enough time has passed since the last execution, run it.
    if (now - lastTime >= wait) {
      func.apply(this, args);
      lastTime = now;
    }
  };
}


// 2. Advanced Throttle (With strict leading & trailing options)
function opThrottle(func, wait, options = {}) {
  let timer = null;
  let lastTime = 0;
  let lastArgs = null;
  let lastThis = null;

  // Strict Level 14 defaults: False unless explicitly provided
  const leading = options && options.leading === true;
  const trailing = options && options.trailing === true;

  return function (...args) {
    const now = Date.now();
    lastArgs = args;
    lastThis = this;

    // If leading is false and it's the first call, spoof the lastTime
    if (!lastTime && !leading) {
      lastTime = now;
    }

    const remaining = wait - (now - lastTime);

    // Condition 1: Time is up (Leading edge or standard cooldown)
    if (remaining <= 0 || remaining > wait) {
      if (timer) {
        clearTimeout(timer);
        timer = null;
      }
      
      lastTime = now;
      
      if (leading) {
        func.apply(lastThis, lastArgs);
        lastArgs = null;
        lastThis = null;
      }
    } 
    // Condition 2: Inside wait period, but trailing is enabled
    else if (!timer && trailing) {
      timer = setTimeout(() => {
        lastTime = leading ? Date.now() : 0;
        timer = null;
        if (lastArgs) {
          func.apply(lastThis, lastArgs);
          lastArgs = null;
          lastThis = null;
        }
      }, remaining);
    }
  };
}