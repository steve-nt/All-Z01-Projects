// src/core/settings.js
// v6: Default sound enabled (on) - always defaults to true on first load
const LS_KEY = 'fth_settings_v6';

export function createSettingsHandlers({
  getSoundEnabled,
  setSoundEnabled,
  getVolume,
  setVolume,
  soundToggle,
  volumeSlider,
}) {
  function loadSettings() {
    try {
      // Clear old settings keys to ensure fresh start with new default
      ['fth_settings_v1', 'fth_settings_v2', 'fth_settings_v3', 'fth_settings_v4', 'fth_settings_v5'].forEach(key => {
        try {
          localStorage.removeItem(key);
        } catch (_) {}
      });
      
      const raw = localStorage.getItem(LS_KEY);
      if (raw) {
        const obj = JSON.parse(raw);
        // Load saved settings - respect user's explicit choice if they've set it
        if (typeof obj.soundEnabled === 'boolean') {
          setSoundEnabled(obj.soundEnabled);
        } else {
          // If soundEnabled is not a boolean, default to true (ON)
          setSoundEnabled(true);
        }
        if (typeof obj.volume === 'number') setVolume(clamp01(obj.volume));
      } else {
        // Default to sound enabled (ON) when no saved settings exist (first load)
        setSoundEnabled(true);
        // Also ensure the checkbox is checked in HTML
        if (soundToggle) {
          soundToggle.checked = true;
          soundToggle.setAttribute('checked', '');
          soundToggle.setAttribute('aria-checked', 'true');
        }
      }
    } catch (_) { 
      // On any error, default to sound enabled
      setSoundEnabled(true);
      if (soundToggle) {
        soundToggle.checked = true;
      }
    }
  }

  function applySettingsToUI() {
    try {
      if (soundToggle) {
        const v = !!getSoundEnabled();
        soundToggle.checked = v;
        soundToggle.setAttribute('aria-checked', String(v));
        // Ensure the checkbox visually reflects the state
        if (v) {
          soundToggle.setAttribute('checked', '');
        } else {
          soundToggle.removeAttribute('checked');
        }
      }
      if (volumeSlider) {
        const vol = clamp01(Number(getVolume()));
        volumeSlider.value = String(vol);
        volumeSlider.setAttribute('aria-valuenow', String(vol));
      }
    } catch (_) { /* ignore */ }
  }

  function saveSettings() {
    try {
      const obj = {
        soundEnabled: !!getSoundEnabled(),
        volume: clamp01(Number(getVolume())),
      };
      localStorage.setItem(LS_KEY, JSON.stringify(obj));
    } catch (_) { /* ignore */ }
  }

  return {
    loadSettings,
    applySettingsToUI,
    saveSettings,
  };
}

function clamp01(n) {
  if (Number.isNaN(n)) return 0;
  return Math.max(0, Math.min(1, n));
}
