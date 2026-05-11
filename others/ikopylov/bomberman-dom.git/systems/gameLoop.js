import state from '../core/state.js';
import { update as updateHud } from '../hud.js';
import {
  update as updatePlayerMovement,
  checkPowerUpCollisions,
  checkExitCollision,
  checkPlayerEnemyCollisions,
} from '../movement/playerMovement.js';
import { update as updatePlayerAnimation } from './playerAnimation.js';
import { update as updateEnemyAI } from '../movement/enemyMovement.js';
import { revealHiddenExitIfAllEnemiesGone } from './mapLoader.js';

function updateFps(now) {
  state.hud.fpsCounter += 1;
  if (now - state.hud.fpsLastTime >= 1000) {
    if (state.hud.fpsDisplay) {
      state.hud.fpsDisplay.textContent = state.hud.fpsCounter;
    }
    state.hud.fpsCounter = 0;
    state.hud.fpsLastTime = now;
  }
}

function step(now) {
  const dt = Math.min(0.032, (now - state.loop.lastFrame) / 1000);
  state.loop.lastFrame = now;

  const mpSpectating =
    state.multiplayer && (state.status.eliminated || state.spectator);
  const canRunMpWorld =
    state.multiplayer &&
    !state.pause.isPaused &&
    !state.status.over &&
    !state.status.won &&
    (mpSpectating || !state.status.eliminated);

  if (!state.pause.isPaused && !state.status.over && !state.status.won) {
    if (!state.multiplayer) {
      updatePlayerMovement(dt);
      updatePlayerAnimation(dt);
      updateEnemyAI(dt);

      state.entities.bombs.forEach((bomb) => bomb.update(dt));
      revealHiddenExitIfAllEnemiesGone();
      state.entities.exit.forEach((exitEntity) => {
        if (typeof exitEntity.update === 'function') {
          exitEntity.update(dt);
        }
      });

      checkPowerUpCollisions();
      checkExitCollision();
      checkPlayerEnemyCollisions();

      if (state.status.score >= 100 && !state.storyFlags.midGameShown) {
        state.storyFlags.midGameShown = true;
        if (state.system && state.system.showStory) {
          state.system.showStory('MIDGAME', true);
        }
      }

      const player = state.player.entity;
      if (player && player.invulnerable) {
        player.invulnerabilityTimer -= dt;
        if (player.invulnerabilityTimer <= 0) {
          player.invulnerable = false;
          player.el.classList.remove('invulnerable');
        } else {
          player.el.classList.add('invulnerable');
        }
      }

      state.status.time += dt;
    }
  }

  if (canRunMpWorld) {
    const player = state.player.entity;
    if (player && !state.spectator) {
      updatePlayerAnimation(dt);
    }
    updateEnemyAI(dt);
    state.entities.exit.forEach((exitEntity) => {
      if (typeof exitEntity.update === 'function') {
        exitEntity.update(dt);
      }
    });
    if (!state.status.eliminated && !state.spectator) {
      checkPlayerEnemyCollisions();
    }
    state.status.time += dt;
  }

  updateFps(now);
  updateHud();

  state.loop.requestId = requestAnimationFrame(step);
}

function start() {
  state.loop.lastFrame = performance.now();
  if (state.loop.requestId) cancelAnimationFrame(state.loop.requestId);
  state.loop.requestId = requestAnimationFrame(step);
}

export const gameLoop = {
  start,
};

export { start };

export default gameLoop;
