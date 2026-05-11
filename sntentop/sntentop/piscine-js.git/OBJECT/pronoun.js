function pronoun(str) {
  const targetPronouns = ['i', 'you', 'he', 'she', 'it', 'they', 'we'];
  
  // Convert string to lowercase, split by any non-word characters (spaces/punctuation), 
  // and remove any empty strings from the resulting array.
  const words = str.toLowerCase().split(/\W+/).filter(Boolean);
  const result = {};

  words.forEach((word, index) => {
    // Check if the current word is in our pronoun list
    if (targetPronouns.includes(word)) {
      
      // Initialize the pronoun in our result object if it doesn't exist yet
      if (!result[word]) {
        result[word] = { word: [], count: 0 };
      }
      
      // Increment the count
      result[word].count++;

      // Look at the next word in the array
      const nextWord = words[index + 1];
      
      // If there is a next word, and it is NOT also a pronoun, push it to the array
      if (nextWord && !targetPronouns.includes(nextWord)) {
        result[word].word.push(nextWord);
      }
    }
  });

  return result;
}