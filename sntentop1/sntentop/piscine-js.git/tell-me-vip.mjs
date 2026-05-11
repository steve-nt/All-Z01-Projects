#!/usr/bin/env node

import { readdir, readFile, writeFile } from 'fs';
import { join } from 'path';

const path = process.argv[2];

readdir(path, (err, files) => {
  if (err) {
    console.error(err);
    return;
  }

  const filePromises = files.map(file => {
    return new Promise((resolve) => {
      readFile(join(path, file), 'utf8', (err, data) => {
        if (err) {
          resolve(null);
          return;
        }
        const guest = JSON.parse(data);
        if (guest.answer === 'yes') {
          const name = file.slice(0, -5);
          const parts = name.split('_');
          resolve(parts[1] + ' ' + parts[0]);
        } else {
          resolve(null);
        }
      });
    });
  });

  Promise.all(filePromises).then(guests => {
    const vips = guests.filter(g => g !== null).sort();
    const output = vips.map((g, i) => `${i + 1}. ${g}`).join('\n');
    writeFile('vip.txt', output, err => {
      if (err) console.error(err);
    });
  });
});
