import { writeFile } from 'fs/promises';

const input = process.argv[2];

if (!input) {
	console.error("Please provide a word or sentence to make very disco.");
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

try {
	await writeFile("verydisco-forever.txt", result);
	console.log("Disco result saved to verydisco-forever.txt!");
} catch (err) {
	console.error("Failed to write file:", err);
}
