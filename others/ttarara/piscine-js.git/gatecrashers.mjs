import http from 'node:http';
import { mkdir, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';

const PORT = 5000;
const BASE_DIR = 'guests';
const AUTH_USERS = new Map([
	['Caleb_Squires', 'abracadabra'],
	['Tyrique_Dalton', 'abracadabra'],
	['Rahima_Young', 'abracadabra'],
]);

mkdir(resolve(process.cwd(), BASE_DIR), { recursive: true })
	.then(() => {
		http.createServer((req, res) => {
			const guest = decodeURIComponent(req.url.slice(1));
			if (req.method !== 'POST' || !guest) {
				res.writeHead(404, { 'Content-Type': 'application/json' });
				return res.end(JSON.stringify({ error: 'not found' }));
			}

			const auth = req.headers.authorization;
			if (!auth || !auth.startsWith('Basic ')) {
				res.writeHead(401, {
					'Content-Type': 'application/json',
					'WWW-Authenticate': 'Basic'
				});
				return res.end(JSON.stringify({ error: 'Authorization Required' }));
			}
			const [user, pass] = Buffer.from(auth.split(' ')[1], 'base64')
				.toString().split(':');
			if (!AUTH_USERS.has(user) || AUTH_USERS.get(user) !== pass) {
				res.writeHead(401, {
					'Content-Type': 'application/json',
					'WWW-Authenticate': 'Basic'
				});
				return res.end(JSON.stringify({ error: 'Authorization Required' }));
			}

			const headerBody = req.headers['body'];
			if (headerBody !== undefined) {
				let data;
				try { data = JSON.parse(headerBody); } catch { data = headerBody; }
				const filePath = resolve(process.cwd(), BASE_DIR, `${guest}.json`);
				writeFile(
					filePath,
					typeof data === 'object' ? JSON.stringify(data) : data,
					'utf8'
				).catch(() => { });
				res.writeHead(200, { 'Content-Type': 'application/json' });
				return res.end(JSON.stringify(data));
			}

			const chunks = [];
			req.on('data', chunk => chunks.push(chunk));
			req.on('end', async () => {
				try {
					const text = Buffer.concat(chunks).toString('utf8');
					let data;
					try { data = JSON.parse(text); } catch { data = text; }

					const filePath = resolve(process.cwd(), BASE_DIR, `${guest}.json`);
					await writeFile(
						filePath,
						typeof data === 'object' ? JSON.stringify(data) : data,
						'utf8'
					);

					res.writeHead(200, { 'Content-Type': 'application/json' });
					res.end(JSON.stringify(data));
				} catch (err) {
					console.error('Error handling POST:', err);
					res.writeHead(500, { 'Content-Type': 'application/json' });
					res.end(JSON.stringify({ error: 'server failed' }));
				}
			});

		}).listen(PORT, () => console.log(`Listening on port ${PORT}`));
	})
	.catch(err => {
		console.error('Initialization failed:', err);
		process.exit(1);
	});
