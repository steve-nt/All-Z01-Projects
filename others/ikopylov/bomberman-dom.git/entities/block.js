import state from '../core/state.js'; 
import map from '../map.js';

class Block {
  constructor(image, x, y, width, height, className = '') {
    this.image = image;
    this.x = x;
    this.y = y;
    this.startX = x;
    this.startY = y;
    this.width = width;
    this.height = height;

    const el = document.createElement('div');
    
    el.className = ['entity', className].filter(Boolean).join(' ');
    
    el.style.position = 'absolute';
    el.style.width = `${width}px`;
    el.style.height = `${height}px`;
    
    if (image) {
      el.style.backgroundImage = `url("${image.src}")`;
      el.style.backgroundRepeat = 'no-repeat';
    } else {
        if (className === 'player') el.style.backgroundColor = 'rgba(255, 0, 0, 0.5)';
    }

    this.el = el;

    this.isDynamic = className === 'player' || className === 'enemy';

    if (className === 'player') {
      this.isPlayer = true;
      this.cols = 3;
      this.rows = 4;
      this.baseFrameWidth = 16;
      this.baseFrameHeight = 18;
      this._col = 1;
      this._row = 0;

      el.style.imageRendering = 'pixelated';
      el.style.backgroundRepeat = 'no-repeat';

      this.updatePlayerDimensions();
      this.updatePlayerFrame();
    } else if (className === 'enemy') {
      this.isEnemy = true;

      this.enemySrcW = 64;
      this.enemySrcH = 64;

      this.enemyGx = 0;
      // Slime sheets (512×256) are a tight 8×4 grid of 64×64 frames — no inter-row gap.
      this.enemyGy = 0;

      this.enemyOffX = 0;
      this.enemyOffY = 0;

      this.eCols = 8;
      this.eRows = 4;

      this._eCol = 0;
      this._eRow = 0;

      this.enemyTrimL = 12;
      this.enemyTrimR = 12;
      this.enemyTrimT = 2;
      this.enemyTrimB = 4;

      el.style.imageRendering = 'pixelated';
      el.style.backgroundRepeat = 'no-repeat';

      if (this.image) {
        const onReady = () => {
          if (!this.image.naturalWidth || !this.image.naturalHeight) return;
          const cols = Math.max(1, Math.floor((this.image.naturalWidth + this.enemyGx) / (this.enemySrcW + this.enemyGx)));
          const rows = Math.max(1, Math.floor((this.image.naturalHeight + this.enemyGy) / (this.enemySrcH + this.enemyGy)));
          this.eCols = cols;
          this.eRows = rows;
          this.eWalkStartCol = 0;
          this.eWalkCols = 8;

          this.eWalkStartPerRow = { 0: 0, 1: 0, 2: 0, 3: 0 };
          this.eWalkColsPerRow = {
            0: Math.min(cols, 8),
            1: Math.min(cols, 8),
            2: Math.min(cols, 8),
            3: Math.min(cols, 8),
          };

          this.updateEnemyDimensions();
          this.updateEnemyFrame();
          this.sync();
        };

        if (this.image.complete) onReady();
        else this.image.addEventListener('load', onReady, { once: true });
      }

      this.updateEnemyDimensions();
      this.updateEnemyFrame();
    } else {
      el.style.backgroundSize = '100% 100%';
      el.style.backgroundRepeat = 'no-repeat';
    }

    if (!state.board) {
      state.board = document.getElementById('board');
    }
    if (state.board) {
      state.board.appendChild(el);
    }
    this.sync();
  }

  updatePlayerDimensions() {
    if (!this.isPlayer) return;

    const tileSize = map.config.tileSize;
    this.lastPlayerTileSize = tileSize;

    const targetWidth = Math.round(tileSize * 0.75);
    const scale = targetWidth / this.baseFrameWidth;

    this.frameW = Math.round(this.baseFrameWidth * scale);
    this.frameH = Math.round(this.baseFrameHeight * scale);

    if (this.el) {
      this.el.style.width = `${this.frameW}px`;
      this.el.style.height = `${this.frameH}px`;
      this.el.style.backgroundSize = `${this.frameW * this.cols}px ${this.frameH * this.rows}px`;
    }
  }

  updatePlayerFrame() {
    if (!this.isPlayer) return;
    const col = this._col ?? 1;
    const row = this._row ?? 0;
    if (this.el) {
      this.el.style.backgroundPosition = `${-col * this.frameW}px ${-row * this.frameH}px`;
    }
  }

  updateEnemyDimensions() {
    if (!this.isEnemy) return;

    const tileSize = map.config.tileSize;

    const visibleSrcW = Math.max(1, this.enemySrcW - (this.enemyTrimL + this.enemyTrimR));
    const targetWidth = Math.min(tileSize, tileSize * 0.75 + 12);
    const scale = targetWidth / visibleSrcW;
    this._enemyScale = scale;

    if (this.image && this.image.naturalWidth && this.image.naturalHeight) {
      const cols = Math.max(1, Math.floor((this.image.naturalWidth + this.enemyGx) / (this.enemySrcW + this.enemyGx)));
      const rows = Math.max(1, Math.floor((this.image.naturalHeight + this.enemyGy) / (this.enemySrcH + this.enemyGy)));
      this.eCols = cols;
      this.eRows = rows;
    }

    const srcW = this.enemySrcW - (this.enemyTrimL + this.enemyTrimR);
    const srcH = this.enemySrcH - (this.enemyTrimT + this.enemyTrimB);

    this.eFrameW = Math.floor(srcW * scale);
    this.eFrameH = Math.floor(srcH * scale);
    this.ePadX = Math.floor(this.enemyGx * scale);
    this.ePadY = Math.floor(this.enemyGy * scale);

    const naturalW = this.image?.naturalWidth || this.eCols * this.enemySrcW + (this.eCols - 1) * this.enemyGx;
    const naturalH = this.image?.naturalHeight || this.eRows * this.enemySrcH + (this.eRows - 1) * this.enemyGy;
    const sheetW = Math.floor(naturalW * scale);
    const sheetH = Math.floor(naturalH * scale);

    if (this.el) {
      this.el.style.width = `${this.eFrameW}px`;
      this.el.style.height = `${this.eFrameH}px`;
      this.el.style.backgroundSize = `${sheetW}px ${sheetH}px`;
    }
  }

  updateEnemyFrame() {
    if (!this.isEnemy) return;

    const col = this._eCol ?? 1;
    const row = this._eRow ?? 0;

    const srcX = col * (this.enemySrcW + this.enemyGx) + this.enemyTrimL;
    const srcY = row * (this.enemySrcH + this.enemyGy) + this.enemyTrimT;

    const scale = this._enemyScale || 1;
    const offsetX = Math.floor(srcX * scale);
    const offsetY = Math.floor(srcY * scale);

    if (this.el) {
      this.el.style.backgroundPosition = `${-offsetX}px ${-offsetY}px`;
    }
  }

  sync() {
    const left = Math.round(this.x);
    const top = Math.round(this.y);
    
    if (this.el) {
      this.el.style.left = '0px';
      this.el.style.top = '0px';
      this.el.style.transform = `translate3d(${left}px, ${top}px, 0)`;
    }
  }

  destroy() {
    if (this.el && this.el.parentNode) {
      this.el.parentNode.removeChild(this.el);
    }
    this.el = null;
  }
}

export { Block };

export default Block;