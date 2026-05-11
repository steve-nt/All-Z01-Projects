/**
 * @file spawning.js
 * @module spawning
 * @description
 * Provides helper functions for spawning decorative and gameplay entities:
 * - Ambient bubbles (purely visual)
 * - Main fish entities (including life fish)
 *
 * Works closely with:
 * - `entities.js` (for positioning and DOM manipulation)
 * - `lifeFish.js` (for spawning bonus fish)
 * - `gameLoop.js` (which calls these spawn functions)
 */

import { positionElement } from '../entities/entities.js';

/**
 * Spawns a decorative bubble element that floats upward and disappears.
 * @param {Object} config - Spawn configuration.
 * @param {HTMLElement} config.bubblesLayer - Layer to append the bubble to.
 * @param {Object} config.world - World dimensions (width, height).
 * @param {Function} config.rand - RNG function.
 */
export function spawnBubble({ bubblesLayer, world, rand }) {
    if (!bubblesLayer) return;

    const bubble = document.createElement('div');
    const sizeClass = rand() < 0.3 ? 'big' : rand() < 0.5 ? 'small' : '';
    bubble.className = `bubble ${sizeClass}`.trim();

    // Randomize x position and rise speed
    const x = 20 + rand() * (world.width - 40);
    bubble.style.left = `${x}px`;
    // Start near the bottom so bubble rises into view
    bubble.style.bottom = '-12px';
    const duration = 4 + rand() * 4;
    bubble.style.animationDuration = `${duration}s`;
    bubblesLayer.appendChild(bubble);

    // Clean up after animation ends
    setTimeout(() => {
        if (bubble.parentNode) bubble.parentNode.removeChild(bubble);
    }, (duration + 0.2) * 1000);
}

/**
 * Spawns a new fish (or life fish) entity when spawn timer expires.
 * Dynamically scales difficulty based on current level.
 */
export function spawnParrot({
    lifeFishManager,
    lives,
    maxLives,
    rand,
    entitiesLayer,
    entities,
    nextEntityId,
    world,
    parrotBaseSpeed,
    baseMaxFish,
    level,
    showCenterNotification,
}) {
    // 1️⃣ Attempt to spawn the rare life fish
    if (lifeFishManager.shouldSpawnLifeFish(lives, maxLives, rand)) {
        if (rand() < lifeFishManager.getSpawnChance()) {
            const entity = lifeFishManager.spawnLifeFish(
                entitiesLayer,
                entities,
                nextEntityId,
                world,
                parrotBaseSpeed,
                rand,
                showCenterNotification
            );
            positionElement(entity);
            lifeFishManager.markSpawnedThisLevel();
            return { nextEntityId };
        }
    }

    // 2️⃣ Regular fish spawning (limited by max active fish)
    const maxFish = baseMaxFish + Math.floor((level - 1) * 0.6);
    const activeFish = entities.filter(e => e.type === 'fish' && e.alive).length;
    if (activeFish >= maxFish) return { nextEntityId };

    // Spawn properties
    const enterFromLeft = rand() < 0.6;
    const hudHeight = 80;
    const y = hudHeight + 40 + rand() * ((world.height - hudHeight) * 0.7);
    const speedScale = 1 + (level - 1) * 0.12;
    const speed = parrotBaseSpeed * speedScale * (0.8 + rand() * 0.8);
    const vx = enterFromLeft ? speed : -speed;
    const x = enterFromLeft ? -60 : world.width + 60;

    // 3️⃣ Choose fish variant set based on level progression
    let fishTypes = [];
    if (level <= 2) {
        fishTypes = ['fish1', 'fish2', 'fish3'];
    } else if (level <= 4) {
        fishTypes = ['fish1', 'fish2', 'fish3', 'fish4'];
    } else if (level <= 6) {
        fishTypes = ['fish1', 'fish2', 'fish3', 'fish4', 'fish5'];
    } else {
        fishTypes = ['fish1', 'fish2', 'fish3', 'fish4', 'fish5', 'fish6', 'fish7'];
    }

    const variant = fishTypes[Math.floor(rand() * fishTypes.length)];
    const el = document.createElement('div');
    el.className = `sprite fish ${variant}`;
    entitiesLayer.appendChild(el);

    // Register new entity in the world state
    const entity = {
        id: nextEntityId++,
        type: 'fish',
        variant,
        x,
        y,
        vx,
        vy: 0,
        width: 80,
        height: 60,
        alive: true,
        dead: false,
        capturing: false,
        scale: 1,
        flip: vx >= 0 ? 1 : -1,
        el,
    };

    entities.push(entity);
    positionElement(entity);
    return { nextEntityId };
}
