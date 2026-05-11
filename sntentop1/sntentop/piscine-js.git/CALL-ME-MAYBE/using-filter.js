// Filter states with name length less than 7 characters
const filterShortStateName = (arr) => arr.filter((state) => state.length < 7);

// Filter strings that start with a vowel
const filterStartVowel = (arr) =>
  arr.filter((str) => /^[aeiou]/i.test(str));

// Filter strings with at least 5 vowels
const filter5Vowels = (arr) =>
  arr.filter((str) => {
    const vowels = str.match(/[aeiou]/gi);
    return vowels && vowels.length >= 5;
  });

// Filter strings with only one distinct vowel
const filter1DistinctVowel = (arr) =>
  arr.filter((str) => {
    const vowels = new Set(str.toLowerCase().match(/[aeiou]/g));
    return vowels.size === 1;
  });

// Filter objects by multiple conditions
const multiFilter = (arr) =>
  arr.filter(
    (obj) =>
      obj.capital.length >= 8 &&
      !/^[aeiou]/i.test(obj.name) &&
      /[aeiou]/i.test(obj.tag) &&
      obj.region !== 'South'
  );
