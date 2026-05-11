function neuron(arr) {
  return arr.reduce((acc, str) => {
    // 1. Split the string into the main parts
    const [typeInfo, response] = str.split(' - Response: ');
    const [rawType, item] = typeInfo.split(': ');
    
    // 2. Format the keys
    const typeKey = rawType.toLowerCase(); // e.g., 'questions'
    const singularType = typeKey.slice(0, -1); // e.g., 'question'
    
    // Lowercase, remove punctuation (anything that isn't a word or space), and replace spaces with underscores
    const itemKey = item.toLowerCase().replace(/[^\w\s]/g, '').replace(/\s+/g, '_');

    // 3. Build the nested object structure if it doesn't exist yet
    acc[typeKey] = acc[typeKey] || {};
    acc[typeKey][itemKey] = acc[typeKey][itemKey] || { 
      [singularType]: item, 
      responses: [] 
    };
    
    // 4. Push the response into the array
    acc[typeKey][itemKey].responses.push(response);
    
    return acc;
  }, {});
}