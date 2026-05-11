let map;
let markersLayer; // Store markers for removal/updating

function openPopup(artist, locations) {
    // Display the popup
    document.getElementById("popup").style.display = "block";

    // Hide the loading indicator and show the map once the locations are ready
    const loadingElement = document.getElementById("loading");
    const mapElement = document.getElementById("map");

    // Hide loading indicator, show map container
    loadingElement.style.display = "none";
    mapElement.style.display = "block";

    // If the map doesn't exist, initialize it
    if (!map) {
                // Initialize the map and set its view to the first location's coordinates
                map = L.map('map', {
                    center: [locations[0].lat, locations[0].lon],
                    zoom: 10,
                    worldCopyJump: true // Prevent map from wrapping around beyond its bounds
                });
      //  map = L.map('map').setView([locations[0].lat, locations[0].lon], 10);
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; OpenStreetMap contributors'
        }).addTo(map);
        
        // Initialize markersLayer
        markersLayer = L.layerGroup().addTo(map);
    } else {
        // Clear existing markers
        markersLayer.clearLayers();
    }

    // Create a bounds object to include all markers
    const bounds = new L.LatLngBounds();

    // Add markers for all locations
    locations.forEach(location => {
        const marker = L.marker([location.lat, location.lon]).addTo(markersLayer)
            .bindPopup(`
                <div style="text-align: center;">
                    <strong style="font-size: 16px;">${artist}</strong><br>
                    <strong>${location.name}</strong><br>
                    <strong>${location.date}</strong>
                </div>
            `);
        // Extend the bounds to include this marker
        bounds.extend(marker.getLatLng());
    });

    // Set the map view to fit all the markers
    map.fitBounds(bounds);
}

function closePopup() {
    // Hide the popup
    document.getElementById("popup").style.display = "none";
}

// Attach click event to all concert links
document.querySelectorAll('.concert-link').forEach(link => {
    link.addEventListener('click', function(event) {
        event.preventDefault();
        
        // Show the popup and the loading message
        const loadingElement = document.getElementById("loading");
        const mapElement = document.getElementById("map");

        document.getElementById("popup").style.display = "block";
        loadingElement.style.display = "block";  // Show loading indicator
        mapElement.style.display = "none";  // Hide the map initially

        let concert = this.getAttribute('data-location');
        fetch(`/geolocations?concert=${encodeURIComponent(concert)}`)
            .then(response => response.json())
            .then(data => {
                if (
                    data.artist &&
                    Array.isArray(data.dates) &&
                    Array.isArray(data.locations) &&
                    Array.isArray(data.latitudes) &&
                    Array.isArray(data.longitudes) &&
                    data.locations.length === data.latitudes.length &&
                    data.locations.length === data.longitudes.length &&
                    data.locations.length === data.dates.length
                ) {
                    // Collect location data
                    let locations = data.locations.map((location, index) => ({
                        lat: parseFloat(data.latitudes[index]),
                        lon: parseFloat(data.longitudes[index]),
                        name: location,
                        date: data.dates[index] // Include concert date
                    }));
                    
                    openPopup(data.artist, locations);
                } else {
                    console.error('Invalid data received:', data);
                    loadingElement.style.display = "none"; // Hide loading if data is invalid
                }
            })
            .catch(error => {
                console.error('Error:', error);
                loadingElement.style.display = "none";  // Hide loading on error
            });
    });
});
