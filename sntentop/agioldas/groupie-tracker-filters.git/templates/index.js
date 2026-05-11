// UPDATE SLIDER LABELS
creation_year_start.oninput = () => (creation_year_start_d.textContent = creation_year_start.value);
creation_year_end.oninput = () => (creation_year_end_d.textContent = creation_year_end.value);
first_album_year_start.oninput = () => (first_album_year_start_d.textContent = first_album_year_start.value);
first_album_year_end.oninput = () => (first_album_year_end_d.textContent = first_album_year_end.value);


// SUGGESTIONS
function fetchSuggestions() {
    const query = document.querySelector('.searchbar').value;
    
    //if searchbar is emtpy, then remove and hide all suggestions
    if (query.trim() === '') {
        document.getElementById('suggestions').innerHTML = '';
        document.getElementById('suggestions').style.display = 'none';
        return;
    }

    //send query to server
    fetch(`/search?query=${encodeURIComponent(query)}`)
    .then(response => response.json())  // Assuming server returns a JSON array
    .then(suggestions => {
        if (suggestions.length === 0) {
            document.getElementById('suggestions').style.display = 'none'; // Hide if no suggestions
        } else {
            //generate suggestions
            const suggestionsList = suggestions.map((suggestion, index) => 
                `<li class="suggestion-element" data-index="${index}">${suggestion}</li>`
            ).join('');
            const suggestionsElement = document.getElementById('suggestions');
            suggestionsElement.innerHTML = suggestionsList;
            suggestionsElement.style.display = 'block'; // Show suggestions
        }
    })
    .catch(error => console.error('Error fetching suggestions:', error));
}

// CLICK ON SUGGESTIONS
document.getElementById('suggestions').addEventListener('click', function(event) {
    const clickedElement = event.target;
    
    // Ensure it's a <li> element
    if (clickedElement.tagName.toLowerCase() === 'li') {
        document.querySelector('.searchbar').value = clickedElement.textContent;  // Set the value in search bar
        document.getElementById('suggestions').style.display = 'none'; // Hide suggestions after selection
    }
});




// RESET BUTTON
const resetButton = document.getElementById('reset_button');  // Make sure you have this button in your HTML
if (resetButton) {
    resetButton.addEventListener('click', resetFilters);
}
function resetFilters() {
    // Get the sliders and their display elements
    const CreationYearStart = document.getElementById('creation_year_start');
    const CreationYearEnd = document.getElementById('creation_year_end');
    const CreationYearStartD = document.getElementById('creation_year_start_d');
    const CreationYearEndD = document.getElementById('creation_year_end_d');

    CreationYearStart.value = 1958;
    CreationYearEnd.value = 2025;
    CreationYearStartD.textContent = CreationYearStart.value;
    CreationYearEndD.textContent = CreationYearEnd.value;

    const firstAlbumYearStart = document.getElementById('first_album_year_start');
    const firstAlbumYearEnd = document.getElementById('first_album_year_end');
    const firstAlbumYearStartD = document.getElementById('first_album_year_start_d');
    const firstAlbumYearEndD = document.getElementById('first_album_year_end_d');

    firstAlbumYearStart.value = 1958;
    firstAlbumYearEnd.value = 2025;
    firstAlbumYearStartD.textContent = firstAlbumYearStart.value;
    firstAlbumYearEndD.textContent = firstAlbumYearEnd.value;

    const checkboxes = document.querySelectorAll('input[name="band_size"]');
    var valuesToCheck = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10"]
    checkboxes.forEach(checkbox => {
        // Set the checkbox checked if its value is in the valuesToCheck array
        checkbox.checked = valuesToCheck.includes(checkbox.value);
    });

    const dropdown = document.getElementById('concert-dropdown');
    const anyOption = dropdown.querySelector('option[value="any"]');
    
    // If "any" option exists, set it as the selected one
    if (anyOption) {
        anyOption.selected = true;
    }

    
}



 // CLICK ON IMAGE
 function setArtistID(id) {
    // Set the value of the hidden input field to the clicked artist's ID
    document.getElementById('artistID').value = id;
    console.log("Clicked artist ID: ", id);

    // Submit the form automatically after setting the ID
    document.getElementById('artist_filters').submit();
}