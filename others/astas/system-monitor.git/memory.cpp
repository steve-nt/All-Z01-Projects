#include "header.h"
#include <imgui.h>
#include <utility>
#include <string>

#if defined(__linux__)
    #include <fstream>
    #include <sstream>
    #include <unordered_map>
#elif defined(_WIN32)
    #include <windows.h>
#elif defined(__APPLE__)
    #include <sys/types.h>
    #include <sys/sysctl.h>
    #include <mach/mach.h>
#endif

// Cross-platform memory usage
std::pair<float, float> getMemoryUsageMB() {
#if defined(__linux__)
    std::ifstream meminfo("/proc/meminfo");
    std::string line;
    std::unordered_map<std::string, long> memValues;

    while (std::getline(meminfo, line)) {
        std::istringstream iss(line);
        std::string key;
        long value;
        std::string unit;
        iss >> key >> value >> unit;
        key = key.substr(0, key.size() - 1);  // remove trailing ':'
        memValues[key] = value;
    }

    long memTotal = memValues["MemTotal"];
    long memAvailable = memValues["MemAvailable"];

    float totalMB = memTotal / 1024.0f;
    float usedMB = (memTotal - memAvailable) / 1024.0f;
    return {usedMB, totalMB};

#elif defined(_WIN32)
    MEMORYSTATUSEX memStatus;
    memStatus.dwLength = sizeof(memStatus);
    if (GlobalMemoryStatusEx(&memStatus)) {
        float totalMB = static_cast<float>(memStatus.ullTotalPhys) / (1024.0f * 1024.0f);
        float usedMB = static_cast<float>(memStatus.ullTotalPhys - memStatus.ullAvailPhys) / (1024.0f * 1024.0f);
        return {usedMB, totalMB};
    } else {
        return {0.0f, 0.0f}; // error
    }

#elif defined(__APPLE__)
    int mib[2];
    int64_t physical_memory = 0;
    size_t length = sizeof(physical_memory);
    mib[0] = CTL_HW;
    mib[1] = HW_MEMSIZE;
    sysctl(mib, 2, &physical_memory, &length, nullptr, 0);

    mach_msg_type_number_t count = HOST_VM_INFO64_COUNT;
    vm_statistics64_data_t vmStats;
    mach_port_t host = mach_host_self();
    if (host_statistics64(host, HOST_VM_INFO64, reinterpret_cast<host_info64_t>(&vmStats), &count) == KERN_SUCCESS) {
        int64_t used = (vmStats.active_count + vmStats.inactive_count + vmStats.wire_count) * sysconf(_SC_PAGESIZE);
        float totalMB = physical_memory / (1024.0f * 1024.0f);
        float usedMB = used / (1024.0f * 1024.0f);
        return {usedMB, totalMB};
    } else {
        return {0.0f, 0.0f}; // error
    }

#else
    // Unsupported OS
    return {-1.0f, -1.0f};  // Sentinel value
#endif
}

void renderRAMWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::SetNextWindowSize(size, ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowPos(position, ImGuiCond_FirstUseEver);

    if (!ImGui::Begin(id)) {
        ImGui::End();
        return;
    }

    float usedMB, totalMB;
    std::tie(usedMB, totalMB) = getMemoryUsageMB();

    if (usedMB < 0.0f || totalMB <= 0.0f) {
        ImGui::Text("This OS is not currently supported for RAM monitoring.");
        ImGui::End();
        return;
    }

    float ramPercent = usedMB / totalMB;

    ImGui::Text("Physical Memory (RAM) Usage:");
    ImGui::ProgressBar(ramPercent, ImVec2(-1.0f, 20.0f));
    ImGui::Text("Used: %.1f MB / Total: %.1f MB (%.1f%%)", usedMB, totalMB, ramPercent * 100.0f);

    ImGui::End();
}
