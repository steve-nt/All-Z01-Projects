import state, { reset } from './core/state.js';
import { loadImages } from './core/assets.js';
import { init as initAudio } from './core/audio.js';
import { init as initHud, update as updateHud } from './hud.js';
import { init as initPause } from './systems/pause.js';
import { initControls as initPlayerControls } from './movement/playerMovement.js';
import { loadMap } from './systems/mapLoader.js';
import map from './map.js';
import { init as initStoryMode } from './systems/storyMode.js';

function bootstrap() {
  console.log("Bootstrap started");
  
 
  reset.all();
  loadImages();

  const board = document.getElementById('board');
  if (!board) {
    throw new Error('Board element with id "board" was not found.');
  }
  state.board = board;
  board.innerHTML = '';

  initAudio();
  initHud();
  initPause();
  initPlayerControls();

  loadMap();
  map.updateTileSize();

  initStoryMode(); 
  updateHud();
}

window.addEventListener('load', bootstrap);
