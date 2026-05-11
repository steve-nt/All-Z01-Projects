        // Toggle dropdown menu
        function toggleDropdown() {
            document.getElementById("filterMenu").classList.toggle("show");
        }

        // Close the dropdown if the user clicks outside of it
        window.onclick = function(event) {
            if (!event.target.matches('.filter-dropdown button') && !event.target.closest('.dropdown-content')) {
                var dropdowns = document.getElementsByClassName("dropdown-content");
                for (var i = 0; i < dropdowns.length; i++) {
                    var openDropdown = dropdowns[i];
                    if (openDropdown.classList.contains('show')) {
                        openDropdown.classList.remove('show');
                    }
                }
            }
        }

        // Populate locations dynamically
        document.addEventListener("DOMContentLoaded", function() {
            fetch('/locations')
                .then(response => response.json())
                .then(locations => {
                    locations.sort(); // Sort locations alphabetically
                    var locationsContainer = document.getElementById("locationsContainer");
                    locations.forEach(function(location) {
                        var checkbox = document.createElement("input");
                        checkbox.type = "checkbox";
                        checkbox.name = "locations";
                        checkbox.value = location;
                        locationsContainer.appendChild(checkbox);
                        locationsContainer.appendChild(document.createTextNode(location));
                        locationsContainer.appendChild(document.createElement("br"));
                    });
                })
                .catch(error => console.error('Error fetching locations:', error));
        });