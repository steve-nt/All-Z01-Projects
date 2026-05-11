import { state } from './core/state.js';

const baseTileMap = [
  [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
  [1, 'P', 0, 0, 2, 2, 2, 2, 0, 0, 0, 2, 2, 2, 0, 2, 2, 0, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 2, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 0, 0, 2, 0, 2, 0, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 0, 2, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 0, 'e', 2, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 0, 0, 2, 0, 2, 0, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 0, 0, 0, 0, 2, 0, 0, 'e', 0, 0, 2, 0, 0, 0, 2, 0, 2, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 0, 2, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 2, 0, 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 'e', 0, 2, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 'e', 1],
  [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1],
  [1, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 1],
  [1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 1],
  [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
];

const BASE_TILE_SIZE = 32;

const config = {
  rowCount: 20,
  colCount: 20,
  tileSize: BASE_TILE_SIZE,
  baseTileSize: BASE_TILE_SIZE,
  tileMap: baseTileMap.map(row => row.slice()),
  collisionMap: null,
};

function resetTileMap() {
  config.tileMap = baseTileMap.map(row => row.slice());
}

function buildCollisionMap() {
  config.collisionMap = config.tileMap.map(row => row.map(tile => (tile === 1 || tile === 2 ? tile : 0)));
}

function markDynamicBlock(row, col) {
  if (!config.collisionMap) return;
  if (row < 0 || col < 0 || row >= config.rowCount || col >= config.colCount) return;
  config.collisionMap[row][col] = 9;
}

function clearDynamicBlock(row, col) {
  if (!config.collisionMap) return;
  if (row < 0 || col < 0 || row >= config.rowCount || col >= config.colCount) return;
  const base = config.tileMap[row][col];
  config.collisionMap[row][col] = base === 1 || base === 2 ? base : 0;
}

function isBlocked(row, col) {
  if (row < 0 || col < 0 || row >= config.rowCount || col >= config.colCount) return true;
  if (!config.collisionMap) return false;
  const tile = config.collisionMap[row][col];
  return tile === 1 || tile === 2 || tile === 9;
}

function isBlockedAtPixel(px, py) {
  const col = Math.floor(px / config.tileSize);
  const row = Math.floor(py / config.tileSize);
  return isBlocked(row, col);
}

function boxCollision(x, y, width, height) {
  const w = width ?? config.tileSize;
  const h = height ?? config.tileSize;
  const left = x;
  const right = x + w - 1;
  const top = y;
  const bottom = y + h - 1;
  return (
    isBlockedAtPixel(left, top) ||
    isBlockedAtPixel(right, top) ||
    isBlockedAtPixel(left, bottom) ||
    isBlockedAtPixel(right, bottom)
  );
}

function clamp(value, min, max) {
  if (value < min) return min;
  if (value > max) return max;
  return value;
}

function snapX(x, deltaX) {
  if (deltaX > 0) {
    const rightTile = Math.floor((x + config.tileSize) / config.tileSize);
    return rightTile * config.tileSize - config.tileSize;
  }
  const leftTile = Math.floor(x / config.tileSize);
  return leftTile * config.tileSize;
}

function snapY(y, deltaY) {
  if (deltaY > 0) {
    const bottomTile = Math.floor((y + config.tileSize) / config.tileSize);
    return bottomTile * config.tileSize - config.tileSize;
  }
  const topTile = Math.floor(y / config.tileSize);
  return topTile * config.tileSize;
}

function updateTileSize() {
  const boardEl = state.board;
  if (!boardEl) return;

  const oldSize = config.tileSize;
  const availableWidth = window.innerWidth;
  const availableHeight = window.innerHeight - 100;
  const tileW = Math.floor(availableWidth / config.colCount);
  const tileH = Math.floor(availableHeight / config.rowCount);
  let newSize = Math.min(tileW, tileH);

  if (!Number.isFinite(newSize) || newSize <= 0) {
    newSize = oldSize;
  } else {
    newSize = Math.max(1, Math.round(newSize));
  }

  if (newSize !== config.tileSize) {
    const scale = newSize / config.tileSize;
    config.tileSize = newSize;

    const player = state.player.entity;
    if (player) {
      player.x = Math.round(player.x * scale);
      player.y = Math.round(player.y * scale);
      if (player.startX != null) player.startX = Math.round(player.startX * scale);
      if (player.startY != null) player.startY = Math.round(player.startY * scale);
    }

    const { floors, walls, softWalls, enemies, bombs, powerUps, exit } = state.entities;
    [floors, walls, softWalls, enemies, bombs, powerUps, exit].forEach(set => {
      set.forEach(entity => {
        entity.x = Math.round(entity.x * scale);
        entity.y = Math.round(entity.y * scale);
        if (entity.startX != null) entity.startX = Math.round(entity.startX * scale);
        if (entity.startY != null) entity.startY = Math.round(entity.startY * scale);
      });
    });
  }

  boardEl.style.width = `${config.tileSize * config.colCount}px`;
  boardEl.style.height = `${config.tileSize * config.rowCount}px`;
  boardEl.style.transform = 'none';
  boardEl.style.transformOrigin = 'top left';

  const { floors, walls, softWalls, enemies, bombs, powerUps, exit } = state.entities;
  const updateBlockSize = block => {
    if (!block) return;
    if (block.isPlayer) {
      if (typeof block.updatePlayerDimensions === 'function') block.updatePlayerDimensions();
      if (typeof block.sync === 'function') block.sync();
      return;
    }
    if (typeof block.updateDimensions === 'function') {
      block.updateDimensions();
      if (typeof block.sync === 'function') block.sync();
      return;
    }
    block.width = config.tileSize;
    block.height = config.tileSize;
    if (block.el) {
      block.el.style.width = `${config.tileSize}px`;
      block.el.style.height = `${config.tileSize}px`;
    }
    if (typeof block.sync === 'function') block.sync();
  };

  [floors, walls, softWalls, enemies, bombs, powerUps, exit].forEach(set => {
    set.forEach(updateBlockSize);
  });

  const player = state.player.entity;
  if (player) {
    if (typeof player.updatePlayerDimensions === 'function') player.updatePlayerDimensions();
    if (typeof player.sync === 'function') player.sync();
  }
}

window.addEventListener('resize', updateTileSize);

export const map = {
  config,
  resetTileMap,
  buildCollisionMap,
  isBlocked,
  isBlockedAtPixel,
  boxCollision,
  clamp,
  snapX,
  snapY,
  markDynamicBlock,
  clearDynamicBlock,
  updateTileSize,
};

export {
  config,
  resetTileMap,
  buildCollisionMap,
  isBlocked,
  isBlockedAtPixel,
  boxCollision,
  clamp,
  snapX,
  snapY,
  markDynamicBlock,
  clearDynamicBlock,
  updateTileSize,
};

export default map;
