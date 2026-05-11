#include <fstream>
#include <string>
#include <sstream>

struct MemoryStats {
    float totalRamMB = 0.0f;
    float usedRamMB = 0.0f;
    float totalSwapMB = 0.0f;
    float usedSwapMB = 0.0f;
};

MemoryStats getMemoryStats() {
    std::ifstream meminfo("/proc/meminfo");
    std::string line;
    long memTotal = 0, memAvailable = 0, swapTotal = 0, swapFree = 0;

    while (std::getline(meminfo, line)) {
        std::istringstream iss(line);
        std::string key;
        long value;
        std::string unit;
        iss >> key >> value >> unit;

        if (key == "MemTotal:") memTotal = value;
        else if (key == "MemAvailable:") memAvailable = value;
        else if (key == "SwapTotal:") swapTotal = value;
        else if (key == "SwapFree:") swapFree = value;
    }

    MemoryStats stats;
    stats.totalRamMB = memTotal / 1024.0f;
    stats.usedRamMB = (memTotal - memAvailable) / 1024.0f;
    stats.totalSwapMB = swapTotal / 1024.0f;
    stats.usedSwapMB = (swapTotal - swapFree) / 1024.0f;
    return stats;
}
