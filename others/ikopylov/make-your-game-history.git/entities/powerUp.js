import { state } from '../core/state.js';
import map from '../map.js';
import { play as playSound } from '../core/audio.js';

class PowerUp {
  constructor(x, y, type) {
    this.x = x;
    this.y = y;
    this.type = type;

    const wrapper = document.createElement('div');
    wrapper.className = 'entity powerup';
    wrapper.style.position = 'absolute';
    wrapper.style.width = `${map.config.tileSize}px`;
    wrapper.style.height = `${map.config.tileSize}px`;
    wrapper.style.zIndex = '2';

    const icon = document.createElement('div');
    icon.className = `powerup-icon powerup-${type}`;
    icon.style.width = '100%';
    icon.style.height = '100%';
    icon.style.display = 'flex';
    icon.style.alignItems = 'center';
    icon.style.justifyContent = 'center';
    icon.style.borderRadius = '50%';
    icon.style.backgroundColor = this.getColor();
    icon.style.border = '2px solid #fff';
    icon.style.fontSize = '12px';
    icon.style.fontWeight = 'bold';
    icon.style.color = '#fff';
    icon.textContent = this.getSymbol();

    wrapper.appendChild(icon);

    this.el = wrapper;
    this.iconEl = icon;
    state.board.appendChild(wrapper);
    this.sync();
  }

  getColor() {
    switch (this.type) {
      case 'bomb':
        return '#ff6b6b';
      case 'radius':
        return '#cd781cff';
      case 'speed':
        return '#45b7d1';
      default:
        return '#ff6b6b';
    }
  }

  getSymbol() {
    switch (this.type) {
      case 'bomb':
        return 'B';
      case 'radius':
        return 'R';
      case 'speed':
        return 'S';
      default:
        return '?';
    }
  }

  collect() {
    const stats = state.player.stats;
    switch (this.type) {
      case 'bomb':
        stats.maxBombs += 1;
        break;
      case 'radius':
        stats.bombRadius += 1;
        break;
      case 'speed':
        stats.speed += 20;
        break;
      default:
        break;
    }

    playSound('powerUp', 600, 0.15);

    state.entities.powerUps.delete(this);
    if (this.el.parentNode) {
      this.el.parentNode.removeChild(this.el);
    }

    state.status.score += 50;
  }

  sync() {
    const left = Math.round(this.x);
    const top = Math.round(this.y);
    this.el.style.left = `${left}px`;
    this.el.style.top = `${top}px`;
    this.el.style.transform = '';
  }
}

export { PowerUp };

export default PowerUp;
