#include "header.h"
#include <cstdlib>  // for getenv
#include <string>   // for std::string
#include <cstring>  // for memset
#ifdef _WIN32
// Windows-specific networking header.
// Required for functions like gethostname(), socket(), etc. on Windows.
// POSIX systems (Linux/macOS) do not need this; they use unistd.h and others.
#include <winsock2.h>

// Tell the linker to include the Winsock library.
// This provides implementations for networking functions used in this code.
// Without it, you'll get "unresolved external symbol" errors for things like gethostname().
#pragma comment(lib, "ws2_32.lib")
#include <tlhelp32.h>
#else
#include <dirent.h>
#endif

// getOsName, this will get the OS of the current computer
const char *getOsName()
{
#ifdef _WIN32
    return "Windows 32-bit";
#elif _WIN64
    return "Windows 64-bit";
#elif __APPLE__ || __MACH__
    return "Mac OSX";
#elif __linux__
    return "Linux";
#elif __FreeBSD__
    return "FreeBSD";
#elif __unix || __unix__
    return "Unix";
#else
    return "Other";
#endif
}

// getLoggedInUser, this will get the user of the current computer💜
std::string getLoggedInUser()
{
#ifdef _WIN32
    const char* user = std::getenv("USERNAME");
#else
    const char* user = std::getenv("USER");
#endif
    if (user)
        return std::string(user);
    else
        return "Unknown";
}

// getComputerName, this will get the namehost of the current computer💜
std::string getComputerName()
{
    char hostname[256];
    memset(hostname, 0, sizeof(hostname));

#ifdef _WIN32
    // Initialize Winsock
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0)
        return "Unknown";

    if (gethostname(hostname, sizeof(hostname)) == SOCKET_ERROR)
    {
        WSACleanup();
        return "Unknown";
    }

    WSACleanup();
#else
    if (gethostname(hostname, sizeof(hostname)) != 0)
    {
        return "Unknown";
    }
#endif

    return std::string(hostname);
}

//  getTaskStats, this will to get the process states cross-platform💜
TaskStats getTaskStats()
{
    TaskStats stats = {0, 0, 0, 0, 0, 0};

#ifdef _WIN32
    HANDLE hProcessSnap;
    PROCESSENTRY32 pe32;
    hProcessSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);

    if (hProcessSnap == INVALID_HANDLE_VALUE)
        return stats;

    pe32.dwSize = sizeof(PROCESSENTRY32);

    if (!Process32First(hProcessSnap, &pe32))
    {
        CloseHandle(hProcessSnap);
        return stats;
    }

    do
    {
        stats.total++;

        // Windows does not expose direct process states like Unix.
        // By default, assume all running.
        stats.running++;

    } while (Process32Next(hProcessSnap, &pe32));

    CloseHandle(hProcessSnap);

#else
    DIR *dir = opendir("/proc");
    if (!dir)
        return stats;

    struct dirent *entry;
    while ((entry = readdir(dir)) != NULL)
    {
        if (!isdigit(entry->d_name[0]))
            continue;

        stats.total++;

        string statPath = string("/proc/") + entry->d_name + "/stat";
        ifstream statFile(statPath);
        if (!statFile.is_open())
            continue;

        string pid, comm, state;
        statFile >> pid >> comm >> state;

        if (state == "R")
            stats.running++;
        else if (state == "S" || state == "D" || state == "I" || state == "W")
            stats.sleeping++;
        else if (state == "D")
            stats.uninterruptible++;
        else if (state == "T" || state == "t")
            stats.stopped++;
        else if (state == "Z")
            stats.zombie++;

        statFile.close();
    }

    closedir(dir);
#endif

    return stats;
}

//  CPUinfo, this will get the CPU cross-platform💜
// get cpu id and information, you can use `proc/cpuinfo`
std::string CPUinfo()
{
#ifdef _WIN32
    int CPUInfo[4] = {-1};
    char CPUBrandString[0x40];
    __cpuid(CPUInfo, 0x80000000);
    unsigned int nExIds = CPUInfo[0];

    memset(CPUBrandString, 0, sizeof(CPUBrandString));

    for (unsigned int i = 0x80000002; i <= nExIds && i <= 0x80000004; ++i)
    {
        __cpuid(CPUInfo, i);
        memcpy(CPUBrandString + (i - 0x80000002) * 16, CPUInfo, sizeof(CPUInfo));
    }

    return std::string(CPUBrandString);

#elif __linux__
    std::ifstream cpuinfo("/proc/cpuinfo");
    std::string line;
    while (std::getline(cpuinfo, line))
    {
        if (line.substr(0, 10) == "model name")
        {
            return line.substr(line.find(":") + 2);
        }
    }
    return "Unknown CPU";

#elif __APPLE__
    char buffer[256];
    size_t bufferlen = sizeof(buffer);
    if (sysctlbyname("machdep.cpu.brand_string", &buffer, &bufferlen, NULL, 0) == 0)
    {
        return std::string(buffer);
    }
    return "Unknown CPU";

#else
    return "Unknown CPU";
#endif
}
