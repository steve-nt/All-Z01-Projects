import http from 'http';
import { readFile } from 'fs/promises';
import { existsSync } from 'fs';
import path from 'path';
import url from 'url';

const PORT = 5000;

const GUESTS_DIR = path.resolve('guests');

const server = http.createServer(async (req, res) => {
	try {
		const parsedUrl = url.parse(req.url, true);
		const pathname = parsedUrl.pathname;

		if (req.method === 'GET' && pathname.length > 1) {
			const guestName = pathname.slice(1);

			const guestFile = path.join(GUESTS_DIR, `${guestName}.json`);

			if (!existsSync(guestFile)) {
				res.writeHead(404, { 'Content-Type': 'application/json' });
				res.end(JSON.stringify({ error: 'guest not found' }));
				return;
			}

			const content = await readFile(guestFile, 'utf-8');

			res.writeHead(200, { 'Content-Type': 'application/json' });
			res.end(content);
		} else {
			res.writeHead(404, { 'Content-Type': 'application/json' });
			res.end(JSON.stringify({ error: 'guest not found' }));
		}
	} catch (err) {
		res.writeHead(500, { 'Content-Type': 'application/json' });
		res.end(JSON.stringify({ error: 'server failed' }));
	}
});

server.listen(PORT, () => {
	console.log(`Server listening on port ${PORT}`);
});
