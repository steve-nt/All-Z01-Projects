#include "header.h"
#include <deque>
#include <vector>
#include <fstream>
#include <string>

// Platform detection for OS-specific includes
#ifdef _WIN32
#include <windows.h>
#elif __APPLE__
#include <mach/mach.h>
#include <mach/host_info.h>
#include <mach/mach_host.h>
#else // Assume Linux
#include <unistd.h>
#endif

// ------------------------------
// UI STATE (for CPU tab controls)
// ------------------------------

// Whether graph updates are paused
static bool pauseCPU = false;

// Graph refresh rate in frames per second
static int fpsCPU = 60;

// Max value on Y-axis (used to scale the CPU usage graph)
static float yScaleCPU = 100.0f;

// ------------------------------
// GRAPH DATA STORAGE
// ------------------------------

// Deque holds recent CPU usage values (acts like a rolling buffer)
static std::deque<float> cpuUsageHistory;

// Maximum number of samples stored in the graph (roughly equivalent to how wide the graph is)
static const int maxSamples = 100;

// ------------------------------
// CPU USAGE FUNCTION (Cross-platform)
// ------------------------------

// This function returns the current CPU usage in percentage
float getCpuUsagePercent() {
#ifdef __linux__
    // Static variables retain values between calls to calculate deltas
    static long long lastIdle = 0, lastTotal = 0;

    // Open /proc/stat (Linux-only virtual file with CPU stats)
    std::ifstream file("/proc/stat");
    if (!file.is_open()) return 0.0f; // Return 0 if file can’t be opened

    std::string cpu;
    CPUStats stat = {}; // Struct declared in header.h to hold values from /proc/stat

    // Read values: cpu user nice system idle iowait irq softirq steal
    file >> cpu >> stat.user >> stat.nice >> stat.system >> stat.idle
         >> stat.iowait >> stat.irq >> stat.softirq >> stat.steal;

    // Calculate total and idle times
    long long idle = stat.idle + stat.iowait;
    long long nonIdle = stat.user + stat.nice + stat.system + stat.irq + stat.softirq + stat.steal;
    long long total = idle + nonIdle;

    // Calculate change since last call
    long long deltaTotal = total - lastTotal;
    long long deltaIdle = idle - lastIdle;

    // Save current values for next comparison
    lastTotal = total;
    lastIdle = idle;

    // Avoid divide-by-zero
    if (deltaTotal == 0) return 0.0f;

    // Calculate CPU usage as percent
    return 100.0f * (deltaTotal - deltaIdle) / deltaTotal;

#elif _WIN32
    // Windows version using GetSystemTimes()

    static FILETIME prevIdle, prevKernel, prevUser;

    FILETIME idleTime, kernelTime, userTime;
    if (!GetSystemTimes(&idleTime, &kernelTime, &userTime)) return 0.0f;

    // Convert FILETIME to ULONGLONG for math
    ULARGE_INTEGER idle, kernel, user;
    idle.LowPart = idleTime.dwLowDateTime;   idle.HighPart = idleTime.dwHighDateTime;
    kernel.LowPart = kernelTime.dwLowDateTime; kernel.HighPart = kernelTime.dwHighDateTime;
    user.LowPart = userTime.dwLowDateTime;   user.HighPart = userTime.dwHighDateTime;

    ULONGLONG sysIdle = idle.QuadPart;
    ULONGLONG sysKernel = kernel.QuadPart;
    ULONGLONG sysUser = user.QuadPart;

    ULONGLONG prevSysIdle = ((ULARGE_INTEGER&)prevIdle).QuadPart;
    ULONGLONG prevSysKernel = ((ULARGE_INTEGER&)prevKernel).QuadPart;
    ULONGLONG prevSysUser = ((ULARGE_INTEGER&)prevUser).QuadPart;

    // Calculate total time and idle time difference
    ULONGLONG sysTotal = (sysKernel + sysUser) - (prevSysKernel + prevSysUser);
    ULONGLONG idleDiff = sysIdle - prevSysIdle;

    // Save current times for next call
    prevIdle = idleTime;
    prevKernel = kernelTime;
    prevUser = userTime;

    if (sysTotal == 0) return 0.0f;
    return 100.0f * (sysTotal - idleDiff) / sysTotal;

#elif __APPLE__
    // macOS version using host_statistics()

    static host_cpu_load_info_data_t prevLoad = {};
    mach_msg_type_number_t count = HOST_CPU_LOAD_INFO_COUNT;
    host_cpu_load_info_data_t load;

    // Query CPU usage stats from macOS kernel
    kern_return_t kr = host_statistics(mach_host_self(), HOST_CPU_LOAD_INFO,
                                       (host_info_t)&load, &count);
    if (kr != KERN_SUCCESS) return 0.0f;

    // Calculate deltas
    uint64_t user = load.cpu_ticks[CPU_STATE_USER] - prevLoad.cpu_ticks[CPU_STATE_USER];
    uint64_t system = load.cpu_ticks[CPU_STATE_SYSTEM] - prevLoad.cpu_ticks[CPU_STATE_SYSTEM];
    uint64_t idle = load.cpu_ticks[CPU_STATE_IDLE] - prevLoad.cpu_ticks[CPU_STATE_IDLE];
    uint64_t nice = load.cpu_ticks[CPU_STATE_NICE] - prevLoad.cpu_ticks[CPU_STATE_NICE];

    prevLoad = load;

    uint64_t total = user + system + idle + nice;
    if (total == 0) return 0.0f;
    return 100.0f * (user + system + nice) / total;

#else
    // Unsupported platform — return 0 but print a clear error
    #include <iostream>
    std::cerr << "Error: CPU usage monitoring is not supported on this platform." << std::endl;
    return 0.0f;
#endif
}

// ------------------------------
// UI RENDERING FUNCTION FOR CPU TAB
// ------------------------------

// This function draws the CPU tab with controls and the graph
void renderCpuTab() {
    ImGui::Text("CPU Usage");       // Label
    ImGui::Separator();             // Horizontal line

    // ------------------
    // UI CONTROLS
    // ------------------

    ImGui::Checkbox("Pause", &pauseCPU);                           // Toggle pause
    ImGui::SliderInt("FPS", &fpsCPU, 1, 144);                      // Adjust graph update speed
    ImGui::SliderFloat("Y Scale", &yScaleCPU, 10.0f, 200.0f, "%.1f%%"); // Adjust graph height

    static ImVec2 graphSize = ImVec2(0, 100); // Full width, 100px height

    // ------------------
    // UPDATE GRAPH DATA
    // ------------------

    if (!pauseCPU) {
        float usage = getCpuUsagePercent();     // Get current CPU usage
        cpuUsageHistory.push_back(usage);       // Store new value

        // Keep deque size within maxSamples
        if (cpuUsageHistory.size() > maxSamples) {
            cpuUsageHistory.pop_front();
        }
    }

    // ------------------
    // DRAW GRAPH
    // ------------------

    std::vector<float> values(cpuUsageHistory.begin(), cpuUsageHistory.end()); // Convert to vector
    if (!values.empty()) {
        ImGui::PlotLines("CPU %", values.data(), values.size(), 0, nullptr, 0.0f, yScaleCPU, graphSize);
    }

    // ------------------
    // CURRENT VALUE TEXT
    // ------------------

    ImGui::Text("Current: %.2f%%", cpuUsageHistory.empty() ? 0.0f : cpuUsageHistory.back());
}
