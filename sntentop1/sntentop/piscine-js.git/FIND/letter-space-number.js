function letterSpaceNumber(str) {
  const matches = str.match(/[a-z] \d(?![a-z\d])/gi);
  return matches || [];
}
