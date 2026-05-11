import { readdir } from 'fs/promises';
import path from 'path';

const dirPath = process.argv[2] || '.';

try {
	const resolvedPath = path.resolve(dirPath);

	const entries = await readdir(resolvedPath);

	console.log(entries.length);
} catch (err) {
	console.error(`Failed to read directory "${dirPath}": ${err.message}`);
	process.exit(1);
}
