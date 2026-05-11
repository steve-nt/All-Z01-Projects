// Extract city names from array of objects
const citiesOnly = (arr) => arr.map((obj) => obj.city);

// Capitalize first letter of each word in state names
const upperCasingStates = (arr) =>
  arr.map((state) =>
    state
      .split(' ')
      .map((word) => word[0].toUpperCase() + word.slice(1))
      .join(' ')
  );

// Convert fahrenheit to celsius (rounded down) and format as string
const fahrenheitToCelsius = (arr) =>
  arr.map((temp) => {
    const fahrenheit = parseFloat(temp);
    const celsius = Math.floor((fahrenheit - 32) * (5 / 9));
    return `${celsius}°C`;
  });

// Remove spaces from temperature strings while preserving object structure
const trimTemp = (arr) =>
  arr.map((obj) => ({
    ...obj,
    temperature: obj.temperature.replace(/\s/g, ''),
  }));

// Format weather data into readable forecast strings
const tempForecasts = (arr) =>
  arr.map((obj) => {
    const fahrenheit = parseFloat(obj.temperature);
    const celsius = Math.floor((fahrenheit - 32) * (5 / 9));
    const capitalizedState = obj.state
      .split(' ')
      .map((word) => word[0].toUpperCase() + word.slice(1))
      .join(' ');
    return `${celsius}°Celsius in ${obj.city}, ${capitalizedState}`;
  });
