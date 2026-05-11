import state from '../core/state.js'; 
import map from '../map.js';
import Bomb from '../entities/bomb.js';
import { play as playSound } from '../core/audio.js';
import { showGameOver } from '../systems/pause.js';

const pressed = new Set();
let isMoving = false;
let targetX = 0;
let targetY = 0;

function isMoveKey(code) {
  return (
    code === 'ArrowUp' || code === 'KeyW' ||
    code === 'ArrowDown' || code === 'KeyS' ||
    code === 'ArrowLeft' || code === 'KeyA' ||
    code === 'ArrowRight' || code === 'KeyD'
  );
}

function keyDirection(code) {
  switch (code) {
    case 'ArrowUp':
    case 'KeyW':
      return { dx: 0, dy: -1 };
    case 'ArrowDown':
    case 'KeyS':
      return { dx: 0, dy: 1 };
    case 'ArrowLeft':
    case 'KeyA':
      return { dx: -1, dy: 0 };
    case 'ArrowRight':
    case 'KeyD':
      return { dx: 1, dy: 0 };
    default:
      return null;
  }
}

function onKeyDown(event) {
  if (state.pause.isPaused) return;

  if (isMoveKey(event.code)) {
    event.preventDefault();
    pressed.add(event.code);
    if (!isMoving) startStepFromPressed();
  }

  if (event.code === 'Space') {
    event.preventDefault();
    placeBomb();
  }
}

function onKeyUp(event) {
  if (isMoveKey(event.code)) {
    pressed.delete(event.code);
  }
}

function startStepFromPressed() {
  if (state.pause.isPaused) return;

  const player = state.player.entity;
  if (!player) return;

  const codes = Array.from(pressed);
  if (!codes.length) return;

  const code = codes[codes.length - 1];
  const dir = keyDirection(code);
  if (!dir) return;

  const { tileSize, colCount, rowCount } = map.config;
  const currCol = Math.round(player.x / tileSize);
  const currRow = Math.round(player.y / tileSize);
  const startX = currCol * tileSize;
  const startY = currRow * tileSize;

  const nextCol = currCol + dir.dx;
  const nextRow = currRow + dir.dy;

  const destX = map.clamp(nextCol * tileSize, 0, colCount * tileSize - tileSize);
  const destY = map.clamp(nextRow * tileSize, 0, rowCount * tileSize - tileSize);

  if (map.boxCollision(destX, destY, player.width, player.height)) {
    player.x = startX;
    player.y = startY;
    player.velocityX = 0;
    player.velocityY = 0;
    player.sync();
    return;
  }

  player.x = startX;
  player.y = startY;

  const baseSpeed = state.player.stats.speed;
  const baseTile = map.config.baseTileSize || map.config.tileSize;
  const scale = map.config.tileSize / baseTile;
  const speed = baseSpeed * scale;
  player.velocityX = dir.dx * speed;
  player.velocityY = dir.dy * speed;

  targetX = destX;
  targetY = destY;
  isMoving = true;
}

function update(dt) {
  const player = state.player.entity;
  if (!player) return;

  if (!isMoving) {
    player.velocityX = 0;
    player.velocityY = 0;
    return;
  }

  let newX = player.x + player.velocityX * dt;
  let newY = player.y + player.velocityY * dt;

  if (player.velocityX > 0) newX = Math.min(newX, targetX);
  if (player.velocityX < 0) newX = Math.max(newX, targetX);
  if (player.velocityY > 0) newY = Math.min(newY, targetY);
  if (player.velocityY < 0) newY = Math.max(newY, targetY);

  player.x = newX;
  player.y = newY;

  const atTarget = player.x === targetX && player.y === targetY;
  if (atTarget) {
    player.velocityX = 0;
    player.velocityY = 0;
    isMoving = false;
    startStepFromPressed();
  }

  player.sync();
}

function placeBomb() {
  if (state.pause.isPaused) return;

  const player = state.player.entity;
  if (!player || state.status.over || state.status.won) return;

  const stats = state.player.stats;
  if (player.bombCount >= stats.maxBombs) return;

  const size = map.config.tileSize;
  const col = Math.round(player.x / size);
  const row = Math.round(player.y / size);
  const x = col * size;
  const y = row * size;

  let bombExists = false;
  state.entities.bombs.forEach(bomb => {
    const bombCol = Math.floor(bomb.x / size);
    const bombRow = Math.floor(bomb.y / size);
    if (bombRow === row && bombCol === col) {
      bombExists = true;
    }
  });

  if (bombExists) return;

  const bomb = new Bomb(x, y, stats.bombRadius, player);
  state.entities.bombs.add(bomb);
  player.bombCount += 1;

  playSound('bombPlace');
}

function checkPowerUpCollisions() {
  const player = state.player.entity;
  if (!player) return;

  const size = map.config.tileSize;
  const playerCol = Math.floor(player.x / size);
  const playerRow = Math.floor(player.y / size);

  state.entities.powerUps.forEach(powerUp => {
    const col = Math.floor(powerUp.x / size);
    const row = Math.floor(powerUp.y / size);
    if (row === playerRow && col === playerCol) {
      powerUp.collect();
    }
  });
}

function checkExitCollision() {
  if (state.status.won || state.status.over) return;
  // Multiplayer: only the server ends the match when someone reaches the goal.
  if (state.multiplayer) return;

  const player = state.player.entity;
  if (!player) return;

  const size = map.config.tileSize;
  const playerWidth = player.width ?? size;
  const playerHeight = player.height ?? size;

  /** Grid cell (col, row) the player occupies — server row/col in MP avoids DOM/transform drift. */
  let playerCol;
  let playerRow;
  if (state.multiplayer) {
    const me = state.lastServerState?.players?.find((p) => p.slot === state.netSlot);
    if (me && Number.isFinite(me.col) && Number.isFinite(me.row)) {
      playerCol = me.col;
      playerRow = me.row;
    } else {
      playerCol = Math.floor((player.x + playerWidth * 0.5) / size);
      playerRow = Math.floor((player.y + playerHeight * 0.5) / size);
    }
  } else {
    playerCol = Math.floor((player.x + playerWidth * 0.5) / size);
    playerRow = Math.floor((player.y + playerHeight * 0.5) / size);
  }

  state.entities.exit.forEach((exitEntity) => {
    const stillBuried =
      exitEntity.buriedUnderSoft === true ||
      (exitEntity.el && exitEntity.el.classList && exitEntity.el.classList.contains('exit-buried'));
    if (stillBuried) return;

    const exitCol = Math.floor(exitEntity.x / size + 1e-9);
    const exitRow = Math.floor(exitEntity.y / size + 1e-9);

    if (playerCol === exitCol && playerRow === exitRow) {
      state.status.won = true;
      showGameOver(true);
    }
  });
}

function checkPlayerEnemyCollisions() {
  const player = state.player.entity;
  if (!player || player.invulnerable || state.status.over || state.status.won) return;

  const size = map.config.tileSize;
  const playerWidth = player.width ?? size;
  const playerHeight = player.height ?? size;
  const playerRect = {
    left: player.x,
    right: player.x + playerWidth,
    top: player.y,
    bottom: player.y + playerHeight,
  };

  let anyOverlap = false;
  state.entities.enemies.forEach((enemy) => {
    if (!enemy) return;

    const enemyWidth = enemy.width ?? size;
    const enemyHeight = enemy.height ?? size;

    const enemyRect = {
      left: enemy.x,
      right: enemy.x + enemyWidth,
      top: enemy.y,
      bottom: enemy.y + enemyHeight,
    };

    const overlapX = Math.min(playerRect.right, enemyRect.right) - Math.max(playerRect.left, enemyRect.left);
    const overlapY = Math.min(playerRect.bottom, enemyRect.bottom) - Math.max(playerRect.top, enemyRect.top);

    const minOverlapX = Math.min(playerWidth, enemyWidth) * 0.35;
    const minOverlapY = Math.min(playerHeight, enemyHeight) * 0.35;

    if (overlapX > minOverlapX && overlapY > minOverlapY) {
      anyOverlap = true;
      if (state.multiplayer && state.netSend) {
        if (!player._slimeHitPending) {
          player._slimeHitPending = true;
          playSound('playerHit');
          state.netSend({ type: 'slime_hit' });
        }
      } else {
        playSound('playerHit');
        state.status.lives -= 1;
        player.invulnerable = true;
        player.invulnerabilityTimer = 2.0;
        if (state.status.lives <= 0) {
          state.status.over = true;
          showGameOver(false);
        }
      }
    }
  });

  if (state.multiplayer && state.netSend && !anyOverlap) {
    player._slimeHitPending = false;
  }
}

function simulateKeyPress(code) {
  if (state.pause.isPaused) return;
  pressed.add(code);
  if (!isMoving) startStepFromPressed();
}

function simulateKeyRelease(code) {
  pressed.delete(code);
}

function reset() {
  pressed.clear();
  isMoving = false;
  targetX = 0;
  targetY = 0;
}

function initControls() {
  document.addEventListener('keydown', onKeyDown);
  document.addEventListener('keyup', onKeyUp);
}


export const playerMovement = {
  initControls,
  update,
  placeBomb,
  checkPowerUpCollisions,
  checkExitCollision,
  checkPlayerEnemyCollisions,
  simulateKeyPress,
  simulateKeyRelease,
  reset,
};

export {
  initControls,
  update,
  placeBomb,
  checkPowerUpCollisions,
  checkExitCollision,
  checkPlayerEnemyCollisions,
  simulateKeyPress,
  simulateKeyRelease,
  reset,
};

export default playerMovement;