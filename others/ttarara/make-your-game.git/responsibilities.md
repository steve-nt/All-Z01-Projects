# 🧾 Responsibilities & Feature Log

This document tracks all agreed-upon features, fixes, and changes discussed by the team based on the recent conversation.

---

### 🎮 **Gameplay & Mechanics**

*  Add **difficulty levels** (“βαθμοί”) that increase as you catch more fish. — *Andy / Sofia*
* **DONE** Add a **bonus fish** that grants an **extra life**. — *Andy / Sofia*
* **DONE** Add **7 new fish** with unique colors, one per level. — *Andy / Sofia*
*  **DONE** Add a **turtle** that moves slowly to distract or block the player.
* **DONE**  Should appear in front of fish.  
*  **DONE** Player can click it — lose points but turtle remains visible. — *Sofia*
* **DONE** Fix **high score logic** to properly save/load via localStorage (no more hardcoded value). — *Andy / Sofia*
* **DONE** Increase the **catch radius** of the net/hook for better gameplay feel. — *Andy / Sofia*
*  Investigate and fix **unstable FPS counter** (shows 80–130 FPS after restart). Must stabilize around 60 FPS and reset properly on restart. — *Andy*
*  Need to comment main.js and story.js

---

### 🧭 **UI & Layout**

*  Convert **Settings**, **How to Play**, and other info panels into **accordion-style menus** (expand/collapse). — *Xaroula*
*  Add **music and sound effects** toggle. — *Xaroula*
* **DONE** Make **menu fish** swim **both directions** (left & right). — *Georgia*
* **DONE** When the player **loses**, show an **octopus with “Game Over”** animation before the menu appears. — *Andy / Sofia*
* **DONE** Make the **Pause menu button** visually match the **main menu** style. — *Georgia*
* **DONE** Standardize **Level**, **Time**, and **Score** HUD styling to match the main menu aesthetic. — *Georgia*
*  Review **accordion motion** — center settings window and smooth accordion open/close. — *Georgia*
*  Add **credits** section beneath the main menu. — *Georgia*

---

### 🌊 **Visuals & Effects**

*  Add **more background bubbles** for richer underwater feel. — *Andy / Sofia*
* **DONE** Display **celebratory bubbles animation** when achieving a new high score. — *Andy / Sofia*
*  Adjust **combo text (x2, x3, etc.)** position or style for clarity. — *Xaroula*
* **DONE** Slow down **score pop-up animation** for better visibility. — *Georgia*
*  Make **mouse cursor** always visible in gameplay. — *Xaroula*
* **DONE** Fix **keyboard input** responsiveness (no need to spam keys). — *Xaroula*
*  **DONE** Make **point loss** (missed shot penalty) visually clearer. 
*  Make the game have an infinity mode — *Team discussion*

---

### 🧩 **Ideas / Theming**

* Possible future **themes**: Pirate 🏴‍☠️, Christmas 🎄, Halloween 🎃, Haunted 👻, Kids 🧸.
* **Pirate theme** idea: include **treasure chest**, **anchor**, and **boat** decorations. — *Sofia / Georgia*
* **DONE** Add decorative **treasure chest GIF** (currently cosmetic only). — *Sofia*
*  Optionally make the **turtle** perform an **ambush**, stay longer, and swim sideways. — *Team discussion*
*  **DONE**  Show **extra life fish** whenever you lose a life — *Team discussion*
*  Investigate **octopus overlap bug** (falls over menu); delay menu appearance until animation ends. — *Georgia / Sofia*

---

### 🛠️ **Technical Notes**

* Project runs locally using:  
  ```bash
  python3 -m http.server 8080
