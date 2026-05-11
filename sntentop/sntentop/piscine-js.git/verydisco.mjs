#!/usr/bin/env node

const argument = process.argv[2];

const words = argument.split(' ');
const discoWords = words.map(word => {
  const halfLength = Math.ceil(word.length / 2);
  const firstHalf = word.slice(0, halfLength);
  const secondHalf = word.slice(halfLength);
  return secondHalf + firstHalf;
});

console.log(discoWords.join(' '));
