import { readFile } from 'fs/promises';

const filename = process.argv[2];

if (!filename) {
	console.error("Please provide a filename to decode from very disco mode.");
	process.exit(1);
}

function undoVeryDisco(word) {
	const mid = Math.ceil(word.length / 2);
	const firstPart = word.slice(0, word.length - mid);
	const secondPart = word.slice(word.length - mid);
	return secondPart + firstPart;
}

try {
	const content = await readFile(filename, 'utf-8');

	const result = content
		.split(" ")
		.map(undoVeryDisco)
		.join(" ");

	console.log(result);
} catch (err) {
	console.error("Could not read or decode the file:", err.message);
}
