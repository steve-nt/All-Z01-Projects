import { createApp, createStore, h } from "./framework/index.js";
import { bootstrap } from "./game.js";
import { state, reset as coreReset } from "./core/state.js";
import { applyServerState } from "./systems/netSync.js";
import { playVictory, stopMainTheme } from "./core/audio.js";

const store = createStore({
  nickname: "",
  /** "intro" = mission splash (HTML); then "nickname" = join lobby form */
  phase: "intro",
  count: 1,
  players: [],
  readyIn: 0,
  /** Seconds left in 20s grace (2–3 players) before the 10s start countdown */
  waitingGraceIn: 0,
  messages: [],
  connected: false,
  error: "",
  winnerName: "",
  gameOver: false,
  spectator: false,
});

let socket = null;
let gameBootstrapped = false;
/** @type {number} Unix ms when the match starts (after 10s countdown); 0 if not scheduled */
let lobbyStartAtMs = 0;
/** @type {number} Unix ms when 20s grace ends (2–3 players); 0 if N/A */
let lobbyGraceEndsAtMs = 0;

function computeLobbyTimers() {
  const readyIn = lobbyStartAtMs
    ? Math.max(0, Math.ceil((lobbyStartAtMs - Date.now()) / 1000))
    : 0;
  const waitingGraceIn =
    lobbyStartAtMs || !lobbyGraceEndsAtMs
      ? 0
      : Math.max(0, Math.ceil((lobbyGraceEndsAtMs - Date.now()) / 1000));
  return { readyIn, waitingGraceIn };
}

function wsUrl() {
  const protocol = location.protocol === "https:" ? "wss" : "ws";
  return `${protocol}://${location.host}/ws`;
}

function send(data) {
  if (!socket || socket.readyState !== WebSocket.OPEN) return;
  socket.send(JSON.stringify(data));
}

function chatListsEqual(a, b) {
  if (!Array.isArray(a) || !Array.isArray(b) || a.length !== b.length) {
    return false;
  }
  return a.every((m, i) => m.from === b[i].from && m.text === b[i].text);
}

function startGameFromPayload(payload) {
  if (gameBootstrapped) return;
  gameBootstrapped = true;
  store.setState({ phase: "playing", gameOver: false, spectator: false });
  bootstrap({
    multiplayer: true,
    spectator: false,
    mySlot: payload.slot,
    tiles: payload.tiles,
    exitRow: payload.exitRow,
    exitCol: payload.exitCol,
    send,
  });
}

function startSpectatorFromPayload(payload) {
  if (gameBootstrapped) return;
  gameBootstrapped = true;
  store.setState({ phase: "playing", gameOver: false, spectator: true });
  bootstrap({
    multiplayer: true,
    spectator: true,
    mySlot: -1,
    tiles: payload.tiles,
    exitRow: payload.exitRow,
    exitCol: payload.exitCol,
    send,
  });
}

function connect(nickname) {
  gameBootstrapped = false;
  socket = new WebSocket(wsUrl());
  socket.addEventListener("open", () => {
    store.setState({ connected: true, error: "", nickname, phase: "waiting" });
    send({ type: "join", nickname });
  });

  socket.addEventListener("error", () => {
    store.setState({
      error:
        "Could not connect. The lobby may be full (maximum 4 players). Try again later.",
    });
  });

  socket.addEventListener("message", (event) => {
    let payload = null;
    try {
      payload = JSON.parse(event.data);
    } catch (_) {
      return;
    }
    if (!payload || !payload.type) return;

    if (payload.type === "lobby") {
      lobbyStartAtMs = Number(payload.startAtMs) || 0;
      lobbyGraceEndsAtMs = Number(payload.graceEndsAtMs) || 0;
      const { readyIn, waitingGraceIn } = computeLobbyTimers();
      const players = payload.players ?? [];
      const lobbyMsgs = Array.isArray(payload.messages)
        ? payload.messages.map((m) => ({
            from: String(m?.from ?? ""),
            text: String(m?.text ?? ""),
          }))
        : null;
      const prev = store.getState();
      const serverPhase = String(payload.phase ?? "waiting");
      const backToLobby = prev.phase === "ended" && serverPhase === "waiting";
      if (backToLobby) {
        gameBootstrapped = false;
        coreReset.status();
      }
      const same =
        !backToLobby &&
        prev.count === payload.count &&
        prev.readyIn === readyIn &&
        prev.waitingGraceIn === waitingGraceIn &&
        prev.players.length === players.length &&
        prev.players.every((n, i) => n === players[i]) &&
        (lobbyMsgs === null || chatListsEqual(prev.messages, lobbyMsgs));
      if (!same || backToLobby) {
        store.setState({
          count: payload.count,
          players,
          readyIn,
          waitingGraceIn,
          ...(backToLobby ? { phase: "waiting", gameOver: false } : {}),
          ...(lobbyMsgs !== null ? { messages: lobbyMsgs } : {}),
        });
      }
    }

    if (payload.type === "game_start") {
      startGameFromPayload(payload);
    }

    if (payload.type === "spectate_start") {
      startSpectatorFromPayload(payload);
    }

    if (payload.type === "state") {
      applyServerState(payload);
    }

    if (payload.type === "game_over") {
      gameBootstrapped = false;
      const w = payload.winner;
      const name = w && w.nickname ? w.nickname : "Nobody";
      const winSlot = w != null && Number.isFinite(Number(w.slot)) ? Number(w.slot) : -1;
      state.status.over = true;
      state.status.won = winSlot >= 0 && winSlot === state.netSlot;
      state.status.eliminated = false;
      store.setState({
        gameOver: true,
        winnerName: name,
        phase: "ended",
      });
      if (state.status.won) {
        playVictory();
      } else {
        stopMainTheme();
      }
    }

    if (payload.type === "chat") {
      store.setState((prev) => ({
        messages: [...prev.messages, { from: payload.from, text: payload.text }],
      }));
    }
  });

  socket.addEventListener("close", () => {
    store.setState((prev) => ({
      connected: false,
      error:
        prev.error ||
        "Connection closed. Refresh to rejoin.",
    }));
  });
}

function submitNickname(e) {
  e.preventDefault();
  const input = document.getElementById("nickname-input");
  const nickname = (input?.value || "").trim().slice(0, 20);
  if (!nickname) {
    store.setState({ error: "Nickname is required." });
    return;
  }
  connect(nickname);
}

function submitChat(e) {
  e.preventDefault();
  const input =
    document.getElementById("chat-input") ||
    document.getElementById("game-chat-input");
  const text = (input?.value || "").trim().slice(0, 160);
  if (!text) return;
  send({ type: "chat", text });
  input.value = "";
}

function nicknameView(state) {
  return h(
    "div",
    { class: "mp-overlay", key: "nick-overlay" },
    h(
      "form",
      { key: "nick-form", class: "mp-card", on: { submit: submitNickname } },
      h("h2", { key: "nick-h2" }, "Bomberman DOM Multiplayer"),
      h("p", { key: "nick-hint" }, "Enter your nickname to join the lobby."),
      h("input", {
        key: "nickname-field",
        id: "nickname-input",
        class: "mp-input",
        type: "text",
        maxLength: 20,
        placeholder: "Nickname",
      }),
      state.error
        ? h("p", { key: "nick-err", class: "mp-error" }, state.error)
        : null,
      h("button", { key: "nick-submit", class: "mp-btn", type: "submit" }, "Join"),
    ),
  );
}

function waitingView(state) {
  return h(
    "div",
    { class: "mp-overlay", key: "wait-overlay" },
    h(
      "div",
      { key: "wait-card", class: "mp-card" },
      h("h2", { key: "wait-title" }, "Waiting Room"),
      h("p", { key: "wait-count" }, `Players: ${state.count} / 4`),
      h(
        "p",
        { key: "wait-status" },
        state.readyIn > 0
          ? `Match starts in ${state.readyIn}s`
          : state.waitingGraceIn > 0
            ? `Waiting for more players — ${state.waitingGraceIn}s until the 10s start countdown`
            : "Waiting for players...",
      ),
      h(
        "ul",
        { key: "wait-players", class: "mp-player-list" },
        ...(state.players || []).map((name, i) =>
          h("li", { key: name || `p-${i}` }, name),
        ),
      ),
      h(
        "form",
        { key: "wait-chat-form", class: "mp-chat", on: { submit: submitChat } },
        h("div", { key: "wait-chat-title", class: "mp-chat-title" }, "Lobby Chat"),
        h(
          "div",
          { key: "wait-chat-box", class: "mp-chat-box" },
          ...(state.messages || [])
            .slice(-10)
            .map((m, idx) =>
              h("p", { key: `${m.from}-${idx}-${m.text?.slice(0, 8)}` }, `${m.from}: ${m.text}`),
            ),
        ),
        h("input", {
          key: "lobby-chat-field",
          id: "chat-input",
          class: "mp-input",
          type: "text",
          maxLength: 160,
          placeholder: "Write a message...",
        }),
      ),
    ),
  );
}

function gameOverView(state) {
  return h(
    "div",
    { class: "mp-overlay", key: "end-overlay" },
    h(
      "div",
      { key: "end-card", class: "mp-card" },
      h("h2", { key: "end-title" }, "Match Over"),
      h("p", { key: "end-winner" }, `Winner: ${state.winnerName || "—"}`),
      h(
        "p",
        { key: "end-hint", class: "mp-hint" },
        "Wait here for the next countdown, or refresh to rejoin.",
      ),
    ),
  );
}

/** In-game WebSocket chat (Hello World for multiplayer); does not block the board. */
function playingChatView(state) {
  return h(
    "div",
    {
      key: "play-chat-root",
      class: "mp-game-chat-float",
      "aria-label": "Multiplayer chat (WebSocket)",
    },
    h("div", { key: "play-chat-hdr", class: "mp-game-chat-header" }, "Chat · WebSocket"),
    h(
      "div",
      { key: "play-chat-msgs", class: "mp-game-chat-msgs" },
      ...(state.messages || [])
        .slice(-12)
        .map((m, idx) =>
          h(
            "p",
            { key: `gc-${m.from}-${idx}` },
            `${m.from}: ${m.text}`,
          ),
        ),
    ),
    h(
      "form",
      { key: "play-chat-form", class: "mp-game-chat-form", on: { submit: submitChat } },
      h("input", {
        key: "game-chat-field",
        id: "game-chat-input",
        class: "mp-input",
        type: "text",
        maxLength: 160,
        placeholder: "Message other players…",
        autocomplete: "off",
      }),
    ),
  );
}

function view({ state }) {
  if (state.phase === "intro") {
    return h("div", { key: "mp-intro-placeholder", class: "mp-intro-placeholder" });
  }
  if (state.phase === "nickname") return nicknameView(state);
  if (state.phase === "waiting") return waitingView(state);
  if (state.phase === "ended" && state.gameOver) return gameOverView(state);
  if (state.phase === "playing") return playingChatView(state);
  return h("div", null);
}

function syncMultiplayerRootVisibility(phase) {
  const root = document.getElementById("multiplayer-root");
  if (!root) return;
  root.style.display = phase === "intro" ? "none" : "";
}

function wireMissionIntroToLobby() {
  const btn = document.getElementById("intro-start-btn");
  const intro = document.getElementById("intro-menu");
  if (!btn || !intro) return;
  btn.addEventListener(
    "click",
    (e) => {
      if (store.getState().phase !== "intro") return;
      e.preventDefault();
      e.stopPropagation();
      intro.classList.add("hidden");
      store.setState({ phase: "nickname" });
    },
    { capture: true },
  );
}

const app = createApp({ root: "#multiplayer-root", view, store });
window.addEventListener("load", () => {
  app.mount();
  store.subscribe(() => {
    syncMultiplayerRootVisibility(store.getState().phase);
  });
  syncMultiplayerRootVisibility(store.getState().phase);
  wireMissionIntroToLobby();
  setInterval(() => {
    const s = store.getState();
    if (s.phase !== "waiting") return;
    const { readyIn, waitingGraceIn } = computeLobbyTimers();
    if (readyIn !== s.readyIn || waitingGraceIn !== s.waitingGraceIn) {
      store.setState({ readyIn, waitingGraceIn });
    }
  }, 250);
});
