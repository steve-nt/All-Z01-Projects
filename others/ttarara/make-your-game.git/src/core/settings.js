// src/core/settings.js
const LS_KEY = 'fth_settings_v1';

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
      const raw = localStorage.getItem(LS_KEY);
      if (raw) {
        const obj = JSON.parse(raw);
        if (typeof obj.soundEnabled === 'boolean') setSoundEnabled(obj.soundEnabled);
        if (typeof obj.volume === 'number') setVolume(clamp01(obj.volume));
      }
    } catch (_) { /* ignore */ }
  }

  function applySettingsToUI() {
    try {
      if (soundToggle) {
        const v = !!getSoundEnabled();
        soundToggle.checked = v;
        soundToggle.setAttribute('aria-checked', String(v));
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
