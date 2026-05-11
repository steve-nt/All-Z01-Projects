/**
 * Builds menu-related effects such as bubbles, celebrations, and pause overlays.
 */

export function createMenuEffects({
    instructionsBtn,
    settingsBtn,
    bubblesLayer,
    centerNotification,
    world,
    spawnBubble,
    rand,
}) {
    let activeCelebrationElements = [];

    function createMenuBubbles() {
        const bubblesContainer = document.getElementById('menu-bubbles');
        if (!bubblesContainer) return;

        bubblesContainer.innerHTML = '';

        for (let i = 0; i < 15; i++) {
            const bubble = document.createElement('div');
            bubble.className = 'menu-bubble';

            const size = 5 + Math.random() * 15;
            bubble.style.width = `${size}px`;
            bubble.style.height = `${size}px`;

            bubble.style.left = `${Math.random() * 100}vw`;
            bubble.style.bottom = `${-10 - Math.random() * 20}px`;

            const delay = Math.random() * 5;
            const duration = 5 + Math.random() * 10;
            bubble.style.animationDelay = `${delay}s`;
            bubble.style.animationDuration = `${duration}s`;

            bubblesContainer.appendChild(bubble);
        }
    }

    let menuBubbleIntervalId = null;

    function initEnhancedMenu() {
        createMenuBubbles();

        // Also spawn in-game style bubbles on the start menu using the same spawner
        const container = document.getElementById('menu-bubbles');
        if (container && typeof spawnBubble === 'function') {
            // Clear any previous interval
            if (menuBubbleIntervalId !== null) {
                clearInterval(menuBubbleIntervalId);
            }
            // Start a gentle ambient spawn loop
            menuBubbleIntervalId = setInterval(() => {
                spawnBubble({ bubblesLayer: container, world, rand });
            }, 700);
        }

        if (settingsBtn) {
            settingsBtn.addEventListener('click', () => {
                const settingsSection = document.querySelector('.menu-settings');
                if (settingsSection) {
                    settingsSection.scrollIntoView({ behavior: 'smooth' });
                }
            });
        }
    }

    function celebrateHighScore(callback) {
        if (!bubblesLayer) return;

        const totalBubbles = 1000;
        const duration = 11;
        const colors = ['#FF6B6B', '#FFD93D', '#6BCB77', '#4D96FF', '#FF7EB6', '#FFA07A', '#40E0D0'];

        activeCelebrationElements.forEach((element) => {
            if (element.parentNode) element.parentNode.removeChild(element);
        });
        activeCelebrationElements = [];

        for (let i = 0; i < totalBubbles; i++) {
            const bubble = document.createElement('div');
            const size = 10 + Math.random() * 30;
            bubble.className = 'bubble celebration';
            bubble.style.width = `${size}px`;
            bubble.style.height = `${size}px`;
            bubble.style.left = `${Math.random() * world.width}px`;
            bubble.style.top = `${Math.random() * world.height}px`;
            bubble.style.setProperty('--dx', `${(Math.random() - 0.5) * 800}px`);
            bubble.style.setProperty('--dy', `${(Math.random() - 0.5) * 800}px`);
            bubble.style.background = colors[Math.floor(Math.random() * colors.length)];
            bubble.style.opacity = 0.4 + Math.random() * 0.6;
            bubble.style.setProperty('--duration', `${duration}s`);
            bubblesLayer.appendChild(bubble);
            activeCelebrationElements.push(bubble);

            setTimeout(() => {
                if (bubble.parentNode) bubble.parentNode.removeChild(bubble);
            }, duration * 1000 + 200);
        }

        const flash = document.createElement('div');
        flash.className = 'celebration-flash';
        flash.style.position = 'absolute';
        flash.style.inset = '0';
        flash.style.background = 'radial-gradient(circle, rgba(255,255,255,0.3), transparent)';
        flash.style.pointerEvents = 'none';
        bubblesLayer.appendChild(flash);
        activeCelebrationElements.push(flash);
        setTimeout(() => flash.remove(), duration * 1000);

        if (centerNotification) {
            centerNotification.textContent = '🏆 NEW HIGH SCORE!';
            centerNotification.className = 'show level-up';
            setTimeout(() => centerNotification.classList.remove('show'), duration * 1000);
        }

        if (typeof callback === 'function') {
            setTimeout(callback, duration * 1000);
        }
    }

    function clearCelebration() {
        activeCelebrationElements.forEach((element) => {
            if (element.parentNode) element.parentNode.removeChild(element);
        });
        activeCelebrationElements = [];
        if (centerNotification) centerNotification.classList.remove('show');
    }

    function createPauseBubbles() {
        const bubblesContainer = document.getElementById('pause-bubbles');
        if (!bubblesContainer) return;

        bubblesContainer.innerHTML = '';

        for (let i = 0; i < 10; i++) {
            const bubble = document.createElement('div');
            bubble.className = 'menu-bubble';

            const size = 5 + Math.random() * 15;
            bubble.style.width = `${size}px`;
            bubble.style.height = `${size}px`;

            bubble.style.left = `${Math.random() * 100}vw`;
            bubble.style.bottom = `${-10 - Math.random() * 20}px`;

            const delay = Math.random() * 5;
            const duration = 5 + Math.random() * 10;
            bubble.style.animationDelay = `${delay}s`;
            bubble.style.animationDuration = `${duration}s`;

            bubblesContainer.appendChild(bubble);
        }
    }

    return {
        initEnhancedMenu,
        celebrateHighScore,
        clearCelebration,
        createPauseBubbles,
    };
}
