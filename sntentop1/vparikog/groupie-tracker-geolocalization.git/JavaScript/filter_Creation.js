document.addEventListener("DOMContentLoaded", function () {
  console.log("🎵 Filter script loaded!");

  // Create year filter elements
  const fromSlider = document.querySelector('#fromSlider');
  const toSlider = document.querySelector('#toSlider');
  const fromInput = document.querySelector('#fromInput');
  const toInput = document.querySelector('#toInput');

  // First album year filter elements
  const fromAlbumSlider = document.querySelector('#fromAlbumSlider');
  const toAlbumSlider = document.querySelector('#toAlbumSlider');
  const fromAlbumInput = document.querySelector('#fromAlbumInput');
  const toAlbumInput = document.querySelector('#toAlbumInput');

  const applyFilterBtn = document.getElementById("applyFilterBtn");
  const filterPopup = document.getElementById("filterPopup");


  // 🌍 Location filter elements
  const countrySelect = document.getElementById("countrySelect");
  const citySelect = document.getElementById("citySelect");
  
  console.log("🎛 Elements found:", {
    fromInput,
    toInput,
    fromAlbumInput,
    toAlbumInput,
    countrySelect,
    citySelect,
    applyFilterBtn
  });

    // 🌍 Fetch locations from backend API
    fetch("/api/locations")  //api address 
        .then(response => response.json())
        .then(locationsMap => {
          console.log("🌍 Locations fetched:", locationsMap);
            // 📌 Populate country dropdown
            Object.keys(locationsMap).forEach(country => {
                let option = document.createElement("option");
                option.value = country;
                option.textContent = country.replace("_", " ").toUpperCase();
                countrySelect.appendChild(option);
            });

            // 📌 Update city dropdown based on selected country
            countrySelect.addEventListener("change", function () {
                citySelect.innerHTML = '<option value="">Select a city</option>'; // Reset cities
                const selectedCountry = countrySelect.value;
                if (selectedCountry && locationsMap[selectedCountry]) {
                    locationsMap[selectedCountry].forEach(city => {
                        let option = document.createElement("option");
                        option.value = city;
                        option.textContent = city.replace("_", " ").toUpperCase();
                        citySelect.appendChild(option);
                    });
                }
            });
        })
        .catch(error => console.error("Error fetching location data:", error));


      



  // Function to synchronize input fields and sliders for the "from" side
  function controlFromInput(fromSlider, fromInput, toInput, controlSlider) {
      const [from, to] = getParsed(fromInput, toInput);
      fillSlider(fromInput, toInput, '#C6C6C6', '#25daa5', controlSlider);
      if (from > to) {
          fromSlider.value = to;
          fromInput.value = to;
      } else {
          fromSlider.value = from;
      }
  }

  // Function to synchronize input fields and sliders for the "to" side
  function controlToInput(toSlider, fromInput, toInput, controlSlider) {
      const [from, to] = getParsed(fromInput, toInput);
      fillSlider(fromInput, toInput, '#C6C6C6', '#25daa5', controlSlider);
      setToggleAccessible(toInput);
      if (from <= to) {
          toSlider.value = to;
          toInput.value = to;
      } else {
          toInput.value = from;
      }
  }

  // Function to synchronize sliders when the "from" slider is changed
  function controlFromSlider(fromSlider, toSlider, fromInput) {
    const [from, to] = getParsed(fromSlider, toSlider);
    fillSlider(fromSlider, toSlider, '#C6C6C6', '#25daa5', toSlider);
    if (from > to) {
      fromSlider.value = to;
      fromInput.value = to;
    } else {
      fromInput.value = from;
    }
  }

  // Function to synchronize sliders when the "to" slider is changed
  function controlToSlider(fromSlider, toSlider, toInput) {
    const [from, to] = getParsed(fromSlider, toSlider);
    fillSlider(fromSlider, toSlider, '#C6C6C6', '#25daa5', toSlider);
    setToggleAccessible(toSlider);
    if (from <= to) {
      toSlider.value = to;
      toInput.value = to;
    } else {
      toInput.value = from;
      toSlider.value = from;
    }
  }

  // Function to parse the values from input fields
  function getParsed(currentFrom, currentTo) {
    const from = parseInt(currentFrom.value, 10);
    const to = parseInt(currentTo.value, 10);
    return [from, to];
  }

  // Function to fill the slider with color gradient based on input values
  function fillSlider(from, to, sliderColor, rangeColor, controlSlider) {
      const rangeDistance = to.max - to.min;
      const fromPosition = from.value - to.min;
      const toPosition = to.value - to.min;
      controlSlider.style.background = `linear-gradient(
        to right,
        ${sliderColor} 0%,
        ${sliderColor} ${(fromPosition) / (rangeDistance) * 100}%,
        ${rangeColor} ${((fromPosition) / (rangeDistance)) * 100}%,
        ${rangeColor} ${(toPosition) / (rangeDistance) * 100}%, 
        ${sliderColor} ${(toPosition) / (rangeDistance) * 100}%, 
        ${sliderColor} 100%)`;
  }

  // Function to toggle the accessibility (Z-index) of the slider based on input value
  function setToggleAccessible(currentTarget) {
    const toSlider = document.querySelector('#toSlider');
    if (Number(currentTarget.value) <= 0) {
      toSlider.style.zIndex = 2;
    } else {
      toSlider.style.zIndex = 0;
    }
  }

  // Function to parse the DD-MM-YYYY string into a valid year format
  function parseAlbumDate(dateString) {
    const [day, month, year] = dateString.split('-');
    return parseInt(year, 10); // Return the year as a number
  }

  // Initialize the creation year slider behavior
  fillSlider(fromSlider, toSlider, '#C6C6C6', '#25daa5', toSlider);
  setToggleAccessible(toSlider);

  // Initialize the first album year slider behavior
  fillSlider(fromAlbumSlider, toAlbumSlider, '#C6C6C6', '#25daa5', toAlbumSlider);
  setToggleAccessible(toAlbumSlider);

  // Event listeners for the creation year range
  fromSlider.oninput = () => controlFromSlider(fromSlider, toSlider, fromInput);
  toSlider.oninput = () => controlToSlider(fromSlider, toSlider, toInput);
  fromInput.oninput = () => controlFromInput(fromSlider, fromInput, toInput, toSlider);
  toInput.oninput = () => controlToInput(toSlider, fromInput, toInput, toSlider);

  // Event listeners for the first album year range
  fromAlbumSlider.oninput = () => controlFromSlider(fromAlbumSlider, toAlbumSlider, fromAlbumInput);
  toAlbumSlider.oninput = () => controlToSlider(fromAlbumSlider, toAlbumSlider, toAlbumInput);
  fromAlbumInput.oninput = () => controlFromInput(fromAlbumSlider, fromAlbumInput, toAlbumInput, toAlbumSlider);
  toAlbumInput.oninput = () => controlToInput(fromAlbumSlider, fromAlbumInput, toAlbumInput, toAlbumSlider);

  // Open the filter popup when the "Filter" button is clicked
  document.getElementById("openFilterBtn").addEventListener("click", () => {
      filterPopup.style.display = "flex";  // Show popup
  });

  // Close the filter popup when the "X" button is clicked
  document.getElementById("closeFilterBtn").addEventListener("click", () => {
      filterPopup.style.display = "none";  // Hide popup
  });


  

     // 🎯 Clear Filter Logic
  clearFilterBtn.addEventListener("click", function () {
    console.log("🗑 Clearing filters...");

    // Reset Sliders & Inputs
    fromInput.value = fromSlider.min;
    toInput.value = toSlider.max;
    fromAlbumInput.value = fromAlbumSlider.min;
    toAlbumInput.value = toAlbumSlider.max;

    fromSlider.value = fromSlider.min;
    toSlider.value = toSlider.max;
    fromAlbumSlider.value = fromAlbumSlider.min;
    toAlbumSlider.value = toAlbumSlider.max;

    // Reset Checkboxes
    document.querySelectorAll(".my-form input[type='checkbox']").forEach(checkbox => {
      checkbox.checked = false;
    });

    // Reset Dropdowns
    countrySelect.value = "";
    citySelect.innerHTML = '<option value="">Select a city</option>'; // Reset city dropdown

    // Show all artist cards
    document.querySelectorAll(".artist-card").forEach(card => {
      card.style.display = "block";
    });
    
    fillSlider(fromSlider, toSlider, '#C6C6C6', '#25daa5', toSlider);
    fillSlider(fromAlbumSlider, toAlbumSlider, '#C6C6C6', '#25daa5', toAlbumSlider);


    console.log("✅ Filters cleared");
  });

  // Apply the filter logic for filters
  applyFilterBtn.addEventListener("click", function () {

    console.log("✅ Apply button clicked!"); //debug

      const minYear = parseInt(fromInput.value, 10); // Get min value from creation year input
      const maxYear = parseInt(toInput.value, 10); // Get max value from creation year input
      console.log(`Applying filter: Min Year = ${minYear}, Max Year = ${maxYear}`);

      const minAlbYear = parseInt(fromAlbumInput.value, 10); // Get min value from first album year input
      const maxAlbYear = parseInt(toAlbumInput.value, 10); // Get max value from first album year input
      console.log(`Applying filter: Min Album Year = ${minAlbYear}, Max Album Year = ${maxAlbYear}`);

         // Collect selected checkboxes
      const selectedMembers = [];
      document.querySelectorAll(".my-form input[type='checkbox']:checked").forEach(checkbox => {
          selectedMembers.push(parseInt(checkbox.nextElementSibling.textContent, 10));
      });

      const selectedCountry = countrySelect.value.toLowerCase();
      const selectedCity = citySelect.value.toLowerCase();

    console.log(`Filtering: Min Year = ${minYear}, Max Year = ${maxYear}, Min Album Year = ${minAlbYear}, Max Album Year = ${maxAlbYear}, Selected Members = ${selectedMembers}`);

      // Loop through all artist cards and apply the filter based on both years
      document.querySelectorAll(".artist-card").forEach(card => {
          const creationDate = parseInt(card.getAttribute("data-year"), 10); // Get artist creation date
          
          const firstAlbumDateStr = card.getAttribute("data-album-year"); // Get artist first album year (as string)
          const firstAlbumDate = parseAlbumDate(firstAlbumDateStr); // Parse the first album date to get the year
          
          const numMembers = parseInt(card.getAttribute("data-members"), 10);

          const artistLocationsRaw = card.getAttribute("data-locations") || "";
          const artistLocations = artistLocationsRaw.replace(/[\[\]]/g, "").split(/\s+/).map(loc => loc.toLowerCase());

          // Show or hide artist card based on both creation year and first album year
        
        // Check if the artist matches the filters
        const matchesYear = creationDate >= minYear && creationDate <= maxYear;
        
        const matchesAlbumYear = firstAlbumDate >= minAlbYear && firstAlbumDate <= maxAlbYear;
        
        const matchesMembers = selectedMembers.length === 0 || selectedMembers.includes(numMembers);

        const matchesCountry = selectedCountry === "" || artistLocations.some(loc => loc.endsWith(selectedCountry));
        const matchesCity = selectedCity === "" || artistLocations.includes(`${selectedCity}-${selectedCountry}`);


        if (matchesYear && matchesAlbumYear && matchesMembers && matchesCountry && matchesCity) {
            card.style.display = "block";
        } else {
            card.style.display = "none";
        }
      });

      // Close filter popup after applying filter
      filterPopup.style.display = "none";
  });
});
