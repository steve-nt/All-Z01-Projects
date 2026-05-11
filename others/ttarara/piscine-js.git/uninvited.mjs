import http from 'node:http';
import { mkdir, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';

const PORT = 5000;
const BASE_DIR = 'guests';

mkdir(resolve(process.cwd(), BASE_DIR), { recursive: true })
	.then(() => {
		const server = http.createServer((req, res) => {
			const guest = decodeURIComponent(req.url.slice(1));
			if (req.method !== 'POST' || !guest) {
				res.writeHead(404, { 'Content-Type': 'application/json' });
				return res.end(JSON.stringify({ error: 'not found' }));
			}

			let body = '';
			req.on('data', chunk => { body += chunk; });
			req.on('end', async () => {
				try {
					const contentType = req.headers['content-type'] || '';
					let data;
					if (contentType.includes('application/json')) {
						try {
							data = JSON.parse(body);
						} catch {
							data = body;
						}
					} else {
						data = body;
					}

					const filePath = resolve(process.cwd(), BASE_DIR, `${guest}.json`);
					await writeFile(
						filePath,
						typeof data === 'object' ? JSON.stringify(data) : data,
						'utf8'
					);

					res.writeHead(201, { 'Content-Type': 'application/json' });
					res.end(JSON.stringify(data));
				} catch (err) {
					console.error('Error handling POST:', err);
					res.writeHead(500, { 'Content-Type': 'application/json' });
					res.end(JSON.stringify({ error: 'server failed' }));
				}
			});
		});

		server.listen(PORT, () => console.log(`Listening on port ${PORT}`));
	})
	.catch(err => {
		console.error('Failed to set up directory:', err);
		process.exit(1);
	});
