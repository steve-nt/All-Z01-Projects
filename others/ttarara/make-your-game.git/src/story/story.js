/**
 * @file story.js
 * @description Controls the optional story overlay tutorial. Manages scene
 *              data, DOM bindings, navigation events, and emits callbacks to
 *              resume gameplay (via window.startNewRun()).
 */

class StoryManager {
    constructor() {
        this.currentStoryScene = 0;
        this.storyScenes = [
            {
                title: "The Underwater Adventure",
                text: "Deep beneath the ocean waves, in a magical coral reef, you are an underwater hunter with a special net. The reef is teeming with colorful fish, each more beautiful than the last."
            },
            {
                title: "The Challenge",
                text: "Your mission is to catch as many fish as possible before time runs out! Different fish are worth different points - the rarer the fish, the higher the score. Build combos for bonus points!"
            },
            {
                title: "The Hunt Begins",
                text: "Use your mouse to aim your net and click to catch the fish. Watch out - if fish escape, you'll lose a life! Catch the special red fish to gain an extra life when you're running low."
            },
            {
                title: "Ready to Dive?",
                text: "The ocean awaits your skills! Can you become the ultimate fish hunter and achieve the highest score? Your underwater adventure starts now!"
            }
        ];

        // DOM elements
        this.storyOverlay = document.getElementById('story-overlay');
        this.storyText = document.getElementById('story-text');
        this.storyNextBtn = document.getElementById('story-next-btn');
        this.storySkipBtn = document.getElementById('story-skip-btn');
        this.storyDots = document.querySelectorAll('.story-dot');

        this.init();
    }

    init() {
        this.setupEventListeners();
    }

    setupEventListeners() {
        // Story event listeners
        if (this.storyNextBtn) {
            this.storyNextBtn.addEventListener('click', () => this.nextStoryScene());
        }

        if (this.storySkipBtn) {
            this.storySkipBtn.addEventListener('click', () => {
                this.hideStory();
                // Call the global startNewRun function
                if (window.startNewRun) {
                    window.startNewRun();
                }
            });
        }

        // Story dot navigation
        this.storyDots.forEach((dot, index) => {
            dot.addEventListener('click', () => {
                this.currentStoryScene = index;
                this.updateStoryDisplay();
            });
        });
    }

    showStory() {
        if (this.storyOverlay) {
            this.storyOverlay.classList.remove('hidden');
            this.currentStoryScene = 0;
            this.updateStoryDisplay();
            this.createStoryBubbles();
        }
    }

    hideStory() {
        if (this.storyOverlay) {
            this.storyOverlay.classList.add('hidden');
        }
    }

    updateStoryDisplay() {
        if (!this.storyText || this.currentStoryScene >= this.storyScenes.length) return;

        const scene = this.storyScenes[this.currentStoryScene];
        const titleEl = document.getElementById('story-title');
        
        if (titleEl) titleEl.textContent = scene.title;
        
        // Animate text change
        this.storyText.style.opacity = '0';
        this.storyText.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
            this.storyText.innerHTML = `<p class="story-paragraph">${scene.text}</p>`;
            this.storyText.style.opacity = '1';
            this.storyText.style.transform = 'translateY(0)';
        }, 300);

        // Update progress dots
        this.storyDots.forEach((dot, index) => {
            dot.classList.remove('active', 'completed');
            if (index < this.currentStoryScene) {
                dot.classList.add('completed');
            } else if (index === this.currentStoryScene) {
                dot.classList.add('active');
            }
        });

        // Update button text
        if (this.storyNextBtn) {
            if (this.currentStoryScene === this.storyScenes.length - 1) {
                this.storyNextBtn.textContent = 'Start Game';
            } else {
                this.storyNextBtn.textContent = 'Continue';
            }
        }
    }

    nextStoryScene() {
        if (this.currentStoryScene < this.storyScenes.length - 1) {
            this.currentStoryScene++;
            this.updateStoryDisplay();
        } else {
            // Start the game
            this.hideStory();
            // Call the global startNewRun function
            if (window.startNewRun) {
                window.startNewRun();
            }
        }
    }

    createStoryBubbles() {
        const bubblesContainer = document.getElementById('story-bubbles');
        if (!bubblesContainer) return;

        // Clear existing bubbles
        bubblesContainer.innerHTML = '';

        // Create bubbles for story
        for (let i = 0; i < 12; i++) {
            const bubble = document.createElement('div');
            bubble.className = 'menu-bubble';

            // Random size
            const size = 5 + Math.random() * 15;
            bubble.style.width = `${size}px`;
            bubble.style.height = `${size}px`;

            // Random position
            bubble.style.left = `${Math.random() * 100}vw`;
            bubble.style.bottom = `${-10 - Math.random() * 20}px`;

            // Random animation delay and duration
            const delay = Math.random() * 5;
            const duration = 5 + Math.random() * 10;
            bubble.style.animationDelay = `${delay}s`;
            bubble.style.animationDuration = `${duration}s`;

            bubblesContainer.appendChild(bubble);
        }
    }
}

// Initialize story manager when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.storyManager = new StoryManager();
});

// Export for module systems if needed
if (typeof module !== 'undefined' && module.exports) {
    module.exports = StoryManager;
}
