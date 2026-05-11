import state from '../core/state.js';
import map from '../map.js';
import { pickFloor } from '../core/assets.js';
import Block from '../entities/block.js';
import Exit from '../entities/exit.js';
import { reset as resetPlayerAnimation } from './playerAnimation.js';

function removeEntityDom(entity) {
  if (entity && entity.el && entity.el.parentNode) {
    entity.el.parentNode.removeChild(entity.el);
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

        const sheets = state.images.enemies || [];
        const img = sheets.length ? sheets[Math.floor(Math.random() * sheets.length)] : null;
        const enemy = new Block(img, x, y, tileSize, tileSize, 'enemy');

        enemy.animElapsed = 0;
        enemy.animFps = 8;

        state.entities.enemies.add(enemy);
        tileMap[row][col] = 0;
      }
    }
  }

  if (!exitAssigned && softWallPool.length) {
    const chosen = softWallPool[Math.floor(Math.random() * softWallPool.length)];
    if (chosen) {
      chosen.containsExit = true;
      exitAssigned = true;
       // cheat for testing
    if (chosen.el) {
          chosen.el.style.border = "4px solid #00ff00"; 
          chosen.el.style.boxShadow = "0 0 15px #00ff00";
          chosen.el.style.zIndex = "100"; 
          console.log(`[CHEAT] Exit hidden at X: ${chosen.x}, Y: ${chosen.y}`);
      }

    }
  }
}

export const mapLoader = {
  loadMap,
};

export { loadMap };

export default mapLoader;