async function isWinner(countryName) {
  try {
    // Attempt to get the country. If it's not found, this will throw an error
    // which immediately skips to the catch block.
    const country = await db.getWinner(countryName);

    // Check continent requirement
    if (country.continent !== 'Europe') {
      return `${countryName} is not what we are looking for because of the continent`;
    }

    // Get the results for the found country
    const results = await db.getResults(country.id);

    // Check number of wins
    if (results.length < 3) {
      return `${countryName} is not what we are looking for because of the number of times it was champion`;
    }

    // Map out the years and scores, joining them with a comma and a space
    const years = results.map(result => result.year).join(', ');
    const scores = results.map(result => result.score).join(', ');

    // Return the final success string
    return `${countryName} won the FIFA World Cup in ${years} winning by ${scores}`;

  } catch (error) {
    // If db.getWinner throws an error (e.g., "Country Not Found"), it lands here
    return `${countryName} never was a winner`;
  }
}