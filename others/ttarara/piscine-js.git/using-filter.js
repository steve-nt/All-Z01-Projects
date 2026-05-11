function filterShortStateName(array){
    return array.filter((array) => array.length < 7);
}

function filterStartVowel(array) {
    return array.filter(word => 'aeiouAEIOU'.includes(word[0].toLowerCase()));
}

function filter5Vowels(array) {
    return array.filter(str => {
        const matches = str.match(/[aeiou]/gi); // match all vowels, case-insensitive
        return matches.length >= 5;
    });
}

function filter1DistinctVowel(words) {
	return words.filter(word => {
		const vowels = (word.match(/[aeiou]/gi) || []).map(v => v.toLowerCase());
		const distinctVowels = new Set(vowels);
		return distinctVowels.size === 1;
	});
}

function multiFilter(data) {
	const vowels = /[aeiou]/i;

	return data.filter(item => {
		return (
			item.capital.length >= 8 &&
			!/^[aeiou]/i.test(item.name) &&
			vowels.test(item.tag) &&
			item.region !== 'South'
		);
	});
}