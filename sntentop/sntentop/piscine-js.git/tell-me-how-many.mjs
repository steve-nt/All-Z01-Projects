#!/usr/bin/env node

import { readdir } from 'fs';

const path = process.argv[2] || '.';

readdir(path, (err, files) => {
  if (err) {
    console.error(err);
    return;
  }
  console.log(files.length);
});
