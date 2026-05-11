function interpolation({ step, start, end, callback, duration }) {
  let currentStep = 0;
  
  // Calculate the amount to increment per step
  const timeInterval = duration / step;
  const distInterval = (end - start) / step;

  const timer = setInterval(() => {
    // Calculate exact x (distance) and y (time)
    // We use multiplication instead of repeated addition to avoid compounding float errors
    const x = start + distInterval * currentStep;
    const y = timeInterval * (currentStep + 1);

    // Clean up JavaScript's notorious floating-point math quirks (e.g., 0.2 + 0.4 = 0.6000000000000001)
    const cleanX = Math.round(x * 100) / 100;
    const cleanY = Math.round(y * 100) / 100;

    // Execute the callback with the formatted [x, y] coordinates
    callback([cleanX, cleanY]);

    currentStep++;

    // Stop the interval once we've reached the required number of steps
    if (currentStep === step) {
      clearInterval(timer);
    }
  }, timeInterval);
}