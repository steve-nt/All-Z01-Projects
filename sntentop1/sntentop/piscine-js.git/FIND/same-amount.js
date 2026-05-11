export function sameAmount(str, regex1, regex2) {
  const matches1 = (str.match(regex1.global ? regex1 : new RegExp(regex1.source, regex1.flags + 'g')) || []).length;
  const matches2 = (str.match(regex2.global ? regex2 : new RegExp(regex2.source, regex2.flags + 'g')) || []).length;
  return matches1 === matches2;
}
