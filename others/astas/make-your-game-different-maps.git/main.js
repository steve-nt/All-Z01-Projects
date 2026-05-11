/**
 * @file main.js
 * Bootstraps Fish Tank Hunt (loop, input, lifecycle, UI, audio).
 */
import {
    GAME_DURATION_S,
    MAX_LEVEL,
    WORLD,
    PARROT_BASE_SPEED,
    BASE_SPAWN_MIN,
    BASE_SPAWN_MAX,
    BASE_MAX_FISH,
    LEVEL_DURATION,
    MAX_LIVES,
    SHOT_RATE_LIMIT_MS,
    COMBO_WINDOW,
} from './src/core/constants.js';
import { createSettingsHandlers } from './src/core/settings.js';
import { createGameLoop } from './src/core/gameLoop.js';
import { createLifecycleSystem } from './src/core/lifecycle.js';
import { createInputSystem } from './src/core/input.js';

import { positionElement, removeEntity, clearEntities } from './src/entities/entities.js';
import { LifeFishManager } from './src/entities/lifeFish.js';
import { TurtleManager } from './src/entities/turtle.js';

import { spawnBubble, spawnParrot } from './src/gameplay/spawning.js';
import { createShootingSystem } from './src/gameplay/shooting.js';

import { addHighScore, updateHighScoreListUI, isHighScore } from './src/ui/highScores.js';
import { createHudSystem } from './src/ui/hud.js';
import { createMenuEffects } from './src/ui/menuEffects.js';
import { createFeedbackSystem } from './src/ui/feedback.js';
import { getStoryManager } from './src/story/story.js';

import { TileMapRenderer } from './src/tilemap/tilemap.js';
import { createGameTileset } from './src/tilemap/tilesetGenerator.js';
import { getMapByIndex } from './src/tilemap/maps.js';

function bootstrapGame() {
    'use strict';

    // DOM refs
    const root = document.getElementById('game-root');
    const gameEl = document.getElementById('game');
    const entitiesLayer = document.getElementById('entities');
    const bubblesLayer = document.getElementById('bubbles-layer');
    const plantsLayer = document.getElementById('plants-layer');
    const tilemapLayer = document.getElementById('tilemap-layer');
    const crosshair = document.getElementById('crosshair');

    // HUD
    const hudTime = document.getElementById('hud-time');
    const hudScore = document.getElementById('hud-score');
    const hudLives = document.getElementById('hud-lives');
    const hudFps = document.getElementById('hud-fps');
    const hudLevel = document.getElementById('hud-level');

    // Pause overlay
    const pauseBtn = document.getElementById('pause-btn');
    const pauseOverlay = document.getElementById('pause-overlay');
    const continueBtn = document.getElementById('continue-btn');
    const restartBtn = document.getElementById('restart-btn');
    const menuTime = document.getElementById('menu-time');
    const menuScore = document.getElementById('menu-score');
    const menuLives = document.getElementById('menu-lives');
    const menuFps = document.getElementById('menu-fps');
    const menuLevel = document.getElementById('menu-level');

    // Coral overlay toggle
    const coralOverlay = document.getElementById('coral-overlay');
    const toggleCoralBtn = document.getElementById('toggle-coral-btn');

    // Combo & notifications
    const comboDisplay = document.getElementById('combo-display');
    const centerNotification = document.getElementById('center-notification');

    // Gameover overlay
    const gameoverOverlay = document.getElementById('gameover-overlay');
    const finalScore = document.getElementById('final-score');
    const finalAccuracy = document.getElementById('final-accuracy');
    const finalCombo = document.getElementById('final-combo');
    const finalLevel = document.getElementById('final-level');
    const finalCaught = document.getElementById('final-caught');
    const restartGameoverBtn = document.getElementById('restart-gameover-btn');
    const mainMenuGameoverBtn = document.getElementById('main-menu-gameover-btn');
    let highScore = 0;

    // Start menu & settings
    const startOverlay = document.getElementById('start-overlay');
    const startBtn = document.getElementById('start-btn');
    const instructionsBtn = document.getElementById('instructions-btn');
    const settingsBtn = document.getElementById('settings-btn');
    const creditsBtn = document.getElementById('credits-btn');
    const soundToggle = document.getElementById('sound-toggle');
    const volumeSlider = document.getElementById('volume-slider');

    // Popup windows
    const howtoPopup = document.getElementById('howto-popup');
    const settingsPopup = document.getElementById('settings-popup');
    const creditsPopup = document.getElementById('credits-popup');
    const closeHowto = document.getElementById('close-howto');
    const closeSettings = document.getElementById('close-settings');
    const closeCredits = document.getElementById('close-credits');

    const {
        initEnhancedMenu,
        celebrateHighScore,
        clearCelebration,
        createPauseBubbles,
    } = createMenuEffects({
        instructionsBtn,
        settingsBtn,
        bubblesLayer,
        centerNotification,
        world: WORLD,
        spawnBubble,
        rand,
    });

    // RNG & state
    let rngSeed = 1337;
    function rand() {
        rngSeed = (rngSeed * 1664525 + 1013904223) % 0xffffffff;
        return (rngSeed >>> 0) / 0xffffffff;
    }
    let entities = [];
    let nextEntityId = 1;
    let running = false;
    let paused = false;
    let timeLeft = GAME_DURATION_S;
    let score = 0;
    let lives = MAX_LIVES;
    let lastFrameTs = 0;
    let spawnTimer = 0;
    let bubbleTimer = 0;
    let lastShotMs = 0;
    let level = 1;
    let levelTimer = LEVEL_DURATION;

    // Combo
    let combo = 0;
    let maxCombo = 0;
    let comboTimer = 0;

    // Stats
    let totalShots = 0;
    let totalHits = 0;
    let fishCaught = 0;
    let missedShots = 0;

    // FPS
    let fps = 60;
    let fpsAccum = 0;
    let fpsFrames = 0;

    // Crosshair
    let crosshairX = window.innerWidth / 2;
    let crosshairY = window.innerHeight / 2;
    let lastCrossTs = window.performance.now();
    const keys = {};
    const CROSSHAIR_HALF = 16;
    const SPEED_PX_PER_SEC = 800;

    // Settings
    let soundEnabled = true;
    let volume = 0.8;

    // Audio: menu music
    const menuMusic = new Audio('sounds/Menu.mp3');
    menuMusic.loop = true;
    menuMusic.preload = 'auto';
    let menuMusicStarted = false;

    async function startMenuMusic() {
        if (!soundEnabled || menuMusicStarted) return;
        menuMusicStarted = true;
        try {
            menuMusic.muted = true;
            menuMusic.currentTime = 0;
            await menuMusic.play();
            setTimeout(() => {
                menuMusic.muted = false;
                menuMusic.volume = Math.max(0, Math.min(1, volume * 0.75));
            }, 400);
        } catch (e) {
            window.addEventListener('pointerdown', () => {
                try {
                    menuMusic.muted = false;
                    menuMusic.volume = Math.max(0, Math.min(1, volume * 0.75));
                    menuMusic.currentTime = 0;
                    menuMusic.play().catch(() => { });
                } catch (_) { }
            }, { once: true });
            menuMusicStarted = false;
        }
    }
    function stopMenuMusic() {
        try { menuMusic.pause(); } catch (_) { }
        menuMusicStarted = false;
    }

    // Audio: background music
    const bgMusic = new Audio('sounds/Background.mp3');
    bgMusic.loop = true;
    bgMusic.volume = 0.25;
    bgMusic.preload = 'auto';

    function startBgMusic() {
        if (!soundEnabled) return;
        try {
            bgMusic.currentTime = 0;
            bgMusic.play().catch(() => { });
        } catch (_) { }
    }
    function stopBgMusic() {
        try { bgMusic.pause(); } catch (_) { }
    }

    // Audio: gameover SFX
    const gameOverSfx = new Audio('sounds/Gameover.mp3');
    let goSfxPlayed = false;
    function resetGoSfx() { goSfxPlayed = false; }
    function playGameOverOnce() {
        if (!soundEnabled || goSfxPlayed) return;
        goSfxPlayed = true;
        try {
            gameOverSfx.currentTime = 0;
            gameOverSfx.volume = volume;
            const p = gameOverSfx.play();
            if (p && typeof p.catch === 'function') p.catch(() => { });
        } catch (_) { }
    }

    if (gameoverOverlay) {
        let wasVisible = false;
        const goObserver = new MutationObserver(() => {
            const visible =
                !gameoverOverlay.hasAttribute('hidden') &&
                getComputedStyle(gameoverOverlay).display !== 'none';
            if (visible && !wasVisible) {
                playGameOverOnce();
                stopBgMusic();
            }
            wasVisible = visible;
        });
        goObserver.observe(gameoverOverlay, {
            attributes: true,
            attributeFilter: ['hidden', 'style', 'class'],
        });
    }

    // Start menu observer
    if (startOverlay) {
        let wasVisible = false;
        const startObserver = new MutationObserver(() => {
            const visible =
                !startOverlay.hasAttribute('hidden') &&
                getComputedStyle(startOverlay).display !== 'none';
            if (visible && !wasVisible) startMenuMusic();
            if (!visible && wasVisible) stopMenuMusic();
            wasVisible = visible;
        });
        startObserver.observe(startOverlay, {
            attributes: true,
            attributeFilter: ['hidden', 'style', 'class'],
        });
        const initialVisible =
            !startOverlay.hasAttribute('hidden') &&
            getComputedStyle(startOverlay).display !== 'none';
        if (initialVisible) startMenuMusic();
    }

    const lifeFishManager = new LifeFishManager();
    const turtleManager = new TurtleManager();

    // Settings handlers
    const {
        loadSettings,
        applySettingsToUI,
        saveSettings,
    } = createSettingsHandlers({
        getSoundEnabled: () => soundEnabled,
        setSoundEnabled: (value) => { soundEnabled = value; },
        getVolume: () => volume,
        setVolume: (value) => { volume = value; },
        soundToggle,
        volumeSlider,
    });

    // Tile Map System - Declare early so resize() can access it
    let tileMapRenderer = null;
    const TILE_SIZE = 96; // Increased from 64 to make tiles more visible
    const TILES_PER_ROW = 3;
    let currentMapIndex = 0;

    // World size
    function resize() {
        WORLD.width = gameEl.clientWidth;
        WORLD.height = gameEl.clientHeight;
        // Update tile map viewport when world resizes
        if (tileMapRenderer) {
            tileMapRenderer.setViewport(0, 0, WORLD.width, WORLD.height);
            tileMapRenderer.update();
        }
    }
    window.addEventListener('resize', resize);
    resize();

    async function initializeTileMap() {
        try {
            // Clear any existing tiles first
            if (tilemapLayer) {
                tilemapLayer.innerHTML = '';
            }

            // Generate tileset programmatically (solid colored tiles)
            const tilesetImage = await createGameTileset(TILE_SIZE);

            // Verify tileset was created
            if (!tilesetImage || tilesetImage.width === 0) {
                throw new Error('Failed to generate tileset image');
            }

            // Create tile map renderer
            tileMapRenderer = new TileMapRenderer(
                tilesetImage,
                TILE_SIZE,
                TILES_PER_ROW,
                tilemapLayer
            );

            // Set initial viewport to match world size
            tileMapRenderer.setViewport(0, 0, WORLD.width, WORLD.height);

            // Load first map with depth level 0 (transparent)
            const initialMap = getMapByIndex(0);
            tileMapRenderer.setMap(initialMap, 0);

            console.log('[Tilemap] System initialized successfully');
        } catch (error) {
            console.error('[Tilemap] Failed to initialize:', error);
        }
    }


    function loadMapForLevel(levelNum) {
        if (!tileMapRenderer) return;

        // Cycle through 4 maps (0-3)
        const mapIndex = Math.floor((levelNum - 1) / 3) % 4;

        // Depth level matches map index
        // Map 0 = depth 0 (transparent)
        // Map 1 = depth 1 (light blue)
        // Map 2 = depth 2 (medium blue)
        // Map 3 = depth 3 (dark blue)
        const depthLevel = mapIndex;

        const newMap = getMapByIndex(mapIndex);
        tileMapRenderer.setMap(newMap, depthLevel);

        console.log(`[Tilemap] Loaded map ${mapIndex} for level ${levelNum} with depth ${depthLevel}`);

    }

    // Function to set tile map by index
    function setTileMap(mapIndex) {
        if (tileMapRenderer && getMapByIndex) {
            currentMapIndex = mapIndex;
            const newMap = getMapByIndex(currentMapIndex);
            tileMapRenderer.setMap(newMap, currentMapIndex);
        }
    }

    // Initialize tile map system
    initializeTileMap();

    // HUD
    const hudSystem = createHudSystem({
        hudTime, hudScore, hudLives, hudFps, hudLevel,
        menuTime, menuScore, menuLives, menuFps, menuLevel,
        documentRef: document,
        getTimeLeft: () => timeLeft,
        getScore: () => score,
        getFps: () => fps,
        getLevel: () => level,
        getLives: () => lives,
        getHighScore: () => highScore,
    });
    const { updateHud, updateMenu } = hudSystem;

    // Story system - initialize early so it's available for gameLoop
    const storyManager = getStoryManager();

    const feedbackSystem = createFeedbackSystem({
        entitiesLayer,
        comboDisplay,
        centerNotification,
        gameoverOverlay,
        finalScore,
        finalAccuracy,
        finalCombo,
        finalLevel,
        finalCaught,
        celebrateHighScore,
        updateHighScoreListUI,
        addHighScore,
        isHighScore,
        storyManager,
        state: {
            getCombo: () => combo,
            getScore: () => score,
            setScore: (value) => { score = value; },
            getHighScore: () => highScore,
            setHighScore: (value) => { highScore = value; },
            getTotalShots: () => totalShots,
            getTotalHits: () => totalHits,
            getMaxCombo: () => maxCombo,
            getLevel: () => level,
            getFishCaught: () => fishCaught,
            getLives: () => lives,
            getTimeLeft: () => timeLeft,
            setRunning: (value) => { running = value; },
            setPaused: (value) => { paused = value; },
            getRunning: () => running,
            getPaused: () => paused,
        },
    });
    const {
        showScorePopup,
        updateComboDisplay,
        showCenterNotification,
        clearCenterNotification,
        showGameOver,
        hideGameOver,
    } = feedbackSystem;

    const originalShowGameOver = showGameOver;
    function showGameOverWithSound(...args) {
        playGameOverOnce();
        stopBgMusic();
        return originalShowGameOver(...args);
    }

    const shootingSystem = createShootingSystem({
        gameEl,
        entitiesLayer,
        lifeFishManager,
        turtleManager,
        removeEntity,
        rand,
        constants: {
            shotRateLimitMs: SHOT_RATE_LIMIT_MS,
            comboWindow: COMBO_WINDOW,
            maxLives: MAX_LIVES,
        },
        state: {
            getEntities: () => entities,
            getLastShotTimestamp: () => lastShotMs,
            setLastShotTimestamp: (value) => { lastShotMs = value; },
            incrementTotalShots: () => { totalShots += 1; },
            incrementTotalHits: () => { totalHits += 1; },
            incrementFishCaught: () => { fishCaught += 1; },
            incrementMissedShots: () => { missedShots += 1; },
            getCombo: () => combo,
            setCombo: (value) => { combo = value; },
            getComboTimer: () => comboTimer,
            setComboTimer: (value) => { comboTimer = value; },
            getMaxCombo: () => maxCombo,
            setMaxCombo: (value) => { maxCombo = value; },
            getScore: () => score,
            setScore: (value) => { score = value; },
            getLives: () => lives,
            setLives: (value) => { lives = value; },
            getLevel: () => level,
        },
        feedback: {
            updateComboDisplay,
            showScorePopup,
            showCenterNotification,
        },
    });
    const { attemptShot } = shootingSystem;

    const gameLoop = createGameLoop({
        lifeFishManager,
        turtleManager,
        spawnBubble,
        spawnParrot,
        positionElement,
        removeEntity,
        entitiesLayer,
        bubblesLayer,
        world: WORLD,
        constants: {
            maxLevel: MAX_LEVEL,
            levelDuration: LEVEL_DURATION,
            baseSpawnMin: BASE_SPAWN_MIN,
            baseSpawnMax: BASE_SPAWN_MAX,
            parrotBaseSpeed: PARROT_BASE_SPEED,
            baseMaxFish: BASE_MAX_FISH,
            maxLives: MAX_LIVES,
        },
        hud: { updateHud },
        feedback: { updateComboDisplay, showCenterNotification, showGameOver: showGameOverWithSound },
        storyManager,
        tileMapRenderer: () => tileMapRenderer, // Pass getter function
        getMapByIndex,
        loadMapForLevel,
        state: {
            getEntities: () => entities,
            getTimeLeft: () => timeLeft,
            setTimeLeft: (value) => { timeLeft = value; },
            getLevel: () => level,
            setLevel: (value) => { level = value; },
            getLevelTimer: () => levelTimer,
            setLevelTimer: (value) => { levelTimer = value; },
            getSpawnTimer: () => spawnTimer,
            setSpawnTimer: (value) => { spawnTimer = value; },
            getBubbleTimer: () => bubbleTimer,
            setBubbleTimer: (value) => { bubbleTimer = value; },
            getCombo: () => combo,
            setCombo: (value) => { combo = value; },
            getComboTimer: () => comboTimer,
            setComboTimer: (value) => { comboTimer = value; },
            getLives: () => lives,
            setLives: (value) => { lives = value; },
            getScore: () => score,
            setScore: (value) => { score = value; },
            getRunning: () => running,
            setRunning: (value) => { running = value; },
            getPaused: () => paused,
            setPaused: (value) => { paused = value; },
            getNextEntityId: () => nextEntityId,
            setNextEntityId: (value) => { nextEntityId = value; },
            getFps: () => fps,
            setFps: (value) => { fps = value; },
            getFpsAccum: () => fpsAccum,
            setFpsAccum: (value) => { fpsAccum = value; },
            getFpsFrames: () => fpsFrames,
            setFpsFrames: (value) => { fpsFrames = value; },
            getLastFrameTs: () => lastFrameTs,
            setLastFrameTs: (value) => { lastFrameTs = value; },
        },
        rand,
    });
    const { frame } = gameLoop;

    const lifecycleSystem = createLifecycleSystem({
        lifeFishManager,
        turtleManager,
        clearEntities,
        clearCelebration,
        updateComboDisplay,
        clearCenterNotification,
        showCenterNotification,
        updateHud,
        updateMenu,
        updateHighScoreListUI,
        hideGameOver,
        createPauseBubbles,
        pauseOverlay,
        startOverlay,
        requestFrame: window.requestAnimationFrame.bind(window),
        frame,
        storyManager,
        constants: {
            gameDuration: GAME_DURATION_S,
            maxLives: MAX_LIVES,
            levelDuration: LEVEL_DURATION,
        },
        state: {
            getEntities: () => entities,
            setRunning: (value) => { running = value; },
            getRunning: () => running,
            setPaused: (value) => { paused = value; },
            getPaused: () => paused,
            setTimeLeft: (value) => { timeLeft = value; },
            setScore: (value) => { score = value; },
            setLives: (value) => { lives = value; },
            setSpawnTimer: (value) => { spawnTimer = value; },
            // Expose bubble timer so lifecycle can initialize ambience
            setBubbleTimer: (value) => { bubbleTimer = value; },
            // Expose entity id counter reset for clean runs
            setNextEntityId: (value) => { nextEntityId = value; },
            setRngSeed: (value) => { rngSeed = value; },
            setLastFrameTs: (value) => { lastFrameTs = value; },
            setLevel: (value) => { level = value; },
            setLevelTimer: (value) => { levelTimer = value; },
            setCombo: (value) => { combo = value; },
            setMaxCombo: (value) => { maxCombo = value; },
            setComboTimer: (value) => { comboTimer = value; },
            setTotalShots: (value) => { totalShots = value; },
            setTotalHits: (value) => { totalHits = value; },
            setFishCaught: (value) => { fishCaught = value; },
            setMissedShots: (value) => { missedShots = value; },
        },
    });
    const {
        pauseGame,
        resumeGame,
        restartGame,
        returnToMainMenu,
        showStartMenu,
        startNewRun: originalStartNewRun,
    } = lifecycleSystem;

    // Wrapper to reset tile map when starting a new run
    function startNewRun() {
        // Reset tile map to first map
        setTileMap(0);
        originalStartNewRun();
    }

    const inputSystem = createInputSystem({
        gameEl,
        crosshair,
        attemptShot,
        documentRef: document,
        windowRef: window,
        constants: {
            crosshairHalf: CROSSHAIR_HALF,
            speedPxPerSec: SPEED_PX_PER_SEC,
        },
        state: {
            getRunning: () => running,
            getPaused: () => paused,
            getCrosshairX: () => crosshairX,
            setCrosshairX: (value) => { crosshairX = value; },
            getCrosshairY: () => crosshairY,
            setCrosshairY: (value) => { crosshairY = value; },
            getLastCrossTs: () => lastCrossTs,
            setLastCrossTs: (value) => { lastCrossTs = value; },
            getKeys: () => keys,
        },
    });
    inputSystem.init();

    // --- Button click SFX (global) --------------------------------------------
    const buttonSfx = new Audio('sounds/button.mp3');
    buttonSfx.preload = 'auto';
    buttonSfx.volume = 0.7;

    function playButtonSound() {
        if (!soundEnabled) return;
        try {
            const sfx = buttonSfx.cloneNode();
            sfx.volume = Math.max(0, Math.min(1, volume * 0.7));
            sfx.play().catch(() => { });
        } catch (_) { }
    }
    document.querySelectorAll('button').forEach(btn => {
        btn.addEventListener('click', playButtonSound);
    });
    // --------------------------------------------------------------------------

    // Popup functionality
    function setupPopup(openBtn, popup, closeBtn) {
        if (openBtn && popup) {
            openBtn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                popup.classList.remove('hidden');
                playButtonSound();
                // Re-apply settings to UI when settings popup is opened
                if (popup === settingsPopup) {
                    applySettingsToUI();
                }
            });
        }

        if (closeBtn && popup) {
            closeBtn.addEventListener('click', () => {
                popup.classList.add('hidden');
                playButtonSound();
            });

            // Close when clicking outside content
            popup.addEventListener('click', (e) => {
                if (e.target === popup) {
                    popup.classList.add('hidden');
                    playButtonSound();
                }
            });

            // Close with Escape key
            document.addEventListener('keydown', (e) => {
                if (e.key === 'Escape' && !popup.classList.contains('hidden')) {
                    popup.classList.add('hidden');
                    playButtonSound();
                }
            });
        }
    }

    // Setup all popups
    setupPopup(instructionsBtn, howtoPopup, closeHowto);
    setupPopup(settingsBtn, settingsPopup, closeSettings);
    setupPopup(creditsBtn, creditsPopup, closeCredits);

    // Init - load settings (defaults to ON if no saved preference)
    loadSettings();
    applySettingsToUI();

    // Ensure sound toggle checkbox state matches soundEnabled value
    // Force update the checkbox to ensure it reflects the current state
    if (soundToggle) {
        const currentState = soundEnabled;
        soundToggle.checked = currentState;
        if (currentState) {
            soundToggle.setAttribute('checked', '');
            soundToggle.setAttribute('aria-checked', 'true');
        } else {
            soundToggle.removeAttribute('checked');
            soundToggle.setAttribute('aria-checked', 'false');
        }
    }
    updateHighScoreListUI();
    initEnhancedMenu();

    const handleStoryComplete = () => {
        if (storyManager.storyMode === 'introduction') {
            startBgMusic();
            startNewRun();
        } else if (storyManager.storyMode === 'development') {
            // Resume gameplay after development scene
            // Game continues automatically
        } else if (storyManager.storyMode === 'conclusion') {
            // Conclusion complete - show game over screen
            // This will be handled by the feedback system
        }
    };

    const handleStorySkip = () => {
        if (storyManager.storyMode === 'introduction') {
            startBgMusic();
            startNewRun();
        } else if (storyManager.storyMode === 'development') {
            // Resume gameplay
        }
    };

    function attachDefaultStoryCallbacks() {
        storyManager.onStoryComplete = handleStoryComplete;
        storyManager.onStorySkip = handleStorySkip;
    }

    attachDefaultStoryCallbacks();

    // Start button - show introduction story
    if (startBtn) {
        startBtn.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
            stopMenuMusic();
            resetGoSfx();
            startOverlay.classList.add('hidden');
            attachDefaultStoryCallbacks();
            storyManager.showIntroduction();
        });
    }


    // Settings: sound toggle
    if (soundToggle) {
        soundToggle.addEventListener('change', () => {
            soundEnabled = !!soundToggle.checked;
            // Save the user's preference so it persists
            saveSettings();

            try {
                menuMusic.muted = !soundEnabled;
                bgMusic.muted = !soundEnabled;
                gameOverSfx.muted = !soundEnabled;

                const menuVisible = startOverlay &&
                    !startOverlay.hasAttribute('hidden') &&
                    getComputedStyle(startOverlay).display !== 'none';

                if (menuVisible) {
                    if (soundEnabled) startMenuMusic(); else stopMenuMusic();
                } else {
                    if (!soundEnabled) stopBgMusic();
                }
            } catch (_) { }
        });
    }

    // Settings: volume slider
    if (volumeSlider) {
        volumeSlider.addEventListener('input', () => {
            const val = Number(volumeSlider.value);
            if (!Number.isNaN(val)) {
                volume = Math.min(1, Math.max(0, val));
                saveSettings();
                try {
                    if (!menuMusic.paused) {
                        menuMusic.volume = Math.max(0, Math.min(1, volume * 0.75));
                    }
                    if (!bgMusic.paused) {
                        bgMusic.volume = Math.max(0, Math.min(1, Math.max(0.15, volume * 0.25)));
                    }
                } catch (_) { }
                buttonSfx.volume = Math.max(0, Math.min(1, volume * 0.7));
            }
        });
    }

    // Pause controls
    if (pauseBtn) {
        pauseBtn.addEventListener('click', () => {
            if (!running) return;
            if (!paused) {
                stopBgMusic();
                pauseGame(true);
            } else {
                startBgMusic();
                resumeGame();
            }
        });
    }

    // Tile map and coral overlay toggle
    if (toggleCoralBtn && coralOverlay && tilemapLayer) {
        let tileMapVisible = false;

        // Make sure button only responds to direct clicks - stop all events at mousedown
        toggleCoralBtn.addEventListener('mousedown', (e) => {
            e.preventDefault();
            e.stopPropagation();
            e.stopImmediatePropagation();
        }, true); // Use capture phase to catch early

        toggleCoralBtn.addEventListener('click', (e) => {
            // Stop event immediately to prevent any propagation
            e.preventDefault();
            e.stopPropagation();
            e.stopImmediatePropagation();

            // STRICT CHECK: Only trigger if clicking on the button or its direct content
            // Do NOT trigger if clicking on fish, entities, or anything outside the button
            const clickedElement = e.target;
            const isButtonClick = clickedElement === toggleCoralBtn ||
                clickedElement.closest('#toggle-coral-btn') === toggleCoralBtn;

            if (!isButtonClick) {
                return;
            }

            // Toggle tile map visibility - ONLY two states: visible or invisible
            tileMapVisible = !tileMapVisible;

            if (tileMapVisible) {
                // Show tile map and hide coral overlay
                tilemapLayer.style.opacity = '0.15'; // 15% opacity - very subtle
                coralOverlay.classList.remove('opacity-100');
                coralOverlay.classList.add('opacity-0');
                toggleCoralBtn.textContent = '🌊';
                toggleCoralBtn.title = 'Hide Tile Map';
            } else {
                // Hide tile map and show coral overlay
                tilemapLayer.style.opacity = '0'; // Hide tiles completely
                coralOverlay.classList.remove('opacity-0');
                coralOverlay.classList.add('opacity-100');
                toggleCoralBtn.textContent = '🌊';
                toggleCoralBtn.title = 'Show Tile Map';
            }

            playButtonSound();
        }, true); // Use capture phase
    }
    if (continueBtn) continueBtn.addEventListener('click', () => {
        startBgMusic();
        resumeGame();
    });
    if (restartBtn) restartBtn.addEventListener('click', () => {
        resetGoSfx();
        startBgMusic();
        restartGame();
    });

    const mainMenuBtn = document.getElementById('main-menu-btn');
    if (mainMenuBtn) mainMenuBtn.addEventListener('click', () => {
        stopBgMusic();
        resetGoSfx();
        startMenuMusic();
        returnToMainMenu();
    });

    if (restartGameoverBtn) restartGameoverBtn.addEventListener('click', () => {
        resetGoSfx();
        stopMenuMusic();
        startBgMusic();
        // Use startNewRun instead of restartGame to match the "Start Game" flow exactly
        // This ensures notifications work properly and the game starts the same way
        startNewRun();
    });
    if (mainMenuGameoverBtn) mainMenuGameoverBtn.addEventListener('click', () => {
        stopBgMusic();
        resetGoSfx();
        startMenuMusic();
        returnToMainMenu();
    });

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            if (!running) return;
            if (paused) {
                startBgMusic();
                resumeGame();
            } else {
                stopBgMusic();
                pauseGame(true);
            }
        }
    });

    // Optional autostart for debugging: visit with ?autostart=1 to start immediately
    const autoStart = (() => {
        try {
            return new URLSearchParams(window.location.search).get('autostart') === '1';
        } catch (_) {
            return false;
        }
    })();

    if (autoStart) {
        stopMenuMusic();
        resetGoSfx();
        startBgMusic();
        startNewRun();
    } else {
        showStartMenu();
        window.startNewRun = startNewRun;
    }
}

bootstrapGame();

