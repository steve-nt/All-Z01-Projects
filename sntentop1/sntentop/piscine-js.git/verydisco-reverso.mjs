#!/usr/bin/env node

import { readFile } from 'fs';

const filename = process.argv[2];

readFile(filename, 'utf8', (err, data) => {
  try {
    if (err) throw err;

    const words = data.split(' ');
    const originalWords = words.map(word => {
      const halfLength = Math.floor(word.length / 2);
      const firstHalf = word.slice(0, halfLength);
      const secondHalf = word.slice(halfLength);
      return secondHalf + firstHalf;
    });

    console.log(originalWords.join(' '));
  } catch (error) {
    console.error(error);
  }
});
