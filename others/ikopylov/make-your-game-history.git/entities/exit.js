import { state } from '../core/state.js';
import map from '../map.js';

class Exit {
  constructor(x, y) {
    this.x = x;
    this.y = y;
    this.width = map.config.tileSize;
    this.height = map.config.tileSize;

    this.frameCount = 6;
    this.frameIndex = 0;
    this.frameDuration = 0.15;
    this.elapsed = 0;

    const el = document.createElement('div');
    el.className = 'entity exit';
    el.style.position = 'absolute';
    el.style.zIndex = '2';
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

  destroy() {
    if (this.el && this.el.parentNode) {
      this.el.parentNode.removeChild(this.el);
    }
    this.el = null;
  }
}

export { Exit };

export default Exit;
