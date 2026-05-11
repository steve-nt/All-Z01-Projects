function sunnySunday(date) {
  const dayNames = [
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday'
  ];

  const epoch = new Date('0001-01-01');
  const daysDiff = Math.floor((date - epoch) / (1000 * 60 * 60 * 24));
  
  return dayNames[daysDiff % 6];
}

export { sunnySunday };
