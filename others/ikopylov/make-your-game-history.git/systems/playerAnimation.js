import { state } from '../core/state.js';

const PLAYER_ROWS = {
  down: 0,
  right: 3,
  up: 1,
  left: 2,
};

const PLAYER_FPS = 10;

function setPlayerFrame(col, row) {
  const player = state.player.entity;
  if (!player) return;

  player._col = col;
  player._row = row;
  if (typeof player.updatePlayerFrame === 'function') {
    player.updatePlayerFrame();
  } else if (player.el) {
    const x = -col * player.frameW;
    const y = -row * player.frameH;
    player.el.style.backgroundPosition = `${x}px ${y}px`;
  }
}

function update(dt) {
  const player = state.player.entity;
  if (!player) return;

  const animation = state.player.animation;
  const dx = player.x - animation.lastX;
  const dy = player.y - animation.lastY;
  const moving = dx !== 0 || dy !== 0;

  let row = player._row ?? PLAYER_ROWS.down;
  if (Math.abs(dx) >= Math.abs(dy)) {
    if (dx > 0) row = PLAYER_ROWS.right;
    else if (dx < 0) row = PLAYER_ROWS.left;
  } else {
    if (dy > 0) row = PLAYER_ROWS.down;
    else if (dy < 0) row = PLAYER_ROWS.up;
  }

  if (moving) {
    animation.elapsed += dt;
    const spf = 1 / PLAYER_FPS;
    if (animation.elapsed >= spf) {
      animation.elapsed -= spf;
      const nextCol = ((player._col ?? 0) + 1) % 3;
      setPlayerFrame(nextCol, row);
    } else {
      setPlayerFrame(player._col ?? 1, row);
    }
  } else {
    setPlayerFrame(1, row);
  }

  animation.lastX = player.x;
  animation.lastY = player.y;
}

function reset() {
  const player = state.player.entity;
  state.player.animation.elapsed = 0;
  if (player) {
    state.player.animation.lastX = player.x;
    state.player.animation.lastY = player.y;
    setPlayerFrame(1, PLAYER_ROWS.down);
  } else {
    state.player.animation.lastX = 0;
    state.player.animation.lastY = 0;
  }
}

export const playerAnimation = {
  setFrame: setPlayerFrame,
  update,
  reset,
};

export { setPlayerFrame, update, reset };

export default playerAnimation;
