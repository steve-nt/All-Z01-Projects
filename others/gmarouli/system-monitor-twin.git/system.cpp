#include "header.h"
#include <pwd.h> // for getpwuid()
#include <algorithm>
#include <regex>
#include <dirent.h>
#include <fstream>

// get cpu id and information, you can use `proc/cpuinfo`
string CPUinfo()
{
    char CPUBrandString[0x40];
    unsigned int CPUInfo[4] = {0, 0, 0, 0};

    // unix system
    // for windoes maybe we must add the following
    // __cpuid(regs, 0);
    // regs is the array of 4 positions
    __cpuid(0x80000000, CPUInfo[0], CPUInfo[1], CPUInfo[2], CPUInfo[3]);
    unsigned int nExIds = CPUInfo[0];

    memset(CPUBrandString, 0, sizeof(CPUBrandString));

    for (unsigned int i = 0x80000000; i <= nExIds; ++i)
    {
        __cpuid(i, CPUInfo[0], CPUInfo[1], CPUInfo[2], CPUInfo[3]);

        if (i == 0x80000002)
            memcpy(CPUBrandString, CPUInfo, sizeof(CPUInfo));
        else if (i == 0x80000003)
            memcpy(CPUBrandString + 16, CPUInfo, sizeof(CPUInfo));
        else if (i == 0x80000004)
            memcpy(CPUBrandString + 32, CPUInfo, sizeof(CPUInfo));
    }
    string str(CPUBrandString);
    return str;
}

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

string getUsername()
{
    struct passwd *pw = getpwuid(getuid());
    return pw ? pw->pw_name : "Unknown";
}

string getHostname()
{
    char hostname[HOST_NAME_MAX];
    if (gethostname(hostname, HOST_NAME_MAX) == 0)
        return hostname;
    return "Unknown";
}

int countProcesses()
{
    int count = 0;
    DIR *dir = opendir("/proc");
    if (dir == nullptr)
        return 0;

    struct dirent *entry;
    while ((entry = readdir(dir)) != nullptr)
    {
        if (entry->d_type == DT_DIR)
        {
            if (std::all_of(entry->d_name, entry->d_name + strlen(entry->d_name), ::isdigit))
            {
                count++;
            }
        }
    }
    closedir(dir);
    return count;
}

// Helper: read an integer from a file, return -1 on failure
static int readInt(const std::string &path)
{
    std::ifstream ifs(path);
    int v = -1;
    if (ifs >> v)
        return v;
    return -1;
}

std::vector<int> getFanSpeeds()
{
    std::vector<int> speeds;
    const std::string base = "/sys/class/hwmon/";
    DIR *d = opendir(base.c_str());
    if (!d)
        return speeds;
    struct dirent *e;
    std::regex r("fan[0-9]+_input");
    while ((e = readdir(d)))
    {
        std::string fname = e->d_name;
        if (std::regex_match(fname, r))
        {
            int rpm = readInt(base + fname);
            if (rpm >= 0)
                speeds.push_back(rpm);
        }
    }
    closedir(d);
    return speeds;
}

std::vector<std::pair<std::string, double>> getThermalZones()
{
    std::vector<std::pair<std::string, double>> zones;
    const std::string base = "/sys/class/thermal/";
    DIR *d = opendir(base.c_str());
    if (!d)
        return zones;
    struct dirent *e;
    std::regex r("thermal_zone[0-9]+");
    while ((e = readdir(d)))
    {
        std::string zone = e->d_name;
        if (std::regex_match(zone, r))
        {
            std::string tfile = base + zone + "/temp";
            int mdeg = readInt(tfile);
            if (mdeg >= 0)
            {
                double deg = mdeg / 1000.0;
                zones.emplace_back(zone, deg);
            }
        }
    }
    closedir(d);
    return zones;
}