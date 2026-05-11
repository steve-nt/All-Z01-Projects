/**
 * @file story.js
 * @description Story mode system with introduction, development (score milestones),
 *              and conclusion (win/loss) scenes. Manages story progression and
 *              integrates with game lifecycle.
 * 
 * Development scenes trigger at score milestones:
 * - 50000 points: "First Light"
 * - 100000 points: "Hope Returns"
 * - 180000 points: "The Turning Tide"
 * 
 * @version 2.2 - Updated score milestones (50000, 100000, 180000)
 */

export class StoryManager {
    constructor() {
        // Story mode state
        this.currentStoryScene = 0;
        this.storyMode = 'introduction'; // 'introduction', 'development', 'conclusion'
        this.shownDevelopmentScenes = new Set(); // Track which development scenes have been shown
        this.lastShownScore = 0; // Track last score milestone reached
        
        // Introduction scenes (shown before game starts)
        this.introductionScenes = [
            {
                title: "The Coral Crisis",
                text: "Deep beneath the ocean waves, a once-vibrant coral reef faces a grave threat. Mysterious dark shadows have begun to appear, causing the fish to scatter and the reef's colors to fade. You are a skilled underwater guardian, armed with a special net that can capture and relocate the fish to safety."
            },
            {
                title: "The Mission Begins",
                text: "Your mission is clear: rescue as many fish as possible before the darkness consumes the reef. Each fish you save strengthens the reef's ecosystem. The more fish you rescue, the brighter the reef becomes. But beware - if too many fish escape, the reef's life force will fade away."
            },
            {
                title: "The Power of Combos",
                text: "As you work, you'll discover that rescuing fish in quick succession creates powerful combos. These combos not only boost your score but also help restore the reef's energy. Watch for the rare red life fish - they carry the essence of the reef itself and can restore your strength when you're running low."
            },
            {
                title: "The Journey Ahead",
                text: "The ocean depths hold many challenges. As you progress, the darkness grows stronger, and the fish become more elusive. But with skill and determination, you can save the reef. Your underwater adventure begins now - the fate of the coral reef rests in your hands!"
            }
        ];

        // Development scenes (shown at specific score milestones during gameplay)
        this.developmentScenes = [
            {
                scoreThreshold: 50000,
                title: "First Light",
                text: "As you rescue the first wave of fish, you notice something remarkable - small patches of the reef are beginning to glow again! The coral responds to your efforts, its colors slowly returning. But this is only the beginning. The darkness still looms, and you must continue your mission."
            },
            {
                scoreThreshold: 100000,
                title: "Hope Returns",
                text: "The reef's recovery is accelerating! More fish are returning to the area, drawn by the restored energy. You can see the coral branches growing stronger, their vibrant colors spreading. However, the shadows are also growing bolder. Stay vigilant - the battle for the reef is far from over."
            },
            {
                scoreThreshold: 180000,
                title: "The Turning Tide",
                text: "Something incredible is happening! The reef is not just recovering - it's thriving! The fish you've rescued are now breeding, creating new life. The darkness is retreating, but it's fighting back with increased intensity. You're making real progress, but the final challenge awaits."
            }
        ];

        // Conclusion scenes (shown on game over)
        this.conclusionScenes = {
            victory: [
                {
                    title: "Triumph of Light",
                    text: "Against all odds, you have succeeded! The coral reef is fully restored, its colors more vibrant than ever. The fish swim freely in crystal-clear waters, and the darkness has been completely banished. You are hailed as the Guardian of the Deep, a hero whose legend will be told for generations."
                },
                {
                    title: "A New Beginning",
                    text: "The reef is not just saved - it's reborn. New coral formations are growing, and the ecosystem is stronger than before. Your dedication and skill have created a sanctuary where marine life can thrive. The ocean thanks you, Guardian. Your mission is complete, but the reef will remember your heroism forever."
                }
            ],
            defeat: [
                {
                    title: "The Darkness Prevails",
                    text: "Despite your valiant efforts, the darkness has spread too far. The reef's colors have faded, and many fish have been lost. The shadows grow stronger, and the coral weakens. But this is not the end - the ocean remembers those who fight for it, and your spirit will inspire future guardians."
                },
                {
                    title: "A Second Chance",
                    text: "Though the reef has fallen, hope remains. The fish you managed to save carry the memory of your courage. The ocean is vast and patient - it will wait for another guardian to rise. Your journey may have ended here, but the fight for the reef continues. Will you try again?"
                }
            ]
        };

        // DOM elements
        this.storyOverlay = document.getElementById('story-overlay');
        this.storyText = document.getElementById('story-text');
        this.storyNextBtn = document.getElementById('story-next-btn');
        this.storySkipBtn = document.getElementById('story-skip-btn');
        this.storyDots = document.querySelectorAll('.story-dot');
        this.storyTitle = document.getElementById('story-title');

        // Callbacks
        this.onStoryComplete = null;
        this.onStorySkip = null;

        this.init();
    }

    init() {
        this.setupEventListeners();
    }

    setupEventListeners() {
        if (this.storyNextBtn) {
            this.storyNextBtn.addEventListener('click', () => this.nextStoryScene());
        }

        if (this.storySkipBtn) {
            this.storySkipBtn.addEventListener('click', () => {
                this.hideStory();
                if (this.onStorySkip) {
                    this.onStorySkip();
                }
            });
        }

        // Story dot navigation (only for introduction)
        if (this.storyDots) {
            this.storyDots.forEach((dot, index) => {
                dot.addEventListener('click', () => {
                    if (this.storyMode === 'introduction' && index < this.introductionScenes.length) {
                        this.currentStoryScene = index;
                        this.updateStoryDisplay();
                    }
                });
            });
        }
    }

    /**
     * Show introduction story (before game starts)
     */
    showIntroduction() {
        this.storyMode = 'introduction';
        this.currentStoryScene = 0;
        this.shownDevelopmentScenes.clear();
        this.lastShownScore = 0;
        this.showStory();
    }

    /**
     * Show development story (at specific score milestones)
     * @param {number} currentScore - Current game score
     * @param {Function} onContinue - Callback when story is dismissed
     */
    showDevelopment(currentScore, onContinue) {
        // Find the scene for the current score if we haven't shown it yet
        let sceneToShow = null;
        for (const scene of this.developmentScenes) {
            // Check if score meets threshold, scene hasn't been shown, and score has increased
            if (currentScore >= scene.scoreThreshold && 
                !this.shownDevelopmentScenes.has(scene.scoreThreshold) &&
                currentScore > this.lastShownScore) {
                sceneToShow = scene;
                break; // Show the first matching scene
            }
        }

        if (!sceneToShow) return false; // No new scene to show

        this.storyMode = 'development';
        this.currentStoryScene = 0;
        this.shownDevelopmentScenes.add(sceneToShow.scoreThreshold);
        this.lastShownScore = currentScore;
        this.onStoryComplete = onContinue;
        this.onStorySkip = onContinue;

        // Set the development scene
        this.currentDevelopmentScene = sceneToShow;
        this.showStory();
        return true;
    }

    /**
     * Show conclusion story (on game over)
     * @param {boolean} isVictory - Whether the player won (completed all levels or high score)
     * @param {Function} onComplete - Callback when story is dismissed
     */
    showConclusion(isVictory, onComplete) {
        this.storyMode = 'conclusion';
        this.currentStoryScene = 0;
        this.isVictory = isVictory;
        this.onStoryComplete = onComplete;
        this.onStorySkip = onComplete;
        this.showStory();
    }

    showStory() {
        if (this.storyOverlay) {
            this.storyOverlay.classList.remove('hidden');
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
        if (!this.storyText) return;

        let scene = null;
        let totalScenes = 0;

        // Get the appropriate scene based on mode
        if (this.storyMode === 'introduction') {
            if (this.currentStoryScene >= this.introductionScenes.length) return;
            scene = this.introductionScenes[this.currentStoryScene];
            totalScenes = this.introductionScenes.length;
        } else if (this.storyMode === 'development') {
            if (!this.currentDevelopmentScene) return;
            scene = this.currentDevelopmentScene;
            totalScenes = 1; // Development scenes are single-page
        } else if (this.storyMode === 'conclusion') {
            const conclusionArray = this.isVictory ? 
                this.conclusionScenes.victory : 
                this.conclusionScenes.defeat;
            if (this.currentStoryScene >= conclusionArray.length) return;
            scene = conclusionArray[this.currentStoryScene];
            totalScenes = conclusionArray.length;
        }

        if (!scene) return;

        // Update title
        if (this.storyTitle) {
            this.storyTitle.textContent = scene.title;
        }

        // Animate text change
        this.storyText.style.opacity = '0';
        this.storyText.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
            this.storyText.innerHTML = `<p class="story-paragraph text-center text-sm md:text-base leading-relaxed text-cyan-200 animate-story">${scene.text}</p>`;
            this.storyText.style.opacity = '1';
            this.storyText.style.transform = 'translateY(0)';
        }, 300);

        // Update progress dots (only for introduction and conclusion)
        if (this.storyDots && (this.storyMode === 'introduction' || this.storyMode === 'conclusion')) {
            this.storyDots.forEach((dot, index) => {
                dot.classList.remove('active', 'completed');
                if (index < totalScenes) {
                    dot.style.display = 'block';
                    if (index < this.currentStoryScene) {
                        dot.classList.add('completed');
                    } else if (index === this.currentStoryScene) {
                        dot.classList.add('active');
                    }
                } else {
                    dot.style.display = 'none';
                }
            });
        } else if (this.storyDots) {
            // Hide dots for development scenes
            this.storyDots.forEach(dot => dot.style.display = 'none');
        }

        // Update button text and visibility
        if (this.storyNextBtn) {
            const isLastScene = this.currentStoryScene >= totalScenes - 1;
            if (isLastScene) {
                if (this.storyMode === 'introduction') {
                    this.storyNextBtn.textContent = 'Start Game';
                } else if (this.storyMode === 'development') {
                    this.storyNextBtn.textContent = 'Continue';
                } else {
                    this.storyNextBtn.textContent = 'Continue';
                }
            } else {
                this.storyNextBtn.textContent = 'Continue';
            }
        }

        // Show/hide skip button (hide for conclusion)
        if (this.storySkipBtn) {
            if (this.storyMode === 'conclusion') {
                this.storySkipBtn.style.display = 'none';
            } else {
                this.storySkipBtn.style.display = 'block';
            }
        }
    }

    nextStoryScene() {
        let totalScenes = 0;
        
        if (this.storyMode === 'introduction') {
            totalScenes = this.introductionScenes.length;
        } else if (this.storyMode === 'development') {
            totalScenes = 1;
        } else if (this.storyMode === 'conclusion') {
            totalScenes = this.isVictory ? 
                this.conclusionScenes.victory.length : 
                this.conclusionScenes.defeat.length;
        }

        if (this.currentStoryScene < totalScenes - 1) {
            this.currentStoryScene++;
            this.updateStoryDisplay();
        } else {
            // Last scene - complete the story
            this.hideStory();
            if (this.onStoryComplete) {
                this.onStoryComplete();
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

    /**
     * Reset story state for a new game
     */
    reset() {
        this.shownDevelopmentScenes.clear();
        this.lastShownScore = 0;
        this.currentStoryScene = 0;
        // Ensure we're ready for new score milestones (50000, 100000, 180000)
    }
}

// Create global instance
let storyManagerInstance = null;

export function getStoryManager() {
    if (!storyManagerInstance) {
        storyManagerInstance = new StoryManager();
    }
    return storyManagerInstance;
}

// Initialize when DOM is ready
if (typeof document !== 'undefined') {
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => {
            window.storyManager = getStoryManager();
        });
    } else {
        window.storyManager = getStoryManager();
    }
}
