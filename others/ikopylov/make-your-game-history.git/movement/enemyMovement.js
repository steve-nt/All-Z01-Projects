import { state } from '../core/state.js';
import map from '../map.js';

function dirVec(direction) {
  switch (direction) {
    case 'U': return { vx: 0, vy: -1 };
    case 'D': return { vx: 0, vy: 1 };
    case 'L': return { vx: -1, vy: 0 };
    case 'R': return { vx: 1, vy: 0 };
    default: return { vx: 0, vy: 0 };
  }
}

function oppositeDir(direction) {
  switch (direction) {
    case 'U': return 'D';
    case 'D': return 'U';
    case 'L': return 'R';
    case 'R': return 'L';
    default: return null;
  }
}

const DEFAULT_ROW_MAP = { down: 0, right: 3, up: 1, left: 2 };
const EXTENDED_ROW_MAP = { down: 1, right: 3, up: 5, left: 2 };

function resolveRowMap(enemy) {
  if (enemy && typeof enemy.eRows === 'number' && enemy.eRows >= 6) {
    return EXTENDED_ROW_MAP;
  }
  return DEFAULT_ROW_MAP;
}

function isAtTileCenter(enemy, tileSize) {
  const cx = enemy.x + enemy.width * 0.5;
  const cy = enemy.y + enemy.height * 0.5;

  const col = Math.floor(cx / tileSize);
  const row = Math.floor(cy / tileSize);
  const centerX = (col + 0.5) * tileSize;
  const centerY = (row + 0.5) * tileSize;

  const EPS = 0.5;
  return Math.abs(cx - centerX) <= EPS && Math.abs(cy - centerY) <= EPS;
}

function getAvailableDirs(row, col) {
  const options = [];
  if (!map.isBlocked(row - 1, col)) options.push('U');
  if (!map.isBlocked(row + 1, col)) options.push('D');
  if (!map.isBlocked(row, col - 1)) options.push('L');
  if (!map.isBlocked(row, col + 1)) options.push('R');
  return options;
}

function update(dt) {
  if (state.entities.enemies.size === 0) return;

  const tileSize = map.config.tileSize;

  state.entities.enemies.forEach(enemy => {
    if (enemy.aiSpeed == null) enemy.aiSpeed = 55;
    if (enemy._baseAiSpeed == null || enemy._baseAiSpeed !== enemy.aiSpeed) {
      enemy._baseAiSpeed = enemy.aiSpeed;
    }
    if (enemy.animFps == null) enemy.animFps = 8;

    if (enemy.eWalkStartPerRow == null) enemy.eWalkStartPerRow = { 0: 0, 1: 0, 2: 0, 3: 0 };
    if (enemy.eWalkColsPerRow == null) enemy.eWalkColsPerRow = { 0: 8, 1: 8, 2: 8, 3: 8 };

    if (!enemy.aiDir) {
      const cx = enemy.x + enemy.width * 0.5;
      const cy = enemy.y + enemy.height * 0.5;
      const col = Math.floor(cx / tileSize);
      const row = Math.floor(cy / tileSize);
      const options = getAvailableDirs(row, col);
      enemy.aiDir = options.length ? options[Math.floor(Math.random() * options.length)] : 'R';
    }

    const rowMap = resolveRowMap(enemy);

    if (isAtTileCenter(enemy, tileSize)) {
      const cx = enemy.x + enemy.width * 0.5;
      const cy = enemy.y + enemy.height * 0.5;
      const col = Math.floor(cx / tileSize);
      const row = Math.floor(cy / tileSize);
      const tileKey = row * map.config.colCount + col;
      const prevTileKey = enemy._lastDecisionTileKey;
      const newTile = prevTileKey == null || prevTileKey !== tileKey;

      const options = getAvailableDirs(row, col);
      const reverse = oppositeDir(enemy.aiDir);

      let forwardOpen = false;
      if (enemy.aiDir) {
        const nextCol = enemy.aiDir === 'L' ? col - 1 : enemy.aiDir === 'R' ? col + 1 : col;
        const nextRow = enemy.aiDir === 'U' ? row - 1 : enemy.aiDir === 'D' ? row + 1 : row;
        forwardOpen = !map.isBlocked(nextRow, nextCol);
      }

      const legalNoReverse = options.filter(d => d !== reverse);
      const mustDecide = newTile || !forwardOpen;

      if (mustDecide && !forwardOpen) {
        if (legalNoReverse.length) {
          enemy.aiDir = legalNoReverse[Math.floor(Math.random() * legalNoReverse.length)];
        } else if (options.length) {
          enemy.aiDir = reverse;
        }
      } else if (mustDecide) {
        const intersectionChoices = legalNoReverse.includes(enemy.aiDir)
          ? legalNoReverse
          : [enemy.aiDir, ...legalNoReverse];
        if (intersectionChoices.length > 1) {
          enemy.aiDir = intersectionChoices[Math.floor(Math.random() * intersectionChoices.length)];
        }
      }

      if (mustDecide) {
        enemy._lastDecisionTileKey = tileKey;
      }

      if (newTile || !forwardOpen) {
        enemy.x = (col + 0.5) * tileSize - enemy.width * 0.5;
        enemy.y = (row + 0.5) * tileSize - enemy.height * 0.5;
      }
    }

    const { vx, vy } = dirVec(enemy.aiDir);
    const baseTile = map.config.baseTileSize || tileSize;
    const speed = enemy._baseAiSpeed * (tileSize / baseTile);
    const stepX = vx * speed * dt;
    const stepY = vy * speed * dt;

    if (stepX) {
      const nextX = map.clamp(enemy.x + stepX, 0, map.config.colCount * tileSize - enemy.width);
      if (!map.boxCollision(nextX, enemy.y, enemy.width, enemy.height)) {
        enemy.x = nextX;
      } else {
        enemy.x = map.snapX(enemy.x, stepX);
      }
    }

    if (stepY) {
      const nextY = map.clamp(enemy.y + stepY, 0, map.config.rowCount * tileSize - enemy.height);
      if (!map.boxCollision(enemy.x, nextY, enemy.width, enemy.height)) {
        enemy.y = nextY;
      } else {
        enemy.y = map.snapY(enemy.y, stepY);
      }
    }

    const moving = vx !== 0 || vy !== 0;
    let nextRow = enemy._eRow ?? rowMap.down;
    if (enemy.aiDir === 'R') nextRow = rowMap.right;
    else if (enemy.aiDir === 'L') nextRow = rowMap.left;
    else if (enemy.aiDir === 'U') nextRow = rowMap.up;
    else if (enemy.aiDir === 'D') nextRow = rowMap.down;

    const prevRow = enemy._eRow;
    enemy._eRow = nextRow;

    const rowStart = (enemy.eWalkStartPerRow && Number.isFinite(enemy.eWalkStartPerRow[nextRow]))
      ? enemy.eWalkStartPerRow[nextRow]
      : (enemy.eWalkStartCol ?? 0);

    const rowCols = (enemy.eWalkColsPerRow && Number.isFinite(enemy.eWalkColsPerRow[nextRow]) && enemy.eWalkColsPerRow[nextRow] > 0)
      ? enemy.eWalkColsPerRow[nextRow]
      : (enemy.eWalkCols ?? 3);

    if (prevRow !== nextRow) {
      enemy._eCol = rowStart;
    }

    if (moving) {
      enemy.animElapsed = (enemy.animElapsed || 0) + dt;
      const spf = 1 / (enemy.animFps || 8);
      if (enemy.animElapsed >= spf) {
        enemy.animElapsed -= spf;
        const rel = ((enemy._eCol ?? rowStart) - rowStart + 1) % rowCols;
        enemy._eCol = rowStart + rel;
      }
    } else {
      const idleRel = Math.floor((rowCols - 1) / 2);
      enemy._eCol = rowStart + idleRel;
    }

    if (typeof enemy.updateEnemyFrame === 'function') {
      enemy.updateEnemyFrame();
    }

    enemy.sync();
  });
}

export const enemyMovement = {
  update,
};

export { update };

export default enemyMovement;
