import { postScore } from '../api/client.js';
import { GAME_DURATION_S } from '../core/constants.js';
import { showScoreboard } from './scoreboard.js';

/**
* Bundles UI feedback helpers for scores, combos, notifications, and Game Over.
*/

export function createFeedbackSystem({
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
   state,
   storyManager,
}) {
   let overlayRevealTimeoutId = null;
   let initialsPromptCleanup = null;
   let notificationTimeoutId = null;

   function cleanupInitialsPrompt() {
       if (typeof initialsPromptCleanup === 'function') {
           initialsPromptCleanup();
           initialsPromptCleanup = null;
       } else {
           const container = document.getElementById('high-score-input-container');
           if (container) container.classList.add('hidden');
       }
   }

   function clearOverlayRevealTimeout() {
       if (overlayRevealTimeoutId !== null) {
           clearTimeout(overlayRevealTimeoutId);
           overlayRevealTimeoutId = null;
       }
   }

   function clearNotificationTimeout() {
       if (notificationTimeoutId !== null) {
           clearTimeout(notificationTimeoutId);
           notificationTimeoutId = null;
       }
   }

   function clearCenterNotification() {
       if (!centerNotification) return;
       clearNotificationTimeout();
       centerNotification.classList.remove('show', 'life-lost', 'level-up');
       centerNotification.style.opacity = '';
       // Ensure opacity-0 class is present
       if (!centerNotification.classList.contains('opacity-0')) {
           centerNotification.classList.add('opacity-0');
       }
       centerNotification.textContent = '';
   }

   function scheduleGameOverOverlay(context) {
       if (!gameoverOverlay) return;
       cleanupInitialsPrompt();
       clearOverlayRevealTimeout();
       overlayRevealTimeoutId = setTimeout(() => {
           overlayRevealTimeoutId = null;
           gameoverOverlay.classList.remove('hidden');
           setupInitialsPrompt(context);
       }, 1600);
   }

   function setupInitialsPrompt({ qualifiesForHighScore, currentScore, submitOnlineScore }) {
       const nameInputContainer = document.getElementById('high-score-input-container');
       const playerInitialsInput = document.getElementById('player-initials');
       const submitInitialsBtn = document.getElementById('submit-initials-btn');
       const canCollectInitials =
           Boolean(nameInputContainer) && Boolean(playerInitialsInput) && Boolean(submitInitialsBtn);

       if (!canCollectInitials) {
           if (nameInputContainer) nameInputContainer.classList.add('hidden');
           submitOnlineScore();
           return;
       }

       nameInputContainer.classList.remove('hidden');
       playerInitialsInput.value = '';
       playerInitialsInput.focus();

       const handleInput = (e) => {
           e.target.value = e.target.value.replace(/[^A-Za-z]/g, '').toUpperCase().slice(0, 3);
       };
       const handleSubmit = () => {
           const initials = playerInitialsInput.value.trim().toUpperCase().slice(0, 3) || '---';
           if (qualifiesForHighScore) {
               addHighScore(currentScore, initials);
               updateHighScoreListUI();
           }
           submitOnlineScore(initials);
           cleanup();
       };
       const handleKeyPress = (e) => {
           if (e.key === 'Enter') {
               handleSubmit();
           }
       };
       const cleanup = () => {
           if (nameInputContainer) nameInputContainer.classList.add('hidden');
           submitInitialsBtn.removeEventListener('click', handleSubmit);
           playerInitialsInput.removeEventListener('keypress', handleKeyPress);
           playerInitialsInput.removeEventListener('input', handleInput);
           initialsPromptCleanup = null;
       };

       initialsPromptCleanup = cleanup;
       playerInitialsInput.addEventListener('input', handleInput);
       submitInitialsBtn.addEventListener('click', handleSubmit);
       playerInitialsInput.addEventListener('keypress', handleKeyPress);
   }

   function showScorePopup(x, y, points, comboCount) {
       if (!entitiesLayer) return;
       const popup = document.createElement('div');
       popup.className = 'score-popup';
       const comboText = comboCount > 1 ? ` (${comboCount}x)` : '';
       popup.textContent = `+${points}${comboText}`;
       // Position popup at click location, centered
       popup.style.position = 'absolute';
       popup.style.left = `${x}px`;
       popup.style.top = `${y}px`;
       popup.style.pointerEvents = 'none';
       popup.style.zIndex = '1000';
       entitiesLayer.appendChild(popup);
       setTimeout(() => {
           if (popup.parentNode) popup.parentNode.removeChild(popup);
       }, 1000);
   }

   function updateComboDisplay() {
       if (!comboDisplay) return;
       if (state.getCombo() >= 2) {
           comboDisplay.textContent = `${state.getCombo()}x COMBO!`;
           comboDisplay.classList.add('active');
       } else {
           comboDisplay.classList.remove('active');
       }
   }

   function showCenterNotification(message, type = '', duration = 2000) {
       if (!centerNotification) {
           console.warn('centerNotification element not found');
           return;
       }

       // Clear any existing notification timeout
       clearNotificationTimeout();

       // Set the message first
       centerNotification.textContent = message;
       
       // Handle type classes
       if (type) {
           centerNotification.classList.remove('life-lost', 'level-up');
           centerNotification.classList.add(type);
       }
       
       // Remove opacity-0 class and show class to reset
       centerNotification.classList.remove('opacity-0', 'show');
       
       // Force animation restart
       const hasAnimation = centerNotification.classList.contains('animate-notif');
       if (hasAnimation) {
           centerNotification.classList.remove('animate-notif');
           // Force a reflow
           void centerNotification.offsetHeight;
           centerNotification.classList.add('animate-notif');
       }
       
       // Use a small delay to ensure DOM is ready, then show
       setTimeout(() => {
           // Remove opacity-0 class if it exists
           centerNotification.classList.remove('opacity-0');
           // Directly set opacity via inline style to override Tailwind
           centerNotification.style.opacity = '1';
           centerNotification.classList.add('show');
       }, 50);

       notificationTimeoutId = setTimeout(() => {
           notificationTimeoutId = null;
           centerNotification.classList.remove('show');
           centerNotification.style.removeProperty('opacity');
           centerNotification.classList.add('opacity-0');
       }, duration);
   }

   function showGameOver() {
       state.setRunning(false);
       state.setPaused(true);

       // Determine if it's a victory (completed all levels or achieved high score)
       const isVictory = state.getLevel() >= 10 || state.getScore() >= 5000;

       const totalShots = state.getTotalShots();
       const totalHits = state.getTotalHits();
       const accuracy = totalShots > 0 ? Math.round((totalHits / totalShots) * 100) : 0;
       const currentScore = state.getScore();
       const qualifiesForHighScore = isHighScore && isHighScore(currentScore);
       const timeSeconds = Math.max(0, GAME_DURATION_S - state.getTimeLeft());
       let onlineScoreSubmitted = false;

       const submitOnlineScore = (initialsValue = '') => {
           if (onlineScoreSubmitted) return;
           onlineScoreSubmitted = true;

           const normalizedInitials = normalizeInitials(initialsValue);
           postScore({
               name: normalizedInitials,
               score: currentScore,
               timeSeconds: Math.round(timeSeconds),
           })
               .then((recentSubmission) => {
                   showScoreboard({ recentSubmission });
               })
               .catch((error) => {
                   console.error('Failed to submit score to server:', error);
                   // Preserve offline fallback silently when API isn't reachable.
                   // But log it so we can debug
               });
       };

       function normalizeInitials(value = '') {
           const base = value || getInitialsFromInput();
           const trimmed = (base || '').trim().toUpperCase().slice(0, 3);
           return trimmed || 'Anon';
       }

       function getInitialsFromInput() {
           const input = document.getElementById('player-initials');
           return input ? input.value : '';
       }

       // Show conclusion story before game over screen
       if (storyManager) {
           storyManager.showConclusion(isVictory, () => {
               // After conclusion story, show game over screen
               showCenterNotification('GAME OVER', 'life-lost', 1500);
               if (!gameoverOverlay) return;

               if (finalScore) finalScore.textContent = String(currentScore);
              
               // Show high score message in menu if score qualifies (beats top 5)
               const highScoreMessage = document.getElementById('high-score-message');
               if (qualifiesForHighScore && highScoreMessage) {
                   highScoreMessage.classList.remove('hidden');
                   if (currentScore > state.getHighScore()) {
                       state.setHighScore(currentScore);
                   }
               } else if (highScoreMessage) {
                   highScoreMessage.classList.add('hidden');
               }

               if (finalAccuracy) finalAccuracy.textContent = `${accuracy}%`;
               if (finalCombo) finalCombo.textContent = `${state.getMaxCombo()}x`;
               if (finalLevel) finalLevel.textContent = String(state.getLevel());
               if (finalCaught) finalCaught.textContent = String(state.getFishCaught());

               scheduleGameOverOverlay({ qualifiesForHighScore, currentScore, submitOnlineScore });

               updateHighScoreListUI();
           });
       } else {
           // Fallback if story manager not available
           showCenterNotification('GAME OVER', 'life-lost', 1500);
           if (!gameoverOverlay) return;

           if (finalScore) finalScore.textContent = String(currentScore);
          
           // Show high score message in menu if score qualifies (beats top 5)
           const highScoreMessage = document.getElementById('high-score-message');
           if (qualifiesForHighScore && highScoreMessage) {
               highScoreMessage.classList.remove('hidden');
               if (currentScore > state.getHighScore()) {
                   state.setHighScore(currentScore);
               }
           } else if (highScoreMessage) {
               highScoreMessage.classList.add('hidden');
           }

           if (finalAccuracy) finalAccuracy.textContent = `${accuracy}%`;
           if (finalCombo) finalCombo.textContent = `${state.getMaxCombo()}x`;
           if (finalLevel) finalLevel.textContent = String(state.getLevel());
           if (finalCaught) finalCaught.textContent = String(state.getFishCaught());

           scheduleGameOverOverlay({ qualifiesForHighScore, currentScore, submitOnlineScore });

           updateHighScoreListUI();
       }
   }

   function hideGameOver() {
       clearOverlayRevealTimeout();
       cleanupInitialsPrompt();
       if (gameoverOverlay) {
           gameoverOverlay.classList.add('hidden');
           // Ensure no inline display styles interfere
           gameoverOverlay.style.display = '';
       }
       // Hide high score message when hiding game over
       const highScoreMessage = document.getElementById('high-score-message');
       if (highScoreMessage) {
           highScoreMessage.classList.add('hidden');
       }
   }

   return {
       showScorePopup,
       updateComboDisplay,
       showCenterNotification,
       clearCenterNotification,
       showGameOver,
       hideGameOver,
   };
}

