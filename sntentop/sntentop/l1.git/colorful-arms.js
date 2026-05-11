// Select the button with id 'arm-color'
const armColorButton = document.querySelector("#arm-color");

// Select the left and right arm elements
const armLeft = document.querySelector("#arm-left");
const armRight = document.querySelector("#arm-right");

// Select the button with id 'speak-button'
const speakButton = document.querySelector("#speak-button");

// Your function that gets triggered when clicking the arm-color button
const handleChangeArmColor = (event) => {
  // Generate a random color
  const randomColor = `#${Math.floor(Math.random() * 16777215).toString(16)}`;

  // Apply the random color to both arms
  armLeft.style.backgroundColor = randomColor;
  armRight.style.backgroundColor = randomColor;
};

// Your function that gets triggered when clicking the speak button
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
    const newDiv = document.createElement("div");
    newDiv.classList.add("words");
    newDiv.textContent = "Hello World";
    torso.append(newDiv);
  }
};

// Attach the handleChangeArmColor function to the 'arm-color' button
armColorButton.addEventListener("click", handleChangeArmColor);

// Attach the handleSpeakClick function to the 'speak-button' button
speakButton.addEventListener("click", handleSpeakClick);
