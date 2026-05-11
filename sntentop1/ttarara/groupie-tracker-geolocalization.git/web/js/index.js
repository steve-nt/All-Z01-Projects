document.addEventListener("DOMContentLoaded", function () {
           
    const searchInput = document.getElementById("search");
    const suggestionsBox = document.getElementById("suggestions");

    let currentFocus = -1; // Track the currently focused suggestion
    let currentSuggestions = []; // Store current suggestions


    // Fetch and display search suggestions and handle input changes (including deletions)
    searchInput.addEventListener("input", function () {

        const query = searchInput.value.trim().toLowerCase(); // Normalize the query

        if (query.length === 0) {
             // If the input is empty, reset the view
            showAllCards(); // Show all cards
            suggestionsBox.innerHTML = ""; // Clear suggestions if search bar is empty

        } else {

            // Fetch search suggestions from the backend
            fetch(`/search?q=${query}`)
                .then(response => response.json())
                .then(results => {
                    suggestionsBox.innerHTML = ""; // Clear previous suggestions
                    currentSuggestions = results; // Store current suggestions

                    // Sort results: Artists first, then Members, then Locations
                    results.sort((a, b) => {
                    // Define priority levels for each type
                    const typePriority = {
                        "artist": 1,    
                        "member": 2,   
                        "location": 3  
                    };

                    // Compare priority levels
                    const aPriority = typePriority[a.type] || 99; // Default to 99 if type not found
                    const bPriority = typePriority[b.type] || 99;

                    if (aPriority < bPriority) return -1;
                    if (aPriority > bPriority) return 1;

                    // If same priority, sort by query match (matches at the start of the string come first)
                    const aStartsWithQuery = a.display.toLowerCase().startsWith(query);
                    const bStartsWithQuery = b.display.toLowerCase().startsWith(query);

                    if (aStartsWithQuery && !bStartsWithQuery) return -1;
                    if (!aStartsWithQuery && bStartsWithQuery) return 1;

                    return 0; // Keep order if both have the same priority and match status
                });


            
            // Add sorted suggestions to the suggestions box
            results.forEach(item => {
                const div = document.createElement("div");
                div.textContent = item.display; // Display the suggestion text
                div.classList.add("suggestion-item");

                // When a suggestion is clicked, filter the cards
                div.addEventListener("click", function () {
                    searchInput.value = item.display; // Fill the search bar with the selected suggestion
                    suggestionsBox.innerHTML = "";

                    handleSearchSelection(item); // Filter cards based on the selected suggestion

                });
                suggestionsBox.appendChild(div); // Add the suggestion to the suggestions box
            });
            currentFocus = -1; // Reset focus when new suggestions are loaded
        })
        .catch(error => console.error("Error:", error));
    }
});


    // Handle Enter key press
    searchInput.addEventListener("keydown", function (e) {
        if (e.key === "Enter") {
            e.preventDefault(); // Prevent form submission

            if (currentFocus > -1) {
                // If a suggestion is selected, simulate a click on it
                const selectedSuggestion = document.querySelector(".suggestion-item.suggestion-active");
                if (selectedSuggestion) {
                    selectedSuggestion.click(); // Trigger the click event on the selected suggestion
                }
            } else {
                // If no suggestion is selected, show cards for the current query
                const query = searchInput.value.trim();
                if (query.length > 0) {
                    showCardsForSuggestions(currentSuggestions); // Show cards for the current suggestions
                    suggestionsBox.innerHTML = ""; // Clear the suggestions box
                }
            }
        }
    });


    // Keyboard Navigation for Suggestions
    searchInput.addEventListener("keydown", function (e) {
        const suggestions = document.querySelectorAll(".suggestion-item");

        if (e.key === "ArrowDown") {
            e.preventDefault(); // Prevent default behavior (moving cursor to the end of the input)
            currentFocus++;
            if (currentFocus >= suggestions.length) currentFocus = 0; // Loop to the first suggestion
            setActiveSuggestion(suggestions);
            scrollSuggestionIntoView(suggestions[currentFocus]);
        } else if (e.key === "ArrowUp") {
            e.preventDefault(); // Prevent default behavior (moving cursor to the start of the input)
            currentFocus--;
            if (currentFocus < 0) currentFocus = suggestions.length - 1; // Loop to the last suggestion
            setActiveSuggestion(suggestions);
            scrollSuggestionIntoView(suggestions[currentFocus]); // Scroll to the selected suggestion
        } else if (e.key === "Enter") {
            e.preventDefault(); // Prevent form submission
            if (currentFocus > -1 && suggestions[currentFocus]) {
                suggestions[currentFocus].click(); 
            }
        }
    });
    

    // Function to highlight the active suggestion
    function setActiveSuggestion(suggestions) {
        suggestions.forEach((suggestion, index) => {
            if (index === currentFocus) {
                suggestion.classList.add("suggestion-active"); 
            } else {
                suggestion.classList.remove("suggestion-active");
            }
        });
    }

    // Function to scroll the selected suggestion into view
    function scrollSuggestionIntoView(suggestion) {
        if (suggestion) {
            suggestion.scrollIntoView({
                behavior: "smooth", // Smooth scrolling
                block: "nearest",  // Align the suggestion to the top or bottom of the container
            });
        }
    }

    // Function to filter cards based on the selected suggestion
    function handleSearchSelection(selected) {
        const artistId = selected.artistId; // Get the artist ID from the selected suggestion
        const cards = document.querySelectorAll('.card'); // Select all artist cards

        cards.forEach(card => {
            const cardArtistId = card.getAttribute('data-id'); // Get the artist ID from the card
            if (cardArtistId == artistId) {
                card.style.display = 'block'; // Show the card if it matches the selected artist ID
            } else {
                card.style.display = 'none'; // Hide other cards
            }
        });
    }

    // Function to show cards for the current suggestions
    function showCardsForSuggestions(suggestions) {
        const cards = document.querySelectorAll('.card'); // Select all artist cards
        cards.forEach(card => card.style.display = 'none'); // Hide all cards initially

        suggestions.forEach(suggestion => {
            const artistId = suggestion.artistId; // Get the artist ID from the suggestion
            const card = document.querySelector(`.card[data-id="${artistId}"]`);
            if (card) {
                card.style.display = 'block'; // Show the card if it matches the suggestion
            }
        });
    }

    // Function to show all cards (reset the view)
    function showAllCards() {
        const cards = document.querySelectorAll('.card');
        cards.forEach(card => {
            card.style.display = 'block'; // Show all cards
        });

        // Clear the search bar
        const searchInput = document.getElementById("search");
        searchInput.value = ""; // Clear the input field

        // Clear the suggestions box
        const suggestionsBox = document.getElementById("suggestions");
        suggestionsBox.innerHTML = ""; // Clear any visible suggestions
    }

    // Add Show All Button (Reset)
    const showAllButton = document.createElement("button");
    showAllButton.textContent = "Reset";
    showAllButton.id = "show-all-button";
    showAllButton.addEventListener("click", showAllCards);
   

    // Append the button to the search container
    const searchContainer = document.querySelector(".searchContainer");
    searchContainer.appendChild(showAllButton); // Add the button next to the search bar
});


const resetButton = document.getElementById("resetFilters");



// Open Filter Modal
function openFilter() {
    document.getElementById("filterModal").style.display = "block";
    document.getElementById("modalBackdrop").style.display = "block"; // Show backdrop

}

// Close Filter Modal
function closeFilter() {
    document.getElementById("filterModal").style.display = "none";
    document.getElementById("modalBackdrop").style.display = "none"; // Hide backdrop

}

  // Reset Filters
  resetButton.addEventListener("click", function () {
  document.getElementById("creationStart").value = "";
  document.getElementById("creationEnd").value = "";
  document.getElementById("albumStart").value = "";
  document.getElementById("albumEnd").value = "";
  document.getElementById("locationsfilter").value = "";

  document.querySelectorAll("input[name='Member']").forEach(cb => cb.checked = false);
});

function attachYearValidation(id, errorId) {
            const input = document.getElementById(id);
            const errorSpan = document.getElementById(errorId);
    
            input.addEventListener("input", () => {
                const value = input.value;
                if (value && !/^\d{4}$/.test(value)) {
                    input.classList.add("invalid-year");
                    errorSpan.textContent = "Please input a year with 4 digits";
                    errorSpan.style.display = "block";
                } else {
                    input.classList.remove("invalid-year");
                    errorSpan.textContent = "";
                    errorSpan.style.display = "none";
                }
            });
        }
    
        attachYearValidation("creationStart", "creationStartError");
        attachYearValidation("creationEnd", "creationEndError");
        attachYearValidation("albumStart", "albumStartError");
        attachYearValidation("albumEnd", "albumEndError");
    

// Fetch all locations for the filter datalist
document.addEventListener("DOMContentLoaded", function () {
    fetch("/api/minmax")
    .then(res => res.json())
    .then(data => { defaultValues = data; })
    .catch(err => console.error("Error fetching min/max data:", err));

    fetch("/all_locations")
        .then(res => res.json())
        .then(locations => {
            const dataList = document.getElementById("search1");
            dataList.innerHTML = "";
            locations.forEach(loc => {
                const option = document.createElement("option");
                option.value = loc;
                dataList.appendChild(option);
            });
        })
        .catch(err => console.error("Error fetching all_locations:", err));
        
// Validation
document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector("form");
    if (form) {
        form.addEventListener("submit", function (e) {
            const fields = [
                "creationStart",
                "creationEnd",
                "albumStart",
                "albumEnd"
            ];
            let isValid = true;

            fields.forEach((id) => {
                const input = document.getElementById(id);
                const value = input.value.trim();
                // Only validate if there's a value (empty is allowed)
                if (value && !/^\d{4}$/.test(value)) {
                    showPopup("Must be a 4-digit year (e.g. 1980).");
                    isValid = false;
                    // Focus on the invalid field
                    input.focus();
                }
            });

            if (!isValid) {
                e.preventDefault();
                return false;
            }
            return true;
        });
    }
});
});


// Load Number of Members dynamically
fetch("/api/maxmembers")
  .then(res => res.json())
  .then(data => {
    const max = data.maxMembers;
    const container = document.getElementById("membersFieldset");

    container.innerHTML = ""; // clear
    const legend = document.createElement("legend");
    legend.textContent = "Number of Members:";
    container.appendChild(legend);

    for (let i = 1; i <= max; i++) {
      const label = document.createElement("label");
      label.innerHTML = `<input type="checkbox" name="Member" value="${i}"> ${i}`;
      label.style.marginRight = "10px";
      container.appendChild(label);
    }
  })
  .catch(err => console.error("Error fetching maxmembers:", err));


// Uncheck all member checkboxes
document.querySelectorAll("input[name='Member']").forEach(checkbox => {
    checkbox.checked = false;
});

window.onload = function() {
  document.getElementById("yourFilterForm").reset();
};
