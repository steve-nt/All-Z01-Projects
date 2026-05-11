# Make Your Game — Bomberman Edition (Plain JS / DOM)

A Bomberman-inspired game built using **only plain JavaScript, HTML, and CSS**, without any frameworks, libraries, or canvas. All rendering is done via DOM manipulation and CSS—with performance and smooth animations in focus.

---

## 🎯 About the Project

This is a single‑player Bomberman‑style game made from scratch with an immersive story-driven narrative. The goals include:

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
- **No frameworks / no canvas** — pure JS/DOM/HTML only  
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

🍉 Maria   
🍇 Thanos  
🍓 Yana
