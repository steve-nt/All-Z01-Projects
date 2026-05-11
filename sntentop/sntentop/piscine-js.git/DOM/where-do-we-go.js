import { places } from './where-do-we-go.data.js';

export function explore() {
  // 1. Sort the places from North to South (descending latitude)
  const sortedPlaces = [...places].sort((a, b) => {
    const parseLat = (coords) => {
      // Get the latitude part (before the space)
      const latPart = coords.split(' ')[0];
      const match = latPart.match(/(\d+)°(\d+)'([\d.]+)"([NS])/);
      
      const deg = parseFloat(match[1]);
      const min = parseFloat(match[2]);
      const sec = parseFloat(match[3]);
      const direction = match[4];
      
      // Convert DMS to decimal degrees
      const decimal = deg + (min / 60) + (sec / 3600);
      
      // North is positive, South is negative
      return direction === 'N' ? decimal : -decimal;
    };
    
    return parseLat(b.coordinates) - parseLat(a.coordinates);
  });

  // 2. Create fullscreen sections for each place
  sortedPlaces.forEach((place) => {
    const section = document.createElement('section');
    
    // Extract the location name before the comma, convert to lowercase, and replace spaces with hyphens
    const imageName = place.name.split(',')[0].toLowerCase().split(' ').join('-');
    
    section.style.background = `url('./where-do-we-go_images/${imageName}.jpg')`;
    section.style.backgroundSize = 'cover';
    section.style.backgroundPosition = 'center';
    
    document.body.appendChild(section);
  });

  // 3. Create the location indicator
  const locationIndicator = document.createElement('a');
  locationIndicator.className = 'location';
  locationIndicator.target = '_blank';
  document.body.appendChild(locationIndicator);

  // 4. Create the direction compass
  const directionCompass = document.createElement('div');
  directionCompass.className = 'direction';
  document.body.appendChild(directionCompass);

  // 5. Handle the scroll events to update location and compass
  let lastScrollY = window.scrollY;

  const updateOnScroll = () => {
    const currentScrollY = window.scrollY;
    
    // Update the compass direction
    if (currentScrollY > lastScrollY) {
      directionCompass.textContent = 'S';
    } else if (currentScrollY < lastScrollY) {
      directionCompass.textContent = 'N';
    }
    lastScrollY = currentScrollY;

    // Determine which section is currently taking up the majority of the viewport.
    // Since each section is 100vh, rounding the scroll position divided by the window height 
    // gives the index of the image in the exact middle of the screen.
    const currentIndex = Math.round(currentScrollY / window.innerHeight);
    const currentPlace = sortedPlaces[currentIndex];

    // Update the location indicator's text, color, and link
    if (currentPlace) {
      locationIndicator.textContent = `${currentPlace.name}\n${currentPlace.coordinates}`;
      locationIndicator.style.color = currentPlace.color;
      locationIndicator.href = `https://www.google.com/maps/place/${currentPlace.coordinates}`;
    }
  };

  // Listen to the scroll event and initialize the first view on page load
  window.addEventListener('scroll', updateOnScroll);
  updateOnScroll(); 
}