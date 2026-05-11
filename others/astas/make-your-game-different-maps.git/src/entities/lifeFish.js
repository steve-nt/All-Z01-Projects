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
        const speed = PARROT_BASE_SPEED * 3.2; // Faster than normal fish - naughty and quick!
        const vx = enterFromLeft ? speed : -speed;
        const x = enterFromLeft ? -60 : WORLD.width + 60;

        // Create DOM node
        const el = document.createElement('div');
        el.className = 'sprite lifeFish';
        el.style.width = '80px';
        el.style.height = '60px';
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
            // Movement parameters for unpredictable behavior
            directionChangeTimer: 0,
            directionChangeInterval: 0.3 + rand() * 0.4, // Change direction every 0.3-0.7 seconds (more frequent)
            verticalDrift: (rand() - 0.5) * 60, // Random vertical drift speed (increased range)
            speedBoostTimer: 0,
            speedBoostInterval: 0.8 + rand() * 0.6, // Speed boost every 0.8-1.4 seconds (more frequent)
            baseSpeed: speed,
            originalVx: vx, // Store original direction for reference
            goalDirection: vx > 0 ? 1 : -1, // Goal: reach opposite side (keep original direction as goal)
            currentHorizontalDirection: vx > 0 ? 1 : -1, // Track current horizontal direction
            horizontalReverseTimer: 0,
            horizontalReverseInterval: 1.5 + rand() * 1.0, // Reverse horizontal direction every 1.5-2.5 seconds
            reverseDuration: 0, // How long to stay reversed
            maxReverseDuration: 0.8 + rand() * 0.6, // Max time to stay reversed (0.8-1.4s)
            lastProgressCheck: 0,
            progressCheckInterval: 2.0, // Check progress every 2 seconds
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
     * Applies unpredictable "naughty" movement patterns to make the life fish
     * harder to catch. Includes frequent direction changes, temporary horizontal reversals,
     * vertical drift, and speed boosts. Ensures fish eventually reaches opposite side.
     */
    updateLifeFishMovement(entity, dt) {
        if (!entity.isLifeFish || !entity.alive || entity.dead) return;

        entity.directionChangeTimer += dt;
        entity.speedBoostTimer += dt;
        entity.horizontalReverseTimer += dt;
        entity.lastProgressCheck += dt;

        // 1. Frequent unpredictable vertical direction changes (zigzag pattern)
        if (entity.directionChangeTimer >= entity.directionChangeInterval) {
            // Randomly change vertical drift direction with wider range
            entity.verticalDrift = (Math.random() - 0.5) * 80; // -40 to +40 pixels/sec (increased)
            entity.directionChangeTimer = 0;
            entity.directionChangeInterval = 0.3 + Math.random() * 0.4; // Next change in 0.3-0.7s (more frequent)
        }

        // 2. Temporary horizontal direction reversals (fish can turn around, but returns to goal!)
        if (entity.reverseDuration > 0) {
            // Currently reversed - count down
            entity.reverseDuration -= dt;
            if (entity.reverseDuration <= 0) {
                // Time to return to goal direction
                entity.currentHorizontalDirection = entity.goalDirection;
                entity.horizontalReverseTimer = 0;
                entity.horizontalReverseInterval = 1.5 + Math.random() * 1.0; // Next reversal in 1.5-2.5s
            }
        } else if (entity.horizontalReverseTimer >= entity.horizontalReverseInterval) {
            // Time for a temporary reversal
            entity.currentHorizontalDirection = -entity.goalDirection; // Reverse from goal
            entity.reverseDuration = entity.maxReverseDuration;
            entity.maxReverseDuration = 0.8 + Math.random() * 0.6; // Next reversal duration 0.8-1.4s
        }

        // 3. Progress check - ensure fish is making progress toward goal
        if (entity.lastProgressCheck >= entity.progressCheckInterval) {
            // If fish has been going wrong direction too long, force it back to goal
            if (entity.currentHorizontalDirection !== entity.goalDirection && entity.reverseDuration <= 0.2) {
                entity.currentHorizontalDirection = entity.goalDirection;
                entity.reverseDuration = 0;
            }
            entity.lastProgressCheck = 0;
        }

        // 4. Occasional speed boosts (makes it dart away)
        let currentSpeed = entity.baseSpeed;
        if (entity.speedBoostTimer >= entity.speedBoostInterval) {
            // Speed boost for a short duration
            currentSpeed = entity.baseSpeed * 1.8; // 80% faster (increased from 60%)
            if (entity.speedBoostTimer >= entity.speedBoostInterval + 0.25) {
                // Reset after boost duration (slightly shorter)
                entity.speedBoostTimer = 0;
                entity.speedBoostInterval = 0.8 + Math.random() * 0.6; // Next boost in 0.8-1.4s (more frequent)
            }
        }

        // 5. Apply movement with current horizontal direction and vertical drift
        entity.vx = entity.currentHorizontalDirection * currentSpeed;
        entity.vy = entity.verticalDrift;

        // Update flip direction based on movement
        entity.flip = entity.vx >= 0 ? 1 : -1;
    }
}
