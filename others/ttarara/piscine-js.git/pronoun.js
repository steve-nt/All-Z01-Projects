function pronoun(str) {
	const pronouns = ['i', 'you', 'he', 'she', 'it', 'they', 'we'];
	const words = str.toLowerCase().split(/\s+/);
	const result = {};

	for (let i = 0; i < words.length; i++) {
		const current = words[i].replace(/[^a-z]/g, '');

		if (pronouns.includes(current)) {
			if (!result[current]) {
				result[current] = { word: [], count: 0 };
			}

			result[current].count += 1;

			const nextWord = words[i + 1]?.replace(/[^a-z]/g, '');

			if (!pronouns.includes(nextWord) && nextWord) {
				result[current].word.push(nextWord);
			}
		}
	}

	return result;
}
