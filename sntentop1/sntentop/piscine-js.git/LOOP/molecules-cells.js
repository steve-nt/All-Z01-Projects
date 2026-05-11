function RNA(dna) {
	let result = '';
	for (let i = 0; i < dna.length; i++) {
		if (dna[i] === 'G') {
			result += 'C';
		} else if (dna[i] === 'C') {
			result += 'G';
		} else if (dna[i] === 'T') {
			result += 'A';
		} else if (dna[i] === 'A') {
			result += 'U';
		}
	}
	return result;
}

function DNA(rna) {
	let result = '';
	for (let i = 0; i < rna.length; i++) {
		if (rna[i] === 'C') {
			result += 'G';
		} else if (rna[i] === 'G') {
			result += 'C';
		} else if (rna[i] === 'A') {
			result += 'T';
		} else if (rna[i] === 'U') {
			result += 'A';
		}
	}
	return result;
}
