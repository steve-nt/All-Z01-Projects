function format(date, formatStr) {
  const monthNames = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];
  const monthAbbr = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
  const dayAbbr = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  const year = Math.abs(date.getFullYear());
  const month = date.getMonth();
  const day = date.getDate();
  const hour24 = date.getHours();
  const hour12 = hour24 % 12 || 12;
  const minute = date.getMinutes();
  const second = date.getSeconds();
  const dayOfWeek = date.getDay();
  const era = date.getFullYear() > 0 ? 'AD' : 'BC';

  const pad = (n) => String(n).padStart(2, '0');

  let result = '';
  let i = 0;

  while (i < formatStr.length) {
    // Try to match 4-letter patterns first
    if (i + 4 <= formatStr.length) {
      const four = formatStr.substr(i, 4);
      if (four === 'GGGG') {
        result += era === 'AD' ? 'Anno Domini' : 'Before Christ';
        i += 4;
        continue;
      }
      if (four === 'MMMM') {
        result += monthNames[month];
        i += 4;
        continue;
      }
      if (four === 'EEEE') {
        result += dayNames[dayOfWeek];
        i += 4;
        continue;
      }
      if (four === 'yyyy') {
        result += String(year).padStart(4, '0');
        i += 4;
        continue;
      }
    }

    // Try to match 3-letter patterns
    if (i + 3 <= formatStr.length) {
      const three = formatStr.substr(i, 3);
      if (three === 'MMM') {
        result += monthAbbr[month];
        i += 3;
        continue;
      }
      if (three === 'hh' || three === 'mm' || three === 'ss') {
        // Will handle these with 2-letter check below
      }
    }

    // Try to match 2-letter patterns
    if (i + 2 <= formatStr.length) {
      const two = formatStr.substr(i, 2);
      if (two === 'MM') {
        result += pad(month + 1);
        i += 2;
        continue;
      }
      if (two === 'dd') {
        result += pad(day);
        i += 2;
        continue;
      }
      if (two === 'HH') {
        result += pad(hour24);
        i += 2;
        continue;
      }
      if (two === 'hh') {
        result += pad(hour12);
        i += 2;
        continue;
      }
      if (two === 'mm') {
        result += pad(minute);
        i += 2;
        continue;
      }
      if (two === 'ss') {
        result += pad(second);
        i += 2;
        continue;
      }
    }

    // Try to match 1-letter patterns
    const char = formatStr[i];
    if (char === 'G') {
      result += era;
      i++;
      continue;
    }
    if (char === 'M') {
      result += String(month + 1);
      i++;
      continue;
    }
    if (char === 'E') {
      result += dayAbbr[dayOfWeek];
      i++;
      continue;
    }
    if (char === 'y') {
      result += String(year);
      i++;
      continue;
    }
    if (char === 'd') {
      result += String(day);
      i++;
      continue;
    }
    if (char === 'H') {
      result += String(hour24);
      i++;
      continue;
    }
    if (char === 'h') {
      result += String(hour12);
      i++;
      continue;
    }
    if (char === 'm') {
      result += String(minute);
      i++;
      continue;
    }
    if (char === 's') {
      result += String(second);
      i++;
      continue;
    }
    if (char === 'a') {
      result += hour24 < 12 ? 'AM' : 'PM';
      i++;
      continue;
    }

    // Not a format token, just add the character
    result += char;
    i++;
  }

  return result;
}

export { format };
