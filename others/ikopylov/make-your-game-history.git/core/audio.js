const TRACKS = {
  music: './audio/main-theme.mp3',
  victory: './audio/victory.mp3',
};

const SOUND_PRESETS = {
  bombPlace: { frequency: 200, duration: 0.1 },
  bombExplode: { frequency: 150, duration: 0.3 },
  playerHit: { frequency: 100, duration: 0.2 },
  enemyDeath: { frequency: 300, duration: 0.2 },
  powerUp: { frequency: 600, duration: 0.15 },
};

const audioState = {
  context: null,
  unlocked: false,
  music: null,
  sounds: {},
  initialized: false,
  musicMuted: false,
  musicButton: null,
};

function createAudioElement(src, { loop = false, volume = 1 } = {}) {
  if (!src) return null;
  const audio = new Audio(src);
  audio.loop = loop;
  audio.preload = 'auto';
  audio.volume = volume;
  audio.crossOrigin = 'anonymous';
  return audio;
}

function primeAudioElement(audio) {
  if (!audio) return;
  try {
    const playPromise = audio.play();
    if (playPromise && typeof playPromise.then === 'function') {
      playPromise
        .then(() => {
          audio.pause();
          audio.currentTime = 0;
        })
        .catch(() => {});
    }
  } catch (err) {
  }
}

function ensureContext() {
  if (!audioState.context) return null;
  if (audioState.context.state === 'suspended') {
    return audioState.context.resume().catch(() => {});
  }
  return null;
}

function syncMusicToggleButton() {
  const btn = audioState.musicButton;
  if (!btn) return;
  btn.setAttribute('aria-pressed', audioState.musicMuted ? 'true' : 'false');
}

function handleMusicButtonClick() {
  toggleMusicMute();
}

function bindMusicToggleButton() {
  if (audioState.musicButton && !document.body.contains(audioState.musicButton)) {
    audioState.musicButton.removeEventListener('click', handleMusicButtonClick);
    audioState.musicButton = null;
  }

  const btn = document.getElementById('ui-music-btn');
  if (!btn) return;

  if (audioState.musicButton !== btn) {
    audioState.musicButton = btn;
    btn.addEventListener('click', handleMusicButtonClick);
  }

  syncMusicToggleButton();
}

function scheduleMusicToggleButtonBinding() {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', bindMusicToggleButton, { once: true });
    return;
  }
  bindMusicToggleButton();
}

function setupUnlockHandler() {
  const resume = () => {
    if (audioState.unlocked) return;
    audioState.unlocked = true;

    if (audioState.context && audioState.context.state === 'suspended') {
      audioState.context.resume().catch(() => {});
    }

    Object.values(audioState.sounds).forEach(primeAudioElement);
    primeAudioElement(audioState.music);
  };

  window.addEventListener('pointerdown', resume, { once: true, capture: true });
  window.addEventListener('keydown', resume, { once: true, capture: true });
}

function loadTracks() {
  audioState.music = createAudioElement(TRACKS.music, { loop: true, volume: 0.4 });
  audioState.sounds.victory = createAudioElement(TRACKS.victory, { volume: 0.5 });

  if (audioState.music) audioState.music.muted = !!audioState.musicMuted;
  syncMusicToggleButton();
}

function init() {
  if (audioState.initialized) return;
  audioState.initialized = true;

  try {
    const saved = localStorage.getItem('musicMuted');
    if (saved != null) audioState.musicMuted = JSON.parse(saved);
  } catch (err) {
    // ignore storage errors
  }

  try {
    audioState.context = new (window.AudioContext || window.webkitAudioContext)();
  } catch (err) {
    console.warn('Web Audio API not supported.', err);
  }

  loadTracks();
  setupUnlockHandler();
  scheduleMusicToggleButtonBinding();
}

function play(typeOrFrequency, maybeFrequency, maybeDuration) {
  if (!audioState.context) return;

  const resumed = ensureContext();

  const oscillator = audioState.context.createOscillator();
  const gain = audioState.context.createGain();

  oscillator.connect(gain);
  gain.connect(audioState.context.destination);

  let frequency = 440;
  let duration = 0.1;

  if (typeof typeOrFrequency === 'number') {
    frequency = typeOrFrequency;
    if (Number.isFinite(maybeFrequency)) {
      duration = maybeFrequency;
    }
  } else {
    const preset = SOUND_PRESETS[typeOrFrequency] || {};
    frequency = Number.isFinite(maybeFrequency) ? maybeFrequency : preset.frequency ?? frequency;
    duration = Number.isFinite(maybeDuration) ? maybeDuration : preset.duration ?? duration;
  }

  if (!Number.isFinite(frequency)) frequency = 440;
  if (!Number.isFinite(duration) || duration <= 0) duration = 0.1;

  const startPlayback = () => {
    const now = audioState.context.currentTime;
    oscillator.frequency.setValueAtTime(frequency, now);
    oscillator.type = 'square';

    gain.gain.setValueAtTime(0.08, now);
    gain.gain.exponentialRampToValueAtTime(0.005, now + duration);

    oscillator.start(now);
    oscillator.stop(now + duration);
  };

  if (resumed && typeof resumed.then === 'function') {
    resumed.then(startPlayback).catch(startPlayback);
  } else {
    startPlayback();
  }
}

function playMainTheme({ restart = false } = {}) {
  const track = audioState.music;
  if (!track) return;
  if (restart) track.currentTime = 0;
  const resumed = ensureContext();
  const startPlayback = () => {
    track.play().catch(() => {});
  };
  if (resumed && typeof resumed.then === 'function') {
    resumed.then(startPlayback).catch(startPlayback);
  } else {
    startPlayback();
  }
}

function pauseMainTheme() {
  const track = audioState.music;
  if (track && !track.paused) {
    track.pause();
  }
}

function resumeMainTheme() {
  const track = audioState.music;
  if (!track) return;
  const resumed = ensureContext();
  const startPlayback = () => {
    track.play().catch(() => {});
  };
  if (resumed && typeof resumed.then === 'function') {
    resumed.then(startPlayback).catch(startPlayback);
  } else {
    startPlayback();
  }
}

function stopMainTheme() {
  const track = audioState.music;
  if (!track) return;
  track.pause();
  track.currentTime = 0;
}

function playVictory() {
  if (audioState.musicMuted) return;
  const sound = audioState.sounds.victory;
  if (!sound) return;
  stopMainTheme();
  sound.currentTime = 0;
  const resumed = ensureContext();
  const startPlayback = () => {
    sound.play().catch(() => {});
  };
  if (resumed && typeof resumed.then === 'function') {
    resumed.then(startPlayback).catch(startPlayback);
  } else {
    startPlayback();
  }
}

function setMusicMuted(muted) {
  audioState.musicMuted = !!muted;
  applyMuteState();
  if (audioState.music) audioState.music.muted = audioState.musicMuted;
  try {
    localStorage.setItem('musicMuted', JSON.stringify(audioState.musicMuted));
  } catch (err) {
    // ignore storage errors
  }
  syncMusicToggleButton();
}

function applyMuteState() {
  const muted = !!audioState.musicMuted;
  if (audioState.music) {
    audioState.music.muted = muted;
  }
  Object.values(audioState.sounds).forEach(sound => {
    if (sound && typeof sound === 'object' && 'muted' in sound) {
      sound.muted = muted;
    }
  });
}

function toggleMusicMute() {
  setMusicMuted(!audioState.musicMuted);
}

function isMusicMuted() {
  return !!audioState.musicMuted;
}

export const audio = {
  init,
  play,
  playMainTheme,
  pauseMainTheme,
  resumeMainTheme,
  stopMainTheme,
  playVictory,
  setMusicMuted,
  toggleMusicMute,
  isMusicMuted,
};

export {
  init,
  play,
  playMainTheme,
  pauseMainTheme,
  resumeMainTheme,
  stopMainTheme,
  playVictory,
  setMusicMuted,
  toggleMusicMute,
  isMusicMuted,
};

export default audio;
