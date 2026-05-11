let brickCount = 0;
let brickInterval = null;

export function build(numberOfBricks) {
  brickCount = 0;
  brickInterval = setInterval(() => {
    if (brickCount < numberOfBricks) {
      brickCount++;
      const brick = document.createElement('div');
      brick.id = `brick-${brickCount}`;
      
      // Add data-foundation attribute for middle column bricks (column 1)
      if ((brickCount - 1) % 3 === 1) {
        brick.dataset.foundation = 'true';
      }
      
      document.body.append(brick);
    } else {
      clearInterval(brickInterval);
    }
  }, 100);
}

export function repair(...ids) {
  ids.forEach(id => {
    const brick = document.getElementById(id);
    if (brick) {
      if (brick.hasAttribute('data-foundation')) {
        brick.dataset.repaired = 'in progress';
      } else {
        brick.dataset.repaired = 'true';
      }
    }
  });
}

export function destroy() {
  const allDivs = document.querySelectorAll('body > div');
  if (allDivs.length > 0) {
    allDivs[allDivs.length - 1].remove();
    brickCount--;
  }
}
