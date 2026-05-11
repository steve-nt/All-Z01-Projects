async function initMap() {
    const map = new google.maps.Map(document.getElementById("map"), {
        center: { lat: 0, lng: 0 },
        zoom: 2,
        mapId: 'DEMO_MAP_ID',
    });

    const geocoder = new google.maps.Geocoder();
    const concerts = concertData; // Use the concertData variable from artist.html
    const locations = [];

    // Parse the concerts data structure to get locations in chronological order
    for (const country in concerts.Countries) {
        for (const town in concerts.Countries[country].Towns) {
            const dates = concerts.Countries[country].Towns[town].Dates;
            for (const date of dates) {
                locations.push({ country, town, date });
            }
        }
    }

    // Sort locations by date
    locations.sort((a, b) => new Date(a.date) - new Date(b.date));

    const markers = [];
    const pathCoordinates = [];

    // Geocode each location and create markers with info windows
    for (let i = 0; i < locations.length; i++) {
        const location = locations[i];
        const address = `${location.town}, ${location.country}`;
        const results = await geocodeAddress(geocoder, address);
        if (results) {
            const position = results[0].geometry.location;
            pathCoordinates.push(position);

            const marker = new google.maps.Marker({
                map: map,
                position: position,
                title: `${location.town}, ${location.country} (${location.date})`,
                label: String.fromCharCode('A'.charCodeAt(0) + i), // Label as A, B, C, etc.
            });

            const infoWindow = new google.maps.InfoWindow({
                content: `<div><strong>${location.town}, ${location.country}</strong></div>`,
            });

            marker.addListener("click", () => {
                infoWindow.open({
                    anchor: marker,
                    map,
                });
            });

            markers.push(marker);
        }
    }

    // Draw lines between the markers
    const path = new google.maps.Polyline({
        path: pathCoordinates,
        geodesic: true,
        strokeColor: '#FF0000',
        strokeOpacity: 1.0,
        strokeWeight: 2,
    });

    path.setMap(map);
}

function geocodeAddress(geocoder, address) {
    return new Promise((resolve, reject) => {
        geocoder.geocode({ address: address }, (results, status) => {
            if (status === "OK") {
                resolve(results);
            } else {
                console.error("Geocode was not successful for the following reason: " + status);
                resolve(null);
            }
        });
    });
}