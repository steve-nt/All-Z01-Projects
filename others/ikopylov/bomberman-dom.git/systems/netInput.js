import state from '../core/state.js';

const keys = {
  up: false,
  down: false,
  left: false,
  right: false,
  bomb: false,
};

function isTypingInFormField(target) {
  if (!target || typeof target !== 'object') return false;
  const el = /** @type {any} */ (target);
  if (typeof el.closest !== 'function') return false;
  return Boolean(
    el.closest('input, textarea, select, [contenteditable="true"]'),
  );
}

function setKey(code, down) {
  switch (code) {
    case 'ArrowUp':
    case 'KeyW':
      keys.up = down;
      break;
    case 'ArrowDown':
    case 'KeyS':
      keys.down = down;
      break;
    case 'ArrowLeft':
    case 'KeyA':
      keys.left = down;
      break;
    case 'ArrowRight':
    case 'KeyD':
      keys.right = down;
      break;
    case 'Space':
      keys.bomb = down;
      break;
    default:
      break;
  }
}

function onKeyDown(event) {
  if (!state.multiplayer || !state.netSend) return;
  if (state.status.eliminated || state.spectator) return;
  if (state.pause.isPaused) return;
  if (isTypingInFormField(event.target)) return;
  if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'Space', 'KeyW', 'KeyA', 'KeyS', 'KeyD'].includes(event.code)) {
    event.preventDefault();
  }
  setKey(event.code, true);
}

function onKeyUp(event) {
  if (!state.multiplayer || !state.netSend) return;
  if (state.status.eliminated || state.spectator) return;
  if (isTypingInFormField(event.target)) return;
  setKey(event.code, false);
}

let rafId = 0;

function tick() {
  if (state.spectator) {
    rafId = requestAnimationFrame(tick);
    return;
  }
  if (state.status.eliminated) {
    rafId = requestAnimationFrame(tick);
    return;
  }
  if (state.netSend && state.multiplayer) {
    state.netSend({
      type: 'input',
      up: keys.up,
      down: keys.down,
      left: keys.left,
      right: keys.right,
      bomb: keys.bomb,
    });
  }
  rafId = requestAnimationFrame(tick);
}

/**
 * @param {(o: object) => void} send
 */
export function initNetInput(send) {
  state.netSend = send;
  document.addEventListener('keydown', onKeyDown, true);
  document.addEventListener('keyup', onKeyUp, true);
  if (rafId) cancelAnimationFrame(rafId);
  rafId = requestAnimationFrame(tick);
}

export function stopNetInput() {
  if (rafId) cancelAnimationFrame(rafId);
  rafId = 0;
  document.removeEventListener('keydown', onKeyDown, true);
  document.removeEventListener('keyup', onKeyUp, true);
  state.netSend = null;
}
