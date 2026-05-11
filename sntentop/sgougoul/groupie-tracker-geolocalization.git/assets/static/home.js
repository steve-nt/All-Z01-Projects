let allArtists = []; // Store artists globally
let allLocations = {}; // Store locations indexed by artist ID

// Fetch all artists
fetch('/artists')
    .then(response => response.json())
    .then(artists => {
        allArtists = artists;
        displayArtists(allArtists); // Display the default artist cards
    })
    .catch(error => console.error('Error fetching artists:', error));

// Fetch locations for an artist (only when searching)
function fetchLocations(artistId) {
    if (allLocations[artistId]) {
        return Promise.resolve(allLocations[artistId]); // Return cached locations if already fetched
    }
    return fetch(`api/locations?id=${artistId}`)
        .then(response => response.json())
        .then(data => {
            allLocations[artistId] = data.locations || []; // Store in cache
            return allLocations[artistId];
        })
        .catch(error => {
            console.error(`Error fetching locations for artist ${artistId}:`, error);
            return [];
        });
}

// Function to display artists dynamically
function displayArtists(artists) {
    const gallery = document.getElementById('gallery');
    gallery.innerHTML = ''; // Clear previous content

    artists.forEach(artist => {
        const card = document.createElement('div');
        card.className = 'col-lg-3 col-md-4 col-sm-6 artist-card';
        card.innerHTML = `
            <div class="card h-100">
                <img src="${artist.image}" alt="${artist.name}" class="card-img-top artist-image">
                <div class="card-body text-center">
                    <h5 class="card-title">${artist.name}</h5>
                    <p class="card-text">Creation Date: ${artist.creationDate}</p>
                    <p class="card-text">First Album: ${artist.firstAlbum}</p>
                    <p class="card-text">Members: ${artist.members.join(', ')}</p>
                    <a href="/locations.html?id=${artist.id}" class="btn locationBtn">View Locations</a>
                    <a href="/dates.html?id=${artist.id}" class="btn datesBtn">View Concert Dates</a>
                    <a href="/relations.html?id=${artist.id}" class="btn relationsBtn">View Relations</a>
                </div>
            </div>
        `;
        gallery.appendChild(card);
    });
}

// Function to extract the year from the first album date (in "DD-MM-YYYY" format)
function extractYearFromDate(dateStr) {
    const parts = dateStr.split("-");
    if (parts.length === 3) {
        return parseInt(parts[2]); // Extract the year part
    }
    return NaN; // Invalid format, return NaN
}

function initFilterListeners() {
    const applyBtn = document.getElementById("applyFilters");
    const clearBtn = document.getElementById("clearFilters")
    
    if (!applyBtn) return;

    const dialog   = document.getElementById("errorDialog");
    const list     = document.getElementById("errorList");
    const closeBtn = document.getElementById("closeError");

    closeBtn.addEventListener("click", () => dialog.close());

    clearBtn.addEventListener("click", ()=> {
        document.getElementById("creationFrom").value = "";
        document.getElementById("creationTo"  ).value = "";
        document.getElementById("albumFrom"   ).value = "";
        document.getElementById("albumTo"     ).value = "";


        document.querySelectorAll('input[name="members"]').forEach(cb => {
            cb.checked = false;
        });

        document.getElementById("locationsFilter").value = "";

    })

   
  

    applyBtn.addEventListener("click", async (event) => {
        console.log("Apply Filters button clicked:", event);
        
        const creationFrom = parseInt(document.getElementById("creationFrom").value) || 0;
        const creationTo = parseInt(document.getElementById("creationTo").value) || 9999;
    
        const albumFrom = parseInt(document.getElementById("albumFrom").value) || 0;
        const albumTo = parseInt(document.getElementById("albumTo").value) || 9999;

        const errors = [];
        if (creationFrom > creationTo) {
          errors.push("“Creation From” must be no larger than “Creation To.”");
        }
        if (albumFrom > albumTo) {
          errors.push("“Album Year From” must be no larger than “Album Year To.”");
        }
        if (errors.length) {
          list.innerHTML = errors.map(e => `<li>${e}</li>`).join("");
          dialog.showModal();
          return;   // stop before fetching/filtering
        }
    
        const selectedMembers = Array.from(document.querySelectorAll('input[name="members"]:checked'))
            .map(cb => parseInt(cb.value));
    
        const selectedLocation = document.getElementById("locationsFilter").value.toLowerCase();
        
        console.log('Filter values:');
        console.log('Creation Date Range:', creationFrom, 'to', creationTo);
        console.log('First Album Year Range:', albumFrom, 'to', albumTo);
        console.log('Selected Members:', selectedMembers);
        console.log('Selected Location:', selectedLocation);

        const filtered = await Promise.all(allArtists.map(async artist => {
            const locations = await fetchLocations(artist.id); // Always get locations
        
            const creationOK = artist.creationDate >= creationFrom && artist.creationDate <= creationTo;
        
            const albumYear = extractYearFromDate(artist.firstAlbum); // Extract the year correctly
            const albumOK = albumYear >= albumFrom && albumYear <= albumTo;
        
            const memberCount = artist.members.length;
            const membersOK = selectedMembers.length === 0 || selectedMembers.includes(memberCount) || (selectedMembers.includes(memberCount) && memberCount >=7);
        
            let locationOK = true;
            if (selectedLocation) {
                locationOK = locations.some(loc => loc.toLowerCase().includes(selectedLocation));
            }
        
            return (creationOK && albumOK && membersOK && locationOK) ? artist : null;
        }));
    
        displayArtists(filtered.filter(Boolean));
    });
}

// 1a) In-memory cache for all location strings
let allLocationOptions = [];

// 1b) Only fetch & cache—no DOM writes
async function loadLocationOptions() {
  try {
    const resp = await fetch("/api/locations-list");
    if (!resp.ok) throw new Error(resp.status);
    const { locations } = await resp.json();
    allLocationOptions = locations.sort((a, b) => a.localeCompare(b));
  } catch (err) {
    console.error("Couldn’t load location options:", err);
  }
}

// 1c) Only DOM writes—no network
function populateLocationDropdown() {
  const select = document.getElementById("locationsFilter");
  if (!select) return;   // not in DOM yet
  // build one HTML string for speed
  const opts = allLocationOptions
    .map(loc => `<option value="${loc}">${loc}</option>`)
    .join("");
  select.innerHTML = `<option value="">-- Select a location --</option>${opts}`;
}


// Collect the filter values from the form and send to the backend for filtering
function applyFilters() {
    const creationFrom = parseInt(document.getElementById("creationFrom").value) || 0;
    const creationTo = parseInt(document.getElementById("creationTo").value) || 9999;
    
    const albumFrom = parseInt(document.getElementById("albumFrom").value) || 0;
    const albumTo = parseInt(document.getElementById("albumTo").value) || 9999;
    
    const selectedMembers = Array.from(document.querySelectorAll('input[name="members"]:checked'))
        .map(cb => parseInt(cb.value));
    
    const selectedLocation = document.getElementById("locationsFilter").value.toLowerCase();
    
    console.log('Filter values:');
    console.log('Creation Date Range:', creationFrom, 'to', creationTo);
    console.log('First Album Year Range:', albumFrom, 'to', albumTo);
    console.log('Selected Members:', selectedMembers);
    console.log('Selected Location:', selectedLocation);

    // Send filter criteria to the backend
    fetch('/filterArtists', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            creationFrom,
            creationTo,
            albumFrom,
            albumTo,
            selectedMembers,
            selectedLocation
        })
    })
    .then(response => response.json())
    .then(filteredArtists => {
        displayArtists(filteredArtists); // Update the displayed artists after filtering
    })
    .catch(error => console.error('Error applying filters:', error));
}

document.addEventListener("DOMContentLoaded", () => {
  const preloadPromise = loadLocationOptions();


  const btn       = document.getElementById("showFiltersBtn");
  const container = document.getElementById("filtersContainer");
  let loaded      = false;

  btn.addEventListener("click", async () => {

    if (!loaded) {
      // ─── First click: load & inject the filters ───
      loaded = true;
      btn.disabled = true;
      await preloadPromise;

      try {
        const resp = await fetch("/static/filters.html");
        if (!resp.ok) throw new Error(resp.statusText);
        container.innerHTML = await resp.text();
      } catch (err) {
        console.error("Failed to load filters:", err);
        btn.disabled = false;
        return;
      }

      // ─── Now that #filtersContainer contains your form, wire it up ───
      populateLocationDropdown();
      initFilterListeners();
      
      const element = document.getElementById("locationsFilter");
      if (element && window.Choices) {
        new Choices(element, {
        searchEnabled: true,
        itemSelectText: '',
        shouldSort: false
        });
        }
      

      // show panel, update button
      container.classList.remove("d-none");
      btn.textContent = "Hide Filters";
      btn.disabled   = false;

    } else {
      // ─── Subsequent clicks: just toggle visibility ───
      container.classList.toggle("d-none");
      btn.textContent = container.classList.contains("d-none")
        ? "Show Filters"
        : "Hide Filters";
    }
  });
});


const searchInput = document.getElementById('searchBar');

searchInput.addEventListener('input', function () {
    const query = this.value;
    debouncedFetchSuggestions(query);
    filterAndDisplayArtists(query);
    
});

function debounce(func, delay) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => func.apply(this, args), delay);
    };
}

const debouncedFetchSuggestions = debounce(async function (query) {
    try {
        const response = await fetch("/suggestions?q=" + encodeURIComponent(query));
        const html = await response.text();
        const suggestionsList = document.getElementById("artistSuggestions");
        if (suggestionsList) {
            suggestionsList.innerHTML = html;
        }
    } catch (err) {
        console.error("Error fetching suggestions:", err);
    }
}, 300);

async function filterAndDisplayArtists(query) {
    if (!query) {
        displayArtists(allArtists);
        return;
    }

    // Use the parser to check for suggestion-type queries.
    const { text: searchText, type } = parseHyphenatedSuggestion(query);
    const lowerSearchText = searchText.toLowerCase();
    console.log("Parsed searchText:", searchText);
    console.log("Parsed type:", type);
    console.log("Lowercase searchText:", lowerSearchText);

    let filteredArtists = await Promise.all(
        allArtists.map(async artist => {
            const locations = await fetchLocations(artist.id);
            let matches = false;

            if (type === "member") {
                matches = artist.members.some(member =>
                    member.toLowerCase().includes(lowerSearchText)
                );
            } else if (type === "location") {
                matches = locations.some(location =>
                    location.toLowerCase().includes(lowerSearchText)
                );
            } else if (type === "firstalbum") {
                const albumDate = artist.firstAlbum.trim().toLowerCase();

                // Match full date
                if (albumDate === lowerSearchText) {
                    matches = true;
                } else {
                    // Match year only
                    const albumYear = extractYearFromDate(artist.firstAlbum);
                    matches = albumYear && albumYear.toString().includes(lowerSearchText);
                }
            } else if (type === "creationdate") {
                matches = artist.creationDate.toString().includes(lowerSearchText);
            } else if (type === "artist/band") {
                matches = artist.name.toLowerCase().includes(lowerSearchText);
            } else {
                // Generic filtering if no type is specified.
                const lowerQuery = query.toLowerCase();

                // Ensure partial matches are found even when typing the full word
                matches =
                    artist.name.toLowerCase().includes(lowerQuery) ||
                    artist.members.some(member => member.toLowerCase().includes(lowerQuery)) ||
                    artist.firstAlbum.toLowerCase().trim().includes(lowerQuery) ||
                    artist.creationDate.toString().includes(lowerQuery) ||
                    locations.some(location => location.toLowerCase().includes(lowerQuery));
            }

            return matches ? artist : null;
        })
    );

    // Filter out null values (artists that didn't match)
    filteredArtists = filteredArtists.filter(artist => artist !== null);

    // Display the filtered artists
    displayArtists(filteredArtists);
}

function parseHyphenatedSuggestion(query) {
    // Trim the query and check if it even contains a hyphen.
    const trimmed = query.trim();
    if (!trimmed.includes("-")) {
        return { text: trimmed, type: null };
    }
    
    // Find the last hyphen.
    const idx = trimmed.lastIndexOf("-");
    // Extract the potential type and the remaining text.
    const potentialType = trimmed.substring(idx + 1).toLowerCase().trim();
    console.log("Parsed type:", potentialType)
    const text = trimmed.substring(0, idx).trim();
    
    // Define allowed types.
    const allowedTypes = new Set(["member", "location", "artist/band","firstalbum", "creationdate"]);
 
    // Only if the potential type is valid, consider it as a type indicator.
    if (allowedTypes.has(potentialType)) {
        return { text, type: potentialType };
    }
    // Otherwise, treat the whole query as generic text.
    return { text: trimmed, type: null };
}
