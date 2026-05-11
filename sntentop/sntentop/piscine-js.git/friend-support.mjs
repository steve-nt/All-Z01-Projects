#!/usr/bin/env node

import { createServer } from 'http';
import { readFile } from 'fs';
import { join } from 'path';

const server = createServer((req, res) => {
  const guest = req.url.slice(1);

  res.setHeader('Content-Type', 'application/json');

  readFile(join('guests', `${guest}.json`), 'utf8', (err, data) => {
    if (err) {
      if (err.code === 'ENOENT') {
        res.writeHead(404);
        res.end(JSON.stringify({ error: 'guest not found' }));
      } else {
        res.writeHead(500);
        res.end(JSON.stringify({ error: 'server failed' }));
      }
      return;
    }

    try {
      const guestData = JSON.parse(data);
      res.writeHead(200);
      res.end(JSON.stringify(guestData));
    } catch {
      res.writeHead(500);
      res.end(JSON.stringify({ error: 'server failed' }));
    }
  });
});

const PORT = 5000;
server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});
