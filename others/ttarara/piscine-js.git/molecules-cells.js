function RNA(dnaStr) {
  let res = '';
  for (let i = 0; i < dnaStr.length; i++) {
    if (dnaStr[i] === 'G') res += 'C';
    if (dnaStr[i] === 'C') res += 'G';
    if (dnaStr[i] === 'T') res += 'A';
    if (dnaStr[i] === 'A') res += 'U';
  }
  return res;
}

function DNA(rnaStr) {
  let res = '';
  for (let i = 0; i < rnaStr.length; i++) {
    if (rnaStr[i] === 'C') res += 'G';
    if (rnaStr[i] === 'G') res += 'C';
    if (rnaStr[i] === 'A') res += 'T';
    if (rnaStr[i] === 'U') res += 'A';
  }
  return res;
}
