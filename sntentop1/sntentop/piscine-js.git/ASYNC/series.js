async function series(asyncFns) {
  const results = [];
  
  for (const fn of asyncFns) {
    // await pauses the loop until the current function resolves, ensuring series execution
    results.push(await fn()); 
  }
  
  return results;
}