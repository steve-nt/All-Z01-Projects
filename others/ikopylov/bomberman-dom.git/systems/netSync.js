import { state } from '../core/state.js';
import map from '../map.js';
import { pickFloor } from '../core/assets.js';
import Block from '../entities/block.js';
import Exit, { findExitAtTile } from '../entities/exit.js';

/** @type {Map<number, any>} */
const playerBySlot = new Map();
/** @type {Map<number, { el: HTMLElement }>} */
const bombById = new Map();
/** @type {Map<number, { el: HTMLElement }>} */
const powerById = new Map();

/** Server sim uses 32px cells; client `map.config.tileSize` scales to fit the viewport. */
function cellToPx(col, row) {
  const ts = map.config.tileSize;
  const c = Number(col);
  const r = Number(row);
  return {
    x: (Number.isFinite(c) ? c : 0) * ts,
    y: (Number.isFinite(r) ? r : 0) * ts,
  };
}

function repositionNetEntitiesFromLastState() {
  const payload = state.lastServerState;
  if (!payload || payload.type !== 'state' || !state.multiplayer) return;
  const ts = map.config.tileSize;
  for (const p of payload.players || []) {
    const ent = playerBySlot.get(p.slot);
    if (!ent) continue;
    const { x, y } = cellToPx(p.col, p.row);
    ent.x = x;
    ent.y = y;
    if (typeof ent.updatePlayerDimensions === 'function') ent.updatePlayerDimensions();
    ent.sync();
  }
  for (const b of payload.bombs || []) {
    const rec = bombById.get(b.id);
    if (!rec?.el) continue;
    const { x, y } = cellToPx(b.col, b.row);
    rec.el.style.width = `${ts}px`;
    rec.el.style.height = `${ts}px`;
    rec.el.style.transform = `translate3d(${Math.round(x)}px, ${Math.round(y)}px, 0)`;
  }
  for (const pu of payload.powerUps || []) {
    const rec = powerById.get(pu.id);
    if (!rec?.el) continue;
    const { x, y } = cellToPx(pu.col, pu.row);
    rec.el.style.width = `${ts}px`;
    rec.el.style.height = `${ts}px`;
    rec.el.style.transform = `translate3d(${Math.round(x)}px, ${Math.round(y)}px, 0)`;
  }
}

if (typeof window !== 'undefined') {
  window.addEventListener('bomberman-tiles-resized', () => {
    repositionNetEntitiesFromLastState();
  });
}

export function resetNetSync() {
  playerBySlot.forEach((ent) => {
    if (ent && ent.destroy) ent.destroy();
    else if (ent && ent.el && ent.el.parentNode) ent.el.parentNode.removeChild(ent.el);
  });
  playerBySlot.clear();
  bombById.forEach((b) => {
    if (b.el && b.el.parentNode) b.el.parentNode.removeChild(b.el);
  });
  bombById.clear();
  powerById.forEach((p) => {
    if (p.el && p.el.parentNode) p.el.parentNode.removeChild(p.el);
  });
  powerById.clear();
}

function destroySoftCell(r, c) {
  const ts = map.config.tileSize;
  const x = c * ts;
  const y = r * ts;
  let hadExit = false;
  state.entities.softWalls.forEach((soft) => {
    if (!soft) return;
    const gc = Math.floor(soft.x / ts + 1e-6);
    const gr = Math.floor(soft.y / ts + 1e-6);
    if (gc === c && gr === r) {
      hadExit = !!soft.containsExit;
      state.entities.softWalls.delete(soft);
      if (soft.el && soft.el.parentNode) soft.el.parentNode.removeChild(soft.el);
    }
  });
  const floor = new Block(pickFloor(r, c), x, y, ts, ts, 'floor');
  state.entities.floors.add(floor);

  if (hadExit) {
    const existing = findExitAtTile(c, r, ts);
    if (existing) {
      existing.unbury();
    } else if (state.entities.exit.size === 0) {
      const exit = new Exit(x, y);
      state.entities.exit.add(exit);
    }
  }
}

function syncTileGrid(tiles) {
  const prev = map.config.tileMap;
  if (!prev || !tiles || !tiles.length) return;
  for (let r = 0; r < map.config.rowCount; r += 1) {
    for (let c = 0; c < map.config.colCount; c += 1) {
      const next = tiles[r][c];
      const was = prev[r][c];
      if (next === was) continue;
      if (was === 2 && next === 0) {
        destroySoftCell(r, c);
      }
      prev[r][c] = next;
    }
  }
  map.buildCollisionMap();
}

/** @param {{ r: number, c: number, v: number }[]} deltas */
function applyTileDeltas(deltas) {
  const prev = map.config.tileMap;
  if (!prev || !Array.isArray(deltas) || !deltas.length) return;
  for (const d of deltas) {
    const r = d.r;
    const c = d.c;
    const v = d.v;
    if (r < 0 || c < 0 || r >= map.config.rowCount || c >= map.config.colCount) continue;
    const was = prev[r][c];
    if (was === v) continue;
    if (was === 2 && v === 0) {
      destroySoftCell(r, c);
    }
    prev[r][c] = v;
  }
  map.buildCollisionMap();
}

function ensurePlayer(p) {
  let ent = playerBySlot.get(p.slot);
  if (!ent) {
    const { x, y } = cellToPx(p.col, p.row);
    ent = new Block(state.images.bomberman, x, y, map.config.tileSize, map.config.tileSize, 'player');
    ent.bombCount = 0;
    ent.velocityX = 0;
    ent.velocityY = 0;
    ent.invulnerable = false;
    ent.invulnerabilityTimer = 0;
    ent.el.classList.add(`mp-slot-${p.slot}`);
    playerBySlot.set(p.slot, ent);
  }
  return ent;
}

function syncPlayers(players) {
  if (!Array.isArray(players)) return;
  const seen = new Set();
  for (const p of players) {
    seen.add(p.slot);
    const ent = ensurePlayer(p);
    const { x, y } = cellToPx(p.col, p.row);
    ent.x = x;
    ent.y = y;
    if (typeof ent.updatePlayerDimensions === 'function') ent.updatePlayerDimensions();
    ent.sync();
    if (p.dead) {
      ent.el.style.opacity = '0.35';
    } else {
      ent.el.style.opacity = '1';
    }
    if (p.inv) {
      ent.el.classList.add('invulnerable');
    } else {
      ent.el.classList.remove('invulnerable');
    }
    ent.invulnerable = !!p.inv;
    if (p.inv) {
      ent._slimeHitPending = false;
    }
    if (p.slot === state.netSlot) {
      state.player.entity = ent;
      ent.bombCount = p.bombCount ?? 0;
      state.player.stats.maxBombs = p.maxBombs ?? 1;
      state.player.stats.bombRadius = p.radius ?? 1;
      state.player.stats.speed = 80 + (p.speed ?? 0) * 20;
    }
  }
  playerBySlot.forEach((ent, slot) => {
    if (!seen.has(slot)) {
      if (ent && ent.destroy) ent.destroy();
      else if (ent && ent.el && ent.el.parentNode) ent.el.parentNode.removeChild(ent.el);
      playerBySlot.delete(slot);
    }
  });
}

function makeBombEl(x, y) {
  const wrapper = document.createElement('div');
  wrapper.className = 'entity bomb';
  wrapper.style.position = 'absolute';
  wrapper.style.width = `${map.config.tileSize}px`;
  wrapper.style.height = `${map.config.tileSize}px`;
  wrapper.style.zIndex = '3';
  const sprite = document.createElement('div');
  sprite.className = 'bomb-sprite';
  sprite.style.position = 'absolute';
  sprite.style.left = '15%';
  sprite.style.top = '15%';
  sprite.style.width = '70%';
  sprite.style.height = '70%';
  if (state.images.bomb) {
    sprite.style.backgroundImage = `url("${state.images.bomb.src}")`;
  }
  sprite.style.backgroundSize = 'contain';
  sprite.style.backgroundRepeat = 'no-repeat';
  sprite.style.imageRendering = 'pixelated';
  wrapper.appendChild(sprite);
  state.board.appendChild(wrapper);
  wrapper.style.left = '0';
  wrapper.style.top = '0';
  wrapper.style.transform = `translate3d(${Math.round(x)}px, ${Math.round(y)}px, 0)`;
  return wrapper;
}

function syncBombs(list) {
  const ids = new Set(list.map((b) => b.id));
  bombById.forEach((b, id) => {
    if (!ids.has(id)) {
      if (b.el && b.el.parentNode) b.el.parentNode.removeChild(b.el);
      bombById.delete(id);
    }
  });
  for (const b of list) {
    let rec = bombById.get(b.id);
    if (!rec) {
      const { x, y } = cellToPx(b.col, b.row);
      const el = makeBombEl(x, y);
      rec = { el };
      bombById.set(b.id, rec);
    } else {
      const { x, y } = cellToPx(b.col, b.row);
      rec.el.style.transform = `translate3d(${Math.round(x)}px, ${Math.round(y)}px, 0)`;
    }
  }
}

function makePowerEl(x, y, type) {
  const wrapper = document.createElement('div');
  wrapper.className = 'entity powerup';
  wrapper.style.position = 'absolute';
  wrapper.style.width = `${map.config.tileSize}px`;
  wrapper.style.height = `${map.config.tileSize}px`;
  wrapper.style.zIndex = '2';
  const icon = document.createElement('div');
  icon.className = `powerup-icon powerup-${type === 'radius' ? 'radius' : type}`;
  icon.style.width = '100%';
  icon.style.height = '100%';
  icon.style.display = 'flex';
  icon.style.alignItems = 'center';
  icon.style.justifyContent = 'center';
  icon.style.borderRadius = '50%';
  icon.style.border = '2px solid #fff';
  icon.style.fontWeight = 'bold';
  icon.style.color = '#fff';
  const colors = { bomb: '#ff6b6b', radius: '#cd781c', speed: '#45b7d1' };
  const syms = { bomb: 'B', radius: 'R', speed: 'S' };
  icon.style.backgroundColor = colors[type] || '#888';
  icon.textContent = syms[type] || '?';
  wrapper.appendChild(icon);
  state.board.appendChild(wrapper);
  wrapper.style.left = '0';
  wrapper.style.top = '0';
  wrapper.style.transform = `translate3d(${Math.round(x)}px, ${Math.round(y)}px, 0)`;
  return wrapper;
}

function syncPowerUps(list) {
  const ids = new Set(list.map((p) => p.id));
  powerById.forEach((p, id) => {
    if (!ids.has(id)) {
      if (p.el && p.el.parentNode) p.el.parentNode.removeChild(p.el);
      powerById.delete(id);
    }
  });
  for (const pu of list) {
    let rec = powerById.get(pu.id);
    if (!rec) {
      const { x, y } = cellToPx(pu.col, pu.row);
      const el = makePowerEl(x, y, pu.type);
      rec = { el };
      powerById.set(pu.id, rec);
    } else {
      const { x, y } = cellToPx(pu.col, pu.row);
      rec.el.style.transform = `translate3d(${Math.round(x)}px, ${Math.round(y)}px, 0)`;
    }
  }
}

/**
 * @param {any} payload
 */
export function applyServerState(payload) {
  if (!payload || payload.type !== 'state') return;
  state.lastServerState = payload;
  if (Array.isArray(payload.tileDeltas)) {
    if (payload.tileDeltas.length) {
      applyTileDeltas(payload.tileDeltas);
    }
  } else if (payload.tiles) {
    syncTileGrid(payload.tiles);
  }
  syncPlayers(payload.players || []);
  syncBombs(payload.bombs || []);
  syncPowerUps(payload.powerUps || []);

  if (state.spectator || state.netSlot < 0) {
    return;
  }
  const me = (payload.players || []).find((p) => p.slot === state.netSlot);
  if (me) {
    state.status.lives = me.lives;
    state.status.eliminated = !!me.dead;
  }
}
