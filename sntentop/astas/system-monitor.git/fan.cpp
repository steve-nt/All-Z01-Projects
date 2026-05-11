#include "fan.h"
#include <imgui.h>

#ifdef __linux__

#include <fstream>      // For file input (reading from sysfs files)
#include <string>       // For std::string manipulation
#include <vector>       // For copying deque to vector for ImGui plotting
#include <deque>        // To store history of fan speed readings
#include <filesystem>   // For directory traversal to find hwmon files
#include <iostream>     // (Optional) For debugging output

namespace fs = std::filesystem;  // Alias to make filesystem calls shorter

// Search through /sys/class/hwmon directories to find the path to the fan speed input file,
// usually named something like "fan1_input"
static std::string findFanInputPath() {
    // Iterate all entries in /sys/class/hwmon
    for (const auto& entry : fs::directory_iterator("/sys/class/hwmon")) {
        if (!entry.is_directory()) continue;  // Only check directories

        // Inside each hwmonX directory, iterate files to find "fan1_input"
        for (const auto& file : fs::directory_iterator(entry.path())) {
            std::string filename = file.path().filename().string();
            if (filename.find("fan1_input") != std::string::npos) {
                // Return the full path to the fan input file when found
                return file.path().string();
            }
        }
    }
    // Return empty string if no fan input file was found
    return "";
}

// Find a related file to control or check fan enable status,
// such as "fan1_enable" or "fan1_status" in the same directory as fanInputPath
static std::string findFanEnablePath(const std::string& fanInputPath) {
    // Get the directory containing fanInputPath (e.g. /sys/class/hwmon/hwmon0/device/)
    fs::path base = fs::path(fanInputPath).parent_path();

    // Try each possible filename; return the first that exists
    for (const char* name : {"fan1_enable", "fan1_status"}) {
        fs::path p = base / name;
        if (fs::exists(p)) return p.string();
    }
    return "";
}

// Find a file that may indicate fan level or PWM control, like "fan1_level", "pwm1", or "pwm1_enable"
static std::string findFanLevelPath(const std::string& fanInputPath) {
    fs::path base = fs::path(fanInputPath).parent_path();

    // Check common fan control or level files, return first found
    for (const char* name : {"fan1_level", "pwm1", "pwm1_enable"}) {
        fs::path p = base / name;
        if (fs::exists(p)) return p.string();
    }
    return "";
}

// Helper function to read an integer value from a given file path
// Returns -1 if file can't be opened or read
static int readIntFromFile(const std::string& path) {
    std::ifstream file(path);
    if (!file.is_open()) return -1; // File couldn't be opened

    int val = -1;
    file >> val; // Read integer from file stream
    return val;
}

// Main function to gather fan information from the hardware monitoring sysfs files
FanInfo getFanInfo() {
    // Cache the paths so we only do filesystem search once
    static std::string fanInputPath = findFanInputPath();
    static std::string fanEnablePath = fanInputPath.empty() ? "" : findFanEnablePath(fanInputPath);
    static std::string fanLevelPath = fanInputPath.empty() ? "" : findFanLevelPath(fanInputPath);

    FanInfo info = {false, 0, 0};  // Default fan info: inactive, zero speed, zero level

    if (fanInputPath.empty()) {
        // No fan input file found; return default info (fan not detected)
        return info;
    }

    // Read current fan speed (RPM) from the fan input file
    int speed = readIntFromFile(fanInputPath);
    if (speed > 0) {
        info.speedRPM = speed;  // Update RPM if a valid reading was obtained
        info.active = true;     // Mark fan as active since it reports speed
    }

    // If available, check whether fan is enabled or active from enable/status file
    if (!fanEnablePath.empty()) {
        int enabled = readIntFromFile(fanEnablePath);
        info.active = (enabled == 1);  // 1 usually means fan is enabled
    }

    // If available, read fan level or PWM control value
    if (!fanLevelPath.empty()) {
        int level = readIntFromFile(fanLevelPath);
        if (level >= 0) info.level = level;
    }

    return info;
}

// Variables used for controlling the UI update and graph parameters
static bool pauseFan = false;       // Pause updating fan data graph
static int fpsFan = 60;             // Frames per second update rate for fan tab
static float yScaleFan = 8000.0f;  // Vertical scale for fan speed graph (max RPM)

// Store recent fan speeds for plotting as a graph, using a deque as a fixed-size buffer
static std::deque<float> fanSpeedHistory;
static constexpr int maxSamples = 100;  // Max number of points in the graph history

// Function to draw the fan tab in the ImGui interface
void renderFanTab() {
    ImGui::Text("Fan Information");
    ImGui::Separator();

    // Query the current fan info from sysfs
    FanInfo fan = getFanInfo();

    // Display fan status and readings
    ImGui::Text("Status: %s", fan.active ? "Active" : "Inactive");
    ImGui::Text("Speed: %d RPM", fan.speedRPM);
    ImGui::Text("Level: %d", fan.level);

    // UI controls for pausing updates, adjusting FPS, and graph Y scale
    ImGui::Checkbox("Pause", &pauseFan);
    ImGui::SliderInt("FPS", &fpsFan, 1, 144);
    ImGui::SliderFloat("Y Scale", &yScaleFan, 100.0f, 16000.0f, "%.0f RPM");

    if (!pauseFan) {
        // Add the current fan speed to the history buffer
        fanSpeedHistory.push_back(static_cast<float>(fan.speedRPM));
        // Keep the buffer size limited to maxSamples by popping oldest values
        if (fanSpeedHistory.size() > maxSamples) {
            fanSpeedHistory.pop_front();
        }
    }

    // Plot the fan speed history as a line graph if we have any data
    if (!fanSpeedHistory.empty()) {
        // Convert deque to vector for ImGui plotting API
        std::vector<float> values(fanSpeedHistory.begin(), fanSpeedHistory.end());
        ImVec2 graphSize = ImVec2(0, 100);  // Width=auto, height=100 pixels

        ImGui::PlotLines("Fan Speed (RPM)", values.data(), static_cast<int>(values.size()), 0, nullptr, 0.0f, yScaleFan, graphSize);

        // Show latest fan speed as text below the graph
        ImGui::Text("Current Speed: %.1f RPM", values.back());
    } else {
        // No fan data available to plot yet
        ImGui::Text("No fan data available.");
    }
}

#else // non-Linux systems

// Stub implementation for other OSs: fan monitoring not supported
FanInfo getFanInfo() {
    return {false, 0, 0};
}

void renderFanTab() {
    ImGui::Text("Fan monitoring is only available on Linux.");
    ImGui::Text("This feature uses /sys/class/hwmon.");
}

#endif











// DUMMY DATA VERSION FOR VISUALS
// #include "fan.h"
// #include <imgui.h>

// #ifdef __linux__

// #include <fstream>
// #include <string>
// #include <vector>
// #include <deque>
// #include <filesystem>
// #include <iostream>

// namespace fs = std::filesystem;

// // (same findFanInputPath, findFanEnablePath, findFanLevelPath as before)

// static std::string findFanInputPath() {
//     // Simulate a path, so the rest of the logic thinks a fan exists
//     return "/sys/class/hwmon/hwmon0/fan1_input";
// }
// static std::string findFanEnablePath(const std::string& fanInputPath) {
//     return "/sys/class/hwmon/hwmon0/fan1_enable";
// }
// static std::string findFanLevelPath(const std::string& fanInputPath) {
//     return "/sys/class/hwmon/hwmon0/fan1_level";
// }

// // Dummy internal state to simulate changing fan data
// static int dummyFanSpeed = 1000;
// static int dummyFanEnabled = 1;
// static int dummyFanLevel = 3;

// // Helper function to simulate reading integers from "files"
// // Instead of reading real files, return dummy values based on the path requested
// static int readIntFromFile(const std::string& path) {
//     // Simulate fan speed that changes over time, for demo
//     if (path.find("fan1_input") != std::string::npos) {
//         // Slowly increase or decrease speed for dynamic effect
//         dummyFanSpeed += (rand() % 21) - 10;  // Random change between -10 and +10 RPM
//         if (dummyFanSpeed < 500) dummyFanSpeed = 500;
//         if (dummyFanSpeed > 3000) dummyFanSpeed = 3000;
//         return dummyFanSpeed;
//     }
//     // Simulate fan enabled status
//     if (path.find("fan1_enable") != std::string::npos || path.find("fan1_status") != std::string::npos) {
//         return dummyFanEnabled;
//     }
//     // Simulate fan level or PWM
//     if (path.find("fan1_level") != std::string::npos || path.find("pwm1") != std::string::npos) {
//         return dummyFanLevel;
//     }

//     // Default fallback for unknown paths
//     return -1;
// }

// FanInfo getFanInfo() {
//     static std::string fanInputPath = findFanInputPath();
//     static std::string fanEnablePath = fanInputPath.empty() ? "" : findFanEnablePath(fanInputPath);
//     static std::string fanLevelPath = fanInputPath.empty() ? "" : findFanLevelPath(fanInputPath);

//     FanInfo info = {false, 0, 0};

//     if (fanInputPath.empty()) return info;

//     int speed = readIntFromFile(fanInputPath);
//     if (speed > 0) {
//         info.speedRPM = speed;
//         info.active = true;
//     }

//     if (!fanEnablePath.empty()) {
//         int enabled = readIntFromFile(fanEnablePath);
//         info.active = (enabled == 1);
//     }

//     if (!fanLevelPath.empty()) {
//         int level = readIntFromFile(fanLevelPath);
//         if (level >= 0) info.level = level;
//     }

//     return info;
// }

// // UI variables & history buffer remain unchanged
// static bool pauseFan = false;
// static int fpsFan = 60;
// static float yScaleFan = 8000.0f;

// static std::deque<float> fanSpeedHistory;
// static constexpr int maxSamples = 100;

// void renderFanTab() {
//     ImGui::Text("Fan Information");
//     ImGui::Separator();

//     FanInfo fan = getFanInfo();

//     ImGui::Text("Status: %s", fan.active ? "Active" : "Inactive");
//     ImGui::Text("Speed: %d RPM", fan.speedRPM);
//     ImGui::Text("Level: %d", fan.level);

//     ImGui::Checkbox("Pause", &pauseFan);
//     ImGui::SliderInt("FPS", &fpsFan, 1, 144);
//     ImGui::SliderFloat("Y Scale", &yScaleFan, 100.0f, 16000.0f, "%.0f RPM");

//     if (!pauseFan) {
//         fanSpeedHistory.push_back(static_cast<float>(fan.speedRPM));
//         if (fanSpeedHistory.size() > maxSamples) {
//             fanSpeedHistory.pop_front();
//         }
//     }

//     if (!fanSpeedHistory.empty()) {
//         std::vector<float> values(fanSpeedHistory.begin(), fanSpeedHistory.end());
//         ImVec2 graphSize = ImVec2(0, 100);
//         ImGui::PlotLines("Fan Speed (RPM)", values.data(), static_cast<int>(values.size()), 0, nullptr, 0.0f, yScaleFan, graphSize);
//         ImGui::Text("Current Speed: %.1f RPM", values.back());
//     } else {
//         ImGui::Text("No fan data available.");
//     }
// }

// #else // non-Linux systems

// FanInfo getFanInfo() {
//     return {false, 0, 0};
// }

// void renderFanTab() {
//     ImGui::Text("Fan monitoring is only available on Linux.");
//     ImGui::Text("This feature uses /sys/class/hwmon.");
// }

// #endif
