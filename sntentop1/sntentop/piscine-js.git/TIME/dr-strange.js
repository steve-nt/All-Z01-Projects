function addWeek(date) {
  const dayNames = [
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday',
    'Sunday',
    'secondMonday',
    'secondTuesday',
    'secondWednesday',
    'secondThursday',
    'secondFriday',
    'secondSaturday',
    'secondSunday'
  ];

  const epoch = new Date('0001-01-01');
  const daysDiff = Math.floor((date - epoch) / (1000 * 60 * 60 * 24));
  
  return dayNames[daysDiff % 14];
}

function timeTravel({ date, hour, minute, second }) {
  const newDate = new Date(date);
  newDate.setHours(hour, minute, second);
  return newDate;
}

export { addWeek, timeTravel };
