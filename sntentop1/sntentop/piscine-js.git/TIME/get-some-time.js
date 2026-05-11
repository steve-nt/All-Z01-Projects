function firstDayWeek(week, year) {
  const y = parseInt(year);
  const yearPadded = String(y).padStart(4, '0');
  
  // Create date differently for year 1 to avoid 2-digit year interpretation
  let jan1;
  if (y <= 100) {
    // Use ISO string format for years <= 100
    jan1 = new Date(`${yearPadded}-01-01`);
  } else {
    jan1 = new Date(y, 0, 1);
  }
  
  const dayOfWeek = jan1.getDay(); // 0 = Sunday, 1 = Monday, etc.
  
  // Calculate days back to Monday from January 1
  const daysBack = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
  
  // Get the Monday of the week containing Jan 1
  let monday1 = new Date(jan1);
  monday1.setDate(1 - daysBack);
  
  // Calculate the date for the requested week
  let resultDate = new Date(monday1);
  resultDate.setDate(resultDate.getDate() + (week - 1) * 7);
  
  // If the result is in the previous year, return Jan 1 instead
  if (resultDate.getFullYear() < y) {
    resultDate = new Date(`${yearPadded}-01-01`);
  }
  
  // Format as dd-mm-yyyy (with padding for years <= 100)
  const day = String(resultDate.getDate()).padStart(2, '0');
  const month = String(resultDate.getMonth() + 1).padStart(2, '0');
  const resultYear = resultDate.getFullYear();
  const resultYearStr = String(resultYear).padStart(4, '0');
  
  return `${day}-${month}-${resultYearStr}`;
}

export { firstDayWeek };
