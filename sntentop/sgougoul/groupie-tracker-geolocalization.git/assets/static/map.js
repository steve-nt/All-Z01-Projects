// map.js
(async function () {
    const params = new URLSearchParams(window.location.search);
    const artistId = params.get("id");
    if (!artistId) {
      alert("No artist specified");
      return;
    }
  
    let data;
    try {
      const res = await fetch(`/api/locations?id=${encodeURIComponent(artistId)}`);
      if (!res.ok) throw new Error(res.statusText);
      data = await res.json();
    } catch (err) {
      console.error("Fetch error:", err);
      alert("Could not load location data");
      return;
    }
  
    // 1) init map
    const map = L.map("map").setView([0, 0], 2);
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: "© OpenStreetMap contributors",
    }).addTo(map);
  
    // 2) drop markers
    const markers = data.locations.map((loc) =>
      L.marker([loc.lat, loc.lon]).bindPopup(loc.name)
    );
    const group = L.featureGroup(markers).addTo(map);
  
    // 3) auto‐zoom
    if (markers.length) {
      map.fitBounds(group.getBounds(), { padding: [20, 20] });
    }
  })();
  