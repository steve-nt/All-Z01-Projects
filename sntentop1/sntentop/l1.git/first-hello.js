// Select the button with id 'speak-button'
const speakButton = document.querySelector("#speak-button");

// Your function that gets triggered when clicking the new button
const handleSpeakClick = (event) => {
  // Select the torso element where the text will be added or removed
  const torso = document.querySelector("#torso");

  // Check if a div with the class 'words' already exists inside the torso
  const existingDiv = document.querySelector("#torso .words");

  if (existingDiv) {
    // If the "Hello World" text exists, remove it from the torso
    torso.removeChild(existingDiv);
  } else {
    // If the "Hello World" text does not exist, create and append it
    // Create a new div element
    const newDiv = document.createElement("div");
    
    // Add the 'words' class to the div
    newDiv.classList.add("words");
    
    // Set the text content to "Hello World"
    newDiv.textContent = "Hello World";
    
    // Append the new div to the torso element
    torso.append(newDiv);
  }
};

// Attach the handleSpeakClick function to the 'speak-button' button
speakButton.addEventListener("click", handleSpeakClick);
