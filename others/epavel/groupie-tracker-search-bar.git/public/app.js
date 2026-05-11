
var globalData;

const params = new URLSearchParams(window.location.search)

window.addEventListener("beforeunload", function () {
    sessionStorage.setItem("scrollPosition", window.scrollY);
});

window.addEventListener("load", async function () {
    const scrollPosition = sessionStorage.getItem("scrollPosition");
    if (scrollPosition && scrollPosition > 200) {
        window.scrollTo(0, parseInt(scrollPosition + 200, 10));
    }
});

function showSection(sectionId) {
    var sections = document.querySelectorAll('.data-section');
    sections.forEach(function (section) {
        section.classList.remove('active');
    });
    document.getElementById(sectionId).classList.add('active');
}

function showModal(id) {
    document.getElementById(id).classList.remove('active');
}
function closeModal(id) {
    document.getElementById(id).classList.add('active');
}

if (window.location.pathname === "/") {
    const clearButton = document.getElementById("clear-filter")
    if (params.size > 3) {
        clearButton.classList.remove('hidden')
        const paramsArray = ["creationDateMin", "creationDateMax", "albumDateMin", "albumDateMax", "locations", "cities", "members"]
        let paramObj = {}
        for (var value of params.keys()) {
            if (!paramsArray.includes(value)) {
                paramObj[value] = params.get(value);
            }
        }

        clearButton.addEventListener('click', function () {
            let queryString = new URLSearchParams(paramObj).toString();
            window.location.search = queryString
        })
    }

    document.getElementById("button-filter").addEventListener('click', function () {
        const container = document.getElementById("show-filters")
        if (container.classList.contains('hidden')) {
            container.classList.remove('hidden')
            this.textContent = "Hide Filters"
        } else {
            container.classList.add('hidden');
            this.textContent = "Show Filters"
        }

        const appendLocations = document.getElementById("locations");
        const appendCities = document.getElementById("cities");

        // Clear existing options
        appendLocations.innerHTML = "";
        appendCities.innerHTML = "";

        // Sort locations by country name
        globalData.Locations.sort((a, b) => a.country.localeCompare(b.country));

        function formatCountryName(country) {
            if (country.toUpperCase() === "USA" || country.toUpperCase() === "UK") {
                return country.toUpperCase();
            } else {
                return country.split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()).join(' ');
            }
        }

        globalData.Locations.forEach((location) => {
            let formattedCountry = formatCountryName(location.country);
            let option = document.createElement("option");
            option.text = formattedCountry;
            option.value = location.country;
            appendLocations.appendChild(option);
        });

        globalData.Locations.forEach((location) => {
            location.towns.sort((a, b) => a.town.localeCompare(b.town));

            location.towns.forEach((town) => {
                let formattedCountry = formatCountryName(location.country);
                let option = document.createElement("option");
                option.text = `${town.town} (${formattedCountry})`;
                option.value = town.town;
                appendCities.appendChild(option);
            });
        });

        console.log(globalData);

        if (params.has("filterMode") && params.get("filterMode") === "true") {
            params.set("filterMode", false);
        } else {
            params.set("filterMode", true);
        }

        window.history.replaceState({}, "", `${window.location.pathname}?${params}`);
    });
}

document.addEventListener("DOMContentLoaded", function () {
    const loadMoreButton = document.getElementById("load-more-button");
    const shuffleButton = document.getElementById("shuffle-button");

    function updateQueryParams() {
        const urlParams = new URLSearchParams(window.location.search);

        // Add filter parameters to the URL
        if (sessionStorage.getItem('creationDateMin')) {
            urlParams.set('creationDateMin', sessionStorage.getItem('creationDateMin'));
        }
        if (sessionStorage.getItem('creationDateMax')) {
            urlParams.set('creationDateMax', sessionStorage.getItem('creationDateMax'));
        }
        if (sessionStorage.getItem('albumDateMin')) {
            urlParams.set('albumDateMin', sessionStorage.getItem('albumDateMin'));
        }
        if (sessionStorage.getItem('albumDateMax')) {
            urlParams.set('albumDateMax', sessionStorage.getItem('albumDateMax'));
        }
        if (sessionStorage.getItem('members')) {
            urlParams.set('members', sessionStorage.getItem('members'));
        }
        if (sessionStorage.getItem('locations')) {
            urlParams.set('locations', sessionStorage.getItem('locations'));
        }
        if (sessionStorage.getItem('cities')) {
            urlParams.set('cities', sessionStorage.getItem('cities'));
        }

        return urlParams.toString();
    }
    if (totalArtists >= 12) {
        loadMoreButton.addEventListener("click", function (event) {
            event.preventDefault();
            const urlParams = updateQueryParams();
            const currentPagination = parseInt(params.get("pagination")) || 12;
            const nextPagination = currentPagination + 12;
            window.location.href = `/?pagination=${nextPagination}&${urlParams}`;
        });
    }

    shuffleButton.addEventListener("click", function (event) {
        event.preventDefault();
        const urlParams = new URLSearchParams(window.location.search);
        urlParams.delete("shuffle");
        const shuffle = params.get("shuffle") === "true" ? "false" : "true";
        urlParams.set("shuffle", shuffle);
        window.location.href = `/?${urlParams.toString()}`;
    });
});

if (window.location.pathname === "/") {
    document.addEventListener("DOMContentLoaded", async function () {
        sessionStorage.removeItem('creationDateMin');
        sessionStorage.removeItem('creationDateMax');
        sessionStorage.removeItem('albumDateMin');
        sessionStorage.removeItem('albumDateMax');
        sessionStorage.removeItem('members');
        sessionStorage.removeItem('locations');
        sessionStorage.removeItem('cities');
    
        const combobox = document.getElementById("combobox");
        const options = document.getElementById("options");
        const searchButton = document.getElementById("searchButton");
    
        try {
            const response = await fetch("/filter/data");
            if (!response.ok) {
                throw new Error(`Response status: ${response.status}`);
            }
            globalData = await response.json();
        } catch (error) {
            console.error("Error fetching data:", error);
            return;
        }
    
        // Retrieve stored filter values from sessionStorage
        if (sessionStorage.getItem('creationDateMin')) {
            minSlider.value = sessionStorage.getItem('creationDateMin');
        }
        if (sessionStorage.getItem('creationDateMax')) {
            maxSlider.value = sessionStorage.getItem('creationDateMax');
        }
        if (sessionStorage.getItem('albumDateMin')) {
            minSlider2.value = sessionStorage.getItem('albumDateMin');
        }
        if (sessionStorage.getItem('albumDateMax')) {
            maxSlider2.value = sessionStorage.getItem('albumDateMax');
        }
        if (sessionStorage.getItem('members')) {
            activateCheckboxesById(sessionStorage.getItem('members'));
        }
        if (sessionStorage.getItem('locations')) {
            const selectedLocation = sessionStorage.getItem('locations');
            const appendLocations = document.getElementById("locations");
            appendLocations.innerHTML = "";
            globalData.Locations.sort((a, b) => a.country.localeCompare(b.country));
            globalData.Locations.forEach((location) => {
                let option = document.createElement("option");
                option.text = location.country;
                option.value = location.country;
                if (selectedLocation.split('+').includes(location.country)) {
                    option.selected = true;
                }
                appendLocations.appendChild(option);
            });
        }
        if (sessionStorage.getItem('cities')) {
            const selectedCity = sessionStorage.getItem('cities');
            const appendCities = document.getElementById("cities");
            appendCities.innerHTML = "";
            globalData.Locations.forEach((location) => {
                location.towns.sort((a, b) => a.town.localeCompare(b.town));
                location.towns.forEach((town) => {
                    let option = document.createElement("option");
                    option.text = `${town.town} (${location.country})`;
                    option.value = town.town;
                    if (selectedCity.split('+').includes(town.town)) {
                        option.selected = true;
                    }
                    appendCities.appendChild(option);
                });
            });
        }
    
        if (params.has("filterMode") && params.get("filterMode") === "true") {
            const container = document.getElementById("show-filters");
            const button = document.getElementById("button-filter");
            container.classList.remove('hidden');
            button.textContent = "Hide Filters";
    
            const appendLocations = document.getElementById("locations");
            const appendCities = document.getElementById("cities");
            const selectedLocation = params.get("locations");
            const selectedCity = params.get("cities");
    
            // Clear existing options
            appendLocations.innerHTML = "";
            appendCities.innerHTML = "";
    
            // Sort locations by country name
            globalData.Locations.sort((a, b) => a.country.localeCompare(b.country));
    
            globalData.Locations.forEach((location) => {
                let option = document.createElement("option");
                option.text = location.country;
                option.value = location.country;
    
                if (selectedLocation === location.country) {
                    option.selected = true;
                }
    
                appendLocations.appendChild(option);
            });
            globalData.Locations.forEach((location) => {
                // Sort towns by town name
                location.towns.sort((a, b) => a.town.localeCompare(b.town));
    
                location.towns.forEach((town) => {
                    let option = document.createElement("option");
                    option.text = `${town.town} (${location.country})`;
                    option.value = town.town;
                    if (selectedCity === town.town) {
                        option.selected = true;
                    }
                    appendCities.appendChild(option);
                });
            });
        }
    
        const searchableData = [];
    
        globalData.Bands.forEach(band => {
            searchableData.push({
                type: "artist",
                value: band.name.toLowerCase(),
                display: `${band.name} (Artist)`
            });
    
            band.members.forEach(member => {
                searchableData.push({
                    type: "member",
                    value: member.toLowerCase(),
                    display: `${member} (Member)`
                });
            });
    
            searchableData.push({
                type: "creationDate",
                value: band.creationDate.toString(),
                display: `${band.creationDate} (Creation Date)`
            });
    
            searchableData.push({
                type: "firstAlbum",
                value: band.firstAlbum.toLowerCase(),
                display: `${band.firstAlbum} (First Album Date)`
            });
        });
    
        globalData.Locations.forEach(location => {
            searchableData.push({
                type: "location",
                value: location.country.toLowerCase(),
                display: `${location.country} (Location)`
            });
    
            location.towns.forEach(town => {
                searchableData.push({
                    type: "town",
                    value: `${town.town.toLowerCase()} (${location.country.toLowerCase()})`,
                    display: `${town.town} (${location.country})`
                });
            });
        });
    
        combobox.addEventListener("input", function () {
            const query = combobox.value.trim().toLowerCase();
            options.innerHTML = "";
    
            if (query.length > 2) {
                const filteredResults = searchableData.filter(item => item.value.includes(query) || item.value.includes(query.split(' ')[0]));
    
                filteredResults.forEach((result, index) => {
                    const option = document.createElement("li");
                    option.classList.add("relative", "cursor-pointer", "select-none", "py-2", "pl-3", "pr-9", "text-gray-900", "hover:bg-gray-400");
                    option.id = `option-${index}`;
                    option.role = "option";
                    option.tabIndex = "-1";
    
                    const text = document.createElement("span");
                    text.classList.add("block", "truncate");
                    text.textContent = result.display;
    
                    option.appendChild(text);
                    options.appendChild(option);
    
                    // Handle selection
                    option.addEventListener("click", () => {
                        combobox.value = result.display;
                        options.classList.add("hidden");
                        redirectToSearch(result.value);
                    });
                });
    
                if (filteredResults.length != 0) {
                    options.classList.remove("hidden");
                } else {
                    options.classList.add("hidden");
                }
            } else {
                options.classList.add("hidden");
            }
        });
    
        combobox.addEventListener("keypress", function (event) {
            if (event.key === "Enter") {
                const query = combobox.value.trim();
                if (query) {
                    redirectToSearch(query, "manual");
                }
            }
        });
    
        searchButton.addEventListener("click", function () {
            const query = combobox.value.trim();
            if (query) {
                redirectToSearch(query, "manual");
            }
        });
    
        function redirectToSearch(query) {
            const formattedQuery = query.replace(/ /g, "_");
            const queryParam = encodeURIComponent(formattedQuery);
            window.location.href = `/?query=${queryParam}`;
        }
    
        document.addEventListener("click", (e) => {
            if (!e.target.closest("#combobox") && !e.target.closest("#options")) {
                options.classList.add("hidden");
            }
        });
    });

    const minSlider = document.getElementById('min-slider');
    const maxSlider = document.getElementById('max-slider');

    if (params.has("creationDateMin") && params.has("creationDateMax")) {
        minSlider.value = params.get("creationDateMin")
        maxSlider.value = params.get("creationDateMax")
    }

    const rangeFill = document.getElementById('range-fill');
    const minValueDisplay = document.getElementById('min-value');
    const maxValueDisplay = document.getElementById('max-value');

    const minSlider2 = document.getElementById('min-slider-2');
    const maxSlider2 = document.getElementById('max-slider-2');

    if (params.has("albumDateMin") && params.has("albumDateMax")) {
        minSlider2.value = params.get("albumDateMin")
        maxSlider2.value = params.get("albumDateMax")
    }
    const rangeFill2 = document.getElementById('range-fill-2');
    const minValueDisplay2 = document.getElementById('min-value-2');
    const maxValueDisplay2 = document.getElementById('max-value-2');

    const applyFiltersButton = document.getElementById('applyFilters');

    function updateRangeFill(slider1, slider2, fillElement, minDisplay, maxDisplay) {
        const min = parseInt(slider1.value);
        const max = parseInt(slider2.value);

        if (min > max) {
            slider1.value = max;
        }

        const minPosition = (slider1.value - slider1.min) / (slider1.max - slider1.min) * 100;
        const maxPosition = (slider2.value - slider2.min) / (slider2.max - slider2.min) * 100;

        fillElement.style.left = minPosition + '%';
        fillElement.style.width = (maxPosition - minPosition) + '%';

        minDisplay.textContent = slider1.value;
        maxDisplay.textContent = slider2.value;
        if (slider1.valueAsNumber > slider2.valueAsNumber - 1) {
            slider1.style.zIndex = "10";
            slider2.style.zIndex = "20";
        } else {
            slider1.style.zIndex = "20";
            slider2.style.zIndex = "10";
        }
    }

    function collectFilters() {
        // Store current filter values in sessionStorage
        sessionStorage.setItem('creationDateMin', minSlider.value);
        sessionStorage.setItem('creationDateMax', maxSlider.value);
        sessionStorage.setItem('albumDateMin', minSlider2.value);
        sessionStorage.setItem('albumDateMax', maxSlider2.value);

        const memberCheckboxes = document.querySelectorAll('input[name="members"]:checked');
        const members = Array.from(memberCheckboxes).map(checkbox => checkbox.value);
        sessionStorage.setItem('members', members.join('+'));

        const locationSelect = document.getElementById('locations');
        const selectedLocations = Array.from(locationSelect.selectedOptions).map(option => option.value);
        sessionStorage.setItem('locations', selectedLocations.join('+'));

        const cities = document.getElementById('cities');
        const selectedCities = Array.from(cities.selectedOptions).map(option => option.value);
        sessionStorage.setItem('cities', selectedCities.join('+'));

        // Set the filters in the URL parameters
        params.set('creationDateMin', minSlider.value);
        params.set('creationDateMax', maxSlider.value);
        params.set('albumDateMin', minSlider2.value);
        params.set('albumDateMax', maxSlider2.value);
        if (members.length > 0) {
            params.set('members', members.join('+'));
        } else {
            params.delete('members');
        }
        if (selectedLocations.length > 0) {
            params.set('locations', selectedLocations.join('+'));
        } else {
            params.delete('locations');
        }
        if (selectedCities.length > 0) {
            const selectedCountries = new Set(params.get('locations') ? params.get('locations').split('+') : []);
            selectedCities.forEach(selectedCity => {
                globalData.Locations.forEach(location => {
                    location.towns.forEach(town => {
                        if (town.town === selectedCity) {
                            selectedCountries.add(location.country);
                        }
                    });
                });
            });
            params.set('locations', Array.from(selectedCountries).join('+'));
            params.set('cities', selectedCities.join('+'));
        } else {
            params.delete('cities');
        }
        if (params.has("pagination")) {
            params.set("pagination", currentPagination);
        }
        if (params.has("shuffle")) {
            params.set("shuffle", shuffle);
        }
        if (params.has("filterMode")) {
            params.set("filterMode", params.get("filterMode"));
        }
        window.location.search = params;
    }

    minSlider.addEventListener('input', () => updateRangeFill(minSlider, maxSlider, rangeFill, minValueDisplay, maxValueDisplay));
    maxSlider.addEventListener('input', () => updateRangeFill(minSlider, maxSlider, rangeFill, minValueDisplay, maxValueDisplay));
    updateRangeFill(minSlider, maxSlider, rangeFill, minValueDisplay, maxValueDisplay);

    minSlider2.addEventListener('input', () => updateRangeFill(minSlider2, maxSlider2, rangeFill2, minValueDisplay2, maxValueDisplay2));
    maxSlider2.addEventListener('input', () => updateRangeFill(minSlider2, maxSlider2, rangeFill2, minValueDisplay2, maxValueDisplay2));
    updateRangeFill(minSlider2, maxSlider2, rangeFill2, minValueDisplay2, maxValueDisplay2);

    applyFiltersButton.addEventListener('click', collectFilters);

    var queryMembers;

    function activateCheckboxesById(params) {
        const membersParam = params

        if (membersParam) {
            const ranges = membersParam.split('+');
            for (let i = 0; i <= 7; i++) {
                const checkbox = document.getElementById(`member-${i}`);
                var num = i
                if (checkbox && ranges.includes(num.toString())) {
                    checkbox.checked = true;
                }
            }
        }
    }

    if (params.has('members')) {
        queryMembers = params.get("members")
        activateCheckboxesById(queryMembers);
    }
}