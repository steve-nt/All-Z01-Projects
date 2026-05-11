const queryServers = (serverName, q) => {
  const url = `/${serverName}?q=${q}`;
  const backup = `/${serverName}_backup?q=${q}`;
  
  // Promise.race returns the result of the first promise to resolve
  return Promise.race([getJSON(url), getJSON(backup)]);
};

const gougleSearch = (q) => {
  // Create a timeout promise that rejects after 80ms
  const timeout = new Promise((_, reject) => {
    setTimeout(() => reject(new Error('timeout')), 80);
  });

  // Fetch all 3 servers concurrently and map the results to an object
  const searchPromises = Promise.all([
    queryServers('web', q),
    queryServers('image', q),
    queryServers('video', q)
  ]).then(([web, image, video]) => ({ web, image, video }));

  // Race the API calls against the 80ms timeout
  return Promise.race([searchPromises, timeout]);
};