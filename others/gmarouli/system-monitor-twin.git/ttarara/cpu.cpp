#include <fstream>
#include <string>
#include <sstream>
#include <unistd.h>
#include "cpu.h"

float getCPUUsage() {
    static long long lastTotalUser, lastTotalUserLow, lastTotalSys, lastTotalIdle;

    std::ifstream file("/proc/stat");
    std::string line;
    std::getline(file, line);
    std::istringstream iss(line);

    std::string cpu;
    long long user, nice, system, idle, iowait, irq, softirq;
    iss >> cpu >> user >> nice >> system >> idle >> iowait >> irq >> softirq;

    long long totalUser = user;
    long long totalUserLow = nice;
    long long totalSys = system;
    long long totalIdle = idle;

    long long total = (totalUser - lastTotalUser) +
                      (totalUserLow - lastTotalUserLow) +
                      (totalSys - lastTotalSys);
    long long totalAll = total + (totalIdle - lastTotalIdle);

    lastTotalUser = totalUser;
    lastTotalUserLow = totalUserLow;
    lastTotalSys = totalSys;
    lastTotalIdle = totalIdle;

    if (totalAll == 0) return 0.0f;
    return (total * 100.0f) / totalAll;
}

std::string getCPUModel() {
    std::ifstream file("/proc/cpuinfo");
    std::string line;
    while (std::getline(file, line)) {
        if (line.find("model name") != std::string::npos) {
            return line.substr(line.find(":") + 2);
        }
    }
    return "Unknown CPU";
}