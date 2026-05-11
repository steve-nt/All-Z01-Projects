import { state, reset as resetState } from '../core/state.js';
import map from '../map.js';
import { loadMap } from './mapLoader.js';
import { reset as resetPlayerMovement } from '../movement/playerMovement.js';
import { reset as resetPlayerAnimation } from './playerAnimation.js';
import { update as updateHud } from '../hud.js';
import { start as startGameLoop } from './gameLoop.js';
import {
  pauseMainTheme,
  resumeMainTheme,
  stopMainTheme,
  playMainTheme,
  playVictory,
} from '../core/audio.js';
import { initScoreboard, handleGameFinished } from './scoreboard.js';

const prevNavCodes = new Set(['ArrowUp', 'KeyW', 'ArrowLeft', 'KeyA']);
const nextNavCodes = new Set(['ArrowDown', 'KeyS', 'ArrowRight', 'KeyD']);

function getVisibleButtons() {
  if (state.pause.introVisible) {
    return [state.pause.introStartBtn].filter(btn => btn && btn.style.display !== 'none');
  }

  if (state.pause.helpOpen) return [];

  if (state.pause.mode === 'name-entry') {
    return [state.pause.nameSubmitBtn, state.pause.nameDefaultBtn].filter(
      btn => btn && btn.style.display !== 'none' && !btn.classList.contains('hidden'),
    );
  }

  if (state.pause.mode === 'scoreboard') {
    return [state.pause.scoreboardPlayAgainBtn, state.pause.scoreboardCloseBtn].filter(
      btn => btn && btn.style.display !== 'none' && !btn.classList.contains('hidden'),
    );
  }

  return [state.pause.continueBtn, state.pause.restartBtn].filter(
    btn => btn && btn.style.display !== 'none' && !btn.classList.contains('hidden'),
  );
}

function focusButton(button) {
  if (!button || typeof button.focus !== 'function') return;
  try {
    button.focus({ preventScroll: true });
  } catch (err) {
    button.focus();
  }
}

function syncSelection() {
  if (state.pause.helpOpen) return;
  if (!state.pause.isPaused) return;
  const buttons = getVisibleButtons();
  if (!buttons.length) return;

  if (state.pause.selectedIndex >= buttons.length) {
    state.pause.selectedIndex = buttons.length - 1;
  }
  if (state.pause.selectedIndex < 0) {
    state.pause.selectedIndex = 0;
  }

  focusButton(buttons[state.pause.selectedIndex]);
}

function resetSelection() {
  if (state.pause.helpOpen) return;
  state.pause.selectedIndex = 0;
  syncSelection();
}

function moveSelection(delta) {
  if (state.pause.helpOpen) return;
  const buttons = getVisibleButtons();
  if (!buttons.length) return;

  const count = buttons.length;
  state.pause.selectedIndex = (state.pause.selectedIndex + delta + count) % count;
  focusButton(buttons[state.pause.selectedIndex]);
}

function activateSelection() {
  if (state.pause.helpOpen) return;
  const buttons = getVisibleButtons();
  if (!buttons.length) return;
  const selected = buttons[state.pause.selectedIndex];
  if (selected) selected.click();
}

function pause(options = {}) {
  const { skipMenu = false } = options;
  if (state.pause.isPaused) return;
  state.pause.isPaused = true;
  if (state.pause.menu && !skipMenu) {
    state.pause.menu.classList.remove('hidden');
    updatePauseContent('PAUSED', true);
    resetSelection();
  }
  if (!state.status.over && !state.status.won && !state.pause.introVisible) {
    pauseMainTheme();
  }
}

function resume() {
  state.pause.isPaused = false;
  if (state.pause.menu) {
    state.pause.menu.classList.add('hidden');
  }
  if (state.pause.helpMenu) {
    state.pause.helpMenu.classList.add('hidden');
  }
  state.pause.helpOpen = false;
  state.pause.wasPausedBeforeHelp = false;
  state.pause.selectedIndex = 0;
  state.pause.mode = 'default';
  if (state.pause.defaultPanel) {
    state.pause.defaultPanel.classList.remove('hidden');
  }
  if (state.pause.scoreboardPanel) {
    state.pause.scoreboardPanel.classList.add('hidden');
  }
  if (state.pause.namePanel) {
    state.pause.namePanel.classList.add('hidden');
  }
  if (state.pause.titleEl) {
    state.pause.titleEl.textContent = 'PAUSED';
  }
  if (state.pause.nameErrorEl) {
    state.pause.nameErrorEl.classList.add('hidden');
  }

  if (!state.status.over && !state.status.won && !state.pause.introVisible) {
    resumeMainTheme();
  }
}

function updatePauseContent(title, showContinue) {
  const pauseContent = document.getElementById('pause-content');
  if (!pauseContent) return;

  const heading = pauseContent.querySelector('h2');
  if (heading) heading.textContent = title;

  if (state.pause.continueBtn) {
    state.pause.continueBtn.style.display = showContinue ? 'block' : 'none';
  }

  syncSelection();
}

function showGameOver(won) {
  state.pause.isPaused = true;
  state.pause.helpOpen = false;
  state.pause.wasPausedBeforeHelp = false;
  if (state.pause.helpMenu) {
    state.pause.helpMenu.classList.add('hidden');
  }

  const launchScoreboardFlow = () => {
    if (state.pause.menu) {
      state.pause.menu.classList.remove('hidden');
    }
    updatePauseContent(won ? 'VICTORY!' : 'GAME OVER', false);
    resetSelection();
    handleGameFinished();
  };

  if (state.system && state.system.showStory) {
    state.system.showStory(won ? 'WIN' : 'LOSS', {
      pauseGame: true,
      onContinue: () => {
        launchScoreboardFlow();
      },
    });
  } else {
    launchScoreboardFlow();
  }

  if (won) {
    playVictory();
  } else {
    stopMainTheme();
  }
}

function showHelp() {
  if (!state.pause.helpMenu) return;
  state.pause.wasPausedBeforeHelp = state.pause.isPaused;
  state.pause.helpOpen = true;
  state.pause.isPaused = true;

  if (state.pause.menu) {
    state.pause.menu.classList.add('hidden');
  }
  state.pause.helpMenu.classList.remove('hidden');

  pauseMainTheme();
}

function hideHelp() {
  if (!state.pause.helpOpen) return;

  const shouldRestorePauseMenu = state.pause.wasPausedBeforeHelp;

  state.pause.helpOpen = false;
  state.pause.wasPausedBeforeHelp = false;

  if (state.pause.helpMenu) {
    state.pause.helpMenu.classList.add('hidden');
  }

  if (shouldRestorePauseMenu) {
    if (state.pause.menu) {
      state.pause.menu.classList.remove('hidden');
    }
    state.pause.isPaused = true;
    resetSelection();
  } else {
    resume();
  }
}

function startFromIntro() {
  state.pause.introVisible = false;
  if (state.pause.introMenu) {
    state.pause.introMenu.classList.add('hidden');
  }
  resume();
}

function restart() {
  stopMainTheme();
  resume();
  resetState.all();
  document.dispatchEvent(new Event('story:restart'));
  map.resetTileMap();
  loadMap();
  resetPlayerMovement();
  resetPlayerAnimation();
  map.updateTileSize();
  updateHud();
  startGameLoop();

  playMainTheme({ restart: true });
}

function bindIntroControls(event) {
  const { code } = event;

  if (prevNavCodes.has(code)) {
    event.preventDefault();
    moveSelection(-1);
    event.stopImmediatePropagation();
    return true;
  }

  if (nextNavCodes.has(code)) {
    event.preventDefault();
    moveSelection(1);
    event.stopImmediatePropagation();
    return true;
  }

  if (code === 'Space' || code === 'Enter') {
    event.preventDefault();
    activateSelection();
    event.stopImmediatePropagation();
    return true;
  }

  if (code === 'Escape' || code === 'KeyP') {
    event.preventDefault();
    event.stopImmediatePropagation();
    return true;
  }

  return false;
}

function init() {
  state.pause.menu = document.getElementById('pause-menu');
  state.pause.continueBtn = document.getElementById('continue-btn');
  state.pause.restartBtn = document.getElementById('restart-btn');
  state.pause.helpMenu = document.getElementById('help-menu');
  state.pause.helpCloseBtn = document.getElementById('help-close-btn');
  state.pause.introMenu = document.getElementById('intro-menu');
  state.pause.introStartBtn = document.getElementById('intro-start-btn');
  state.pause.titleEl = document.getElementById('pause-title');
  state.pause.defaultPanel = document.getElementById('pause-panel-default');
  state.pause.scoreboardPanel = document.getElementById('scoreboard-panel');
  state.pause.scoreboardPlayAgainBtn = document.getElementById('scoreboard-play-again');
  state.pause.scoreboardCloseBtn = document.getElementById('scoreboard-close');
  state.pause.namePanel = document.getElementById('name-entry-panel');
  state.pause.nameInput = document.getElementById('name-entry-input');
  state.pause.nameSubmitBtn = document.getElementById('name-entry-submit');
  state.pause.nameDefaultBtn = document.getElementById('name-entry-default');
  state.pause.nameErrorEl = document.getElementById('name-entry-error');
  state.pause.mode = 'default';

  initScoreboard({ onRestart: restart });

  const helpBtn = document.getElementById('ui-help-btn');
  if (helpBtn) {
    helpBtn.addEventListener('click', showHelp);
  }
  if (state.pause.helpCloseBtn) {
    state.pause.helpCloseBtn.addEventListener('click', hideHelp);
  }

  if (state.pause.continueBtn) {
    state.pause.continueBtn.addEventListener('click', resume);
  }
  if (state.pause.restartBtn) {
    state.pause.restartBtn.addEventListener('click', restart);
  }

  if (state.pause.introMenu) {
    if (state.multiplayer) {
      state.pause.introMenu.classList.add('hidden');
      state.pause.introVisible = false;
      state.pause.isPaused = false;
    } else {
      state.pause.introMenu.classList.remove('hidden');
      state.pause.introVisible = true;
      state.pause.isPaused = true;
      resetSelection();
      pauseMainTheme();
    }
  } else {
    state.pause.introVisible = false;
  }

  document.addEventListener('pause:reset-selection', resetSelection);

  document.addEventListener('keydown', event => {
    const { code } = event;

    if (state.story && state.story.active) {
      event.preventDefault();
      event.stopImmediatePropagation();
      return;
    }

    if (state.pause.introVisible && bindIntroControls(event)) {
      return;
    }

    if (state.pause.helpOpen && (code === 'Escape' || code === 'KeyP' || code === 'Space')) {
      event.preventDefault();
      hideHelp();
      event.stopImmediatePropagation();
      return;
    }

    const isGameFinished = state.status.over || state.status.won;
    const isDuringNameEntry = state.pause.mode === 'name-entry';
    const isShowingScoreboard = state.pause.mode === 'scoreboard';
    const pauseToggleDisabled =
      isGameFinished || isDuringNameEntry || isShowingScoreboard || state.multiplayer;


    if (code === 'Escape' || code === 'KeyP') {
       if (pauseToggleDisabled) {
        return;
      }
      event.preventDefault();
      if (state.pause.isPaused) {
        resume();
      } else {
        pause();
      }
      event.stopImmediatePropagation();
      return;
    }

    if (!state.pause.isPaused) return;

    if (prevNavCodes.has(code)) {
      event.preventDefault();
      moveSelection(-1);
      event.stopImmediatePropagation();
      return;
    }

    if (nextNavCodes.has(code)) {
      event.preventDefault();
      moveSelection(1);
      event.stopImmediatePropagation();
      return;
    }

    if (code === 'Space') {
      event.preventDefault();
      activateSelection();
      event.stopImmediatePropagation();
    }
  });
}

export const pauseSystem = {
  init,
  pause,
  resume,
  showGameOver,
  restart,
  showHelp,
  hideHelp,
  startFromIntro,
};

export {
  init,
  pause,
  resume,
  showGameOver,
  restart,
  showHelp,
  hideHelp,
  startFromIntro,
};

export default pauseSystem;
