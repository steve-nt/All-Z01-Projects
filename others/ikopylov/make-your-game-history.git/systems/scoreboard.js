import { state } from '../core/state.js';

const scoreboardState = {
  titleEl: null,
  defaultPanel: null,
  scoreboardPanel: null,
  namePanel: null,
  messageEl: null,
  errorEl: null,
  tableBody: null,
  prevBtn: null,
  nextBtn: null,
  pageLabel: null,
  playAgainBtn: null,
  closeBtn: null,
  nameInput: null,
  nameSubmitBtn: null,
  nameDefaultBtn: null,
  nameErrorEl: null,
  currentPage: 1,
  totalPages: 1,
  pageSize: 5,
  includeId: null,
  processing: false,
  wired: false,
  onRestart: null,
  lastTitle: 'PAUSED',
};

export function formatTime(milliseconds) {
  if (!Number.isFinite(milliseconds) || milliseconds < 0) {
    return '00:00';
  }
  const totalSeconds = Math.floor(milliseconds / 1000);
  const minutes = Math.floor(totalSeconds / 60)
    .toString()
    .padStart(2, '0');
  const seconds = (totalSeconds % 60).toString().padStart(2, '0');
  return `${minutes}:${seconds}`;
}

export async function submitScore(name, score, timeStr) {
  const response = await fetch('/scores', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, score, time: timeStr }),
  });

  if (!response.ok) {
    throw new Error(await extractError(response, 'Unable to submit score.'));
  }

  return response.json();
}

export async function fetchScores(page = 1, pageSize = 5, includeId) {
  const params = new URLSearchParams();
  params.set('page', page);
  params.set('pageSize', pageSize);
  params.set('sort', 'desc');
  if (Number.isFinite(includeId) && includeId > 0) {
    params.set('includePercentileForId', includeId);
  }

  const response = await fetch(`/scores?${params.toString()}`);
  if (!response.ok) {
    throw new Error(await extractError(response, 'Unable to fetch scores.'));
  }

  return response.json();
}

export function initScoreboard({ onRestart } = {}) {
  const pause = state.pause ?? {};

  scoreboardState.titleEl = resolveElement(pause.titleEl, 'pause-title');
  scoreboardState.defaultPanel = resolveElement(pause.defaultPanel, 'pause-panel-default');
  scoreboardState.scoreboardPanel = resolveElement(pause.scoreboardPanel, 'scoreboard-panel');
  scoreboardState.namePanel = resolveElement(pause.namePanel, 'name-entry-panel');
  scoreboardState.messageEl = resolveElement(null, 'scoreboard-message');
  scoreboardState.errorEl = resolveElement(null, 'scoreboard-error');
  scoreboardState.tableBody = document.querySelector('#scoreboard-table tbody');
  scoreboardState.prevBtn = resolveElement(null, 'scoreboard-prev');
  scoreboardState.nextBtn = resolveElement(null, 'scoreboard-next');
  scoreboardState.pageLabel = resolveElement(null, 'scoreboard-page');
  scoreboardState.playAgainBtn = resolveElement(pause.scoreboardPlayAgainBtn, 'scoreboard-play-again');
  scoreboardState.closeBtn = resolveElement(pause.scoreboardCloseBtn, 'scoreboard-close');
  scoreboardState.nameInput = resolveElement(pause.nameInput, 'name-entry-input');
  scoreboardState.nameSubmitBtn = resolveElement(pause.nameSubmitBtn, 'name-entry-submit');
  scoreboardState.nameDefaultBtn = resolveElement(pause.nameDefaultBtn, 'name-entry-default');
  scoreboardState.nameErrorEl = resolveElement(pause.nameErrorEl, 'name-entry-error');
  scoreboardState.onRestart = typeof onRestart === 'function' ? onRestart : null;

  wireControls();
}

export async function handleGameFinished() {
  if (scoreboardState.processing) return;
  scoreboardState.processing = true;

  try {
    scoreboardState.lastTitle =
      (scoreboardState.titleEl && scoreboardState.titleEl.textContent) || 'PAUSED';

    const name = await requestPlayerName();
    const elapsedMs = Math.max(0, performance.now() - state.status.startTime);
    const score = Math.max(0, Math.floor(state.status.score));
    const timeString = formatTime(elapsedMs);

    const submission = await submitScore(name, score, timeString);
    const submittedId = Number(submission.id);
    scoreboardState.includeId =
      Number.isFinite(submittedId) && submittedId > 0 ? submittedId : null;
    await loadPage(1);
  } catch (error) {
    showScoreboardPanel();
    showError(error?.message || 'Unexpected error while submitting score.');
    console.error(error);
  } finally {
    scoreboardState.processing = false;
  }
}

export function renderScoreboard(serverState) {
  if (!scoreboardState.tableBody || !scoreboardState.messageEl) return;

  scoreboardState.currentPage = coercePositiveInt(serverState?.page, 1);
  scoreboardState.totalPages = Math.max(coercePositiveInt(serverState?.totalPages, 1), 1);

  populateTable(serverState?.items);

  scoreboardState.messageEl.textContent = resolveMessage(serverState);

  updatePaginationControls();

  hideError();
  showScoreboardPanel();
}

function resolveElement(candidate, fallbackId) {
  return candidate ?? (fallbackId ? document.getElementById(fallbackId) : null);
}

function coercePositiveInt(value, fallback) {
  const num = Number(value);
  return Number.isFinite(num) && num > 0 ? Math.trunc(num) : fallback;
}

function populateTable(items) {
  if (!scoreboardState.tableBody) return;

  scoreboardState.tableBody.innerHTML = '';

  if (!Array.isArray(items) || items.length === 0) {
    scoreboardState.tableBody.append(createEmptyRow());
    return;
  }

  const fragment = document.createDocumentFragment();
  items.forEach(item => {
    fragment.append(buildScoreRow(item));
  });
  scoreboardState.tableBody.append(fragment);
}

function buildScoreRow(item) {
  const row = document.createElement('tr');
  const itemId = Number(item?.id);
  if (Number.isFinite(itemId) && itemId === scoreboardState.includeId) {
    row.classList.add('scoreboard-row-highlight');
  }

  row.append(
    createCell(item?.position ?? '', 'row'),
    createCell(item?.name ?? ''),
    createCell(item?.score ?? ''),
    createCell(item?.time ?? '')
  );
  return row;
}

function createCell(text, scope) {
  const cell = document.createElement('td');
  cell.textContent = text;
  if (scope) {
    cell.scope = scope;
  }
  return cell;
}

function createEmptyRow() {
  const emptyRow = document.createElement('tr');
  const emptyCell = document.createElement('td');
  emptyCell.colSpan = 4;
  emptyCell.textContent = 'No scores yet.';
  emptyCell.style.textAlign = 'center';
  emptyRow.append(emptyCell);
  return emptyRow;
}

function resolveMessage(serverState) {
  if (serverState?.subject?.message) {
    return serverState.subject.message;
  }
  if (Array.isArray(serverState?.items) && serverState.items.length > 0) {
    return 'Great run! Keep climbing the leaderboard.';
  }
  return 'Submit a score to see where you stand!';
}

function updatePaginationControls() {
  if (scoreboardState.pageLabel) {
    scoreboardState.pageLabel.textContent = `Page ${scoreboardState.currentPage} / ${scoreboardState.totalPages}`;
  }
  if (scoreboardState.prevBtn) {
    scoreboardState.prevBtn.disabled = scoreboardState.currentPage <= 1;
  }
  if (scoreboardState.nextBtn) {
    scoreboardState.nextBtn.disabled =
      scoreboardState.currentPage >= scoreboardState.totalPages;
  }
}

async function loadPage(page) {
  try {
    const data = await fetchScores(page, scoreboardState.pageSize, scoreboardState.includeId);
    renderScoreboard(data);
  } catch (error) {
    showScoreboardPanel();
    showError(error?.message || 'Unable to load scores.');
    console.error(error);
  }
}

function wireControls() {
  if (scoreboardState.wired) return;

  if (scoreboardState.prevBtn) {
    scoreboardState.prevBtn.addEventListener('click', () => {
      if (scoreboardState.currentPage > 1) {
        loadPage(scoreboardState.currentPage - 1);
      }
    });
  }

  if (scoreboardState.nextBtn) {
    scoreboardState.nextBtn.addEventListener('click', () => {
      if (scoreboardState.currentPage < scoreboardState.totalPages) {
        loadPage(scoreboardState.currentPage + 1);
      }
    });
  }

  if (scoreboardState.playAgainBtn) {
    scoreboardState.playAgainBtn.addEventListener('click', () => {
      hideScoreboardPanel();
      if (typeof scoreboardState.onRestart === 'function') {
        scoreboardState.onRestart();
      }
    });
  }

  if (scoreboardState.closeBtn) {
    scoreboardState.closeBtn.addEventListener('click', () => {
      hideScoreboardPanel();
    });
  }

  scoreboardState.wired = true;
}

function showScoreboardPanel() {
  if (!scoreboardState.scoreboardPanel) return;
  state.pause.mode = 'scoreboard';

  scoreboardState.defaultPanel?.classList.add('hidden');
  scoreboardState.namePanel?.classList.add('hidden');
  scoreboardState.scoreboardPanel.classList.remove('hidden');

  if (scoreboardState.titleEl) {
    scoreboardState.titleEl.textContent = 'HIGH SCORES';
  }

  requestSelectionReset();
}

function hideScoreboardPanel() {
  scoreboardState.scoreboardPanel?.classList.add('hidden');
  scoreboardState.namePanel?.classList.add('hidden');
  if (scoreboardState.defaultPanel) {
    scoreboardState.defaultPanel.classList.remove('hidden');
  }
  if (scoreboardState.titleEl) {
    scoreboardState.titleEl.textContent = scoreboardState.lastTitle || 'PAUSED';
  }

  state.pause.mode = 'default';
  hideError();
  hideNameError();
  requestSelectionReset();
}

async function requestPlayerName() {
  if (!scoreboardState.namePanel) {
    return 'Player';
  }

  showNamePanel();

  return new Promise(resolve => {
    const cleanup = () => {
      scoreboardState.nameSubmitBtn?.removeEventListener('click', onSubmit);
      scoreboardState.nameDefaultBtn?.removeEventListener('click', onDefault);
      scoreboardState.nameInput?.removeEventListener('keydown', onKeydown);
    };

    const onSubmit = () => {
      const submitted = getSanitizedName(scoreboardState.nameInput?.value);
      if (!submitted) {
        showNameError('Please enter a name (1-20 characters) or choose Player.');
        return;
      }
      cleanup();
      resolve(submitted);
    };

    const onDefault = () => {
      cleanup();
      resolve('Player');
    };

    const onKeydown = event => {
      if (event.key === 'Enter') {
        event.preventDefault();
        onSubmit();
      } else if (event.key === 'Escape') {
        event.preventDefault();
        onDefault();
      }
    };

    scoreboardState.nameSubmitBtn?.addEventListener('click', onSubmit);
    scoreboardState.nameDefaultBtn?.addEventListener('click', onDefault);
    scoreboardState.nameInput?.addEventListener('keydown', onKeydown);
  });
}

function showNamePanel() {
  state.pause.mode = 'name-entry';
  scoreboardState.defaultPanel?.classList.add('hidden');
  scoreboardState.scoreboardPanel?.classList.add('hidden');
  scoreboardState.namePanel?.classList.remove('hidden');
  hideNameError();

  if (scoreboardState.titleEl) {
    scoreboardState.titleEl.textContent = 'ENTER NAME';
  }

  if (scoreboardState.nameInput) {
    scoreboardState.nameInput.value = '';
    scoreboardState.nameInput.focus({ preventScroll: true });
    requestAnimationFrame(() => {
      scoreboardState.nameInput?.select();
    });
  }

  state.pause.selectedIndex = 0;
}

function getSanitizedName(value) {
  if (typeof value !== 'string') return '';
  const trimmed = value.trim();
  if (!trimmed) return '';
  return trimmed.slice(0, 20);
}

function showError(message) {
  if (!scoreboardState.errorEl) return;
  scoreboardState.errorEl.textContent = message;
  scoreboardState.errorEl.classList.remove('hidden');
}

function hideError() {
  if (!scoreboardState.errorEl) return;
  scoreboardState.errorEl.classList.add('hidden');
  scoreboardState.errorEl.textContent = '';
}

function showNameError(message) {
  if (!scoreboardState.nameErrorEl) return;
  scoreboardState.nameErrorEl.textContent = message;
  scoreboardState.nameErrorEl.classList.remove('hidden');
}

function hideNameError() {
  if (!scoreboardState.nameErrorEl) return;
  scoreboardState.nameErrorEl.classList.add('hidden');
  scoreboardState.nameErrorEl.textContent = '';
}

async function extractError(response, fallback) {
  try {
    const payload = await response.json();
    if (payload && typeof payload.error === 'string') {
      return payload.error;
    }
  } catch (err) {
    // ignore
  }
  return fallback;
}

function requestSelectionReset() {
  document.dispatchEvent(new Event('pause:reset-selection'));
}
