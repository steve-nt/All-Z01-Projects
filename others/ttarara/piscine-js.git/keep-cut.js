
export function cutFirst(a) {
  return a.slice(2);
}


export function cutLast(a) {
  return a.slice(0, -2);
}

export function cutFirstLast(a) {
  return a.slice(2, -2);
}

export function keepFirst(a) {
  return a.slice(0, 2);
}


export function keepLast(a) {
  return a.slice(-2);
}

export function keepFirstLast(a) {
  return a.length > 4 ? keepFirst(a) + keepLast(a) : a;
}
