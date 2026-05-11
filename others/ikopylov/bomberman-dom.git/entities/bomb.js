import { state } from '../core/state.js';
import map from '../map.js';
import { play as playSound } from '../core/audio.js';
import { pickFloor } from '../core/assets.js';
import Block from './block.js';
import PowerUp from './powerUp.js';
import Exit, { findExitAtTile } from './exit.js';
import { showGameOver } from '../systems/pause.js';

const directions = ['up', 'down', 'left', 'right'];

class Bomb {
  constructor(x, y, radius, owner) {
    this.x = x;
    this.y = y;
    this.radius = radius;
    this.owner = owner;
    this.timer = 3.0;
    this.exploded = false;

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
    sprite.style.backgroundImage = `url("${state.images.bomb.src}")`;
    sprite.style.backgroundSize = 'contain';
    sprite.style.backgroundRepeat = 'no-repeat';
    sprite.style.imageRendering = 'pixelated';

    wrapper.appendChild(sprite);

    this.el = wrapper;
    this.spriteEl = sprite;
    state.board.appendChild(wrapper);

    const tileSize = map.config.tileSize;
    this._row = Math.floor(this.y / tileSize);
    this._col = Math.floor(this.x / tileSize);
    map.markDynamicBlock(this._row, this._col);

    this.sync();
  }

  update(dt) {
    if (this.exploded) return;
    this.timer -= dt;
    if (this.timer <= 0) {
      this.explode();
    }
  }

  explode() {
    if (this.exploded) return;
    this.exploded = true;

    this.clearCollision();

    playSound('bombExplode');

    if (this.owner && this.owner.bombCount) {
      this.owner.bombCount -= 1;
    }

    this.createExplosion();

    if (this.el.parentNode) {
      this.el.parentNode.removeChild(this.el);
    }

    state.entities.bombs.delete(this);
  }

  clearCollision() {
    if (this._row == null || this._col == null) return;
    map.clearDynamicBlock(this._row, this._col);
    this._row = null;
    this._col = null;
  }

  createExplosion() {
    const { tileSize, rowCount, colCount } = map.config;
    const centerCol = Math.floor(this.x / tileSize);
    const centerRow = Math.floor(this.y / tileSize);

    this.explodeTile(centerRow, centerCol);

    directions.forEach(dir => {
      for (let i = 1; i <= this.radius; i += 1) {
        let row = centerRow;
        let col = centerCol;

        if (dir === 'up') row -= i;
        if (dir === 'down') row += i;
        if (dir === 'left') col -= i;
        if (dir === 'right') col += i;

        if (row < 0 || row >= rowCount || col < 0 || col >= colCount) break;

        if (map.isBlocked(row, col)) {
          if (map.config.collisionMap[row][col] === 2) {
            const result = this.destroySoftWall(row, col) || {};
            const exitSpawned = !!result.exitSpawned;
            const hadExit = !!result.hadExit;
            if (!exitSpawned) {
              const chance = hadExit ? 1 : 0.5;
              if (Math.random() < chance) {
                this.spawnPowerUp(row, col);
              }
            }
          }
          break;
        }

        this.explodeTile(row, col);
      }
    });
  }

  explodeTile(row, col) {
    const size = map.config.tileSize;
    const x = col * size;
    const y = row * size;

    this.createExplosionEffect(x, y);

    const player = state.player.entity;
    if (player && !player.invulnerable) {
      const playerCol = Math.floor(player.x / size);
      const playerRow = Math.floor(player.y / size);
      if (playerRow === row && playerCol === col) {
        this.hitPlayer();
      }
    }

    state.entities.enemies.forEach(enemy => {
      const enemyCol = Math.floor(enemy.x / size);
      const enemyRow = Math.floor(enemy.y / size);
      if (enemyRow === row && enemyCol === col) {
        this.hitEnemy(enemy);
      }
    });

    state.entities.bombs.forEach(bomb => {
      if (bomb !== this && !bomb.exploded) {
        const bombCol = Math.floor(bomb.x / size);
        const bombRow = Math.floor(bomb.y / size);
        if (bombRow === row && bombCol === col) {
          bomb.explode();
          bomb.clearCollision();
        }
      }
    });
  }

  createExplosionEffect(x, y) {
    const explosion = document.createElement('div');
    explosion.className = 'entity explosion';
    explosion.style.position = 'absolute';
    explosion.style.width = `${map.config.tileSize}px`;
    explosion.style.height = `${map.config.tileSize}px`;
    explosion.style.zIndex = '5';
    explosion.style.left = `${Math.round(x)}px`;
    explosion.style.top = `${Math.round(y)}px`;
    explosion.style.transform = '';

    state.board.appendChild(explosion);

    setTimeout(() => {
      if (explosion.parentNode) {
        explosion.parentNode.removeChild(explosion);
      }
    }, 300);
  }

  hitPlayer() {
    playSound('playerHit');

    state.status.lives -= 1;
    const player = state.player.entity;
    if (player) {
      player.invulnerable = true;
      player.invulnerabilityTimer = 2.0;
    }

    if (state.status.lives <= 0) {
      state.status.over = true;
      showGameOver(false);
    }
  }

  hitEnemy(enemy) {
    state.entities.enemies.delete(enemy);
    if (enemy.el && enemy.el.parentNode) {
      enemy.el.parentNode.removeChild(enemy.el);
    }
    state.status.score += 200;

    playSound('enemyDeath');
  }

  destroySoftWall(row, col) {
    const { tileSize } = map.config;
    let target = null;
    let hadExit = false;
    let exitSpawned = false;
    state.entities.softWalls.forEach(soft => {
      const softCol = Math.floor(soft.x / tileSize);
      const softRow = Math.floor(soft.y / tileSize);
      if (softRow === row && softCol === col) {
        target = soft;
      }
    });

    if (target) {
      hadExit = !!target.containsExit;
      target.containsExit = false;
      state.entities.softWalls.delete(target);
      if (target.el && target.el.parentNode) {
        target.el.parentNode.removeChild(target.el);
      }
    }

    map.config.tileMap[row][col] = 0;
    state.status.score += 100;

    const x = col * tileSize;
    const y = row * tileSize;
    const floor = new Block(pickFloor(row, col), x, y, tileSize, tileSize, 'floor');
    state.entities.floors.add(floor);

    if (map.config.collisionMap) {
      map.config.collisionMap[row][col] = 0;
    }

    if (hadExit) {
      exitSpawned = this.spawnExit(row, col);
    }

    return { exitSpawned, hadExit };
  }

  spawnPowerUp(row, col) {
    const options = ['bomb', 'radius', 'speed'];
    const type = options[Math.floor(Math.random() * options.length)];
    const size = map.config.tileSize;
    const powerUp = new PowerUp(col * size, row * size, type);
    state.entities.powerUps.add(powerUp);
  }

  spawnExit(row, col) {
    const size = map.config.tileSize;
    const x = col * size;
    const y = row * size;

    const existing = findExitAtTile(col, row, size);
    if (existing) {
      existing.unbury();
      return true;
    }

    if (state.entities.exit.size > 0) return false;

    const exit = new Exit(x, y);
    state.entities.exit.add(exit);
    return true;
  }

  sync() {
    const left = Math.round(this.x);
    const top = Math.round(this.y);
    this.el.style.left = '0px';
    this.el.style.top = '0px';
    this.el.style.transform = `translate3d(${left}px, ${top}px, 0)`;
  }
}

export { Bomb };

export default Bomb;
