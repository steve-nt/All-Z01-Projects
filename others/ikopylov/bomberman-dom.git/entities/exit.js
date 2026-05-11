import { state } from '../core/state.js';
import map from '../map.js';

class Exit {
  /**
   * @param {number} x
   * @param {number} y
   * @param {{ buried?: boolean }} [options] — if buried, goal sits on a soft block until it is destroyed
   */
  constructor(x, y, options = {}) {
    this.x = x;
    this.y = y;
    this.width = map.config.tileSize;
    this.height = map.config.tileSize;

    this.buriedUnderSoft = !!options.buried;

    this.frameCount = 6;
    this.frameIndex = 0;
    this.frameDuration = 0.15;
    this.elapsed = 0;

    const el = document.createElement('div');
    el.className = ['entity', 'exit', this.buriedUnderSoft ? 'exit-buried' : ''].filter(Boolean).join(' ');
    el.style.position = 'absolute';
    el.style.zIndex = '5';
    el.style.backgroundImage = `url("${state.images.exit?.src}")`;
    el.style.backgroundRepeat = 'no-repeat';
    el.style.imageRendering = 'pixelated';

    this.el = el;

    this.updateDimensions();
    state.board.appendChild(el);
    this.sync();
  }

  updateDimensions() {
    const size = map.config.tileSize;
    this.width = size;
    this.height = size;

    if (!this.el) return;
    this.el.style.width = `${size}px`;
    this.el.style.height = `${size}px`;
    this.el.style.backgroundSize = `${size}px ${size * this.frameCount}px`;
    this.updateFrame();
  }

  updateFrame() {
    if (!this.el) return;
    const size = map.config.tileSize;
    this.el.style.backgroundPosition = `0px -${this.frameIndex * size}px`;
  }

  update(dt) {
    this.elapsed += dt;
    while (this.elapsed >= this.frameDuration) {
      this.elapsed -= this.frameDuration;
      this.frameIndex = (this.frameIndex + 1) % this.frameCount;
      this.updateFrame();
    }
  }

  sync() {
    if (!this.el) return;
    const left = Math.round(this.x);
    const top = Math.round(this.y);
    this.el.style.left = '0px';
    this.el.style.top = '0px';
    this.el.style.transform = `translate3d(${left}px, ${top}px, 0)`;
  }

  unbury() {
    this.buriedUnderSoft = false;
    if (this.el) {
      this.el.classList.remove('exit-buried');
    }
  }

  destroy() {
    if (this.el && this.el.parentNode) {
      this.el.parentNode.removeChild(this.el);
    }
    this.el = null;
  }
}

/** @returns {Exit | null} */
function findExitAtTile(col, row, tileSize) {
  let found = null;
  state.entities.exit.forEach((ex) => {
    const c = Math.floor(ex.x / tileSize + 1e-6);
    const r = Math.floor(ex.y / tileSize + 1e-6);
    if (c === col && r === row) found = ex;
  });
  return found;
}

export { Exit, findExitAtTile };

export default Exit;
