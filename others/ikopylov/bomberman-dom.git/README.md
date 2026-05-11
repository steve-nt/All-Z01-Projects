# Make Your Game — Bomberman Edition (Plain JS / DOM)

A Bomberman-inspired game built with **plain JavaScript, HTML, and CSS** — **no Canvas, no WebGL**. The **game board** is rendered with DOM elements and CSS only.

For the **multiplayer UI** (nickname, waiting room, match chat, game over), this project uses the **mini-framework** you build in the separate *mini-framework* repo: the sources are vendored under `framework/` from `mini-framework/packages/mini-framework/src/` (`h`, `createApp`, `createStore`, delegated events, `requestAnimationFrame` scheduling for renders). No React, Vue, or other UI libraries.

**WebSocket chat:** players can send messages over the same `/ws` connection as the game (`{ "type": "chat", "text": "..." }`). This is the lightweight “Hello World” for multiplayer: lobby chat + a floating in-game chat panel.

**Lobby size (bomberman-dom):** the waiting room accepts **at most four** connections at a time. If the lobby is full, new clients receive **HTTP 503** and the WebSocket handshake does not complete. While a match is **in progress**, additional connections may join as **spectators**. When the match ends, everyone returns to the waiting pool; if more than four clients are connected, the server **closes the newest connections** until only four remain, so the next lobby stays within the cap.

---

## 🎯 About the Project

This is a Bomberman‑style game made from scratch: **single‑player** includes an immersive story-driven narrative; **multiplayer** (WebSockets, 2–4 players per match, last player standing) uses the mini-framework shell described above. Shared goals include:

- Maintaining consistent **60 FPS** with no frame drops  
- Using `requestAnimationFrame` effectively  
- Measuring performance to ensure optimization  
- Implementing a **pause menu** (Continue, Restart)  
- Displaying a **scoreboard** with:
  - Timer / countdown  
  - Current score  
  - Remaining lives  
- Smooth keyboard input (holding keys should maintain movement)  
- Minimal layering in the DOM to reduce reflows and repaint cost  
- **No Canvas / WebGL** — board is DOM/CSS only; **mini-framework** is used only for the multiplayer shell UI (see above)  
- **Story system** with narrative segments (Intro, Mid-game, Win, Loss)

---

## 📖 Story

**Year 20XX.** Malicious agents from the 'Centralization' corporation have stolen the core of the peer-to-peer learning system from the school. All students have been taken captive. You are the only one left free.

**Your Mission:** Destroy the enemies, free your friends, and find the encryption key to restore the P2P learning system!

The game features story segments that appear at key moments:
- **Intro**: Mission briefing when you start
- **Mid-game**: Progress update when you break through the first line of defense
- **Win**: Victory message when you complete your mission
- **Loss**: Failure message when you're caught  

---


## 🛠️ Getting Started

1. **Clone the repository**  
   ```bash
   git clone https://platform.zone01.gr/git/aziagaki/make-your-game.git
   cd make-your-game
   ```

2. **Run the project**  
   Start the local server with:  
   ```bash
   go run main.go
   ```

3. **Open in browser**  
   Navigate to [http://localhost:8080](http://localhost:8080) (or the port specified in your `main.go`).  

---

## Multiplayer

- **Players per match:** 2–4 (Bomberman-style elimination).  
- **Waiting lobby:** up to **four** slots; when full, additional browsers cannot connect until a slot frees (see **Lobby size** in the intro). During an active match, extra clients may **spectate**.  
- **Flow:** choose a nickname → waiting room with player count (×/4) and optional lobby chat → countdown rules are enforced on the server (`internal/server/lobby.go`) → in-game state and chat share the **`/ws`** connection.

---

## ⌨️ Controls

- **Arrow keys / WASD** → Move player  
- **Spacebar** → Place bomb  
- **P** or **ESC**→ Pause / Resume / bring up pause menu  

---

## ⚡ Performance & Rendering

- Game loop driven by `requestAnimationFrame`, targeting **60 FPS**  
- Input handling ensures smooth continuous motion (no stutter, no need to spam keys)  
- Rendering done entirely with DOM and CSS (no canvas, no WebGL)  
- Minimal layering and selective DOM updates to avoid unnecessary repaints/reflows  
- Performance measured and monitored via browser dev tools (FPS metrics, profiler, paint flashing, etc.)

---

## 📚 Learning Goals

With this project, we explore and solidify:
- The JavaScript **event loop**  
- Using `requestAnimationFrame` for animations and game loops  
- Efficient **DOM manipulation** for game visuals  
- Identifying and eliminating **jank / stutter animations**  
- Performance profiling techniques in browser dev tools  

---

## 👥 Authors

- Andriana  
- Georgia  
- Theochara  
- Iana
