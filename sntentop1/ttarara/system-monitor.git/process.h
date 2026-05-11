#pragma once
#include <string>
#include <vector>

struct ProcessInfo {
    int pid;                // Process ID
    std::string name;       // Name from /proc/[pid]/status
    std::string state;      // Process state (R/S/Z/...)
    float cpuUsage = 0.0f;  // To be calculated
    float memUsage = 0.0f;  // To be calculated
};

std::vector<ProcessInfo> getProcesses();
