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


//  Turn a “DD-MM-YYYY” string into a JS Date.
function parseDMY(s) {
  const [d, m, y] = s.split("-").map(Number);
  return new Date(y, m - 1, d);
}

async function addMarkers(locations) {
  if (!Array.isArray(locations)) {
    console.error("Invalid locations data");
    return;
  }

  // for each address, keep only its earliest date
  // deduplicate by building a map of address → {address, date}, always picking the lowest date.
  const earliestByAddress = {};
  locations.forEach(loc => {
    const existing = earliestByAddress[loc.address];
    if (
      !existing ||
      parseDMY(loc.date) < parseDMY(existing.date)
    ) {
      earliestByAddress[loc.address] = loc;
    }
  });
  // get deduped list
  const onlyOne = Object.values(earliestByAddress);

  // true chronological sort
  onlyOne.sort((a, b) => parseDMY(a.date) - parseDMY(b.date));

  // geocode every entry in parallel
  const enriched = await Promise.all(
    onlyOne.map(async loc => {
      try {
        const res  = await fetch(
          `/geolocations?address=${encodeURIComponent(loc.address)}`
        );
        if (!res.ok) throw new Error(res.statusText);
        const data = await res.json();
        return { ...loc, position: { lat: data.latitude, lng: data.longitude } };
      } catch (e) {
        console.error("Geocode failed for", loc.address, e);
        return null;
      }
    })
  );

  // place markers in sorted order
  const path = [];
  enriched
    .filter(e => e && e.position)
    .forEach((loc, i) => {
      const marker = new google.maps.Marker({
        position: loc.position,
        map,
        title:    loc.address,
        label:    { 
            text: `${i+1}`, 
            color: "#fff", 
            fontSize: "12px", 
            fontWeight: "bold" 
        }
      });

      const iw = new google.maps.InfoWindow({
        content: `<div style="padding:5px; color:#000; font-weight:bold;">${loc.address}</div>`
      });
      marker.addListener("click", () => {
        infoWindows.forEach(x=>x.close());
        iw.open({ map, anchor: marker });
      });

      markers.push(marker);
      infoWindows.push(iw);
      path.push(loc.position);
    });

  // draw the connecting line
  new google.maps.Polyline({
    path,
    geodesic:     true,
    strokeColor:  "#e94560",
    strokeOpacity: 1,
    strokeWeight:  2,
    map
  });

  // fit map to markers
  if (markers.length) {
    const bounds = new google.maps.LatLngBounds();
    markers.forEach(m => bounds.extend(m.getPosition()));
    map.fitBounds(bounds);
  }
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

function drawPolyline(pathCoordinates) {
    const line = new google.maps.Polyline({
        path: pathCoordinates.filter(Boolean),
        geodesic: true,
        strokeColor: "#e94560",
        strokeOpacity: 1.0,
        strokeWeight: 2
    });
    line.setMap(map);
}

