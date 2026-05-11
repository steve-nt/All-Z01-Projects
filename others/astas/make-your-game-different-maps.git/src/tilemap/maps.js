/**
 * @file maps.js
 * @module maps
 * @description
 * Defines multiple tile map layouts for the game.
 * Each map has a unique pattern and layout.
 */

import { createMap } from './tilemap.js';

/**
 * Map 1: Ocean Floor Pattern
 * Features a varied ocean floor with different colored tiles creating a natural pattern.
 */
export function createOceanFloorMap() {
    const columns = 30;
    const rows = 20;
    // Generate tiles programmatically for larger map
    const tiles = [];
    for (let row = 0; row < rows; row++) {
        for (let col = 0; col < columns; col++) {
            // Create a repeating pattern that shifts each row
            const patternIndex = (col + row) % 5;
            tiles.push(patternIndex + 1); // Tile IDs 1-5
        }
    }

    return createMap(columns, rows, tiles);
}

/**
 * Map 2: Checkerboard Pattern
 * Features an alternating checkerboard pattern with different colors.
 */
export function createCheckerboardMap() {
    const columns = 30;
    const rows = 20;
    const tiles = [];
    
    for (let row = 0; row < rows; row++) {
        for (let col = 0; col < columns; col++) {
            const isEven = (row + col) % 2 === 0;
            // Alternate between different tile types
            if (isEven) {
                tiles.push((row % 3) + 1); // Cycle through tiles 1-3
            } else {
                tiles.push((col % 3) + 4); // Cycle through tiles 4-6
            }
        }
    }

    return createMap(columns, rows, tiles);
}

/**
 * Map 3: Coral Reef Pattern
 * Features a more organic pattern with clusters and varied tile distribution.
 */
export function createCoralReefMap() {
    const columns = 30;
    const rows = 20;
    const tiles = [];
    
    // Create a more organic, clustered pattern
    for (let row = 0; row < rows; row++) {
        for (let col = 0; col < columns; col++) {
            // Create clusters using distance-based patterns
            const centerX1 = 8;
            const centerY1 = 5;
            const centerX2 = 22;
            const centerY2 = 13;
            const centerX3 = 15;
            const centerY3 = 10;
            
            const dist1 = Math.sqrt((col - centerX1) ** 2 + (row - centerY1) ** 2);
            const dist2 = Math.sqrt((col - centerX2) ** 2 + (row - centerY2) ** 2);
            const dist3 = Math.sqrt((col - centerX3) ** 2 + (row - centerY3) ** 2);
            
            let tileId;
            if (dist1 < 3) {
                tileId = 1; // Blue cluster
            } else if (dist2 < 3) {
                tileId = 2; // Brown cluster
            } else if (dist3 < 2.5) {
                tileId = 3; // Green cluster
            } else if (row < 3) {
                tileId = 4; // Red for top area
            } else if (row > rows - 3) {
                tileId = 5; // Yellow for bottom area
            } else {
                // Create a wave pattern
                tileId = ((col + row * 2) % 6) + 1;
            }
            
            tiles.push(tileId);
        }
    }

    return createMap(columns, rows, tiles);
}

/**
 * Map 4: Striped Pattern
 * Features horizontal and vertical stripes for a different visual style.
 */
export function createStripedMap() {
    const columns = 30;
    const rows = 20;
    const tiles = [];
    
    for (let row = 0; row < rows; row++) {
        for (let col = 0; col < columns; col++) {
            // Create horizontal stripes
            if (row % 3 === 0) {
                tiles.push(1); // Blue stripe
            } else if (row % 3 === 1) {
                tiles.push(2); // Brown stripe
            } else {
                // Vertical pattern in this row
                tiles.push((col % 3) + 3); // Green, Red, Yellow alternating
            }
        }
    }

    return createMap(columns, rows, tiles);
}

/**
 * Gets a map by index (for cycling through maps).
 * @param {number} index - Map index (0-3)
 * @returns {object} Map object
 */
export function getMapByIndex(index) {
    const maps = [
        createOceanFloorMap,
        createCheckerboardMap,
        createCoralReefMap,
        createStripedMap,
    ];
    
    const mapIndex = index % maps.length;
    return maps[mapIndex]();
}

