export function generateLetters() {
  const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
  const minSize = 11;
  const maxSize = 130;
  const totalLetters = 120;

  for (let i = 0; i < totalLetters; i++) {
    const div = document.createElement('div');
    
    // Random uppercase letter
    const randomLetter = alphabet[Math.floor(Math.random() * alphabet.length)];
    div.textContent = randomLetter;
    
    // Calculate font-size: linear growth from 11 to 130
    const fontSize = minSize + (maxSize - minSize) * (i / (totalLetters - 1));
    div.style.fontSize = `${fontSize}px`;
    
    // Set font-weight based on third
    if (i < 40) {
      div.style.fontWeight = '300';
    } else if (i < 80) {
      div.style.fontWeight = '400';
    } else {
      div.style.fontWeight = '600';
    }
    
    document.body.append(div);
  }
}
