import state from '../core/state.js';
import map from '../map.js';
import { pickFloor } from '../core/assets.js';
import Block from '../entities/block.js';
import Exit, { findExitAtTile } from '../entities/exit.js';
import { reset as resetPlayerAnimation } from './playerAnimation.js';

function removeEntityDom(entity) {
  if (entity && entity.el && entity.el.parentNode) {
    entity.el.parentNode.removeChild(entity.el);
  }
}

/** Matches server `mapdata.go` spawn corners — keep enemies out of player start zones. */
const MP_SPAWN_CORNERS = [
  [1, 1],
  [1, 18],
  [18, 1],
  [18, 18],
];

function inMultiplayerSpawnSafeZone(row, col) {
  for (let i = 0; i < MP_SPAWN_CORNERS.length; i += 1) {
    const sr = MP_SPAWN_CORNERS[i][0];
    const sc = MP_SPAWN_CORNERS[i][1];
    const d = Math.abs(row - sr) + Math.abs(col - sc);
    if (d <= 3) return true;
  }
  return false;
}

/**
 * True if at least one N/E/S/W neighbor is walkable (not hard wall or soft block).
 * Prevents spawning slimes in a 1×1 floor cell fully surrounded by obstacles.
 *
 * @param {any[][]} tileMap — 0 floor, 1 wall, 2 soft; or loadMap tiles ('P','e',…).
 */
/** Prefer goal under soft blocks at least this far (Manhattan) from any corner spawn. */
const MIN_EXIT_DISTANCE_FROM_SPAWN = 10;

function minManhattanToNearestSpawn(row, col) {
  let minD = 999;
  for (let i = 0; i < MP_SPAWN_CORNERS.length; i += 1) {
    const [sr, sc] = MP_SPAWN_CORNERS[i];
    const d = Math.abs(row - sr) + Math.abs(col - sc);
    if (d < minD) minD = d;
  }
  return minD;
}

/**
 * Pick a destructible block for the exit: among cells far from spawns, choose the hardest (maximin).
 */
function pickSoftBlockForHardestExit(softWallPool, tileSize) {
  if (!softWallPool.length) return null;
  const preferred = softWallPool.filter((soft) => {
    const col = Math.floor(soft.x / tileSize + 1e-6);
    const row = Math.floor(soft.y / tileSize + 1e-6);
    return minManhattanToNearestSpawn(row, col) >= MIN_EXIT_DISTANCE_FROM_SPAWN;
  });
  const pool = preferred.length ? preferred : softWallPool;
  let best = -1;
  const candidates = [];
  for (let i = 0; i < pool.length; i += 1) {
    const soft = pool[i];
    const col = Math.floor(soft.x / tileSize + 1e-6);
    const row = Math.floor(soft.y / tileSize + 1e-6);
    const md = minManhattanToNearestSpawn(row, col);
    if (md > best) {
      best = md;
      candidates.length = 0;
      candidates.push(soft);
    } else if (md === best) {
      candidates.push(soft);
    }
  }
  return candidates[Math.floor(Math.random() * candidates.length)];
}

/** Orange soft walls fade in from map center outward (collision is already active). */
function applyStaggeredSoftWallReveal(softWallPool, tileSize, rowCount, colCount) {
  if (!softWallPool.length) return;
  const midR = (rowCount - 1) / 2;
  const midC = (colCount - 1) / 2;
  const withD2 = softWallPool.map((soft) => {
    const col = Math.floor(soft.x / tileSize + 1e-6);
    const row = Math.floor(soft.y / tileSize + 1e-6);
    const d2 = (row - midR) ** 2 + (col - midC) ** 2;
    return { soft, d2 };
  });
  withD2.sort((a, b) => a.d2 - b.d2);
  const stepMs = 24;
  withD2.forEach(({ soft }, i) => {
    if (soft.el) {
      soft.el.classList.add('soft-wall-enter');
      soft.el.style.setProperty('--soft-enter-delay', `${i * stepMs}ms`);
    }
  });
}

function enemyHasCardinalFloorExit(tileMap, row, col, rowCount, colCount) {
  const dirs = [
    [-1, 0],
    [1, 0],
    [0, -1],
    [0, 1],
  ];
  for (let i = 0; i < dirs.length; i += 1) {
    const r = row + dirs[i][0];
    const c = col + dirs[i][1];
    if (r < 0 || c < 0 || r >= rowCount || c >= colCount) continue;
    const t = tileMap[r][c];
    if (t === 1 || t === 2) continue;
    return true;
  }
  return false;
}

/** Same cells as `map.js` single-player `e` tiles — spawn slimes here when the floor is open. */
const PREFERRED_ENEMY_CELLS = [
  [6, 16],
  [10, 8],
  [14, 1],
  [14, 18],
];

/**
 * Multiplayer maps from the server only contain 0 / 1 / 2 — no enemy layer.
 * Place slimes on walkable cells so story-mode-style enemies appear in VS matches too.
 */
function spawnMultiplayerEnemies(tileSize, tileMap) {
  const sheets = state.images.enemies || [];
  const placed = new Set();
  const { rowCount, colCount } = map.config;

  const trySpawnAt = (row, col) => {
    if (row < 0 || col < 0 || row >= rowCount || col >= colCount) return false;
    if (tileMap[row][col] !== 0) return false;
    if (inMultiplayerSpawnSafeZone(row, col)) return false;
    if (!enemyHasCardinalFloorExit(tileMap, row, col, rowCount, colCount)) return false;
    const key = `${row},${col}`;
    if (placed.has(key)) return false;
    placed.add(key);

    const x = col * tileSize;
    const y = row * tileSize;
    const img = sheets.length ? sheets[Math.floor(Math.random() * sheets.length)] : null;
    const enemy = new Block(img, x, y, tileSize, tileSize, 'enemy');
    enemy.animElapsed = 0;
    enemy.animFps = 8;
    state.entities.enemies.add(enemy);
    return true;
  };

  for (let i = 0; i < PREFERRED_ENEMY_CELLS.length; i += 1) {
    const [row, col] = PREFERRED_ENEMY_CELLS[i];
    trySpawnAt(row, col);
  }

  const targetCount = 4;
  if (placed.size >= targetCount) return;

  const extras = [];
  for (let row = 0; row < rowCount; row += 1) {
    for (let col = 0; col < colCount; col += 1) {
      if (tileMap[row][col] !== 0) continue;
      if (inMultiplayerSpawnSafeZone(row, col)) continue;
      if (!enemyHasCardinalFloorExit(tileMap, row, col, rowCount, colCount)) continue;
      const key = `${row},${col}`;
      if (placed.has(key)) continue;
      extras.push([row, col]);
    }
  }
  for (let i = extras.length - 1; i > 0; i -= 1) {
    const j = Math.floor(Math.random() * (i + 1));
    const t = extras[i];
    extras[i] = extras[j];
    extras[j] = t;
  }
  for (let k = 0; k < extras.length && placed.size < targetCount; k += 1) {
    trySpawnAt(extras[k][0], extras[k][1]);
  }
}

function clearBoard() {
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

function loadMap() {
  clearBoard();
  const board = state.board || document.getElementById('board');
  if (!board) {
    throw new Error('Cannot load map without #board element in the DOM.');
  }
  state.board = board;
  board.innerHTML = '';
  map.resetTileMap();
  map.buildCollisionMap();

  const { tileSize, tileMap } = map.config;

  state.entities.exit.forEach(exitEntity => exitEntity.destroy && exitEntity.destroy());
  state.entities.exit.clear();

  const softWallPool = [];
  let exitAssigned = false;

  for (let row = 0; row < map.config.rowCount; row += 1) {
    for (let col = 0; col < map.config.colCount; col += 1) {
      const tile = tileMap[row][col];
      const x = col * tileSize;
      const y = row * tileSize;

      if (tile === 0) {
        const floor = new Block(pickFloor(row, col), x, y, tileSize, tileSize, 'floor');
        state.entities.floors.add(floor);
        continue;
      }

      if (tile === 1) {
        const wall = new Block(state.images.wall, x, y, tileSize, tileSize, 'wall');
        state.entities.walls.add(wall);
        continue;
      }

      if (tile === 2) {
        const soft = new Block(state.images.soft, x, y, tileSize, tileSize, 'soft');
        soft.containsExit = false;
        state.entities.softWalls.add(soft);
        softWallPool.push(soft);
        continue;
      }

      if (tile === 'x' || tile === 'X') {
        const floor = new Block(pickFloor(row, col), x, y, tileSize, tileSize, 'floor');
        state.entities.floors.add(floor);

        const exit = new Exit(x, y);
        state.entities.exit.add(exit);
        exitAssigned = true;

        tileMap[row][col] = 0;
        continue;
      }

      if (tile === 'P') {
        const floor = new Block(pickFloor(row, col), x, y, tileSize, tileSize, 'floor');
        state.entities.floors.add(floor);

        
        const player = new Block(state.images.bomberman, x, y, tileSize, tileSize, 'player');
        player.bombCount = 0;
        player.velocityX = 0;
        player.velocityY = 0;
        player.invulnerable = false;
        player.invulnerabilityTimer = 0;
        state.player.entity = player;

        state.player.animation.lastX = player.x;
        state.player.animation.lastY = player.y;
        resetPlayerAnimation();

        tileMap[row][col] = 0;
        continue;
      }

      if (tile === 'e') {
        const floor = new Block(pickFloor(row, col), x, y, tileSize, tileSize, 'floor');
        state.entities.floors.add(floor);

        const rc = map.config.rowCount;
        const cc = map.config.colCount;
        if (enemyHasCardinalFloorExit(tileMap, row, col, rc, cc)) {
          const sheets = state.images.enemies || [];
          const img = sheets.length ? sheets[Math.floor(Math.random() * sheets.length)] : null;
          const enemy = new Block(img, x, y, tileSize, tileSize, 'enemy');

          enemy.animElapsed = 0;
          enemy.animFps = 8;

          state.entities.enemies.add(enemy);
        }
        tileMap[row][col] = 0;
      }
    }
  }

  if (!exitAssigned && softWallPool.length) {
    const chosen = pickSoftBlockForHardestExit(softWallPool, tileSize);
    if (chosen) {
      chosen.containsExit = true;
      exitAssigned = true;
      const exit = new Exit(chosen.x, chosen.y, { buried: true });
      state.entities.exit.add(exit);
    }
  }

  applyStaggeredSoftWallReveal(softWallPool, tileSize, map.config.rowCount, map.config.colCount);
}

/**
 * @param {number[][]} tiles - 0 floor, 1 wall, 2 soft (from server)
 * @param {number} [exitRow] - server-authoritative goal cell (soft tile until blasted)
 * @param {number} [exitCol]
 */
function loadMultiplayerMap(tiles, exitRow, exitCol) {
  clearBoard();
  const board = state.board || document.getElementById('board');
  if (!board) {
    throw new Error('Cannot load map without #board element in the DOM.');
  }
  state.board = board;
  board.innerHTML = '';
  map.config.tileMap = tiles.map((row) => row.slice());
  map.buildCollisionMap();

  const { tileSize, tileMap } = map.config;
  const softWallPool = [];
  const er = Number(exitRow);
  const ec = Number(exitCol);
  const useServerExit =
    Number.isFinite(er) &&
    Number.isFinite(ec) &&
    er >= 0 &&
    ec >= 0 &&
    er < map.config.rowCount &&
    ec < map.config.colCount &&
    tileMap[er][ec] === 2;

  for (let row = 0; row < map.config.rowCount; row += 1) {
    for (let col = 0; col < map.config.colCount; col += 1) {
      const tile = tileMap[row][col];
      const x = col * tileSize;
      const y = row * tileSize;

      if (tile === 0) {
        const floor = new Block(pickFloor(row, col), x, y, tileSize, tileSize, 'floor');
        state.entities.floors.add(floor);
        continue;
      }
      if (tile === 1) {
        const wall = new Block(state.images.wall, x, y, tileSize, tileSize, 'wall');
        state.entities.walls.add(wall);
        continue;
      }
      if (tile === 2) {
        const soft = new Block(state.images.soft, x, y, tileSize, tileSize, 'soft');
        soft.containsExit = useServerExit && row === er && col === ec;
        state.entities.softWalls.add(soft);
        softWallPool.push(soft);
      }
    }
  }

  if (!useServerExit && softWallPool.length) {
    const chosen = pickSoftBlockForHardestExit(softWallPool, tileSize);
    if (chosen) {
      chosen.containsExit = true;
      const exit = new Exit(chosen.x, chosen.y, { buried: true });
      state.entities.exit.add(exit);
    }
  } else if (useServerExit) {
    const match = softWallPool.find(
      (s) =>
        Math.floor(s.y / tileSize + 1e-6) === er && Math.floor(s.x / tileSize + 1e-6) === ec,
    );
    if (match) {
      const exit = new Exit(match.x, match.y, { buried: true });
      state.entities.exit.add(exit);
    }
  }

  spawnMultiplayerEnemies(tileSize, tileMap);

  applyStaggeredSoftWallReveal(softWallPool, tileSize, map.config.rowCount, map.config.colCount);
}

/**
 * When the last slime is gone, reveal the goal portal if it was still hidden under a soft block.
 * (Classic Bomberman-style: exit appears after clearing enemies.)
 */
function revealHiddenExitIfAllEnemiesGone() {
  if (state.multiplayer) return;
  if (state.entities.enemies.size > 0) return;

  const ts = map.config.tileSize;
  let target = null;
  state.entities.softWalls.forEach((soft) => {
    if (soft && soft.containsExit) target = soft;
  });
  if (!target) return;

  const col = Math.floor(target.x / ts + 1e-6);
  const row = Math.floor(target.y / ts + 1e-6);
  target.containsExit = false;
  state.entities.softWalls.delete(target);
  if (target.el && target.el.parentNode) {
    target.el.parentNode.removeChild(target.el);
  }

  map.config.tileMap[row][col] = 0;
  if (map.config.collisionMap) {
    map.config.collisionMap[row][col] = 0;
  }

  const x = col * ts;
  const y = row * ts;
  const floor = new Block(pickFloor(row, col), x, y, ts, ts, 'floor');
  state.entities.floors.add(floor);

  const existing = findExitAtTile(col, row, ts);
  if (existing) {
    existing.unbury();
  } else if (state.entities.exit.size === 0) {
    const exit = new Exit(x, y);
    state.entities.exit.add(exit);
  }
}

export const mapLoader = {
  loadMap,
  loadMultiplayerMap,
  revealHiddenExitIfAllEnemiesGone,
};

export { loadMap, loadMultiplayerMap, revealHiddenExitIfAllEnemiesGone };

export default mapLoader;