# 🖥️ System Monitor (C++ · Dear ImGui · SDL2 · OpenGL)

## 🎯 Project Objective
The goal of this project is to demonstrate **programming logic** in C++ and the ability to **adapt to a new language and libraries**.  
Instead of building an application from scratch, the task is to **extend and fix a given desktop system monitor** using **Dear ImGui** (immediate-mode GUI library) on top of **SDL2 + OpenGL**.

The application monitors **system resources and performance** by reading Linux kernel data directly from the **`/proc`** and **`/sys`** filesystems.

---

## ✨ Features (as required)

### 🔹 System Information
- Operating System type (from `/etc/os-release`)
- Current logged user
- Hostname (computer name)
- CPU model
- Total number of processes by state (running, sleeping, zombie, stopped, etc.)

### 🔹 CPU / Fan / Thermal (Tabbed Section)
- **CPU tab**  
  - Real-time CPU usage graph  
  - Overlay text showing current CPU %  
  - Sliders:  
    - FPS (graph refresh rate)  
    - Y-axis scale  
  - Checkbox/button to pause/resume animation
- **Fan tab**  
  - Fan status (enabled/disabled)  
  - Current speed (RPM)  
  - Level (auto/manual)  
  - Same graph system as CPU
- **Thermal tab**  
  - Graph of CPU temperature (from `/sys/class/thermal`)  
  - Overlay with current °C  

### 🔹 Memory and Processes
- **Memory section**  
  - RAM usage (progress bar)  
  - SWAP usage (progress bar)  
  - Disk usage (progress bar)
- **Processes table**  
  - Columns: PID, Name, State, CPU%, Memory%  
  - Search/filter box (by name or PID)  
  - Multi-row selection  

### 🔹 Network
- **IPv4 addresses** per interface (lo, eth0, wlan0, etc.)  
- **RX stats table**: bytes, packets, errors, drops, fifo, frame, compressed, multicast  
- **TX stats table**: bytes, packets, errors, drops, fifo, colls, carrier, compressed  
- **Progress bars per interface** showing usage (scaled 0 GB → 2 GB)  
- Automatic unit conversion (bytes → KB, MB, GB depending on range)

---

## 🗂 Project Structure
```
system-monitor/
├── Makefile
├── main.cpp              # main loop (ImGui + SDL2 + OpenGL)
├── header.h              # declarations, structs
├── cpu.*                 # CPU usage & model
├── fan.*                 # Fan info
├── temp.*                # Thermal sensors
├── mem.*                 # RAM/SWAP usage
├── process.*             # Process parsing from /proc/[pid]
├── network.*             # RX/TX stats and IPv4
├── system.cpp            # OS, Host, User info
└── imgui/lib/            # Dear ImGui core + backends + gl3w loader
```

---

## ⚙️ Tech Stack
- **C++17**  
- **Dear ImGui** (immediate-mode GUI)  
- **SDL2** (window/input backend)  
- **OpenGL (GL3w)** (rendering)  
- Linux **`/proc`** and **`/sysfs`** for system data  

---

## ✅ Prerequisites
On Ubuntu:
```bash
sudo apt update
sudo apt install -y build-essential pkg-config libsdl2-dev libgl1-mesa-dev
```

On Arch:
```bash
sudo pacman -S base-devel sdl2 mesa
```

On Fedora:
```bash
sudo dnf install @development-tools SDL2-devel mesa-libGL-devel
```

---

## 🏗️ Build
```bash
make
```
Produces the executable:
```bash
./system-monitor
```

Clean build:
```bash
make clean
```

---

## 🧭 Usage
- **Top bar**: OS / User / Hostname / CPU model  
- **CPU, Fan, Thermal tabs**: real-time graphs with controls (pause, FPS, scale)  
- **Memory**: RAM, SWAP, Disk progress bars  
- **Processes**: live table with filter & multi-select  
- **Network**: RX/TX tables + usage bars  

---

## 🛠️ Troubleshooting
- `sdl2-config: not found` → Install `libsdl2-dev`
- `fatal error: SDL2/SDL.h` → Missing SDL2 headers
- OpenGL linking errors (`-lGL`) → Install `libgl1-mesa-dev`
- No temperature/fan data → Some systems don’t expose `/sys/class/hwmon`; check with:
  ```bash
  ls -R /sys/class/hwmon
  grep . -R /sys/class/hwmon/* 2>/dev/null | head
  ```

---

## 📜 License
Educational use only. **Dear ImGui** and **SDL2** remain under their own licenses.