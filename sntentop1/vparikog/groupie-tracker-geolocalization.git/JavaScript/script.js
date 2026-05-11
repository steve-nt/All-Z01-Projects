// Function to clear the search bar and reset filtering
function clearSearchBar() {
  // Get the search input element
  const searchBar = document.getElementById('search-bar');
  
  // Clear the input value
  searchBar.value = '';

  // Reset the search suggestions dropdown
  const dropdownMenu = document.getElementById('search-suggestions');
  dropdownMenu.innerHTML = '';
  dropdownMenu.style.display = 'none';

  // Reset the artist cards display
  const artistCards = document.querySelectorAll('.artist-card');
  artistCards.forEach(card => {
      card.style.display = 'block'; // Show all artist cards
  });
}

// Existing search function
async function searchArtists() {
  const query = document.getElementById('search-bar').value;
  const dropdownMenu = document.getElementById('search-suggestions');

  // If the query is empty, hide the dropdown and reset artist cards
  if (!query.trim()) {
    dropdownMenu.innerHTML = '';
    dropdownMenu.style.display = 'none';
    const artistCards = document.querySelectorAll('.artist-card');
    artistCards.forEach(card => {
        card.style.display = 'block'; // Show all artist cards when no search query
    });
    return;
}
  try {
    const response = await fetch(`/search?q=${encodeURIComponent(query)}`);
    const results = await response.json();

    // Clear previous suggestions
    dropdownMenu.innerHTML = '';
    dropdownMenu.style.display = 'none';

    if (results.length > 0) {
      results.forEach(result => {
        const listItem = document.createElement('div');
        listItem.classList.add('result-item');
        listItem.textContent = result.name + " - " + result.type;
        listItem.onclick = () => {
          // Redirect the user to the artist's page
          window.location.href = `/artist?id=${result.id}`; // Assuming result.id contains the artist's ID
        };
        dropdownMenu.appendChild(listItem);
      });
      dropdownMenu.style.display = 'block'; // Show the dropdown with results
    }
  } catch (error) {
    console.error('Error fetching search results:', error);
    dropdownMenu.innerHTML = '<div class="result-item">No Matches Found</div>';
    dropdownMenu.style.display = 'block';
  }
}

// Filter artist cards as user types
document.addEventListener("DOMContentLoaded", function () {
  const searchBar = document.getElementById("search-bar");

  searchBar.addEventListener("input", () => {
      const query = searchBar.value.toLowerCase();
      const artistCards = document.querySelectorAll(".artist-card");
      let anyMatch = false; // ✅ Fixed: Now properly declared

      artistCards.forEach(card => {
          const artistName = card.querySelector("h2").textContent.toLowerCase();
          const namesMembers = card.getAttribute("name-members").toLowerCase();

         
          if (artistName.includes(query) || namesMembers.includes(query)) {
              card.style.display = "block";
              anyMatch = true; // ✅ Track if at least one match is found
          } else {
              card.style.display = "none";
          }
      });

      // ✅ If no match is found, show all cards again
      if (!anyMatch) {
          artistCards.forEach(card => {
              card.style.display = "block";
          });
      }
  });
});


// Hide dropdown if clicked outside
document.addEventListener('click', (event) => {
  const dropdownMenu = document.getElementById('search-suggestions');
  const searchBar = document.getElementById('search-bar');
  if (!searchBar.contains(event.target)) {
    dropdownMenu.style.display = 'none';
  }
});
