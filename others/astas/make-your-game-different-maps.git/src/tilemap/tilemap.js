/**
 * @file tilemap.js
 * @module tilemap
 * @description
 * Tile map engine for rendering tile-based maps from a tileset.
 * Supports efficient rendering with viewport culling and multiple map layouts.
 */

/**
 * Creates a tile map object with grid-based tile data.
 * @param {number} columns - Number of columns in the grid
 * @param {number} rows - Number of rows in the grid
 * @param {number[]} tiles - Array of tile IDs (row-major order)
 * @returns {object} Map object with properties and helper methods
 */
export function createMap(columns, rows, tiles) {
    const size = columns * rows;

    return {
        columns,
        rows,
        size,
        tiles,
        /**
         * Gets the tile ID at the specified column and row.
         * @param {number} col - Column index (0-based)
         * @param {number} row - Row index (0-based)
         * @returns {number} Tile ID, or 0 if out of bounds
         */
        getTile(col, row) {
            if (col < 0 || col >= this.columns || row < 0 || row >= this.rows) {
                return 0;
            }
            return this.tiles[row * this.columns + col] || 0;
        },
        /**
         * Sets the tile ID at the specified column and row.
         * @param {number} col - Column index (0-based)
         * @param {number} row - Row index (0-based)
         * @param {number} tileId - Tile ID to set
         */
        setTile(col, row, tileId) {
            if (col >= 0 && col < this.columns && row >= 0 && row < this.rows) {
                this.tiles[row * this.columns + col] = tileId;
            }
        }
    };
}

/**
 * Tile map renderer that efficiently renders tiles from a tileset.
 */
export class TileMapRenderer {
    /**
     * @param {HTMLImageElement} tilesetImage - The tileset image containing all tiles
     * @param {number} tileSize - Size of each tile in pixels (assumes square tiles)
     * @param {number} tilesPerRow - Number of tiles per row in the tileset
     * @param {HTMLElement} container - DOM element to render tiles into
     */
    constructor(tilesetImage, tileSize, tilesPerRow, container) {
        this.tilesetImage = tilesetImage;
        this.tileSize = tileSize;
        this.tilesPerRow = tilesPerRow;
        this.container = container;
        this.currentMap = null;
        this.tileElements = new Map(); // Cache of DOM elements for each tile position
        this.viewport = { x: 0, y: 0, width: 0, height: 0 };
        this.mapDepthLevel = 0; // 0 = transparent, 1 = light, 2 = medium, 3 = dark
    }

    /**
     * Sets the current map to render.
     * @param {object} map - Map object created with createMap()
     */
    setMap(map, depthLevel = 0) {
        this.currentMap = map;
        this.mapDepthLevel = depthLevel;
        this.clear();
        this.render();

        const depthNames = ['Transparent', 'Light Depth', 'Medium Depth', 'Deep Waters'];
        console.log(`[Tilemap] Loaded map: ${map.name || 'Map'} (${depthNames[depthLevel]})`);
    }

    /**
     * Updates the viewport for culling optimization.
     * @param {number} x - Viewport X position
     * @param {number} y - Viewport Y position
     * @param {number} width - Viewport width
     * @param {number} height - Viewport height
     */
    setViewport(x, y, width, height) {
        this.viewport = { x, y, width, height };
    }

    /**
     * Clears all rendered tiles from the container.
     */
    clear() {
        // Remove all existing tile elements
        this.tileElements.forEach((element) => {
            if (element.parentNode) {
                element.parentNode.removeChild(element);
            }
        });
        this.tileElements.clear();
    }

    /**
     * Renders the current map, only drawing visible tiles for performance.
     */
    render() {
        if (!this.currentMap) {
            // Don't log warning - this is normal during initialization
            return;
        }
        if (!this.tilesetImage || !this.tilesetImage.complete) return;

        const map = this.currentMap;
        const tileSize = this.tileSize;

        // Calculate visible tile range based on viewport
        // Don't add extra tiles outside viewport to prevent artifacts
        const startCol = Math.max(0, Math.floor(this.viewport.x / tileSize));
        const endCol = Math.min(map.columns, Math.ceil((this.viewport.x + this.viewport.width) / tileSize));
        const startRow = Math.max(0, Math.floor(this.viewport.y / tileSize));
        const endRow = Math.min(map.rows, Math.ceil((this.viewport.y + this.viewport.height) / tileSize));

        // Render only visible tiles
        for (let row = startRow; row < endRow; row++) {
            for (let col = startCol; col < endCol; col++) {
                const tileId = map.getTile(col, row);

                // Skip empty tiles only (0)
                if (tileId === 0) continue;

                const key = `${col},${row}`;
                let tileElement = this.tileElements.get(key);

                if (!tileElement) {
                    // Create new tile element
                    tileElement = document.createElement('div');
                    tileElement.className = 'tile';
                    tileElement.style.position = 'absolute';
                    tileElement.style.width = `${tileSize}px`;
                    tileElement.style.height = `${tileSize}px`;
                    tileElement.style.imageRendering = 'pixelated';
                    tileElement.style.pointerEvents = 'none';
                    tileElement.style.zIndex = '-1';
                    this.container.appendChild(tileElement);
                    this.tileElements.set(key, tileElement);
                }

                // Position tile
                const tileLeft = col * tileSize;
                const tileTop = row * tileSize;
                tileElement.style.left = `${tileLeft}px`;
                tileElement.style.top = `${tileTop}px`;

                // Apply depth-based styling
                this.applyDepthStyle(tileElement, tileId);

                tileElement.style.display = 'block';
            }
        }

        // Clean up tiles outside viewport
        this.tileElements.forEach((element, key) => {
            const [col, row] = key.split(',').map(Number);
            if (col < startCol || col >= endCol || row < startRow || row >= endRow) {
                if (element.parentNode) {
                    element.parentNode.removeChild(element);
                }
                this.tileElements.delete(key);
            }
        });
    }


    /**
     * Applies visual styling based on the map's depth level
     */
    applyDepthStyle(tileElement, tileId) {
        // Calculate source position in tileset
        const tilesetCol = (tileId - 1) % this.tilesPerRow;
        const tilesetRow = Math.floor((tileId - 1) / this.tilesPerRow);
        const sourceX = tilesetCol * this.tileSize;
        const sourceY = tilesetRow * this.tileSize;

        // Apply tileset image with background-position
        tileElement.style.backgroundImage = `url("${this.tilesetImage.src}")`;
        tileElement.style.backgroundPosition = `-${sourceX}px -${sourceY}px`;

        const tilesetRows = Math.ceil(6 / this.tilesPerRow);
        tileElement.style.backgroundSize = `${this.tilesPerRow * this.tileSize}px ${tilesetRows * this.tileSize}px`;
        tileElement.style.backgroundRepeat = 'no-repeat';

        // Apply depth-based tinting
        switch (this.mapDepthLevel) {
            case 0: // Levels 1-3: Transparent with very light blue tint
                tileElement.style.opacity = '0.65';
                tileElement.style.filter = 'sepia(0.2) hue-rotate(180deg) saturate(1.1) brightness(1.2)';
                break;

            case 1: // Levels 4-6: Transparent with light blue tint
                tileElement.style.opacity = '0.45';
                tileElement.style.filter = 'sepia(0.3) hue-rotate(185deg) saturate(1.3) brightness(1.0)';
                break;

            case 2: // Levels 7-9: Transparent with medium blue tint
                tileElement.style.opacity = '0.25';
                tileElement.style.filter = 'sepia(0.4) hue-rotate(195deg) saturate(1.5) brightness(0.85)';
                break;

            case 3: // Levels 10+: Transparent with dark blue tint
                tileElement.style.opacity = '0';
                tileElement.style.filter = 'sepia(0.5) hue-rotate(205deg) saturate(1.7) brightness(0.7)';
                break;

            default:
                tileElement.style.opacity = '0.75';
                tileElement.style.filter = 'sepia(0.2) hue-rotate(180deg) saturate(1.1) brightness(1.2)';
        }
    }
    update() {
        // Re-render tiles when viewport changes (e.g., on window resize)
        // This ensures tiles cover the entire visible area after resize

        this.render();
    }
}


/**
 * Loads a tileset image and returns a promise that resolves when loaded.
 * @param {string} tilesetPath - Path to the tileset image
 * @returns {Promise<HTMLImageElement>} Promise that resolves with the loaded image
 */
export function loadTileset(tilesetPath) {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = () => reject(new Error(`Failed to load tileset: ${tilesetPath}`));
        img.src = tilesetPath;
    });
}

