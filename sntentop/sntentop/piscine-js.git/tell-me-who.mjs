#!/usr/bin/env node

import { readdir } from 'fs';

const path = process.argv[2];

readdir(path, (err, files) => {
  if (err) {
    console.error(err);
    return;
  }

  const guests = files
    .map(file => {
      const name = file.slice(0, -5);
      const parts = name.split('_');
      return parts[1] + ' ' + parts[0];
    })
    .sort();

  guests.forEach((guest, index) => {
    console.log(`${index + 1}. ${guest}`);
  });
});
