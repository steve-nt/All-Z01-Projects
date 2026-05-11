// Select necessary elements
const inputBox = document.querySelector("input[type='text']");
const suggBox = document.querySelector("#autocomplete-suggestions");
let currentIndex = -1; // Tracks which suggestion is selected

// Fetch & show suggestions
inputBox.addEventListener("keyup", async (e) => {
    const query = e.target.value.trim();
    if (!query) {
        suggBox.style.display = "none";
        return;
    }

    try {
        const res = await fetch(`/suggestions?query=${query}`);
        
        if (!res.ok) {
            throw new Error(`HTTP error! Status: ${res.status}`);
        }

        const data = await res.json();

        if (!data || !Array.isArray(data.suggestions)) {
            throw new Error("Invalid response format");
        }

        const suggestions = data.suggestions;

        // If no matches, hide the box
        if (suggestions.length === 0) {
            suggBox.style.display = "none";
            return;
        }

        // Generate suggestions HTML
        suggBox.innerHTML = suggestions.map(s => 
            `<div class="suggestion-item" tabindex="0">${s}</div>`
        ).join('');

        suggBox.style.display = "block";
        currentIndex = -1; // Reset index when new suggestions appear

        // Click event for suggestions
        document.querySelectorAll(".suggestion-item").forEach((div) => {
            div.addEventListener("click", () => selectSuggestion(div.textContent));
        });

    } catch (error) {
        console.error("Error fetching suggestions:", error);
        suggBox.style.display = "none"; // Hide on error
    }
});

// Handle keyboard navigation
inputBox.addEventListener("keydown", (e) => {
    const suggestions = document.querySelectorAll(".suggestion-item");
    if (suggestions.length === 0) return;

    if (e.key === "ArrowDown") {
        e.preventDefault();
        currentIndex = (currentIndex + 1) % suggestions.length; // Cycle down
        updateHighlight(suggestions);
        suggestions[currentIndex].focus(); // Move focus to suggestion
    } else if (e.key === "ArrowUp") {
        e.preventDefault();
        if (currentIndex === 0) {
            inputBox.focus(); // Move focus back to input when at the top
        } else {
            currentIndex = (currentIndex - 1 + suggestions.length) % suggestions.length;
            updateHighlight(suggestions);
            suggestions[currentIndex].focus();
        }
    } 
});

// Handle keyboard events inside suggestions
suggBox.addEventListener("keydown", (e) => {
    const suggestions = document.querySelectorAll(".suggestion-item");
    if (suggestions.length === 0) return;

    if (e.key === "ArrowDown") {
        e.preventDefault();
        currentIndex = (currentIndex + 1) % suggestions.length;
        updateHighlight(suggestions);
        suggestions[currentIndex].focus();
    } else if (e.key === "ArrowUp") {
        e.preventDefault();
        if (currentIndex === 0) {
            inputBox.focus();
        } else {
            currentIndex = (currentIndex - 1 + suggestions.length) % suggestions.length;
            updateHighlight(suggestions);
            suggestions[currentIndex].focus();
        }
    } else if (e.key === "Enter") {
        e.preventDefault();
        selectSuggestion(suggestions[currentIndex].textContent);
    }
});

// Hide suggestions when clicking outside
document.addEventListener("click", (e) => {
    if (!suggBox.contains(e.target) && e.target !== inputBox) suggBox.style.display = "none";
});

// Function to select a suggestion and submit the form
function selectSuggestion(value) {
    inputBox.value = value;
    suggBox.style.display = "none";
    inputBox.form.submit();
}

// Function to highlight selected suggestion
function updateHighlight(suggestions) {
    suggestions.forEach((s, i) => {
        s.classList.toggle("highlighted", i === currentIndex);
    });
}
