#include "header.h"
#include <imgui.h>
#include <cstdio>

#if defined(_WIN32)
    #include <windows.h>
#elif defined(__APPLE__)
    #include <sys/mount.h>
#elif defined(__linux__)
    #include <sys/statvfs.h>
#endif

void renderDiskWindow(const char* id, ImVec2 size, ImVec2 position)
{
    ImGui::SetNextWindowSize(size, ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowPos(position, ImGuiCond_FirstUseEver);

    if (!ImGui::Begin(id)) {
        ImGui::End();
        return;
    }

    float usedPercent = 0.0f;
    float totalGB = 0.0f, usedGB = 0.0f, availGB = 0.0f;

#if defined(_WIN32)
    // Windows implementation
    ULARGE_INTEGER freeBytesAvailable, totalBytes, freeBytes;
    if (GetDiskFreeSpaceExW(L"C:\\", &freeBytesAvailable, &totalBytes, &freeBytes)) {
        unsigned long long total = totalBytes.QuadPart;
        unsigned long long free = freeBytes.QuadPart;
        unsigned long long used = total - free;

        usedPercent = (float)used / (float)total;
        totalGB = total / (1024.0f * 1024.0f * 1024.0f);
        usedGB = used / (1024.0f * 1024.0f * 1024.0f);
        availGB = free / (1024.0f * 1024.0f * 1024.0f);
    } else {
        ImGui::Text("Failed to get disk stats (Windows)");
        ImGui::End();
        return;
    }

#elif defined(__APPLE__)
    // macOS implementation
    struct statfs stats;
    if (statfs("/", &stats) == 0) {
        unsigned long long total = stats.f_blocks * stats.f_bsize;
        unsigned long long free = stats.f_bfree * stats.f_bsize;
        unsigned long long available = stats.f_bavail * stats.f_bsize;
        unsigned long long used = total - free;

        usedPercent = (used + available > 0) ? (float)used / (float)(used + available) : 0.0f;
        totalGB = total / (1024.0f * 1024.0f * 1024.0f);
        usedGB = used / (1024.0f * 1024.0f * 1024.0f);
        availGB = available / (1024.0f * 1024.0f * 1024.0f);
    } else {
        ImGui::Text("Failed to get disk stats (macOS)");
        ImGui::End();
        return;
    }

#elif defined(__linux__)
    // Linux implementation
    struct statvfs stats;
    if (statvfs("/", &stats) == 0) {
        unsigned long long blockSize = stats.f_frsize;
        unsigned long long total = stats.f_blocks * blockSize;
        unsigned long long free = stats.f_bfree * blockSize;
        unsigned long long available = stats.f_bavail * blockSize;
        unsigned long long used = total - free;

        usedPercent = (used + available > 0) ? (float)used / (float)(used + available) : 0.0f;
        totalGB = total / (1024.0f * 1024.0f * 1024.0f);
        usedGB = used / (1024.0f * 1024.0f * 1024.0f);
        availGB = available / (1024.0f * 1024.0f * 1024.0f);
    } else {
        ImGui::Text("Failed to get disk stats (Linux)");
        ImGui::End();
        return;
    }

#else
    // Unsupported platform
    ImGui::Text("This OS is not currently supported for disk usage monitoring.");
    ImGui::End();
    return;
#endif

    // Common rendering
    ImGui::Text("Disk Usage for /");
    ImGui::ProgressBar(usedPercent, ImVec2(-1.0f, 20.0f));
    ImGui::Text("Used: %.1f GB / Total: %.1f GB (%.1f%%)", usedGB, totalGB, usedPercent * 100.0f);
    ImGui::Text("Available: %.1f GB", availGB);

    ImGui::End();
}
