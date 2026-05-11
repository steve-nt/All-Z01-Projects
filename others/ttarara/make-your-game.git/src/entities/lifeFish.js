/**
 * @file lifeFish.js
 * @module lifeFish
 * @description
 * Defines the `LifeFishManager` class which controls spawning and
 * behavior of the special “life fish” — a rare entity that grants
 * an extra life when caught. Handles spawn timing, movement patterns,
 * and interaction callbacks.
 *
 * Dependencies:
 * - Used by `gameLoop.js` (spawn/update logic)
 * - Used by `shooting.js` (life gain handling)
 */

export class LifeFishManager {
    constructor() {
        // Runtime flags controlling spawn state
        this.lifeFishSpawned = false;
        this.lifeFishActive = false;
        this.lifeFishSpawnedThisLevel = false;
        this.lifeFishSpawnTimer = 0;
        this.lifeFishSpawnDelay = 0;
    }

    /** Reset all internal life fish state (used between runs) */
    reset() {
        this.lifeFishSpawned = false;
        this.lifeFishActive = false;
        this.lifeFishSpawnedThisLevel = false;
        this.lifeFishSpawnTimer = 0;
        this.lifeFishSpawnDelay = 0;
    }

    /** Reset level-specific state without affecting global flags */
    resetForNewLevel() {
        this.lifeFishSpawnedThisLevel = false;
        this.lifeFishSpawnTimer = 0;
        this.lifeFishSpawnDelay = 0;
    }

    /**
     * Advance internal timer; determines when the life fish should appear.
     * Returns `true` if spawn conditions are met.
     */
    updateSpawnTimer(dt, lives, maxLives, rand) {
        // Only trigger if player lost a life and one hasn't appeared this level
        if (lives < maxLives && !this.lifeFishSpawnedThisLevel && !this.lifeFishActive) {
            if (this.lifeFishSpawnDelay === 0) {
                // Randomize delay window (3–8 s)
                this.lifeFishSpawnDelay = 3 + rand() * 5;
            }

            this.lifeFishSpawnTimer += dt;

            if (this.lifeFishSpawnTimer >= this.lifeFishSpawnDelay) {
                return true; // Ready to spawn
            }
        }
        return false;
    }

    /** Check if conditions are right for spawning a life fish */
    shouldSpawnLifeFish(lives, maxLives, rand) {
        return lives < maxLives && !this.lifeFishSpawnedThisLevel && !this.lifeFishActive;
    }

    /** Fixed probability of actually spawning once conditions met */
    getSpawnChance() {
        return 0.3; // 30 % chance
    }

    /** Mark the life fish as having appeared this level */
    markSpawnedThisLevel() {
        this.lifeFishSpawnedThisLevel = true;
    }

    /** Update whether a life fish is currently active */
    setActive(active) {
        this.lifeFishActive = active;
    }

    /** Check whether a life fish is currently active in the scene */
    isActive() {
        return this.lifeFishActive;
    }

    /**
     * Create and register a new life fish entity DOM node.
     * @returns {object} Entity object pushed to `entities`.
     */
    spawnLifeFish(entitiesLayer, entities, nextEntityId, WORLD, PARROT_BASE_SPEED, rand, showCenterNotification) {
        const enterFromLeft = rand() < 0.5;
        const hudHeight = 80;

        // Spawn around mid-screen area to make it visible
        const y = hudHeight + 80 + rand() * ((WORLD.height - hudHeight) * 0.5);
        const speed = PARROT_BASE_SPEED * 1.8; // Faster than normal fish
        const vx = enterFromLeft ? speed : -speed;
        const x = enterFromLeft ? -60 : WORLD.width + 60;

        // Create DOM node
        const el = document.createElement('div');
        el.className = 'sprite fish lifeFish';
        entitiesLayer.appendChild(el);

        const entity = {
            id: nextEntityId++,
            type: 'fish',
            variant: 'lifeFish',
            isLifeFish: true,
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
            // Additional motion parameters for complex movement
            originalVx: vx,
            movementTimer: 0,
            pulseSpeed: 1.5 + rand() * 0.5,
            waveAmplitude: 30 + rand() * 20,
            waveFrequency: 0.5 + rand() * 0.3,
            zigzagTimer: 0,
            zigzagDirection: 1,
            scaleVariation: 0.2 + rand() * 0.1,
        };

        entities.push(entity);
        this.lifeFishActive = true;

        // Display alert message to player
        showCenterNotification('Catch the red fish', 'level-up', 2500);
        return entity;
    }

    /** Called when a life fish is caught — awards +1 life */
    handleLifeFishCaught(lives, maxLives, showCenterNotification) {
        const newLives = Math.min(maxLives, lives + 1);
        this.lifeFishActive = false;
        showCenterNotification('💖 +1 LIFE!', 'level-up', 2000);
        return newLives;
    }

    /** Called when the life fish leaves the screen without being caught */
    handleLifeFishEscaped() {
        this.lifeFishActive = false;
        // No penalty for missing it
    }

    /**
     * Applies the life fish’s special animated movement pattern.
     * Includes sine wave, zigzag drift, pulsing scale, and varying speed.
     */
    updateLifeFishMovement(entity, dt) {
        if (!entity.isLifeFish || !entity.alive || entity.dead) return;

        entity.movementTimer += dt;
        entity.zigzagTimer += dt;

        const time = entity.movementTimer;

        // 1. Vertical sine-wave motion
        const waveY = Math.sin(time * entity.waveFrequency) * entity.waveAmplitude;

        // 2. Horizontal zigzag (direction changes every 1–2 s)
        if (entity.zigzagTimer > 1 + Math.random() * 1) {
            entity.zigzagDirection *= -1;
            entity.zigzagTimer = 0;
        }
        const zigzagX = Math.sin(time * 2) * 20 * entity.zigzagDirection;

        // 3. Pulsing scale animation (“breathing”)
        const pulseScale = 1 + Math.sin(time * entity.pulseSpeed) * entity.scaleVariation;

        // 4. Speed variation synchronized with pulse
        const speedMultiplier = 0.8 + 0.4 * Math.sin(time * 1.5);

        // Apply velocities and scale
        entity.vx = entity.originalVx * speedMultiplier + zigzagX * 0.5;
        entity.vy = waveY * 0.3;
        entity.scale = pulseScale;
    }
}
