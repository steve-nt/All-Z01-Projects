import { gossips } from './gossip-grid.data.js';

export function grid() {
  // 1. Create the controls (ranges)
  const rangesWrapper = document.createElement('div');
  rangesWrapper.className = 'ranges';

  const rangeConfigs = [
    { id: 'width', min: 200, max: 800, value: 250 },
    { id: 'fontSize', min: 20, max: 40, value: 20 },
    { id: 'background', min: 20, max: 75, value: 50 }
  ];

  rangeConfigs.forEach(config => {
    // The CSS implies the wrapper needs the "range" class (.range label)
    const rangeContainer = document.createElement('div');
    rangeContainer.className = 'range';

    const label = document.createElement('label');
    label.setAttribute('for', config.id);
    label.textContent = config.id;

    const input = document.createElement('input');
    input.type = 'range';
    input.id = config.id;
    input.className = 'range'; // Added to input as requested by instructions
    input.min = config.min;
    input.max = config.max;
    input.value = config.value;

    const span = document.createElement('span');
    span.textContent = config.value;

    // Listen for changes and apply dynamically to all gossip cards
    input.addEventListener('input', (e) => {
      span.textContent = e.target.value;
      const cards = document.querySelectorAll('.gossip');
      
      cards.forEach(card => {
        if (config.id === 'width') {
          card.style.width = `${e.target.value}px`;
        } else if (config.id === 'fontSize') {
          card.style.fontSize = `${e.target.value}px`;
        } else if (config.id === 'background') {
          card.style.background = `hsl(280, 50%, ${e.target.value}%)`;
        }
      });
    });

    rangeContainer.appendChild(label);
    rangeContainer.appendChild(input);
    rangeContainer.appendChild(span);
    rangesWrapper.appendChild(rangeContainer);
  });

  document.body.appendChild(rangesWrapper);

  // 2. Create the first gossip card (the form)
  const form = document.createElement('form');
  form.className = 'gossip';

  const textarea = document.createElement('textarea');
  textarea.placeholder = "Got a gossip?";

  const button = document.createElement('button');
  button.type = 'submit';
  button.textContent = "Share gossip!";

  form.appendChild(textarea);
  form.appendChild(button);

  // Add new gossip right after the form on submit
  form.addEventListener('submit', (e) => {
    e.preventDefault();
    const text = textarea.value.trim();
    if (!text) return;

    const newGossip = document.createElement('div');
    // Including the fade-in class from your CSS for a smooth entrance
    newGossip.className = 'gossip fade-in'; 
    newGossip.textContent = text;

    // Capture the *current* state of the sliders so the new card matches
    newGossip.style.width = `${document.getElementById('width').value}px`;
    newGossip.style.fontSize = `${document.getElementById('fontSize').value}px`;
    newGossip.style.background = `hsl(280, 50%, ${document.getElementById('background').value}%)`;

    form.insertAdjacentElement('afterend', newGossip);
    textarea.value = '';
  });

  document.body.appendChild(form);

  // 3. Render all gossips from the imported array
  gossips.forEach(gossipText => {
    const gossipCard = document.createElement('div');
    gossipCard.className = 'gossip';
    gossipCard.textContent = gossipText;
    document.body.appendChild(gossipCard);
  });
}