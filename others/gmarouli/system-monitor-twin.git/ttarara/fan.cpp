#include "fan.h"
#include <filesystem>
#include <fstream>
#include <string>

int getFanSpeed() {
    namespace fs = std::filesystem;

    // Πρώτα κοιτάμε στο /sys/class/hwmon
    for (const auto& hwmon : fs::directory_iterator("/sys/class/hwmon")) {
        std::string base = hwmon.path().string();
        for (int i = 1; i <= 5; ++i) {
            std::string path = base + "/fan" + std::to_string(i) + "_input";
            std::ifstream in(path);
            int rpm;
            if (in >> rpm && rpm > 0) return rpm;
        }
    }

    // ASUS-specific fallback
    for (const auto& hwmon : fs::directory_iterator("/sys/devices/platform/asus-nb-wmi/hwmon")) {
        std::string path = hwmon.path().string() + "/fan1_input";
        std::ifstream in(path);
        int rpm;
        if (in >> rpm && rpm > 0) return rpm;
    }

    return 0;
}


bool isFanActive() {
    return getFanSpeed() > 0;
}
