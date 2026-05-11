function initArtistMap() {
  const map = L.map('map').setView([20, 0], 2); // Center map globally

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
  }).addTo(map);

  return map;
}

function addMarker(map, name, lat, lon) {
  const marker = L.marker([lat, lon]).addTo(map);
  marker.bindPopup(`<b>${name}</b>`);
  return [lat, lon];
}

function showLoadingIndicator() {
  const mapContainer = document.getElementById('map');
  const loader = document.createElement('div');
  loader.id = 'map-loader';
  loader.textContent = 'Loading locations...';
  loader.style.cssText = `
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: rgba(0,0,0,0.7);
    color: white;
    padding: 15px 25px;
    border-radius: 10px;
    font-size: 1.2em;
    z-index: 1000;
  `;
  mapContainer.appendChild(loader);
}

function hideLoadingIndicator() {
  const loader = document.getElementById('map-loader');
  if (loader) loader.remove();
}

window.onload = async function () {
  const artistId = document.body.getAttribute('data-artist-id');
  if (!artistId) {
    console.error('No artist ID found on page.');
    return;
  }

  const map = initArtistMap(); // Show map immediately
  showLoadingIndicator();

  try {
    const response = await fetch(`/coordinates?id=${artistId}`);
    if (!response.ok) throw new Error('Network response was not OK');

    const geoCoordinates = await response.json();
    const bounds = [];

    for (const [loc, coord] of Object.entries(geoCoordinates)) {
      const lat = parseFloat(coord.lat);
      const lon = parseFloat(coord.lon);

      if (!isNaN(lat) && !isNaN(lon)) {
        const point = addMarker(map, loc, lat, lon);
        bounds.push(point);

        // Optional: delay for progressive display
        await new Promise(res => setTimeout(res, 150)); // 150ms delay
      }
    }

    if (bounds.length > 0) {
      map.fitBounds(bounds, { animate: true, duration: 1.5 });
    }
  } catch (err) {
    console.error('Failed to load coordinates:', err);
  } finally {
    hideLoadingIndicator();
  }
};
