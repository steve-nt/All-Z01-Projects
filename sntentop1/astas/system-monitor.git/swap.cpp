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

struct SwapStats {
    float usedMB = 0.0f;
    float totalMB = 0.0f;
    std::string errorMessage; // empty if no error
};

// Returns SwapStats with errorMessage set if unsupported or failed
SwapStats getSwapInfo()
{
#if defined(__linux__)
    std::ifstream meminfo("/proc/meminfo");
    if (!meminfo.is_open())
        return {0, 0, "Failed to open /proc/meminfo"};

    std::string line;
    std::unordered_map<std::string, long> memValues;

    while (std::getline(meminfo, line)) {
        std::istringstream iss(line);
        std::string key;
        long value;
        std::string unit;
        iss >> key >> value >> unit;
        key = key.substr(0, key.size() - 1); // remove trailing ':'
        memValues[key] = value;
    }

    long swapTotal = memValues["SwapTotal"];
    long swapFree = memValues["SwapFree"];

    return {
        (swapTotal - swapFree) / 1024.0f,
        swapTotal / 1024.0f,
        ""
    };

#elif defined(_WIN32)
    MEMORYSTATUSEX memStatus;
    memStatus.dwLength = sizeof(memStatus);
    if (GlobalMemoryStatusEx(&memStatus)) {
        float totalMB = static_cast<float>(memStatus.ullTotalPageFile) / (1024.0f * 1024.0f);
        float availMB = static_cast<float>(memStatus.ullAvailPageFile) / (1024.0f * 1024.0f);
        float usedMB = totalMB - availMB;
        return { usedMB, totalMB, "" };
    } else {
        return { 0.0f, 0.0f, "Failed to get Windows memory status" };
    }

#elif defined(__APPLE__)
    mach_msg_type_number_t count = HOST_VM_INFO64_COUNT;
    vm_statistics64_data_t vmStats;
    mach_port_t host = mach_host_self();
    if (host_statistics64(host, HOST_VM_INFO64, reinterpret_cast<host_info64_t>(&vmStats), &count) == KERN_SUCCESS) {
        long pageSize = sysconf(_SC_PAGESIZE);
        float usedMB = (float)(vmStats.pageouts * pageSize) / (1024.0f * 1024.0f);
        // Total swap unknown on macOS
        return { usedMB, 0.0f, "" };
    } else {
        return { 0.0f, 0.0f, "Failed to get macOS vm statistics" };
    }

#else
    return { 0.0f, 0.0f, "Swap monitoring is not supported on this OS." };
#endif
}

void renderSwapWindow(const char* id, ImVec2 size, ImVec2 position)
{
    ImGui::SetNextWindowSize(size, ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowPos(position, ImGuiCond_FirstUseEver);

    if (!ImGui::Begin(id)) {
        ImGui::End();
        return;
    }

    SwapStats swap = getSwapInfo();

    if (!swap.errorMessage.empty()) {
        ImGui::TextColored(ImVec4(1,0,0,1), "%s", swap.errorMessage.c_str());
        ImGui::End();
        return;
    }

    if (swap.totalMB <= 0.0f) {
        ImGui::Text("Swap Used: %.1f MB (Total swap unknown)", swap.usedMB);
    } else {
        float swapPercent = (swap.totalMB > 0.0f) ? (swap.usedMB / swap.totalMB) : 0.0f;
        ImGui::Text("Swap (Virtual Memory) Usage:");
        ImGui::ProgressBar(swapPercent, ImVec2(-1.0f, 20.0f));
        ImGui::Text("Used: %.1f MB / Total: %.1f MB (%.1f%%)", swap.usedMB, swap.totalMB, swapPercent * 100.0f);
    }

    ImGui::End();
}
