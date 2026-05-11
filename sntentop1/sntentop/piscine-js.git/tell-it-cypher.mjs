#!/usr/bin/env node

import { readFile, writeFile } from 'fs';

const inputFile = process.argv[2];
const action = process.argv[3];
const outputFile = process.argv[4] || (action === 'encode' ? 'cypher.txt' : 'clear.txt');

readFile(inputFile, (err, data) => {
  if (err) {
    console.error(err);
    return;
  }

  let result;
  if (action === 'encode') {
    result = Buffer.from(data).toString('base64');
  } else if (action === 'decode') {
    result = Buffer.from(data.toString(), 'base64').toString();
  }

  writeFile(outputFile, result, err => {
    if (err) console.error(err);
  });
});
