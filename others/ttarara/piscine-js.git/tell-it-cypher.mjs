import { readFile, writeFile } from 'fs/promises';

const [inputFile, mode, outputName] = process.argv.slice(2);

if (!inputFile || !mode || !['encode', 'decode'].includes(mode)) {
	console.error('Usage: node tell-it-cypher.mjs <inputFile> <encode|decode> [outputFilename.ext]');
	process.exit(1);
}

try {
	const inputContent = await readFile(inputFile);

	let output;
	let defaultName;

	if (mode === 'encode') {
		output = Buffer.from(inputContent).toString('base64');
		defaultName = 'cypher.txt';
	} else {
		output = Buffer.from(inputContent.toString(), 'base64').toString('utf-8');
		defaultName = 'clear.txt';
	}

	const outFile = outputName || defaultName;

	await writeFile(outFile, output, 'utf-8');

	console.log(`File ${mode}d and saved as ${outFile}`);
} catch (err) {
	console.error(`Error: ${err.message}`);
	process.exit(1);
}
