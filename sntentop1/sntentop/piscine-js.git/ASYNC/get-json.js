async function getJSON(path, params) {
  // 1. Construct the valid URL with optional query parameters
  let url = path;
  
  if (params) {
    const searchParams = new URLSearchParams(params).toString();
    if (searchParams) {
      // Append the parameters to the path
      url += `?${searchParams}`;
    }
  }

  // 2. Fulfil the request using fetch
  const response = await fetch(url);

  // 3. If response is not OK, throw an error with the status text
  if (!response.ok) {
    throw new Error(response.statusText);
  }

  // 4. Read and parse the response body from JSON
  const parsedObject = await response.json();

  // 5. Handle the parsed object properties
  if (parsedObject.error) {
    throw new Error(parsedObject.error);
  }

  return parsedObject.data;
}