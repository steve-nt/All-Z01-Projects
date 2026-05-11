#!/usr/bin/env node

import { writeFile } from 'fs';

const argument = process.argv[2];

const words = argument.split(' ');
const discoWords = words.map(word => {
  const halfLength = Math.ceil(word.length / 2);
  const firstHalf = word.slice(0, halfLength);
  const secondHalf = word.slice(halfLength);
  return secondHalf + firstHalf;
});

const result = discoWords.join(' ');

writeFile('verydisco-forever.txt', result, err => {
  if (err) console.error(err);
});
