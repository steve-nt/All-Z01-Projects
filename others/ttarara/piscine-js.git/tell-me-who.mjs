import { readdir } from 'fs/promises';
import path from 'path';

const dirPath = process.argv[2] || '.';

try {
	const resolvedPath = path.resolve(dirPath);
	const files = await readdir(resolvedPath);

	const guests = files
		.map(file => {
			const base = path.parse(file).name;
			const [first, last] = base.split('_');
			if (!first || !last) return null;
			return { first, last };
		})
		.filter(Boolean)
		.sort((a, b) => {
			const lastCmp = a.last.localeCompare(b.last);
			return lastCmp !== 0 ? lastCmp : a.first.localeCompare(b.first);
		});

	guests.forEach((guest, index) => {
		console.log(`${index + 1}. ${guest.last} ${guest.first}`);
	});

} catch (err) {
	console.error(`Could not read directory "${dirPath}": ${err.message}`);
	process.exit(1);
}
