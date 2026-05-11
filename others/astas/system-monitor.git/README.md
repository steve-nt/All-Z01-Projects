# 🖥️ System Monitor

## Overview

This project is a **System Monitor Desktop Application** written in **C++** using the **Dear ImGui** GUI framework. Its goal is to monitor real-time hardware and process statistics from a Linux system by reading directly from the `/proc` and `/sys` filesystems.

This app is meant to demonstrate programming logic, system-level programming skills, and the ability to adapt to a new language (C++). Rather than creating an application from scratch, the task involves extending and fixing a provided base codebase.

---

## 🧰 Technologies Used

- **C++**  
- **Dear ImGui** (immediate mode GUI)  
- **SDL2** (for input/rendering backend)  
- **OpenGL** (for rendering)  
- **Linux `/proc` and `/sys`** (for monitoring data)

---

## 📂 File System Structure

The provided project is structured as follows:

```

system-monitor/
├── header.h
├── imgui/                 # Dear ImGui source and backends
│   └── lib/
│       ├── backend/       # SDL2/OpenGL backends
│       ├── gl3w/          # OpenGL loader
│       └── \*.cpp/h        # Core ImGui components
├── main.cpp               # Main application loop
├── Makefile
├── mem.cpp                # Memory and process monitor
├── network.cpp            # Network monitor
└── system.cpp             # System info monitor

````

Make sure to install SDL2 via:

```bash
sudo apt install libsdl2-dev
````

---

## 💡 Project Features

### 🔧 System Monitor Tab

Displays high-level system information:

* Operating System type
* Logged-in user
* Hostname
* Total number of processes, categorized by state
* CPU model name

Tabbed section includes:

* **CPU Tab**:

  * Live graph of CPU usage with overlay percentage
  * Sliders to adjust FPS and Y-axis scaling
  * Play/pause animation control

* **Fan Tab**:

  * Fan status (active/enabled)
  * Speed and level
  * Graph similar to CPU tab

* **Thermal Tab**:

  * Current CPU temperature
  * Graph with real-time overlay

---

### 🧠 Memory and Processes Tab

Monitors memory usage and running processes:

* **RAM usage** with visual indicator
* **SWAP usage** with visual indicator
* **Disk usage** with visual indicator

**Processes Table**:

* Columns:

  * `PID`, `Name`, `State`, `CPU Usage`, `Memory Usage`
* Search bar to filter processes
* Multi-select rows supported

---

### 🌐 Network Tab

Displays all network interfaces (`lo`, `wlp5s0`, etc.):

* **IPv4 address display**

* **RX & TX tables** with:

  * Bytes, packets, errors, dropped packets, etc.

* Tabbed sections show **usage bars** per interface:

  * Converted from bytes to MB/GB appropriately
  * Graph scale from 0 GB to 2 GB

---

## ❗Known Limitation — CPU Usage per Process

The **CPU usage per process** column in the process table always shows **`0.00%`**.

### Why it happens:

CPU usage is calculated using:

* `currProcCpuTime`: from `/proc/[pid]/stat`
* `totalCpuTime`: from `/proc/stat`
* `deltaProc = currProcCpuTime - lastProcCpuTime`
* `deltaTotal = totalCpuTime - lastTotalCpuTime`

Then:
`cpuPercent = (deltaProc / deltaTotal) * 100.0f * cpuCount`

### Why it's broken:

Despite correct math and parsing:

* **Delta values are too small** due to:

  * High update rate (very short intervals between samples)
  * Lack of sufficient CPU-bound processes during testing
* **Linux `/proc/[pid]/stat`** fields are **not always reliable** for precise real-time per-process usage
* Real-time sampling requires **smoother timing and possibly OS-specific APIs** for accurate deltas

### Is it justifiable?

Yes. Process-level CPU usage is **non-trivial** to measure accurately in Linux using `/proc`, especially without:

* Sampling over longer time deltas
* Smoothing / averaging calculations
* Root access for high-precision system timers (in some distros)

This is acceptable for a student-level system monitor, especially when the rest of the monitoring (memory, network, system info) works as intended.

---

## 🧠 What You'll Learn

* C++ and immediate-mode UI logic
* Linux system internals via `/proc` and `/sys`
* How to read system-level info programmatically
* Dear ImGui UI/UX and rendering logic
* Basics of real-time graph rendering and data smoothing

---

## 🏁 How to Run

```bash
make
./monitor
```

---

## 📜 License

This project is provided for educational purposes.
