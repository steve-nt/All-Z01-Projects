function matchCron(cron, date) {
  const [minute, hour, day, month, dayOfWeek] = cron.split(' ');
  
  const dateMinute = date.getMinutes();
  const dateHour = date.getHours();
  const dateDay = date.getDate();
  const dateMonth = date.getMonth() + 1;
  const dateDayOfWeek = date.getDay() === 0 ? 7 : date.getDay();
  
  const matchesField = (pattern, value) => {
    if (pattern === '*') return true;
    return parseInt(pattern) === value;
  };
  
  return (
    matchesField(minute, dateMinute) &&
    matchesField(hour, dateHour) &&
    matchesField(day, dateDay) &&
    matchesField(month, dateMonth) &&
    matchesField(dayOfWeek, dateDayOfWeek)
  );
}

export { matchCron };
