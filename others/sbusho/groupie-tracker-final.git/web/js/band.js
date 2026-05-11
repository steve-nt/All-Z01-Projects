// Map functions
let map;
let markers = [];
let infoWindows = [];

// Initialize map and year display
document.getElementById("year").innerText = new Date().getFullYear();

// Modal functions
function openModal(modalId) {
    document.getElementById(modalId).style.display = 'block';
    document.getElementById('modalBg').style.display = 'block';

    if (modalId === 'relationsModal' && typeof google !== 'undefined') {
        setTimeout(() => {
            google.maps.event.trigger(map, "resize");
        }, 300);
    }
}

function closeModal() {
    document.querySelectorAll('.modal').forEach(modal => modal.style.display = 'none');
    document.getElementById('modalBg').style.display = 'none';
}

function initMap() {
    const defaultCenter = { lat: 20, lng: 0 };
    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 1,
        center: defaultCenter,
        styles: [
            { elementType: "geometry", stylers: [{ color: "#1a1a2e" }] },
            { elementType: "labels.text.stroke", stylers: [{ color: "#16213e" }] },
            { elementType: "labels.text.fill", stylers: [{ color: "#e94560" }] },
            { featureType: "administrative", elementType: "geometry", stylers: [{ visibility: "off" }] },
            { featureType: "poi", stylers: [{ visibility: "off" }] },
            { featureType: "road", stylers: [{ visibility: "off" }] },
            { featureType: "transit", stylers: [{ visibility: "off" }] },
            { featureType: "water", stylers: [{ color: "#0f3460" }] }
        ]
    });

    // Clear any existing markers
    clearMarkers();
    
    try {
        const locations = JSON.parse(document.getElementById('map').getAttribute('data-locations'));
        addMarkers(locations);
    } catch (error) {
        console.error("Error parsing locations:", error);
    }
}

function addMarkers(locations) {
    if (!locations || !Array.isArray(locations)) {
        console.error("Invalid locations data");
        return;
    }

    locations.forEach((location, index) => {
        fetch(`/geolocations?address=${encodeURIComponent(location)}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (!data.latitude || !data.longitude) {
                    throw new Error("Invalid geocoding data");
                }

                const position = { lat: data.latitude, lng: data.longitude };
                
                const marker = new google.maps.Marker({
                    position,
                    map,
                    title: location,
                    label: {
                        text: `${index + 1}`,
                        color: "#ffffff",
                        fontSize: "12px",
                        fontWeight: "bold"
                    }
                });

                const infoWindow = new google.maps.InfoWindow({
                    content: `<div style="color: black; padding: 5px;"><strong>${location}</strong></div>`
                });

                marker.addListener('click', () => {
                    infoWindows.forEach(iw => iw.close());
                    infoWindow.open(map, marker);
                });

                markers.push(marker);
                infoWindows.push(infoWindow);

                // Auto-center map to show all markers
                if (markers.length === locations.length) {
                    autoCenterMap();
                }
            })
            .catch(error => {
                console.error(`Failed to geocode location "${location}":`, error);
            });
    });
}

function autoCenterMap() {
    if (markers.length === 0) return;

    const bounds = new google.maps.LatLngBounds();
    markers.forEach(marker => {
        bounds.extend(marker.getPosition());
    });
    map.fitBounds(bounds);
}

function clearMarkers() {
    markers.forEach(marker => marker.setMap(null));
    markers = [];
    infoWindows.forEach(iw => iw.close());
    infoWindows = [];
}
