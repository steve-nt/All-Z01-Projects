import { readdir, readFile, writeFile } from 'fs/promises';
import path from 'path';

const dirPath = process.argv[2] || '.';

try {
	const resolvedPath = path.resolve(dirPath);
	const files = await readdir(resolvedPath);

	const vipGuests = [];

	for (const file of files) {
		if (!file.endsWith('.json')) continue;

		const filePath = path.join(resolvedPath, file);
		const content = await readFile(filePath, 'utf-8');

		let data;
		try {
			data = JSON.parse(content);
		} catch {
			continue;
		}

		if (typeof data.answer === 'string' && data.answer.toLowerCase() === 'yes') {
			const base = path.parse(file).name;
			const [first, last] = base.split('_');
			if (first && last) {
				vipGuests.push({ first, last });
			}
		}
	}

	vipGuests.sort((a, b) => {
		const lastCmp = a.last.localeCompare(b.last);
		return lastCmp !== 0 ? lastCmp : a.first.localeCompare(b.first);
	});

	const lines = vipGuests.map((g, i) => `${i + 1}. ${g.last} ${g.first}`);

	await writeFile('vip.txt', lines.join('\n'), 'utf-8');

	lines.forEach(line => console.log(line));

} catch (err) {
	console.error(`Error: ${err.message}`);
	process.exit(1);
}
