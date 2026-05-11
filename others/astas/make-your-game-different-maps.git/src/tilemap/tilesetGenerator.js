/**
 * @file tilesetGenerator.js
 * @module tilesetGenerator
 * @description
 * Generates a tileset image by combining individual tile images into a single image.
 * This improves performance and memory usage by using a single image resource.
 */

/**
 * Creates a tileset image from individual tile images.
 * @param {string[]} tilePaths - Array of paths to individual tile images
 * @param {number} tilesPerRow - Number of tiles per row in the tileset
 * @param {number} tileSize - Size of each tile in pixels (assumes square tiles)
 * @returns {Promise<HTMLImageElement>} Promise that resolves with the generated tileset image
 */
export async function generateTileset(tilePaths, tilesPerRow, tileSize) {
    return new Promise((resolve, reject) => {
        // Calculate tileset dimensions
        const totalTiles = tilePaths.length;
        const tilesetRows = Math.ceil(totalTiles / tilesPerRow);
        const tilesetWidth = tilesPerRow * tileSize;
        const tilesetHeight = tilesetRows * tileSize;

        // Create canvas for tileset
        const canvas = document.createElement('canvas');
        canvas.width = tilesetWidth;
        canvas.height = tilesetHeight;
        const ctx = canvas.getContext('2d');

        // Load all tile images
        const imagePromises = tilePaths.map((path) => {
            return new Promise((resolveImg, rejectImg) => {
                const img = new Image();
                img.onload = () => resolveImg(img);
                img.onerror = () => rejectImg(new Error(`Failed to load tile: ${path}`));
                img.src = path;
            });
        });

        Promise.all(imagePromises)
            .then((images) => {
                // Draw each tile onto the canvas
                images.forEach((img, index) => {
                    const col = index % tilesPerRow;
                    const row = Math.floor(index / tilesPerRow);
                    const x = col * tileSize;
                    const y = row * tileSize;

                    // Draw the tile image, scaling it to fit the tile size
                    ctx.drawImage(img, x, y, tileSize, tileSize);
                });

                // Convert canvas to image
                const tilesetImage = new Image();
                tilesetImage.onload = () => resolve(tilesetImage);
                tilesetImage.onerror = () => reject(new Error('Failed to create tileset image'));
                tilesetImage.src = canvas.toDataURL('image/png');
            })
            .catch(reject);
    });
}

/**
 * Creates solid colored tiles programmatically.
 * @param {number} tileSize - Size of each tile in pixels
 * @returns {Promise<HTMLImageElement>} Promise that resolves with the tileset image
 */
function createSolidColorTileset(tileSize = 64) {
    return new Promise((resolve) => {
        const tilesPerRow = 3;
        const colors = [
            { r: 64, g: 164, b: 223, name: 'Light Water' },      // Tile 1: Light blue water
            { r: 41, g: 128, b: 185, name: 'Deep Water' },       // Tile 2: Darker blue water
            { r: 180, g: 100, b: 50, name: 'Sand' },             // Tile 3: Sandy brown
            { r: 100, g: 100, b: 100, name: 'Rock' },            // Tile 4: Gray rock
            { r: 255, g: 107, b: 53, name: 'Coral' },            // Tile 5: Orange coral
            { r: 46, g: 125, b: 50, name: 'Seaweed' }            // Tile 6: Green seaweed
        ];

        const tilesetRows = Math.ceil(colors.length / tilesPerRow);
        const tilesetWidth = tilesPerRow * tileSize;
        const tilesetHeight = tilesetRows * tileSize;

        const canvas = document.createElement('canvas');
        canvas.width = tilesetWidth;
        canvas.height = tilesetHeight;
        const ctx = canvas.getContext('2d');

        // Clear canvas first
        ctx.clearRect(0, 0, tilesetWidth, tilesetHeight);

        // Draw each colored tile
        colors.forEach((color, index) => {
            const col = index % tilesPerRow;
            const row = Math.floor(index / tilesPerRow);
            const x = col * tileSize;
            const y = row * tileSize;

            // Fill with solid color - make it more vibrant
            ctx.fillStyle = `rgb(${color.r}, ${color.g}, ${color.b})`;
            ctx.fillRect(x, y, tileSize, tileSize);

            // Add a darker border for better definition
            ctx.strokeStyle = `rgba(0, 0, 0, 0.4)`;
            ctx.lineWidth = 2;
            ctx.strokeRect(x + 1, y + 1, tileSize - 2, tileSize - 2);
            
            // Add a subtle highlight for depth
            ctx.strokeStyle = `rgba(255, 255, 255, 0.2)`;
            ctx.lineWidth = 1;
            ctx.strokeRect(x, y, tileSize, tileSize);
        });

        // Convert canvas to image
        const tilesetImage = new Image();
        tilesetImage.onload = () => {
            resolve(tilesetImage);
        };
        tilesetImage.onerror = () => {
            reject(new Error('Failed to create tileset image from canvas'));
        };
        // Data URLs cannot have query parameters - remove the timestamp
        tilesetImage.src = canvas.toDataURL('image/png');
    });
}

/**
 * Creates a tileset using programmatically generated solid colored tiles.
 * This ensures consistent tile appearance regardless of available image files.
 * @param {number} tileSize - Size of each tile in pixels
 * @returns {Promise<HTMLImageElement>} Promise that resolves with the tileset image
 */
export async function createGameTileset(tileSize = 64) {
    // Use programmatically generated solid colored tiles
    // This ensures we get proper colored tiles, not fish images
    return createSolidColorTileset(tileSize);
}

