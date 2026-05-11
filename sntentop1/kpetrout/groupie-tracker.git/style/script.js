let currentAudioElement = null; // Keep track of the currently playing audio element
let globalVolume = 1.0;

// VOLUME SLIDER
function showSlider() {
  const slider = document.getElementById("volume-slider");
  const output = document.getElementById("volume-value");

  // Toggle display using a helper function
  toggleDisplay(slider);
  toggleDisplay(output);

  const inputEvent = new Event('input');
  slider.dispatchEvent(inputEvent);
}

function toggleDisplay(element) {
  if (element.style.display === "none" || element.style.display === "") {
    element.style.display = "block";
  } else {
    element.style.display = "none";
  }
}

// Update audio volume when slider value changes
document.addEventListener('DOMContentLoaded', function () {
  const volumeSlider = document.getElementById("volume-slider");
  const volumeIcon = document.getElementById("sound-icon");
  const output = document.getElementById("volume-value");
  
  const savedVolume = localStorage.getItem('volume') || '100';
  
  volumeSlider.value = savedVolume;
  output.textContent = `${savedVolume}%`;
  
  function updateVolume(volume) {
    let iconSrc;

    if (volume === '0') {
      iconSrc = 'style/media/icons/volume-min.png';
    } else if (volume <= 45) {
      iconSrc = 'style/media/icons/volume-1.png';
    } else if (volume <= 85) {
      iconSrc = 'style/media/icons/volume-2.png';
    } else {
      iconSrc = 'style/media/icons/volume-max.png';
    }
    output.textContent = `${volume}%`;
    globalVolume = volume / 100;


    const sliderStyle = getComputedStyle(volumeSlider);
    const outputStyle = getComputedStyle(output);
    let sliderWasHidden = sliderStyle.display === 'none';
    let outputWasHidden = outputStyle.display === 'none';

    if (sliderWasHidden) volumeSlider.style.display = 'block';
    if (outputWasHidden) output.style.display = 'block';

    volumeSlider.offsetHeight;

    const slideRect = volumeSlider.getBoundingClientRect();
    const sliderWidth = slideRect.width;
    const thumbWidth = volumeSlider.offsetHeight;
    const position = ((volume - volumeSlider.min) / (volumeSlider.max - volumeSlider.min));
    const offset = position * (sliderWidth - thumbWidth) + (thumbWidth / 2);

    volumeIcon.src = iconSrc;
    output.style.left = `${offset}px`;

    if (sliderWasHidden) volumeSlider.style.display = 'none';
    if (outputWasHidden) output.style.display = 'none';

    localStorage.setItem('volume', volume);
  }

  updateVolume(savedVolume);

  volumeSlider.addEventListener('input', function() {
    const volume = this.value;
    updateVolume(volume);
  }
  );
});

// MUSIC
function playTrack(artistID) {
  console.log("Playing track for artist", artistID);
  const audioElement = document.getElementById(`audio-${artistID}`);
  if (audioElement) {
    currentAudioElement = audioElement; // Set the current audio element
    audioElement.volume = globalVolume; // Apply the global volume
    audioElement.play().catch(function (error) {
      console.error("Playback failed:", error);
    });
  } else {
    console.error("Audio element not found for artist", artistID);
  }
}

function stopTrack(artistID) {
  const overlay = document.getElementById(`artist-${artistID}`);
  if (overlay && overlay.classList.contains("show")) {
    return;
  }
  const audioElement = document.getElementById(`audio-${artistID}`);
  if (audioElement) {
    audioElement.pause();
    audioElement.currentTime = 0; // Reset the track to the beginning
    if (currentAudioElement === audioElement) {
      currentAudioElement = null; // Clear the current audio element
    }
  }
}

function togglePlayPause(artistId) {
  const audioElement = document.getElementById(`audio-${artistId}`);
  const pauseIcon = document.getElementById(`pause-icon-${artistId}`);
  const playIcon = document.getElementById(`play-icon-${artistId}`);

  if (audioElement.paused) {
    audioElement.play();
    pauseIcon.style.display = 'inline';
    playIcon.style.display = 'none';
  } else {
    audioElement.pause();
    pauseIcon.style.display = 'none';
    playIcon.style.display = 'inline';
  }
}

function toggleHeader() {
  const header = document.querySelector('.main-header');
  const btn = document.querySelector('.close-header');
  header.classList.toggle('hidden');

  if (header.classList.contains('hidden')) {
    btn.innerHTML = "&#x25BC;";
  } else {
    btn.innerHTML = "&#x25B2;";
  }
}

// Re-show the header when scrolled to the top.
window.addEventListener("scroll", function() {
  const closeBtn = document.querySelector('.close-header');

  if (window.pageYOffset > 100) {
    closeBtn.style.display = "block";
  } else if (window.pageYOffset === 0) {
    closeBtn.innerHTML = "&#x25B2;";
    closeBtn.style.display = "none";
    document.querySelector('.main-header').classList.remove('hidden');
  }
});

// Toggle map display
document.addEventListener("DOMContentLoaded", function () {
  const mapBtn = document.getElementById("map-btn");
  const map = document.getElementById("map");
  
  mapBtn.addEventListener("click", function () {
    toggleDisplay(map);
  });
});

// Artist Names
function bigNames() {
  const overlays = document.querySelectorAll(".overlay");
  
  overlays.forEach(overlay => {
    const heading = overlay.querySelector("h2");
    const btn = overlay.querySelector('.pause-btn');

    if (heading.textContent.length > 21) {
      btn.classList.add("big-name");
    }
  });
}

document.addEventListener("DOMContentLoaded", function () {
    // Select all artist cards
    const artistCards = document.getElementById("artists");
    bigNames();

    // Open overlay on click
    if (artistCards) {
      artistCards.addEventListener("click", function (event) {
        const card = event.target.closest(".card");
        if (card) {
          const artistId = card.getAttribute("data-artist-id");
          const overlay = document.getElementById(`artist-${artistId}`);
          if (overlay) {
            overlay.classList.add("show");
            const pauseIcon = document.getElementById(`pause-icon-${artistId}`);
            const playIcon = document.getElementById(`play-icon-${artistId}`);
            if (pauseIcon && playIcon) {
              pauseIcon.style.display = "inline";
              playIcon.style.display = "none";
            }
          }
        }
      });
    }

    // Close overlay on clicking close button or outside the content
    document.addEventListener("click", function (event) {
      if (
        event.target.classList.contains("close-button") || // Click close button
        event.target.classList.contains("overlay") // Click outside content
      ) {
        const overlay = event.target.closest(".overlay");
        if (overlay) {
          overlay.classList.remove("show");
          // Extract artistID from overlay's id (assuming id format "artist-{{.ID}}")
          const artistID = overlay.id.substring(7);
          stopTrack(artistID);
        }
      }
    });

    // button to show filter
    const minDateRange = document.getElementById("min-date-range");
    const minDateValue = document.getElementById("min-date-value");
  
    if (minDateRange && minDateValue) {
        minDateRange.addEventListener("input", function () {
            minDateValue.textContent = this.value;
        });
  
        // Set initial value on page load
        minDateValue.textContent = minDateRange.value;
    }
    
    const minAlbumDateRange = document.getElementById("min-album-date-range");
    const minAlbumDateValue = document.getElementById("min-album-date-value");
  
    if (minAlbumDateRange && minAlbumDateValue) {
        minAlbumDateRange.addEventListener("input", function () {
            minAlbumDateValue.textContent = this.value;
        });
  
        // Set initial value on page load
        minAlbumDateValue.textContent = minAlbumDateRange.value;
    }
});

// SEARCH BAR

document.addEventListener("DOMContentLoaded", () => {
  const searchForm = document.getElementById("search-form");
  const searchInput = document.getElementById("search-bar");
  const suggestionsDiv = document.getElementById("suggestions");
  const resultsContainer = document.getElementById("artists"); // Your div container for artists
  
function debounce(func, delay) {
  let timeout;
  return function (...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), delay);
  };
}
  async function performSearch(query) {
    try {
      const response = await fetch(`/homepage?query=${encodeURIComponent(query)}`, {
        headers: {
          'X-Requested-With': 'XMLHttpRequest'
        }
      });
      
      if (response.ok) {
        const htmlFragment = await response.text();
        resultsContainer.innerHTML = htmlFragment;
        suggestionsDiv.style.display = 'none';
        bigNames();
      }  
    } catch (error) {
      console.error("Error fetching search results: ", error);
    }
  }
 
  // SUGGESTIONS

  async function getSuggestions(query) {
    try {
      const response = await fetch(`/api/suggestions?q=${encodeURIComponent(query)}`);
      return await response.json();
    } catch (error) {
      console.error('Fetch error:', error);
      return [];
    }
  }

  function showSuggestions(items) {
    suggestionsDiv.innerHTML = ""; // Clear previous suggestions
    const query = searchInput.value.trim().toLowerCase();

    if (items.length === 0) {
      suggestionsDiv.style.display = 'none';
      return;
    }

    items.forEach(item => {
      const div = document.createElement("div");
      div.className = "suggestion-item";

      // Highlight matching text
      const displayText = item.Display;
      const matchStart = displayText.toLowerCase().indexOf(query);
      
      div.innerHTML = query.length >= 2 && matchStart >= 0
        ? highlightMatch(displayText, matchStart, query.length)
        : displayText;

      div.innerHTML += `<span class="entity type"> - ${item.Type}</span>`;
      div.onclick = () => {
        searchInput.value = item.Display;
        performSearch(item.Display);
      };
      suggestionsDiv.appendChild(div);
    });
    suggestionsDiv.style.display = 'block';
  }

  function highlightMatch(text, startIndex, length) {
    return [
      text.slice(0, startIndex),
      '<span class="match-highlight">',
      text.slice(startIndex, startIndex + length),
      '</span>',
      text.slice(startIndex + length)
    ].join('');
  }

  const handleInput = debounce(async () => {
    const query = searchInput.value.trim();
    const [results, suggestions] = await Promise.all([
      performSearch(query),
      getSuggestions(query)
    ]);
    if (query.length > 0) {
      showSuggestions(suggestions);
    } else {
      suggestionsDiv.style.display = 'none';
    }
  }
  , 300); // Adjust the debounce delay as needed

  if (resultsContainer) {
    searchForm.addEventListener("submit", (e) => {
      e.preventDefault();
      handleInput();
    });
  }

  searchInput.addEventListener("input", () => {
    handleInput();
    // Show loading state if needed
  });

  document.addEventListener("click", (e) => {
    if (!e.target.closest('.search-container')) {
      suggestionsDiv.style.display = 'none';
    }
  });

  searchInput.addEventListener("keypress", (e) => {
    const items = suggestionsDiv.querySelectorAll('.suggestion-item');
    if (!items.length) return;

    switch(e.key) {
      case 'ArrowDown':
        console.log('ArrowDown pressed');
        e.preventDefault();
        items[1].focus();
        break;
      case 'Escape':
        suggestionsDiv.style.display = 'none';
        break;
    }
  });
});

// FILTER
function showFilter() {
  const filter = document.getElementById('filter');
  filter.classList.toggle('show');
  if (filter.classList.contains('show')) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}
