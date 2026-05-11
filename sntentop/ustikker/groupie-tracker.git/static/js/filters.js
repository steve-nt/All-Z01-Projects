document.addEventListener("DOMContentLoaded", () => {
    // Setup for date ranges
    const creationDates = Array.from(document.querySelectorAll('.artist-card .artist-info p:nth-child(1) strong')).map(el => parseInt(el.textContent));
    const firstAlbumDates = Array.from(document.querySelectorAll('.artist-card .artist-info p:nth-child(2) strong')).map(el => parseInt(el.textContent.split('-')[2]));

    const minCreationDate = Math.min(...creationDates);
    const maxCreationDate = Math.max(...creationDates);
    const minFirstAlbumDate = Math.min(...firstAlbumDates);
    const maxFirstAlbumDate = Math.max(...firstAlbumDates);

    setupRangeSlider('creationDate', minCreationDate, maxCreationDate);
    setupRangeSlider('firstAlbumDate', minFirstAlbumDate, maxFirstAlbumDate);

    const filterToggleBtn = document.getElementById('filterToggleBtn');
    const filterOptions = document.getElementById('filterOptions');
    const clearFiltersBtn = document.getElementById('clearFilters');
    const artistCards = document.querySelectorAll('.artist-card');

    // Toggle filter options visibility
    filterToggleBtn.addEventListener('click', function () {
        filterOptions.style.display = filterOptions.style.display === 'none' ? 'block' : 'none';
    });

    // Clear all filters
    clearFiltersBtn.addEventListener('click', function () {
        resetFilters();
        applyFilters();
    });

    // Populate member checkboxes
    const memberCounts = [...new Set(Array.from(artistCards).map(card => parseInt(card.querySelector('.artist-info p:nth-child(3) strong').textContent)))].sort((a, b) => a - b);
    const memberCheckboxes = document.getElementById('memberCheckboxes');
    memberCounts.forEach((count, index) => {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.id = `member-${count}`;
        checkbox.value = count;
        checkbox.addEventListener('change', applyFilters);

        const label = document.createElement('label');
        label.htmlFor = `member-${count}`;
        label.textContent = count;

        if (index === memberCounts.length - 1) {
            label.style.margin = '0';
        }

        memberCheckboxes.appendChild(checkbox);
        memberCheckboxes.appendChild(label);
    });

    // Populate country checkboxes
    const countries = getUniqueCountries(artistCards);
    const countryCheckboxes = document.getElementById('countryCheckboxes');
    countries.forEach(country => {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.id = `country-${country}`;
        checkbox.value = country;
        checkbox.addEventListener('change', applyFilters);

        const label = document.createElement('label');
        label.htmlFor = `country-${country}`;
        label.textContent = country;

        const div = document.createElement('div');
        div.appendChild(checkbox);
        div.appendChild(label);

        countryCheckboxes.appendChild(div);
    });

    // Setup event listeners for range sliders
    document.getElementById('creationDateMin').addEventListener('input', applyFilters);
    document.getElementById('creationDateMax').addEventListener('input', applyFilters);
    document.getElementById('firstAlbumDateMin').addEventListener('input', applyFilters);
    document.getElementById('firstAlbumDateMax').addEventListener('input', applyFilters);

    // Close filter options when clicking outside
    document.addEventListener('click', function (event) {
        if (!filterOptions.contains(event.target) && event.target !== filterToggleBtn) {
            filterOptions.style.display = 'none';
        }
    });

    // Prevent propagation for filter options and toggle button
    filterOptions.addEventListener('click', function (event) {
        event.stopPropagation();
    });

    filterToggleBtn.addEventListener('click', function (event) {
        event.stopPropagation();
    });

    // Setup dropdown for countries
    const dropdownToggle = document.getElementById('countriesDropdownToggle');
    dropdownToggle.addEventListener('click', function () {
        countryCheckboxes.style.display = countryCheckboxes.style.display === 'none' || countryCheckboxes.style.display === '' ? 'block' : 'none';
    });

    // Initial filter application
    applyFilters();
});

function setupRangeSlider(id, min, max) {
    const minSlider = document.getElementById(`${id}Min`);
    const maxSlider = document.getElementById(`${id}Max`);
    const minValue = document.getElementById(`${id}MinValue`);
    const maxValue = document.getElementById(`${id}MaxValue`);

    minSlider.min = min;
    minSlider.max = max;
    maxSlider.min = min;
    maxSlider.max = max;
    minSlider.value = min;
    maxSlider.value = max;
    minValue.textContent = min;
    maxValue.textContent = max;

    minSlider.addEventListener('input', function () {
        const minVal = parseInt(minSlider.value);
        const maxVal = parseInt(maxSlider.value);

        if (minVal > maxVal) {
            minSlider.value = maxVal;
        }
        minValue.textContent = minSlider.value;
        applyFilters();
    });

    maxSlider.addEventListener('input', function () {
        const minVal = parseInt(minSlider.value);
        const maxVal = parseInt(maxSlider.value);

        if (maxVal < minVal) {
            maxSlider.value = minVal;
        }
        maxValue.textContent = maxSlider.value;
        applyFilters();
    });
}

function getUniqueCountries(artists) {
    const countries = new Set();
    artists.forEach(artist => {
        const locations = artist.querySelector('.card-locations ul').querySelectorAll('li');
        locations.forEach(location => {
            const locationText = location.textContent.trim();
            const parts = locationText.split(', ');
            if (parts.length > 1) {
                const country = parts[parts.length - 1].trim();
                countries.add(country);
            }
        });
    });
    return Array.from(countries).sort();
}

function resetFilters() {
    document.getElementById('creationDateMin').value = document.getElementById('creationDateMin').min;
    document.getElementById('creationDateMax').value = document.getElementById('creationDateMax').max;
    document.getElementById('firstAlbumDateMin').value = document.getElementById('firstAlbumDateMin').min;
    document.getElementById('firstAlbumDateMax').value = document.getElementById('firstAlbumDateMax').max;

    document.getElementById('creationDateMinValue').textContent = document.getElementById('creationDateMin').min;
    document.getElementById('creationDateMaxValue').textContent = document.getElementById('creationDateMax').max;
    document.getElementById('firstAlbumDateMinValue').textContent = document.getElementById('firstAlbumDateMin').min;
    document.getElementById('firstAlbumDateMaxValue').textContent = document.getElementById('firstAlbumDateMax').max;

    document.querySelectorAll('#memberCheckboxes input[type="checkbox"]').forEach(cb => cb.checked = false);
    document.querySelectorAll('#countryCheckboxes input[type="checkbox"]').forEach(cb => cb.checked = false);
}

function applyFilters() {
    const creationDateMin = parseInt(document.getElementById('creationDateMin').value);
    const creationDateMax = parseInt(document.getElementById('creationDateMax').value);
    const firstAlbumDateMin = parseInt(document.getElementById('firstAlbumDateMin').value);
    const firstAlbumDateMax = parseInt(document.getElementById('firstAlbumDateMax').value);
    const selectedMembers = Array.from(document.querySelectorAll('#memberCheckboxes input[type="checkbox"]:checked')).map(cb => parseInt(cb.value));
    const selectedCountries = Array.from(document.querySelectorAll('#countryCheckboxes input[type="checkbox"]:checked')).map(cb => cb.value.toLowerCase());

    const artistCards = document.querySelectorAll('.artist-card');

    artistCards.forEach(card => {
        const creationDate = parseInt(card.querySelector('.artist-info p:nth-child(1) strong').textContent);
        const firstAlbumDateFull = card.querySelector('.artist-info p:nth-child(2) strong').textContent;
        const firstAlbumDate = parseInt(firstAlbumDateFull.split('-')[2]);
        const membersCount = parseInt(card.querySelector('.artist-info p:nth-child(3) strong').textContent);

        const artistCountries = Array.from(card.querySelectorAll('.card-locations ul li'))
            .map(li => li.textContent.split(', ').pop().trim().toLowerCase());

        const matchesCreationDate = creationDate >= creationDateMin && creationDate <= creationDateMax;
        const matchesFirstAlbumDate = firstAlbumDate >= firstAlbumDateMin && firstAlbumDate <= firstAlbumDateMax;
        const matchesMembers = selectedMembers.length === 0 || selectedMembers.includes(membersCount);
        const matchesCountries = selectedCountries.length === 0 || artistCountries.some(country => selectedCountries.includes(country));

        if (matchesCreationDate && matchesFirstAlbumDate && matchesMembers && matchesCountries) {
            card.style.display = 'block';
        } else {
            card.style.display = 'none';
        }
    });
}