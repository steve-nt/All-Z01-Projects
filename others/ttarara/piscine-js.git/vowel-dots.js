const vowels = /[aeiou]/gi;

export function vowelDots(str) {
  return str.replace(vowels, m => m + '.');
}