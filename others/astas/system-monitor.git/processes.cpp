#include "header.h"
#include <imgui.h>
#include <dirent.h>
#include <cstring>
#include <cctype>
#include <vector>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>
#include <cstdio>
#include <cstdlib>
#include <ctime>
#include <unistd.h>

// -----------------------------
// Structs & Globals
// -----------------------------

// Holds individual process information
struct ProcInfo {
    std::string name;
    char state = '?';
    unsigned long long lastCpuTime = 0;
    float cpuPercent = 0.0f;
    float memPercent = 0.0f;
};

// Global state
static std::unordered_map<int, ProcInfo> processesCpuData;
static std::unordered_set<int> selectedPids;
static unsigned long long lastTotalCpu = 0;
static double lastSampleTime = 0.0;
static const int clockTicksPerSecond = sysconf(_SC_CLK_TCK);
static const int cpuCount = sysconf(_SC_NPROCESSORS_ONLN); // ✅ Add this


// -----------------------------
// Utilities
// -----------------------------

// Return current monotonic time in seconds
double getTimeSeconds() {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return ts.tv_sec + ts.tv_nsec / 1e9;
}

// Parse /proc/stat and return the sum of all CPU time fields
unsigned long long readTotalCpuTime() {
    FILE* file = fopen("/proc/stat", "r");
    if (!file) return 0;

    char line[512];
    if (!fgets(line, sizeof(line), file)) {
        fclose(file);
        return 0;
    }
    fclose(file);

    unsigned long long user, nice, system, idle, iowait, irq, softirq, steal;
    int scanned = sscanf(line, "cpu  %llu %llu %llu %llu %llu %llu %llu %llu",
                         &user, &nice, &system, &idle, &iowait, &irq, &softirq, &steal);
    return (scanned >= 8) ? user + nice + system + idle + iowait + irq + softirq + steal : 0;
}

// Parse /proc/[pid]/stat to extract total CPU time used by a process
unsigned long long readProcessCpuTime(int pid) {
    char path[64];
    snprintf(path, sizeof(path), "/proc/%d/stat", pid);
    FILE* file = fopen(path, "r");
    if (!file) return 0;

    char buffer[1024];
    if (!fgets(buffer, sizeof(buffer), file)) {
        fclose(file);
        return 0;
    }
    fclose(file);

    char* openParen = strchr(buffer, '(');
    char* closeParen = strrchr(buffer, ')');
    if (!openParen || !closeParen || closeParen <= openParen) return 0;

    char* afterName = closeParen + 2;
    const int utimeIndex = 13, stimeIndex = 14;
    unsigned long long utime = 0, stime = 0;
    int index = 0;

    char* saveptr = nullptr;
    char* token = strtok_r(afterName, " ", &saveptr);
    while (token) {
        if (index == utimeIndex) utime = strtoull(token, nullptr, 10);
        if (index == stimeIndex) {
            stime = strtoull(token, nullptr, 10);
            break;
        }
        token = strtok_r(nullptr, " ", &saveptr);
        index++;
    }

    return utime + stime;
}

// Parse /proc/[pid]/status to read memory usage in kB, return percent of total system memory
float readProcessMemoryPercent(int pid) {
    char path[64];
    snprintf(path, sizeof(path), "/proc/%d/status", pid);
    FILE* file = fopen(path, "r");
    if (!file) return 0.0f;

    char line[256];
    unsigned long vmrss = 0;
    while (fgets(line, sizeof(line), file)) {
        if (sscanf(line, "VmRSS: %lu kB", &vmrss) == 1) break;
    }
    fclose(file);

    long totalMemKb = sysconf(_SC_PHYS_PAGES) * sysconf(_SC_PAGE_SIZE) / 1024;
    return (totalMemKb > 0) ? (float)vmrss * 100.0f / totalMemKb : 0.0f;
}

// -----------------------------
// Main UI: Process Table
// -----------------------------

void renderProcessesWindow(const char* id, ImVec2 size, ImVec2 position) {
    static char filter[256] = "";

    ImGui::SetNextWindowSize(size, ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowPos(position, ImGuiCond_FirstUseEver);

    if (!ImGui::Begin(id)) {
        ImGui::End();
        return;
    }

    ImGui::InputText("Filter", filter, sizeof(filter));

    unsigned long long totalCpu = readTotalCpuTime();
    double now = getTimeSeconds();
    double dt = now - lastSampleTime;
    bool canCalculate = (lastSampleTime > 0.0 && dt >= 0.5);

    if (ImGui::BeginTabBar("ProcessTabs")) {
        if (ImGui::BeginTabItem("Processes")) {
            if (ImGui::BeginTable("ProcessTable", 5, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_ScrollY)) {
                ImGui::TableSetupColumn("PID");
                ImGui::TableSetupColumn("Name");
                ImGui::TableSetupColumn("State");
                ImGui::TableSetupColumn("CPU %");
                ImGui::TableSetupColumn("Memory %");
                ImGui::TableHeadersRow();

                DIR* proc = opendir("/proc");
                if (proc) {
                    struct dirent* entry;
                    while ((entry = readdir(proc)) != nullptr) {
                        if (entry->d_type != DT_DIR) continue;
                        if (!std::all_of(entry->d_name, entry->d_name + strlen(entry->d_name), ::isdigit)) continue;

                        int pid = atoi(entry->d_name);

                        // Read process name
                        std::string name = "unknown";
                        char commPath[64];
                        snprintf(commPath, sizeof(commPath), "/proc/%d/comm", pid);
                        if (FILE* commFile = fopen(commPath, "r")) {
                            char buf[256];
                            if (fgets(buf, sizeof(buf), commFile)) {
                                buf[strcspn(buf, "\n")] = '\0';
                                name = buf;
                            }
                            fclose(commFile);
                        }

                        // Read process state
                        char statPath[64];
                        snprintf(statPath, sizeof(statPath), "/proc/%d/stat", pid);
                        char state = '?';
                        if (FILE* statFile = fopen(statPath, "r")) {
                            int dummy;
                            char dummyComm[256];
                            fscanf(statFile, "%d %255s %c", &dummy, dummyComm, &state);
                            fclose(statFile);
                        }

                        // Apply filter
                        std::string nameLower = name;
                        std::transform(nameLower.begin(), nameLower.end(), nameLower.begin(), ::tolower);
                        std::string filterLower = filter;
                        std::transform(filterLower.begin(), filterLower.end(), filterLower.begin(), ::tolower);

                        if (!filterLower.empty() &&
                            nameLower.find(filterLower) == std::string::npos &&
                            std::to_string(pid).find(filterLower) == std::string::npos)
                            continue;

                        // Read CPU & memory
                        unsigned long long currCpu = readProcessCpuTime(pid);
                        float cpuPercent = 0.0f;
                        float memPercent = readProcessMemoryPercent(pid);

                        auto& procInfo = processesCpuData[pid];
                        if (canCalculate && currCpu >= procInfo.lastCpuTime && totalCpu > lastTotalCpu) {
                        unsigned long long deltaProc = currCpu - procInfo.lastCpuTime;
                        unsigned long long deltaTotal = totalCpu - lastTotalCpu;
                        cpuPercent = (deltaProc / (float)deltaTotal) * 100.0f * cpuCount;
                    }



                        // Update proc info
                        procInfo.name = name;
                        procInfo.state = state;
                        procInfo.cpuPercent = cpuPercent;
                        procInfo.memPercent = memPercent;
                        procInfo.lastCpuTime = currCpu;

                        // Render table row
                        ImGui::TableNextRow();
                        ImGui::TableSetColumnIndex(0);

                        bool isSelected = selectedPids.count(pid) > 0;
                        if (ImGui::Selectable(std::to_string(pid).c_str(), isSelected, ImGuiSelectableFlags_SpanAllColumns)) {
                            if (isSelected)
                                selectedPids.erase(pid);
                            else
                                selectedPids.insert(pid);
                        }

                        ImGui::TableSetColumnIndex(1); ImGui::Text("%s", name.c_str());
                        ImGui::TableSetColumnIndex(2); ImGui::Text("%c", state);
                        ImGui::TableSetColumnIndex(3); ImGui::Text("%.2f%%", cpuPercent);
                        ImGui::TableSetColumnIndex(4); ImGui::Text("%.2f%%", memPercent);
                    }
                    closedir(proc);
                } else {
                    ImGui::Text("Failed to open /proc");
                }

                ImGui::EndTable();
            }
            ImGui::EndTabItem();
        }
        ImGui::EndTabBar();
    }

    ImGui::End();

    // ✅ Fixed: Only update after rendering
    lastTotalCpu = totalCpu;
    lastSampleTime = now;
}
