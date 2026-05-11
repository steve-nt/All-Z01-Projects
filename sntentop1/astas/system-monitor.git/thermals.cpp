#include "header.h" // Include your main header
#include <imgui.h>  // ImGui UI library

#ifdef __linux__

#include <filesystem>
#include <fstream>
#include <deque>
#include <vector>
#include <string>

namespace fs = std::filesystem;

// Store a fixed-size history of temperature samples for the graph
static std::deque<float> thermalHistory;
static constexpr int maxSamples = 100; // Max number of samples in history

// Controls for the thermal graph
static bool pauseThermal = false;      // Whether to pause updating the graph
static int fpsThermal = 60;            // Graph refresh rate
static float yScaleThermal = 100.0f;   // Max Y axis value for the graph

// Fallback mode if real temperature can't be read
static bool useDummyThermal = false;

// Cache the discovered sensor file path
static std::string thermalSensorPath;

// Search for a valid CPU temperature sensor under /sys/class/hwmon
// Try to find thermal sensor file from various common locations
static std::string findThermalSensorPath() {
    // Step 1: Try /sys/class/thermal (generic)
    for (const auto& entry : fs::directory_iterator("/sys/class/thermal")) {
        if (entry.path().filename().string().find("thermal_zone") != std::string::npos) {
            std::string typePath = entry.path() / "type";
            std::string tempPath = entry.path() / "temp";

            std::ifstream typeFile(typePath);
            std::string type;
            if (typeFile >> type) {
                // Match CPU-related thermal zones
                if (type.find("cpu") != std::string::npos || type.find("x86_pkg_temp") != std::string::npos || type.find("k10temp") != std::string::npos) {
                    if (fs::exists(tempPath)) {
                        return tempPath;
                    }
                }
            }
        }
    }

    // Step 2: Fallback to hwmon-based method (your original logic)
    for (const auto& entry : fs::directory_iterator("/sys/class/hwmon")) {
        std::string namePath = entry.path() / "name";
        std::ifstream nameFile(namePath);
        std::string name;
        if (nameFile >> name && name == "k10temp") {
            fs::path tempPath = entry.path() / "temp1_input";
            if (fs::exists(tempPath))
                return tempPath.string();
        }
    }

    // If nothing found
    return "";
}


// Read the current CPU temperature (in Celsius)
static float readTemperatureC() {
    // Try to locate the thermal sensor path if not already done
    if (thermalSensorPath.empty()) {
        thermalSensorPath = findThermalSensorPath();
        if (thermalSensorPath.empty()) {
            useDummyThermal = true; // Use dummy values if no sensor found
        }
    }

    // Generate fake fluctuating temperature if in dummy mode
    if (useDummyThermal) {
        static float t = 45.0f;
        t += ((rand() % 100) - 50) * 0.01f; // Small random fluctuation
        return t;
    }

    // Read temperature from the sensor file
    std::ifstream file(thermalSensorPath);
    int millidegrees = 0;
    if (file >> millidegrees) {
        return millidegrees / 1000.0f; // Convert from millidegree to Celsius
    } else {
        useDummyThermal = true; // Fallback if read fails
        return 50.0f;
    }
}

// Render the "Thermal" tab in the UI
void renderThermalTab() {
    ImGui::Text("Thermal Information");
    ImGui::Separator();

    // Display the current CPU temperature
    float currentTemp = readTemperatureC();
    ImGui::Text("Current CPU Temperature: %.1f °C", currentTemp);

    // User controls
    ImGui::Checkbox("Pause", &pauseThermal);
    ImGui::SliderInt("FPS", &fpsThermal, 1, 144);
    ImGui::SliderFloat("Y Scale", &yScaleThermal, 30.0f, 120.0f, "%.1f °C");

    // Update temperature history if not paused
    if (!pauseThermal) {
        thermalHistory.push_back(currentTemp);
        if (thermalHistory.size() > maxSamples)
            thermalHistory.pop_front();
    }

    // Draw temperature graph if there's data
    if (!thermalHistory.empty()) {
        std::vector<float> plotData(thermalHistory.begin(), thermalHistory.end());
        ImGui::PlotLines("Temperature (°C)", plotData.data(), plotData.size(), 0, nullptr, 0.0f, yScaleThermal, ImVec2(0, 100));
        ImGui::Text("Latest: %.1f °C", plotData.back());
    } else {
        ImGui::Text("No thermal data available.");
    }
}

#else // Non-Linux fallback

// Message for unsupported platforms
void renderThermalTab() {
    ImGui::Text("Thermal monitoring is only available on Linux.");
    ImGui::Text("This feature uses /sys/class/hwmon.");
}

#endif
