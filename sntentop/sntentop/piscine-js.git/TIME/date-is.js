function isValid(date) {
  if (date instanceof Date) {
    return !isNaN(date.getTime());
  }
  if (typeof date === 'number') {
    return !isNaN(date) && isFinite(date);
  }
  return false;
}

function isAfter(date1, date2) {
  return date1 > date2;
}

function isBefore(date1, date2) {
  return date2 > date1;
}

function isFuture(date) {
  return isValid(date) && date > new Date();
}

function isPast(date) {
  return isValid(date) && date < new Date();
}

export { isValid, isAfter, isBefore, isFuture, isPast };
