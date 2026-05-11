/**
 * @file turtle.js
 * @module turtle
 * @description
 * Manages the spawning and behavior of turtle “hazard” entities.
 * Turtles act as obstacles — they can appear suddenly near the center
 * and end the player’s combo if hit or escaped. Controlled via
 * `TurtleManager` used in gameLoop and shooting systems.
 */

export class TurtleManager {
    constructor() {
        this.turtleSpawnedThisLevel = false;
        this.turtleActive = false;
        this.spawnTimer = 0;
        this.spawnDelay = 0;
        this.readyToSpawn = false;
    }

    /** Reset all internal state between runs or levels */
    reset() {
        this.turtleSpawnedThisLevel = false;
        this.turtleActive = false;
        this.spawnTimer = 0;
        this.spawnDelay = 0;
        this.readyToSpawn = false;
    }

    resetForNewLevel() {
        this.reset();
    }

    /**
     * Advance turtle spawn timer; when elapsed, marks readyToSpawn.
     * Called every frame from `gameLoop.update()`.
     */
    updateSpawnTimer(dt, rand) {
        if (this.turtleSpawnedThisLevel || this.turtleActive) return false;
        if (this.spawnDelay === 0) {
            // Random delay between 5–12 s
            this.spawnDelay = 5 + rand() * 7;
        }
        this.spawnTimer += dt;
        if (this.spawnTimer >= this.spawnDelay) {
            this.readyToSpawn = true;
            return true;
        }
        return false;
    }

    markSpawnedThisLevel() {
        this.turtleSpawnedThisLevel = true;
    }

    setActive(active) {
        this.turtleActive = active;
    }

    /**
     * Creates and inserts a new turtle entity into the DOM and entity list.
     * @returns {object} Spawned turtle entity
     */
    spawnTurtle(entitiesLayer, entities, nextEntityId, WORLD, PARROT_BASE_SPEED, rand) {
        // Turtle “ambush” behavior: spawn near center and dash diagonally
        const speed = PARROT_BASE_SPEED * 1.0;
        const diag = Math.SQRT1_2; // 1/sqrt(2) for even diagonal speed
        let x, y, vx, vy;

        // Always ambush in this version
        x = WORLD.width / 2 + (rand() - 0.5) * 100;
        y = WORLD.height / 2 + (rand() - 0.5) * 100;
        const dx = rand() < 0.5 ? 1 : -1;
        const dy = rand() < 0.5 ? 1 : -1;
        vx = dx * speed * diag;
        vy = dy * speed * diag;

        // Build DOM element
        const el = document.createElement('div');
        el.className = 'sprite turtle';
        el.style.opacity = '0';
        el.style.transition = 'opacity 160ms ease-out';
        entitiesLayer.appendChild(el);

        // Slightly larger sprite for emphasis
        const width = 150;
        const height = 120;
        el.style.width = width + 'px';
        el.style.height = height + 'px';

        const entity = {
            id: nextEntityId,
            type: 'fish',
            variant: 'turtle',
            isTurtle: true,
            x,
            y,
            vx,
            vy,
            width,
            height,
            alive: true,
            dead: false,
            capturing: false,
            scale: 1,
            flip: vx >= 0 ? -1 : 1, // Flip orientation for direction
            el,
        };
        entities.push(entity);
        this.turtleActive = true;

        // Small pop/splash animation when spawning
        setTimeout(() => {
            el.style.opacity = '1';

            const splash = document.createElement('div');
            splash.className = 'turtle-splash';
            splash.style.left = `${Math.max(0, Math.min(WORLD.width - 200, entity.x - 100))}px`;
            splash.style.top = `${Math.max(0, Math.min(WORLD.height - 200, entity.y - 100))}px`;
            entitiesLayer.appendChild(splash);

            // Auto-remove after animation completes
            setTimeout(() => {
                if (splash.parentNode) splash.parentNode.removeChild(splash);
            }, 520);
        }, 20);

        return entity;
    }

    /** Called when turtle is hit — resets its active flag and combo */
    handleTurtleHit(comboResetCallback) {
        this.turtleActive = false;
        if (typeof comboResetCallback === 'function') comboResetCallback();
    }

    /** Called when turtle leaves the screen */
    handleTurtleEscaped() {
        this.turtleActive = false;
    }
}
