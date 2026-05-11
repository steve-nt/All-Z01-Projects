import { state } from './core/state.js';

function init() {
  state.hud.fpsDisplay = document.getElementById('hud-fps');
  state.hud.timeDisplay = document.getElementById('hud-time');
  state.hud.scoreDisplay = document.getElementById('hud-score');
  state.hud.livesDisplay = document.getElementById('hud-lives');
  state.hud.bombsDisplay = document.getElementById('hud-bombs');
}

function update() {
  const time = Number.isFinite(state.status.time) ? state.status.time : 0;
  const minutes = Math.floor(time / 60).toString().padStart(2, '0');
  const seconds = Math.floor(time % 60).toString().padStart(2, '0');
  
  if (state.hud.timeDisplay) {
    state.hud.timeDisplay.textContent = `${minutes}:${seconds}`;
  }


  if (state.hud.scoreDisplay) {
    const currentScore = Math.floor(state.status.score);
    state.hud.scoreDisplay.textContent = currentScore;

    if (currentScore >= 500 && state.storyFlags && !state.storyFlags.midGameShown && state.system && state.system.showStory) {
      state.storyFlags.midGameShown = true; 
      state.system.showStory('MIDGAME', true); 
    }
  }

  if (state.hud.livesDisplay) {
    state.hud.livesDisplay.textContent = state.status.lives;
    

    if (state.status.lives <= 0 && state.system && state.system.showStory) {
 
    }
  }

  if (state.hud.bombsDisplay) {
    const player = state.player.entity;
    const activeBombs = player && typeof player.bombCount === 'number' ? player.bombCount : 0;
    const available = state.player.stats.maxBombs - activeBombs;
    state.hud.bombsDisplay.textContent = Math.max(available, 0);
  }
}

export const hud = {
  init,
  update,
};

export { init, update };

export default hud;