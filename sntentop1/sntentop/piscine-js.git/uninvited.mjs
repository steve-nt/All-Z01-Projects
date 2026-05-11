#!/usr/bin/env node

import { createServer } from 'http';
import { writeFile } from 'fs';
import { join } from 'path';

const server = createServer((req, res) => {
  const guest = req.url.slice(1);

  res.setHeader('Content-Type', 'application/json');

  if (req.method === 'POST') {
    let body = '';

    req.on('data', chunk => {
      body += chunk.toString();
    });

    req.on('end', () => {
      try {
        let guestData;
        try {
          guestData = JSON.parse(body);
        } catch {
          guestData = body;
        }

        writeFile(join('guests', `${guest}.json`), typeof guestData === 'string' ? guestData : JSON.stringify(guestData, null, '\t'), (err) => {
          if (err) {
            res.writeHead(500);
            res.end(JSON.stringify({ error: 'server failed' }));
            return;
          }

          res.writeHead(201);
          res.end(JSON.stringify(typeof guestData === 'string' ? { data: guestData } : guestData));
        });
      } catch {
        res.writeHead(500);
        res.end(JSON.stringify({ error: 'server failed' }));
      }
    });
  }
});

const PORT = 5000;
server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});
