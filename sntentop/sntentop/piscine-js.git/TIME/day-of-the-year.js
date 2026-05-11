function dayOfTheYear(date) {
  const year = date.getFullYear();
  let jan1;
  
  // For early years, use ISO string format to avoid 2-digit year interpretation
  if (year <= 100) {
    const yearStr = String(year).padStart(4, '0');
    jan1 = new Date(`${yearStr}-01-01`);
  } else {
    jan1 = new Date(year, 0, 1);
  }
  
  const diff = date - jan1;
  return Math.floor(diff / (1000 * 60 * 60 * 24)) + 1;
}

export { dayOfTheYear };
