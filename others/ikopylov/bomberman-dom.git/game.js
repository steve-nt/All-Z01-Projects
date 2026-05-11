import state, { reset } from './core/state.js';
import { loadImages } from './core/assets.js';
import { init as initAudio, resumeMainTheme } from './core/audio.js';
import { init as initHud, update as updateHud } from './hud.js';
import { init as initPause } from './systems/pause.js';
import { initControls as initPlayerControls } from './movement/playerMovement.js';
import { loadMap, loadMultiplayerMap } from './systems/mapLoader.js';
import map from './map.js';
import { init as initStoryMode } from './systems/storyMode.js';
import { initNetInput } from './systems/netInput.js';
import { start as startGameLoop } from './systems/gameLoop.js';
import { resetNetSync } from './systems/netSync.js';

/**
 * @param {{ multiplayer?: boolean, spectator?: boolean, mySlot?: number, tiles?: number[][], exitRow?: number, exitCol?: number, send?: (o: object) => void }} [options]
 */
function bootstrap(options = {}) {
  console.log('Bootstrap started', options);

  reset.all();
  state.multiplayer = !!options.multiplayer;
  state.spectator = !!options.spectator;
  state.netSlot = options.mySlot ?? 0;
  state.netSend = options.send ?? null;

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

  if (state.multiplayer) {
    resetNetSync();
    if (state.spectator) {
      const host = document.getElementById('game-container') || document.body;
      let bar = document.getElementById('spectate-banner');
      if (!bar) {
        bar = document.createElement('div');
        bar.id = 'spectate-banner';
        bar.className = 'spectate-banner';
        bar.textContent = 'Spectating — you cannot control players';
        host.insertBefore(bar, host.firstChild);
      }
    } else {
      initNetInput(options.send);
    }
    loadMultiplayerMap(options.tiles, options.exitRow, options.exitCol);
    map.updateTileSize();
    updateHud();
    startGameLoop();
    resumeMainTheme();
  } else {
    initPlayerControls();
    loadMap();
    map.updateTileSize();
    initStoryMode();
    updateHud();
  }
}

export { bootstrap };
