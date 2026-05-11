#include "process.h"
#include <filesystem>
#include <fstream>
#include <sstream>
#include <string>
#include <vector>
#include <cctype>
#include <algorithm>


namespace fs = std::filesystem;

std::vector<ProcessInfo> getProcesses() {
    std::vector<ProcessInfo> processes;

    for (const auto& entry : fs::directory_iterator("/proc")) {
        if (!entry.is_directory())
            continue;

        std::string filename = entry.path().filename().string();
        bool isNumber = std::all_of(filename.begin(), filename.end(), [](char c) { return std::isdigit(c); });
        if (!isNumber)
            continue;

        int pid = std::stoi(filename);
        std::ifstream statusFile(entry.path() / "status");
        if (!statusFile.is_open())
            continue;

        ProcessInfo proc;
        proc.pid = pid;

        std::string line;
        while (std::getline(statusFile, line)) {
            if (line.find("Name:") == 0) {
                proc.name = line.substr(6);
            } else if (line.find("State:") == 0) {
                proc.state = line.substr(7);
            }
        }

        // Υπολογισμοί CPU/Memory usage (θα υλοποιηθούν στο επόμενο βήμα)
        proc.cpuUsage = 0.0f;
        proc.memUsage = 0.0f;

        processes.push_back(proc);
    }

    return processes;
}
