import state, { reset } from '../core/state.js';
import { loadMap } from './mapLoader.js';
import map from '../map.js';
import { start as startGameLoop } from './gameLoop.js';
import {
  pause as pauseGame,
  resume as resumeGame,
  startFromIntro,
  restart as restartGame,
} from './pause.js';

const STORY_TEXTS = {
  INTRO: {
    title: 'INCOMING TRANSMISSION...',
    text: "Year 20XX. Malicious agents from the 'Centralization' corporation have stolen the core of the peer-to-peer learning system from the school. All students have been taken captive. You are the only one left free.\n Your mission: \n destroy the enemies and find the key!",
    btn: 'INITIATE MISSION',
  },
  MIDGAME: {
    title: 'SIGNAL INTERCEPTED!',
    text: "Excellent! You've broken through the first line of defense. Scanners show that the key to the learning system is somewhere nearby. But be careful — the enemies are becoming more aggressive!",
    btn: 'RESUME BATTLE',
  },
  WIN: {
    title: 'MISSION ACCOMPLISHED!',
    text: 'The P2P system has been restored! You destroyed the kidnappers and freed the students. The planet can once again freely exchange knowledge.',
    btn: 'CONTINUE',
  },
  LOSS: {
    title: 'CONNECTION LOST...',
    text: "You've been caught... The learning system has been destroyed. Try again, the future depends on you!",
    btn: 'RETRY SYSTEM',
  },
};

let storyTextTimeoutId = null;
let resetListenerBound = false;

function createDefaultFlags() {
  return {
    introShown: false,
    midGameShown: false,
  };
}

function ensureStoryState() {
  if (!state.story) {
    state.story = {};
  }
  const existingFlags = state.story.flags || state.storyFlags;
  state.story.flags = existingFlags || createDefaultFlags();
  state.story.active = state.story.active || false;
  state.storyFlags = state.story.flags;
}

function resetStoryModeState() {
  if (!state.story) {
    state.story = {};
  }
  state.story.flags = createDefaultFlags();
  state.storyFlags = state.story.flags;
  state.story.active = false;
  stopStoryTextTyping();
  if (state.story.elements && state.story.elements.menu) {
    state.story.elements.menu.classList.add('hidden');
  }
}

function cacheElements() {
  const menu = document.getElementById('story-menu');
  const title = document.getElementById('story-title');
  const text = document.getElementById('story-text');
  const button = document.getElementById('story-continue-btn');

  state.story.elements = { menu, title, text, button };
}

function stopStoryTextTyping() {
  if (storyTextTimeoutId) {
    clearTimeout(storyTextTimeoutId);
    storyTextTimeoutId = null;
  }
}

function typeStoryText(element, fullText, speed = 35) {
  if (!element) return;

  stopStoryTextTyping();
  element.innerText = '';

  let index = 0;
  const typeNextCharacter = () => {
    element.innerText = fullText.slice(0, index);
    index += 1;

    if (index <= fullText.length) {
      storyTextTimeoutId = setTimeout(typeNextCharacter, speed);
    } else {
      storyTextTimeoutId = null;
    }
  };

  typeNextCharacter();
}

function showIntroStory(onContinue) {
  showStoryModal('INTRO', {
    pauseGame: true,
    onContinue: () => {
      if (state.storyFlags) {
        state.storyFlags.introShown = true;
      }
      if (typeof onContinue === 'function') {
        onContinue();
      }
    },
  });
}

function handleStoryContinue(type, shouldPause, onContinue) {
  const { menu } = state.story.elements;
  if (menu) {
    menu.classList.add('hidden');
  }

  stopStoryTextTyping();
  state.story.active = false;

  if (typeof onContinue === 'function') {
    onContinue();
    return;
  }

  if (type === 'WIN') {
    restartGame();
    return;
  }

  if (shouldPause && type !== 'WIN' && type !== 'LOSS') {
    resumeGame();
  }
}

function showStoryModal(type, options = {}) {
  const storyData = STORY_TEXTS[type];
  if (!storyData) return;
  if (state.story && state.story.active) return;

  const { pauseGame: shouldPause = true, onContinue } = options;

  if (!state.story || !state.story.elements) return;
  const { menu, title, text, button } = state.story.elements;
  if (!menu || !title || !text || !button) return;

  title.innerText = storyData.title;
  typeStoryText(text, storyData.text);
  button.innerText = storyData.btn;

  state.story.active = true;
  menu.classList.remove('hidden');

  if (shouldPause) {
    pauseGame({ skipMenu: true });
  }

  button.onclick = () => handleStoryContinue(type, shouldPause, onContinue);
}

function handleIntroStart(event) {
  event.preventDefault();
  event.stopImmediatePropagation();

  if (state.pause.introMenu) {
    state.pause.introMenu.classList.add('hidden');
  }
  state.pause.introVisible = false;

  showIntroStory(() => {
    reset.all();
    loadMap();
    map.updateTileSize();
    resumeGame();
    startGameLoop();
  });
}

function bindIntroStartButton() {
  const startBtn = document.getElementById('intro-start-btn');
  if (!startBtn) return;
  if (typeof startFromIntro === 'function') {
    startBtn.removeEventListener('click', startFromIntro);
  }
  startBtn.addEventListener('click', handleIntroStart, { capture: true });
}

function handleStoryRestart() {
  resetStoryModeState();
  showIntroStory(() => {
    resumeGame();
  });
}

function init() {
  ensureStoryState();
  cacheElements();
  bindIntroStartButton();
  if (!resetListenerBound) {
    document.addEventListener('story:restart', handleStoryRestart);
    resetListenerBound = true;
  }

  if (!state.system) state.system = {};
  state.system.showStory = (type, config = true) => {
    if (typeof config === 'boolean') {
      showStoryModal(type, { pauseGame: config });
      return;
    }
    showStoryModal(type, config);
  };
}

export const storyMode = {
  init,
  show: showStoryModal,
};

export { init, showStoryModal as show };

export default storyMode;
