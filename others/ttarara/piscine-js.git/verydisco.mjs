const input = process.argv[2];

if (!input) {
	console.log("Please provide a word or sentence to make very disco.");
	process.exit(1);
}

function makeVeryDisco(word) {
	const mid = Math.ceil(word.length / 2);
	const firstPart = word.slice(0, mid);
	const secondPart = word.slice(mid);
	return secondPart + firstPart;
}

const result = input
	.split(" ")
	.map(makeVeryDisco)
	.join(" ");

console.log(result);
