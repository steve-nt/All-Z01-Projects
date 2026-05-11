#include "temp.h"
#include <filesystem>
#include <fstream>
#include <string>

float getCPUTemperature() {
    namespace fs = std::filesystem;
    try {
        for (const auto& hwmon : fs::directory_iterator("/sys/class/hwmon")) {
            std::string base = hwmon.path().string();

            // Προσπαθούμε να εντοπίσουμε έναν αισθητήρα CPU
            for (int i = 1; i <= 5; ++i) {
                std::string labelPath = base + "/temp" + std::to_string(i) + "_label";
                std::ifstream labelFile(labelPath);
                std::string label;

                if (labelFile && std::getline(labelFile, label)) {
                    if (label.find("Package") != std::string::npos || label.find("CPU") != std::string::npos) {
                        std::string inputPath = base + "/temp" + std::to_string(i) + "_input";
                        std::ifstream inputFile(inputPath);
                        int millideg;
                        if (inputFile && inputFile >> millideg) {
                            return millideg / 1000.0f;
                        }
                    }
                }
            }

            // fallback χωρίς label (αν δεν υπάρχει)
            for (int i = 1; i <= 5; ++i) {
                std::string inputPath = base + "/temp" + std::to_string(i) + "_input";
                std::ifstream inputFile(inputPath);
                int millideg;
                if (inputFile && inputFile >> millideg && millideg > 0) {
                    return millideg / 1000.0f;
                }
            }
        }
    } catch (...) {}
    return 0.0f; // δεν βρέθηκε θερμοκρασία
}
