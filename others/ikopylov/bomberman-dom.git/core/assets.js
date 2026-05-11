import { state } from './state.js';

function loadImages() {
  state.board = document.getElementById('board');

  state.images.floors = [new Image(), new Image()];
  state.images.floors[0].src = './textures/empty-space-1.png';
  state.images.floors[1].src = './textures/empty-space-2.png';

  state.images.wall = new Image();
  state.images.wall.src = './textures/wall.png';

  state.images.soft = new Image();
  state.images.soft.src = './textures/soft.png';

  state.images.bomberman = new Image();
  state.images.bomberman.src = './textures/player.png';

  state.images.bomb = new Image();
  state.images.bomb.src = './textures/bomb.png';

  state.images.exit = new Image();
  state.images.exit.src = './textures/exit.png';

  state.images.enemies = [new Image(), new Image(), new Image()];
  state.images.enemies[0].src = './textures/Slime1_Walk_full.png';
  state.images.enemies[1].src = './textures/Slime2_Walk_full.png';
  state.images.enemies[2].src = './textures/Slime3_Walk_full.png';
}

function pickFloor() {
  if (!state.images.floors.length) return null;
  return Math.random() < 0.35 ? state.images.floors[1] : state.images.floors[0];
}

export const assets = {
  loadImages,
  pickFloor,
};

export { loadImages, pickFloor };

export default assets;
