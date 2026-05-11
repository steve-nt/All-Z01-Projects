function flags(obj) {
  const result = {
    alias: { h: 'help' }
  };

  // Get all flags from the object, excluding the 'help' array
  const availableFlags = Object.keys(obj).filter(key => key !== 'help');

  // Build the alias mappings (e.g., 'i' for 'invert')
  availableFlags.forEach(key => {
    result.alias[key[0]] = key;
  });

  // Determine which flags to describe: the ones in the help array, or all available
  const flagsToDescribe = obj.help || availableFlags;

  // Map the chosen flags to their formatted strings and join them with a newline
  result.description = flagsToDescribe
    .map(key => `-${key[0]}, --${key}: ${obj[key]}`)
    .join('\n');

  return result;
}