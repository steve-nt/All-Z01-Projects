    // Navigate to artist page
    function goToPage(element) {
        const id = element.getAttribute('data-id');
        if (id) {
            window.location.href = "/Artist/" + id;
        } else {
            console.error("ID not found for the clicked card.");
        }
    }
    
    document.getElementById("year").innerText = new Date().getFullYear();