const cardStates = new Map();

document.addEventListener("DOMContentLoaded", () => {

    const cards = document.querySelectorAll(".artist-card");
    cards.forEach((card, index) => {
        // Initialize each card's state to 0 (front face)
        cardStates.set(card, 0);

        // Initialize card to show front face
        const cardInner = card.querySelector(".card-inner");
        const front = card.querySelector(".card-front");
        const members = card.querySelector(".card-members");
        const locations = card.querySelector(".card-locations");
        const dates = card.querySelector(".card-dates");
        const relations = card.querySelector(".card-relations");
        const progressBar = document.createElement('div');
        progressBar.classList.add('card-progress');
        card.appendChild(progressBar);

        cardInner.style.transform = "rotateY(0deg)";
        front.style.opacity = "1";
        members.style.opacity = "0";
        locations.style.opacity = "0";
        dates.style.opacity = "0";

        card.addEventListener("click", () => flipCard(card));
    });
});

function flipCard(cardElement) {
    const artistName = cardElement.getAttribute('data-artist-name');
    const currentState = cardStates.get(cardElement);
    const nextState = (currentState + 1) % 5; // Cycle through 5 states
    cardStates.set(cardElement, nextState);

    const cardInner = cardElement.querySelector(".card-inner");
    const front = cardElement.querySelector(".card-front");
    const members = cardElement.querySelector(".card-members");
    const locations = cardElement.querySelector(".card-locations");
    const dates = cardElement.querySelector(".card-dates");
    const relations = cardElement.querySelector(".card-relations");
    const progressBar = cardElement.querySelector('.card-progress');

    // Reset opacity and z-index for all faces
    front.style.opacity = "0";
    members.style.opacity = "0";
    locations.style.opacity = "0";
    dates.style.opacity = "0";
    relations.style.opacity = "0";

    front.style.zIndex = "1";
    members.style.zIndex = "1";
    locations.style.zIndex = "1";
    dates.style.zIndex = "1";
    relations.style.zIndex = "1";

    // Update the card based on the next state
    if (nextState === 0) {
        console.log(`Flipping ${artistName}'s card to: Front Side`);
        cardInner.style.transform = "rotateY(0deg)";
        front.style.opacity = "1";
        front.style.zIndex = "5";
        progressBar.style.width = "0"; // Hide progress bar
    } else if (nextState === 1) {
        console.log(`Flipping ${artistName}'s card to: Members Section`);
        cardInner.style.transform = "rotateY(180deg)";
        members.style.opacity = "1";
        members.style.zIndex = "5";
        progressBar.style.width = "25%"; // Update progress bar

        const membersTitle = cardElement.querySelector('.card-members .category-title');
        const membersList = cardElement.querySelectorAll('.card-members li');
        membersTitle.textContent = membersList.length === 1 ? 'Artist' : 'Members';
    } else if (nextState === 2) {
        console.log(`Flipping ${artistName}'s card to: Locations Section`);
        cardInner.style.transform = "rotateY(360deg)";
        locations.style.opacity = "1";
        locations.style.zIndex = "5";
        progressBar.style.width = "50%"; // Update progress bar

        const locationsTitle = cardElement.querySelector('.card-locations .category-title');
        const locationsList = cardElement.querySelectorAll('.card-locations li');
        locationsTitle.textContent = locationsList.length === 1 ? 'Location' : 'Locations';

    } else if (nextState === 3) {
        console.log(`Flipping ${artistName}'s card to: Dates Section`);
        cardInner.style.transform = "rotateY(540deg)";
        dates.style.opacity = "1";
        dates.style.zIndex = "5";
        progressBar.style.width = "75%"; // Update progress bar

        const datesTitle = cardElement.querySelector('.card-dates .category-title');
        const datesList = cardElement.querySelectorAll('.card-dates li');
        datesTitle.textContent = datesList.length === 1 ? 'Date' : 'Dates';
    } else if (nextState === 4) {
        console.log(`Flipping ${artistName}'s card to: Relations Section`);
        cardInner.style.transform = "rotateY(720deg)";
        relations.style.opacity = "1";
        relations.style.zIndex = "5";
        progressBar.style.width = "98%"; // Update progress bar

        const relationsTitle = cardElement.querySelector('.card-relations .category-title');
        const relationsList = cardElement.querySelectorAll('.card-relations li');
        relationsTitle.textContent = relationsList.length === 1 ? 'Relation' : 'Relations';
    }
}

function flipAllCards() {
    const cards = document.querySelectorAll(".artist-card");
    cards.forEach(card => flipCard(card));
}

// Add this event listener to handle the keyboard shortcut
document.addEventListener('keydown', function (event) {
    // Check if the pressed key is 'F' (you can change this to any key you prefer)
    if (event.key === 'f' || event.key === 'F') {
        flipAllCards();
    }
});
