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
  preferProceduralMusic: false,
  preferProceduralVictory: false,
  procedural: {
    running: false,
    intervals: [],
    gainNode: null,
    step: 0,
  },
};

function createAudioElement(src, { loop = false, volume = 1 } = {}) {
  if (!src) return null;
  const audio = new Audio(src);
  audio.loop = loop;
  audio.preload = 'auto';
  audio.volume = volume;
  return audio;
}

function primeAudioElement(audio) {
  if (!audio) return;
  // Do not interrupt an already-playing element (e.g. main theme after resumeMainTheme);
  // the unlock handler's play→pause would silence the loop on the first movement key.
  if (!audio.paused) return;
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

const BASS_PATTERN = [98, 98, 123, 98, 87, 98, 123, 110];

function getProceduralMusicBus() {
  if (!audioState.context) return null;
  if (!audioState.procedural.gainNode) {
    const g = audioState.context.createGain();
    g.gain.value = audioState.musicMuted ? 0 : 0.11;
    g.connect(audioState.context.destination);
    audioState.procedural.gainNode = g;
  }
  return audioState.procedural.gainNode;
}

function playMusicTone(frequency, duration, type = 'triangle') {
  if (!audioState.context || audioState.musicMuted) return;
  const ctx = audioState.context;
  const dest = getProceduralMusicBus();
  if (!dest) return;

  const resumed = ensureContext();
  const startPlayback = () => {
    const osc = ctx.createOscillator();
    const g = ctx.createGain();
    osc.type = type;
    const now = ctx.currentTime;
    osc.frequency.setValueAtTime(frequency, now);
    g.gain.setValueAtTime(0, now);
    g.gain.linearRampToValueAtTime(0.07, now + 0.02);
    g.gain.exponentialRampToValueAtTime(0.0008, now + duration);
    osc.connect(g);
    g.connect(dest);
    osc.start(now);
    osc.stop(now + duration + 0.06);
  };

  if (resumed && typeof resumed.then === 'function') {
    resumed.then(startPlayback).catch(startPlayback);
  } else {
    startPlayback();
  }
}

function proceduralMusicTick() {
  if (!audioState.procedural.running || audioState.musicMuted) return;
  const i = audioState.procedural.step % BASS_PATTERN.length;
  playMusicTone(BASS_PATTERN[i], 0.16, 'triangle');
  if (i % 4 === 0) {
    playMusicTone(392, 0.07, 'square');
  }
  audioState.procedural.step += 1;
}

function stopProceduralMainTheme() {
  audioState.procedural.running = false;
  audioState.procedural.intervals.forEach(id => clearInterval(id));
  audioState.procedural.intervals = [];
}

function startProceduralMainTheme({ restart = false } = {}) {
  if (!audioState.context || audioState.musicMuted) return;
  stopProceduralMainTheme();
  ensureContext();
  if (restart) audioState.procedural.step = 0;
  audioState.procedural.running = true;
  proceduralMusicTick();
  const id = setInterval(proceduralMusicTick, 185);
  audioState.procedural.intervals.push(id);
}

/** Short victory fanfare: sine melody on master out — not the procedural *theme* bus (triangle bass + square). */
function playProceduralVictory() {
  if (!audioState.context || audioState.musicMuted) return;
  const ctx = audioState.context;
  const phrase = [
    { hz: 523.25, len: 0.16, gap: 0.07 },
    { hz: 659.25, len: 0.16, gap: 0.07 },
    { hz: 783.99, len: 0.2, gap: 0.08 },
    { hz: 1046.5, len: 0.32, gap: 0.1 },
    { hz: 783.99, len: 0.14, gap: 0.06 },
    { hz: 1046.5, len: 0.42, gap: 0 },
  ];

  const run = () => {
    let acc = 0;
    const t0 = ctx.currentTime + 0.02;
    for (const n of phrase) {
      const start = t0 + acc;
      const osc = ctx.createOscillator();
      const g = ctx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(n.hz, start);
      g.gain.setValueAtTime(0, start);
      g.gain.linearRampToValueAtTime(0.13, start + 0.025);
      g.gain.exponentialRampToValueAtTime(0.0008, start + n.len);
      osc.connect(g);
      g.connect(ctx.destination);
      osc.start(start);
      osc.stop(start + n.len + 0.04);
      acc += n.len + n.gap;
    }
  };

  const resumed = ensureContext();
  if (resumed && typeof resumed.then === 'function') {
    resumed.then(run).catch(run);
  } else {
    run();
  }
}

function shouldUseHtmlMusic() {
  const track = audioState.music;
  return track && !audioState.preferProceduralMusic;
}

function shouldUseHtmlVictory() {
  const sound = audioState.sounds.victory;
  return sound && !audioState.preferProceduralVictory;
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

  if (audioState.music) {
    audioState.music.muted = !!audioState.musicMuted;
    audioState.music.addEventListener(
      'error',
      () => {
        audioState.preferProceduralMusic = true;
      },
      { once: true },
    );
  }
  if (audioState.sounds.victory) {
    audioState.sounds.victory.addEventListener(
      'error',
      () => {
        audioState.preferProceduralVictory = true;
      },
      { once: true },
    );
  }

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
  applyMuteState();
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
  stopProceduralMainTheme();
  if (shouldUseHtmlMusic()) {
    const track = audioState.music;
    if (restart) track.currentTime = 0;
    const resumed = ensureContext();
    const startPlayback = () => {
      track
        .play()
        .catch(() => {
          audioState.preferProceduralMusic = true;
          startProceduralMainTheme({ restart });
        });
    };
    if (resumed && typeof resumed.then === 'function') {
      resumed.then(startPlayback).catch(startPlayback);
    } else {
      startPlayback();
    }
    return;
  }
  startProceduralMainTheme({ restart });
}

function pauseMainTheme() {
  stopProceduralMainTheme();
  const track = audioState.music;
  if (track && !track.paused) {
    track.pause();
  }
}

function resumeMainTheme() {
  if (audioState.musicMuted) return;
  if (shouldUseHtmlMusic()) {
    const track = audioState.music;
    const resumed = ensureContext();
    const startPlayback = () => {
      track
        .play()
        .catch(() => {
          audioState.preferProceduralMusic = true;
          startProceduralMainTheme({ restart: false });
        });
    };
    if (resumed && typeof resumed.then === 'function') {
      resumed.then(startPlayback).catch(startPlayback);
    } else {
      startPlayback();
    }
    return;
  }
  startProceduralMainTheme({ restart: false });
}

function stopMainTheme() {
  stopProceduralMainTheme();
  const track = audioState.music;
  if (!track) return;
  track.pause();
  track.currentTime = 0;
}

function playVictory() {
  if (audioState.musicMuted) return;
  stopMainTheme();
  if (shouldUseHtmlVictory()) {
    const sound = audioState.sounds.victory;
    sound.currentTime = 0;
    const resumed = ensureContext();
    const startPlayback = () => {
      // Do not set preferProceduralVictory on play() rejection (autoplay / suspend);
      // only the audio `error` event means the file is missing or unloadable.
      sound.play().catch(() => {
        playProceduralVictory();
      });
    };
    if (resumed && typeof resumed.then === 'function') {
      resumed.then(startPlayback).catch(startPlayback);
    } else {
      startPlayback();
    }
    return;
  }
  playProceduralVictory();
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
  if (audioState.procedural.gainNode) {
    audioState.procedural.gainNode.gain.value = muted ? 0 : 0.11;
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
