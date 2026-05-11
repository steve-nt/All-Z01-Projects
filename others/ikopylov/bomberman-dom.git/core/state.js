const DEFAULT_PLAYER_STATS = Object.freeze({
  maxBombs: 1,
  bombRadius: 1,
  speed: 80,
});

function createPlayerStats() {
  return {
    maxBombs: DEFAULT_PLAYER_STATS.maxBombs,
    bombRadius: DEFAULT_PLAYER_STATS.bombRadius,
    speed: DEFAULT_PLAYER_STATS.speed,
  };
}

const state = {
  board: null,
  images: {
    floors: [],
    wall: null,
    soft: null,
    bomberman: null,
    bomb: null,
    enemy: null,
    exit: null,
  },
  entities: {
    floors: new Set(),
    walls: new Set(),
    softWalls: new Set(),
    bombs: new Set(),
    powerUps: new Set(),
    enemies: new Set(),
    exit: new Set(),
  },
  player: {
    entity: null,
    stats: createPlayerStats(),
    animation: {
      lastX: 0,
      lastY: 0,
      elapsed: 0,
      frame: { col: 1, row: 0 },
    },
  },
  status: {
    time: 0,
    score: 0,
    lives: 3,
    startTime: 0,
    over: false,
    won: false,
    /** Multiplayer: you are eliminated but the match continues for others */
    eliminated: false,
  },
  pause: {
    isPaused: false,
    menu: null,
    continueBtn: null,
    restartBtn: null,
    selectedIndex: 0,
    helpMenu: null,
    helpCloseBtn: null,
    helpOpen: false,
    wasPausedBeforeHelp: false,
    introMenu: null,
    introStartBtn: null,
    introVisible: true,
    titleEl: null,
    defaultPanel: null,
    scoreboardPanel: null,
    scoreboardPlayAgainBtn: null,
    scoreboardCloseBtn: null,
    namePanel: null,
    nameInput: null,
    nameSubmitBtn: null,
    nameDefaultBtn: null,
    nameErrorEl: null,
    mode: 'default',
  },
  hud: {
    fpsCounter: 0,
    fpsLastTime: performance.now(),
    fpsDisplay: null,
    timeDisplay: null,
    scoreDisplay: null,
    livesDisplay: null,
    bombsDisplay: null,
  },
  loop: {
    lastFrame: performance.now(),
    requestId: null,
  },
  /** Multiplayer (server-authoritative) */
  multiplayer: false,
  spectator: false,
  netSlot: 0,
  netSend: null,
  lastServerState: null,
};

function resetPlayerStats() {
  state.player.stats = createPlayerStats();
  if (state.player.entity) {
    removeEntityDom(state.player.entity);
  }
  state.player.entity = null;
  state.player.animation.lastX = 0;
  state.player.animation.lastY = 0;
  state.player.animation.elapsed = 0;
  state.player.animation.frame.col = 1;
  state.player.animation.frame.row = 0;
}

function resetGameStatus() {
  state.status.time = 0;
  state.status.score = 0;
  state.status.lives = 3;
  state.status.startTime = performance.now();
  state.status.over = false;
  state.status.won = false;
  state.status.eliminated = false;
  state.pause.isPaused = false;
  state.pause.selectedIndex = 0;
}

function removeEntityDom(entity) {
  if (entity && entity.el && entity.el.parentNode) {
    entity.el.parentNode.removeChild(entity.el);
  }
}

function clearEntities() {
  const { floors, walls, softWalls, bombs, powerUps, enemies, exit } = state.entities;
  [floors, walls, softWalls, bombs, powerUps, enemies, exit].forEach(set => {
    set.forEach(entity => {
      if (entity && typeof entity.clearCollision === 'function') {
        entity.clearCollision();
      }
      if (entity && typeof entity.destroy === 'function') {
        entity.destroy();
        return;
      }
      removeEntityDom(entity);
    });
    set.clear();
  });

  if (state.player.entity) {
    removeEntityDom(state.player.entity);
    state.player.entity = null;
  }
}

export const reset = {
  player: resetPlayerStats,
  status: resetGameStatus,
  entities: clearEntities,
  all() {
    clearEntities();
    resetPlayerStats();
    resetGameStatus();
    state.multiplayer = false;
    state.spectator = false;
    state.netSlot = 0;
    state.netSend = null;
    state.lastServerState = null;
  },
};

export { state, resetPlayerStats, resetGameStatus, clearEntities };

export default state;
