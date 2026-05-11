#include "fan.h"
#include <filesystem>
#include <fstream>
#include <string>

int getFanSpeed() {
    namespace fs = std::filesystem;
    try {
        for (const auto& hwmon : fs::directory_iterator("/sys/class/hwmon")) {
            std::string base = hwmon.path().string();
            for (int i = 1; i <= 5; ++i) {
                std::string path = base + "/fan" + std::to_string(i) + "_input";
                std::ifstream in(path);
                int rpm;
                if (in >> rpm && rpm > 0) return rpm;
            }
        }
    } catch (...) {}
    // ASUS-specific fallback
    try {
        for (const auto& hwmon : fs::directory_iterator("/sys/devices/platform/asus-nb-wmi/hwmon")) {
            std::string path = hwmon.path().string() + "/fan1_input";
            std::ifstream in(path);
            int rpm;
            if (in >> rpm && rpm > 0) return rpm;
        }
    } catch (...) {}
    return 0;
}

bool isFanActive() {
    return getFanSpeed() > 0;
}
