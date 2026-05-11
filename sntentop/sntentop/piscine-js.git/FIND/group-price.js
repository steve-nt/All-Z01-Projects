export function groupPrice(str) {
  const regex = /[A-Za-z$]+(\d+)\.(\d+)/g;
  const matches = Array.from(str.matchAll(regex));
  return matches.map(match => [match[0], match[1], match[2]]);
}

