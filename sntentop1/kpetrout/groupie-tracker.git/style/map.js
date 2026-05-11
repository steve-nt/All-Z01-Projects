window.initMap = function () {
    const artistArray = window.artistsData;
    if (!artistArray || artistArray.length === 0) {
        console.error("No artist data found");
        return;
    }

    const artist = artistArray[0];
    console.log(artist);
    const locations = artist.locations;

    if (!locations || locations.length === 0) {
        console.error("No locations found for artist");
        return;
    }
    console.log(artistArray);
    const map = new google.maps.Map(document.getElementById("map"), {
        zoom: 3,
        center: { lat: locations[0].lat, lng: locations[0].lng }, // Center on the first location
    });

    locations.forEach(loc => {
        if (loc.lat && loc.lng) {
            new google.maps.Marker({
                position: { lat: loc.lat, lng: loc.lng },
                map: map,
                title: loc.name
            });
        }
    });
};
