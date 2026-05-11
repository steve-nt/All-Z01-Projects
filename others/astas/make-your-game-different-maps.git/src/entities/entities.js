/**
 * @file entities.js
 * @module entities
 * @description
 * Provides low-level DOM manipulation helpers for game entities (fish, turtles, etc.).
 * These functions are shared between multiple gameplay systems (spawning, gameLoop)
 * and handle element positioning and cleanup in the scene graph.
 *
 * Exports:
 * - `positionElement(entity)` → update element transform.
 * - `removeEntity(entity)` → detach entity from DOM.
 * - `clearEntities(entities)` → batch-remove all.
 */

/**
 * Applies position and scaling transforms to a given entity’s DOM element.
 * Used every frame by the render() step in `gameLoop.js`.
 */
export function positionElement(entity) {
    const tx = entity.x - entity.width / 2;
    const ty = entity.y - entity.height / 2;
    const sx = entity.flip || 1;
    const scale = entity.scale || 1;

    entity.el.style.transform = `translate3d(${tx}px, ${ty}px, 0) scale(${scale}) scaleX(${sx})`;
}

/**
 * Safely removes an entity’s DOM element from its parent layer.
 * @param {object} entity - Entity object with `el` property.
 */
export function removeEntity(entity) {
    if (entity.el && entity.el.parentNode) {
        entity.el.parentNode.removeChild(entity.el);
    }
}

/**
 * Clears and removes all entity DOM nodes from the provided array.
 * Also empties the array to reset state.
 */
export function clearEntities(entities) {
    for (const entity of entities) {
        removeEntity(entity);
    }
    entities.length = 0;
}
