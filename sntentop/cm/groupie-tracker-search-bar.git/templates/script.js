// document.addEventListener('DOMContentLoaded', function() {
//     console.log('JavaScript Loaded'); // Check if script is loaded and DOM is ready
  
//     const searchInput = document.getElementById('query');
//     if (!searchInput) {
//       console.log('Input element not found!');
//       return;
//     }
  
//     searchInput.addEventListener('input', function(event) {
//       const query = event.target.value;
//       console.log('Input detected:', query); // This will log the value of input as you type
  
//       // Avoid making requests when input is empty
//       if (query.trim() === "") {
//         return;
//       }
  
//       fetch(`/searchSuggestions?query=${query}`)
//         .then(response => response.json())
//         .then(suggestions => {
//           console.log('Suggestions received:', suggestions);
//           const datalist = document.getElementById('data');
//           datalist.innerHTML = ''; // Clear previous suggestions
  
//           suggestions.forEach(artist => {
//             const option = document.createElement('option');
//             option.value = `${artist.name} - ${artist.firstAlbum}`;
//             datalist.appendChild(option);
//           });
//         })
//         .catch(error => {
//           console.error('Error fetching suggestions:', error);
//         });
//     });
//   });
  