#include "header.h"
#include <SDL.h>
#include <algorithm>
#include <cctype>
#include <iomanip>
#include <map>
#include <pwd.h>
#include <set>
#include <sstream>

/*
NOTE : You are free to change the code as you wish, the main objective is to make the
       application work and pass the audit.

       It will be provided the main function with the following functions :

       - `void systemWindow(const char *id, ImVec2 size, ImVec2 position)`
            This function will draw the system window on your screen
       - `void memoryProcessesWindow(const char *id, ImVec2 size, ImVec2 position)`
            This function will draw the memory and processes window on your screen
       - `void networkWindow(const char *id, ImVec2 size, ImVec2 position)`
            This function will draw the network window on your screen
*/

// About Desktop OpenGL function loaders:
//  Modern desktop OpenGL doesn't have a standard portable header file to load OpenGL function pointers.
//  Helper libraries are often used for this purpose! Here we are supporting a few common ones (gl3w, glew, glad).
//  You may use another loader/header of your choice (glext, glLoadGen, etc.), or chose to manually implement your own.
#if defined(IMGUI_IMPL_OPENGL_LOADER_GL3W)
#include <GL/gl3w.h> // Initialize with gl3wInit()
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLEW)
#include <GL/glew.h> // Initialize with glewInit()
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD)
#include <glad/glad.h> // Initialize with gladLoadGL()
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD2)
#include <glad/gl.h> // Initialize with gladLoadGL(...) or gladLoaderLoadGL()
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING2)
#define GLFW_INCLUDE_NONE      // GLFW including OpenGL headers causes ambiguity or multiple definition errors.
#include <glbinding/Binding.h> // Initialize with glbinding::Binding::initialize()
#include <glbinding/gl/gl.h>
using namespace gl;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING3)
#define GLFW_INCLUDE_NONE        // GLFW including OpenGL headers causes ambiguity or multiple definition errors.
#include <glbinding/glbinding.h> // Initialize with glbinding::initialize()
#include <glbinding/gl/gl.h>
using namespace gl;
#else
#include IMGUI_IMPL_OPENGL_LOADER_CUSTOM
#endif

struct TaskCounts
{
    int total;
    int running;
    int sleeping;
    int uninterruptible;
    int zombie;
    int stopped;
    int interrupted;
};

struct MemoryInfo
{
    unsigned long long memTotalKb;
    unsigned long long memAvailableKb;
    unsigned long long swapTotalKb;
    unsigned long long swapFreeKb;
};

struct DiskInfo
{
    unsigned long long total;
    unsigned long long used;
    unsigned long long available;
};

struct ProcessRow
{
    int pid;
    string name;
    char state;
    unsigned long long ticks;
    long long rssPages;
    float cpuPercent;
    float memPercent;
};

struct NetDevice
{
    string name;
    RX rx;
    TX tx;
};

struct IPv4Address
{
    string name;
    string address;
};

struct FanInfo
{
    string status;
    string level;
    float speed;
};

struct GraphState
{
    vector<float> values;
    int offset;
    double lastUpdate;
    bool paused;
    float fps;
    float scale;

    GraphState(float initialScale = 100.0f)
        : values(120, 0.0f), offset(0), lastUpdate(0.0), paused(false), fps(5.0f), scale(initialScale)
    {
    }
};

static string trim(const string &value)
{
    size_t start = value.find_first_not_of(" \t\n\r");
    if (start == string::npos)
        return "";
    size_t end = value.find_last_not_of(" \t\n\r");
    return value.substr(start, end - start + 1);
}

static bool readFileLines(const string &path, vector<string> &lines)
{
    ifstream file(path.c_str());
    if (!file)
        return false;

    string line;
    while (getline(file, line))
        lines.push_back(line);
    return true;
}

static bool readUnsignedFromFile(const string &path, unsigned long long &value)
{
    ifstream file(path.c_str());
    if (!file)
        return false;
    file >> value;
    return !file.fail();
}

static string currentUser()
{
    const char *login = getlogin();
    if (login && *login)
        return login;

    const char *user = getenv("USER");
    if (user && *user)
        return user;

    passwd *pwd = getpwuid(getuid());
    return pwd ? pwd->pw_name : "unknown";
}

static string currentHostname()
{
    char hostname[HOST_NAME_MAX + 1];
    if (gethostname(hostname, sizeof(hostname)) == 0)
    {
        hostname[HOST_NAME_MAX] = '\0';
        return hostname;
    }
    return "unknown";
}

static bool readProcessStat(int pid, ProcessRow &process)
{
    string path = "/proc/" + to_string(pid) + "/stat";
    ifstream file(path.c_str());
    if (!file)
        return false;

    string line;
    getline(file, line);
    size_t left = line.find('(');
    size_t right = line.rfind(')');
    if (left == string::npos || right == string::npos || right + 2 >= line.size())
        return false;

    process.pid = pid;
    process.name = line.substr(left + 1, right - left - 1);

    string rest = line.substr(right + 2);
    istringstream fields(rest);
    vector<string> parts;
    string value;
    while (fields >> value)
        parts.push_back(value);

    if (parts.size() < 22)
        return false;

    process.state = parts[0][0];
    unsigned long long utime = strtoull(parts[11].c_str(), NULL, 10);
    unsigned long long stime = strtoull(parts[12].c_str(), NULL, 10);
    process.ticks = utime + stime;
    process.rssPages = strtoll(parts[21].c_str(), NULL, 10);
    process.cpuPercent = 0.0f;
    process.memPercent = 0.0f;
    return true;
}

static vector<ProcessRow> readProcesses()
{
    vector<ProcessRow> processes;
    DIR *proc = opendir("/proc");
    if (!proc)
        return processes;

    dirent *entry;
    while ((entry = readdir(proc)) != NULL)
    {
        if (!isdigit(entry->d_name[0]))
            continue;
        int pid = atoi(entry->d_name);
        ProcessRow process;
        if (readProcessStat(pid, process))
            processes.push_back(process);
    }
    closedir(proc);

    sort(processes.begin(), processes.end(), [](const ProcessRow &a, const ProcessRow &b) {
        return a.pid < b.pid;
    });
    return processes;
}

static TaskCounts readTaskCounts()
{
    TaskCounts counts = {0, 0, 0, 0, 0, 0, 0};
    vector<ProcessRow> processes = readProcesses();
    counts.total = (int)processes.size();

    for (size_t i = 0; i < processes.size(); ++i)
    {
        switch (processes[i].state)
        {
        case 'R':
            counts.running++;
            break;
        case 'S':
            counts.sleeping++;
            break;
        case 'D':
            counts.uninterruptible++;
            break;
        case 'Z':
            counts.zombie++;
            break;
        case 'T':
        case 't':
            counts.stopped++;
            break;
        case 'I':
            counts.interrupted++;
            break;
        default:
            break;
        }
    }
    return counts;
}

static bool readCpuTimes(unsigned long long &idle, unsigned long long &total)
{
    ifstream file("/proc/stat");
    if (!file)
        return false;

    string cpu;
    unsigned long long user = 0, nice = 0, system = 0, idleTime = 0, iowait = 0;
    unsigned long long irq = 0, softirq = 0, steal = 0, guest = 0, guestNice = 0;
    file >> cpu >> user >> nice >> system >> idleTime >> iowait >> irq >> softirq >> steal >> guest >> guestNice;
    idle = idleTime + iowait;
    total = user + nice + system + idleTime + iowait + irq + softirq + steal;
    return cpu == "cpu";
}

static float currentCpuUsage()
{
    static unsigned long long previousIdle = 0;
    static unsigned long long previousTotal = 0;

    unsigned long long idle = 0;
    unsigned long long total = 0;
    if (!readCpuTimes(idle, total))
        return 0.0f;

    unsigned long long idleDelta = idle - previousIdle;
    unsigned long long totalDelta = total - previousTotal;
    previousIdle = idle;
    previousTotal = total;

    if (totalDelta == 0)
        return 0.0f;
    return 100.0f * (float)(totalDelta - idleDelta) / (float)totalDelta;
}

static unsigned long long totalCpuJiffies()
{
    unsigned long long idle = 0;
    unsigned long long total = 0;
    readCpuTimes(idle, total);
    return total;
}

static bool readMemoryInfo(MemoryInfo &info)
{
    info = {0, 0, 0, 0};
    ifstream file("/proc/meminfo");
    if (!file)
        return false;

    string key;
    unsigned long long value;
    string unit;
    while (file >> key >> value >> unit)
    {
        if (key == "MemTotal:")
            info.memTotalKb = value;
        else if (key == "MemAvailable:")
            info.memAvailableKb = value;
        else if (key == "SwapTotal:")
            info.swapTotalKb = value;
        else if (key == "SwapFree:")
            info.swapFreeKb = value;
    }
    return info.memTotalKb > 0;
}

static bool readDiskInfo(DiskInfo &info)
{
    struct statvfs stat;
    if (statvfs("/", &stat) != 0)
        return false;

    info.total = (unsigned long long)stat.f_blocks * stat.f_frsize;
    info.available = (unsigned long long)stat.f_bavail * stat.f_frsize;
    info.used = info.total - ((unsigned long long)stat.f_bfree * stat.f_frsize);
    return info.total > 0;
}

static string formatBytes(double bytes)
{
    const char *units[] = {"B", "KB", "MB", "GB", "TB"};
    int unit = 0;
    while (bytes >= 1024.0 && unit < 4)
    {
        bytes /= 1024.0;
        unit++;
    }

    ostringstream out;
    out << fixed << setprecision(unit == 0 ? 0 : 2) << bytes << " " << units[unit];
    return out.str();
}

static string formatNetworkBytes(unsigned long long bytes)
{
    double value = (double)bytes;
    const char *unit = "B";
    if (bytes >= 1024ULL * 1024ULL * 1024ULL)
    {
        value = value / 1024.0 / 1024.0 / 1024.0;
        unit = "GB";
    }
    else if (bytes >= 1024ULL * 1024ULL)
    {
        value = value / 1024.0 / 1024.0;
        unit = "MB";
    }
    else if (bytes >= 1024ULL)
    {
        value = value / 1024.0;
        unit = "KB";
    }

    ostringstream out;
    out << fixed << setprecision(strcmp(unit, "B") == 0 ? 0 : 2) << value << " " << unit;
    return out.str();
}

static float readThermalTemperature()
{
    vector<string> lines;
    if (readFileLines("/proc/acpi/ibm/thermal", lines))
    {
        for (size_t i = 0; i < lines.size(); ++i)
        {
            size_t colon = lines[i].find(':');
            string values = colon == string::npos ? lines[i] : lines[i].substr(colon + 1);
            istringstream input(values);
            float temperature;
            while (input >> temperature)
            {
                if (temperature > 0.0f)
                    return temperature;
            }
        }
    }

    DIR *thermal = opendir("/sys/class/thermal");
    if (!thermal)
        return 0.0f;

    dirent *entry;
    while ((entry = readdir(thermal)) != NULL)
    {
        string name = entry->d_name;
        if (name.find("thermal_zone") != 0)
            continue;

        unsigned long long raw = 0;
        string path = "/sys/class/thermal/" + name + "/temp";
        if (readUnsignedFromFile(path, raw) && raw > 0)
        {
            closedir(thermal);
            return raw > 1000 ? (float)raw / 1000.0f : (float)raw;
        }
    }
    closedir(thermal);
    return 0.0f;
}

static FanInfo readFanInfo()
{
    FanInfo info = {"unavailable", "unknown", 0.0f};
    vector<string> lines;
    if (readFileLines("/proc/acpi/ibm/fan", lines))
    {
        for (size_t i = 0; i < lines.size(); ++i)
        {
            size_t colon = lines[i].find(':');
            if (colon == string::npos)
                continue;
            string key = trim(lines[i].substr(0, colon));
            string value = trim(lines[i].substr(colon + 1));
            if (key == "status")
                info.status = value;
            else if (key == "speed")
                info.speed = (float)atof(value.c_str());
            else if (key == "level")
                info.level = value;
        }
        return info;
    }

    DIR *hwmon = opendir("/sys/class/hwmon");
    if (!hwmon)
        return info;

    dirent *hwmonEntry;
    while ((hwmonEntry = readdir(hwmon)) != NULL)
    {
        string hwmonName = hwmonEntry->d_name;
        if (hwmonName == "." || hwmonName == "..")
            continue;
        string hwmonPath = "/sys/class/hwmon/" + hwmonName;
        DIR *device = opendir(hwmonPath.c_str());
        if (!device)
            continue;

        dirent *deviceEntry;
        while ((deviceEntry = readdir(device)) != NULL)
        {
            string fileName = deviceEntry->d_name;
            if (fileName.find("fan") == 0 && fileName.find("_input") != string::npos)
            {
                unsigned long long speed = 0;
                if (readUnsignedFromFile(hwmonPath + "/" + fileName, speed))
                {
                    info.speed = (float)speed;
                    info.status = speed > 0 ? "active" : "inactive";
                    info.level = "unknown";
                    closedir(device);
                    closedir(hwmon);
                    return info;
                }
            }
        }
        closedir(device);
    }
    closedir(hwmon);
    return info;
}

static void pushGraphValue(GraphState &graph, float value)
{
    if (graph.values.empty())
        graph.values.assign(120, 0.0f);
    graph.values[graph.offset] = value;
    graph.offset = (graph.offset + 1) % (int)graph.values.size();
}

static void updateGraph(GraphState &graph, float value)
{
    double now = ImGui::GetTime();
    float fps = max(1.0f, graph.fps);
    if (!graph.paused && now - graph.lastUpdate >= 1.0 / fps)
    {
        pushGraphValue(graph, value);
        graph.lastUpdate = now;
    }
}

static void drawGraphControls(const char *name, GraphState &graph)
{
    string pauseLabel = string("Pause##") + name;
    string fpsLabel = string("FPS##") + name;
    string scaleLabel = string("Y scale##") + name;

    ImGui::Checkbox(pauseLabel.c_str(), &graph.paused);
    ImGui::SliderFloat(fpsLabel.c_str(), &graph.fps, 1.0f, 30.0f, "%.0f");
    ImGui::SliderFloat(scaleLabel.c_str(), &graph.scale, 1.0f, 5000.0f, "%.0f");
}

static void drawGraph(const char *label, GraphState &graph, float current, const string &overlay)
{
    updateGraph(graph, current);
    ImGui::PlotLines(label, graph.values.data(), (int)graph.values.size(), graph.offset, overlay.c_str(), 0.0f, graph.scale, ImVec2(0, 110));
}

static vector<ProcessRow> readProcessRowsWithUsage()
{
    static map<int, unsigned long long> previousTicks;
    static unsigned long long previousTotal = 0;

    vector<ProcessRow> processes = readProcesses();
    unsigned long long total = totalCpuJiffies();
    unsigned long long totalDelta = previousTotal == 0 || total < previousTotal ? 0 : total - previousTotal;
    long pageSize = sysconf(_SC_PAGESIZE);

    MemoryInfo memory;
    readMemoryInfo(memory);
    unsigned long long totalBytes = memory.memTotalKb * 1024ULL;
    int cpuCount = (int)max(1L, sysconf(_SC_NPROCESSORS_ONLN));

    map<int, unsigned long long> currentTicks;
    for (size_t i = 0; i < processes.size(); ++i)
    {
        currentTicks[processes[i].pid] = processes[i].ticks;
        map<int, unsigned long long>::iterator previous = previousTicks.find(processes[i].pid);
        if (previous != previousTicks.end() && totalDelta > 0 && processes[i].ticks >= previous->second)
        {
            unsigned long long processDelta = processes[i].ticks - previous->second;
            processes[i].cpuPercent = 100.0f * (float)processDelta / (float)totalDelta * (float)cpuCount;
        }

        unsigned long long processBytes = (unsigned long long)max(0LL, processes[i].rssPages) * (unsigned long long)pageSize;
        if (totalBytes > 0)
            processes[i].memPercent = 100.0f * (float)processBytes / (float)totalBytes;
    }

    previousTicks.swap(currentTicks);
    previousTotal = total;
    return processes;
}

static vector<IPv4Address> readIPv4Addresses()
{
    vector<IPv4Address> addresses;
    ifaddrs *interfaces = NULL;
    if (getifaddrs(&interfaces) == -1)
        return addresses;

    for (ifaddrs *ifa = interfaces; ifa != NULL; ifa = ifa->ifa_next)
    {
        if (!ifa->ifa_addr || ifa->ifa_addr->sa_family != AF_INET)
            continue;

        char buffer[INET_ADDRSTRLEN];
        sockaddr_in *addr = (sockaddr_in *)ifa->ifa_addr;
        if (inet_ntop(AF_INET, &addr->sin_addr, buffer, sizeof(buffer)))
            addresses.push_back({ifa->ifa_name, buffer});
    }

    freeifaddrs(interfaces);
    sort(addresses.begin(), addresses.end(), [](const IPv4Address &a, const IPv4Address &b) {
        return a.name < b.name;
    });
    return addresses;
}

static vector<NetDevice> readNetworkDevices()
{
    vector<NetDevice> devices;
    ifstream file("/proc/net/dev");
    if (!file)
        return devices;

    string line;
    getline(file, line);
    getline(file, line);
    while (getline(file, line))
    {
        size_t colon = line.find(':');
        if (colon == string::npos)
            continue;

        NetDevice device;
        device.name = trim(line.substr(0, colon));
        istringstream values(line.substr(colon + 1));
        values >> device.rx.bytes >> device.rx.packets >> device.rx.errs >> device.rx.drop >> device.rx.fifo >> device.rx.frame >> device.rx.compressed >> device.rx.multicast;
        values >> device.tx.bytes >> device.tx.packets >> device.tx.errs >> device.tx.drop >> device.tx.fifo >> device.tx.colls >> device.tx.carrier >> device.tx.compressed;
        if (!values.fail())
            devices.push_back(device);
    }
    return devices;
}

static void drawUsageBar(const char *label, unsigned long long used, unsigned long long total)
{
    float fraction = total == 0 ? 0.0f : min(1.0f, (float)((double)used / (double)total));
    string overlay = formatBytes((double)used) + " / " + formatBytes((double)total);
    ImGui::Text("%s", label);
    ImGui::ProgressBar(fraction, ImVec2(-1, 0), overlay.c_str());
}

// systemWindow, display information for the system monitorization
void systemWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::Begin(id);
    ImGui::SetWindowSize(id, size);
    ImGui::SetWindowPos(id, position);

    TaskCounts counts = readTaskCounts();
    ImGui::Text("Operating System: %s", getOsName());
    ImGui::Text("User: %s", currentUser().c_str());
    ImGui::Text("Hostname: %s", currentHostname().c_str());
    ImGui::Separator();
    ImGui::Text("Tasks: %d total", counts.total);
    ImGui::Text("Running: %d  Sleeping: %d  Uninterruptible: %d", counts.running, counts.sleeping, counts.uninterruptible);
    ImGui::Text("Zombie: %d  Traced/Stopped: %d  Interrupted/Idle: %d", counts.zombie, counts.stopped, counts.interrupted);
    ImGui::Separator();
    ImGui::TextWrapped("CPU: %s", CPUinfo().c_str());
    ImGui::Separator();

    static GraphState cpuGraph(100.0f);
    static GraphState fanGraph(5000.0f);
    static GraphState thermalGraph(100.0f);

    if (ImGui::BeginTabBar("SystemTabs"))
    {
        if (ImGui::BeginTabItem("CPU"))
        {
            float cpu = currentCpuUsage();
            ostringstream overlay;
            overlay << fixed << setprecision(1) << cpu << "% CPU";
            drawGraph("CPU usage", cpuGraph, cpu, overlay.str());
            drawGraphControls("cpu", cpuGraph);
            ImGui::EndTabItem();
        }
        if (ImGui::BeginTabItem("Fan"))
        {
            FanInfo fan = readFanInfo();
            ImGui::Text("Status: %s", fan.status.c_str());
            ImGui::Text("Speed: %.0f RPM", fan.speed);
            ImGui::Text("Level: %s", fan.level.c_str());
            ostringstream overlay;
            overlay << fixed << setprecision(0) << fan.speed << " RPM";
            drawGraph("Fan speed", fanGraph, fan.speed, overlay.str());
            drawGraphControls("fan", fanGraph);
            ImGui::EndTabItem();
        }
        if (ImGui::BeginTabItem("Thermal"))
        {
            float temperature = readThermalTemperature();
            ostringstream overlay;
            overlay << fixed << setprecision(1) << temperature << " C";
            drawGraph("Temperature", thermalGraph, temperature, overlay.str());
            drawGraphControls("thermal", thermalGraph);
            ImGui::EndTabItem();
        }
        ImGui::EndTabBar();
    }

    ImGui::End();
}

// memoryProcessesWindow, display information for the memory and processes information
void memoryProcessesWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::Begin(id);
    ImGui::SetWindowSize(id, size);
    ImGui::SetWindowPos(id, position);

    MemoryInfo memory;
    if (readMemoryInfo(memory))
    {
        unsigned long long memUsed = (memory.memTotalKb - memory.memAvailableKb) * 1024ULL;
        unsigned long long memTotal = memory.memTotalKb * 1024ULL;
        unsigned long long swapUsed = (memory.swapTotalKb - memory.swapFreeKb) * 1024ULL;
        unsigned long long swapTotal = memory.swapTotalKb * 1024ULL;

        drawUsageBar("RAM", memUsed, memTotal);
        drawUsageBar("SWAP", swapUsed, swapTotal);
    }
    else
    {
        ImGui::Text("Memory information unavailable");
    }

    DiskInfo disk;
    if (readDiskInfo(disk))
        drawUsageBar("Disk /", disk.used, disk.total);
    else
        ImGui::Text("Disk information unavailable");

    ImGui::Separator();
    static char filter[128] = "";
    static set<int> selectedPids;
    ImGui::InputText("Filter processes", filter, IM_ARRAYSIZE(filter));

    vector<ProcessRow> processes = readProcessRowsWithUsage();
    string needle = filter;
    transform(needle.begin(), needle.end(), needle.begin(), [](unsigned char c) { return (char)tolower(c); });

    if (ImGui::BeginTable("ProcessTable", 5, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_ScrollY | ImGuiTableFlags_Resizable, ImVec2(0, 0)))
    {
        ImGui::TableSetupColumn("PID");
        ImGui::TableSetupColumn("Name");
        ImGui::TableSetupColumn("State");
        ImGui::TableSetupColumn("CPU usage");
        ImGui::TableSetupColumn("Memory usage");
        ImGui::TableHeadersRow();

        for (size_t i = 0; i < processes.size(); ++i)
        {
            string lowerName = processes[i].name;
            transform(lowerName.begin(), lowerName.end(), lowerName.begin(), [](unsigned char c) { return (char)tolower(c); });
            if (!needle.empty() && lowerName.find(needle) == string::npos && to_string(processes[i].pid).find(needle) == string::npos)
                continue;

            ImGui::TableNextRow();
            ImGui::TableNextColumn();
            bool selected = selectedPids.find(processes[i].pid) != selectedPids.end();
            string rowLabel = to_string(processes[i].pid) + "##proc" + to_string(processes[i].pid);
            if (ImGui::Selectable(rowLabel.c_str(), selected, ImGuiSelectableFlags_SpanAllColumns))
            {
                if (selected)
                    selectedPids.erase(processes[i].pid);
                else
                    selectedPids.insert(processes[i].pid);
            }
            ImGui::TableNextColumn();
            ImGui::TextUnformatted(processes[i].name.c_str());
            ImGui::TableNextColumn();
            ImGui::Text("%c", processes[i].state);
            ImGui::TableNextColumn();
            ImGui::Text("%.1f%%", processes[i].cpuPercent);
            ImGui::TableNextColumn();
            ImGui::Text("%.1f%%", processes[i].memPercent);
        }
        ImGui::EndTable();
    }

    ImGui::End();
}

// network, display information network information
void networkWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::Begin(id);
    ImGui::SetWindowSize(id, size);
    ImGui::SetWindowPos(id, position);

    vector<IPv4Address> addresses = readIPv4Addresses();
    ImGui::Text("IPv4 addresses");
    for (size_t i = 0; i < addresses.size(); ++i)
        ImGui::BulletText("%s: %s", addresses[i].name.c_str(), addresses[i].address.c_str());

    vector<NetDevice> devices = readNetworkDevices();
    ImGui::Separator();
    if (ImGui::BeginTabBar("NetworkTables"))
    {
        if (ImGui::BeginTabItem("RX"))
        {
            if (ImGui::BeginTable("RXTable", 9, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_Resizable))
            {
                const char *columns[] = {"Interface", "bytes", "packets", "errs", "drop", "fifo", "frame", "compressed", "multicast"};
                for (int i = 0; i < 9; ++i)
                    ImGui::TableSetupColumn(columns[i]);
                ImGui::TableHeadersRow();
                for (size_t i = 0; i < devices.size(); ++i)
                {
                    ImGui::TableNextRow();
                    ImGui::TableNextColumn();
                    ImGui::TextUnformatted(devices[i].name.c_str());
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.bytes);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.packets);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.errs);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.drop);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.fifo);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.frame);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.compressed);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].rx.multicast);
                }
                ImGui::EndTable();
            }
            ImGui::EndTabItem();
        }
        if (ImGui::BeginTabItem("TX"))
        {
            if (ImGui::BeginTable("TXTable", 9, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_Resizable))
            {
                const char *columns[] = {"Interface", "bytes", "packets", "errs", "drop", "fifo", "colls", "carrier", "compressed"};
                for (int i = 0; i < 9; ++i)
                    ImGui::TableSetupColumn(columns[i]);
                ImGui::TableHeadersRow();
                for (size_t i = 0; i < devices.size(); ++i)
                {
                    ImGui::TableNextRow();
                    ImGui::TableNextColumn();
                    ImGui::TextUnformatted(devices[i].name.c_str());
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.bytes);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.packets);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.errs);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.drop);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.fifo);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.colls);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.carrier);
                    ImGui::TableNextColumn();
                    ImGui::Text("%llu", devices[i].tx.compressed);
                }
                ImGui::EndTable();
            }
            ImGui::EndTabItem();
        }
        ImGui::EndTabBar();
    }

    ImGui::Separator();
    if (ImGui::BeginTabBar("NetworkUsage"))
    {
        const double maxBytes = 2.0 * 1024.0 * 1024.0 * 1024.0;
        if (ImGui::BeginTabItem("RX usage"))
        {
            for (size_t i = 0; i < devices.size(); ++i)
            {
                float fraction = min(1.0f, (float)((double)devices[i].rx.bytes / maxBytes));
                string overlay = formatNetworkBytes(devices[i].rx.bytes) + " / 2.00 GB";
                ImGui::Text("%s", devices[i].name.c_str());
                ImGui::ProgressBar(fraction, ImVec2(-1, 0), overlay.c_str());
            }
            ImGui::EndTabItem();
        }
        if (ImGui::BeginTabItem("TX usage"))
        {
            for (size_t i = 0; i < devices.size(); ++i)
            {
                float fraction = min(1.0f, (float)((double)devices[i].tx.bytes / maxBytes));
                string overlay = formatNetworkBytes(devices[i].tx.bytes) + " / 2.00 GB";
                ImGui::Text("%s", devices[i].name.c_str());
                ImGui::ProgressBar(fraction, ImVec2(-1, 0), overlay.c_str());
            }
            ImGui::EndTabItem();
        }
        ImGui::EndTabBar();
    }

    ImGui::End();
}

// Main code
int main(int, char **)
{
    // Setup SDL
    // (Some versions of SDL before <2.0.10 appears to have performance/stalling issues on a minority of Windows systems,
    // depending on whether SDL_INIT_GAMECONTROLLER is enabled or disabled.. updating to latest version of SDL is recommended!)
    if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_TIMER | SDL_INIT_GAMECONTROLLER) != 0)
    {
        printf("Error: %s\n", SDL_GetError());
        return -1;
    }

    // GL 3.0 + GLSL 130
    const char *glsl_version = "#version 130";
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_FLAGS, 0);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_PROFILE_MASK, SDL_GL_CONTEXT_PROFILE_CORE);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MAJOR_VERSION, 3);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MINOR_VERSION, 0);

    // Create window with graphics context
    SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER, 1);
    SDL_GL_SetAttribute(SDL_GL_DEPTH_SIZE, 24);
    SDL_GL_SetAttribute(SDL_GL_STENCIL_SIZE, 8);
    SDL_WindowFlags window_flags = (SDL_WindowFlags)(SDL_WINDOW_OPENGL | SDL_WINDOW_RESIZABLE | SDL_WINDOW_ALLOW_HIGHDPI);
    SDL_Window *window = SDL_CreateWindow("Dear ImGui SDL2+OpenGL3 example", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, 1280, 720, window_flags);
    SDL_GLContext gl_context = SDL_GL_CreateContext(window);
    SDL_GL_MakeCurrent(window, gl_context);
    SDL_GL_SetSwapInterval(1); // Enable vsync

    // Initialize OpenGL loader
#if defined(IMGUI_IMPL_OPENGL_LOADER_GL3W)
    bool err = gl3wInit() != 0;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLEW)
    bool err = glewInit() != GLEW_OK;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD)
    bool err = gladLoadGL() == 0;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD2)
    bool err = gladLoadGL((GLADloadfunc)SDL_GL_GetProcAddress) == 0; // glad2 recommend using the windowing library loader instead of the (optionally) bundled one.
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING2)
    bool err = false;
    glbinding::Binding::initialize();
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING3)
    bool err = false;
    glbinding::initialize([](const char *name) { return (glbinding::ProcAddress)SDL_GL_GetProcAddress(name); });
#else
    bool err = false; // If you use IMGUI_IMPL_OPENGL_LOADER_CUSTOM, your loader is likely to requires some form of initialization.
#endif
    if (err)
    {
        fprintf(stderr, "Failed to initialize OpenGL loader!\n");
        return 1;
    }

    // Setup Dear ImGui context
    IMGUI_CHECKVERSION();
    ImGui::CreateContext();
    // render bindings
    ImGuiIO &io = ImGui::GetIO();

    // Setup Dear ImGui style
    ImGui::StyleColorsDark();

    // Setup Platform/Renderer backends
    ImGui_ImplSDL2_InitForOpenGL(window, gl_context);
    ImGui_ImplOpenGL3_Init(glsl_version);

    // background color
    // note : you are free to change the style of the application
    ImVec4 clear_color = ImVec4(0.0f, 0.0f, 0.0f, 0.0f);

    // Main loop
    bool done = false;
    while (!done)
    {
        // Poll and handle events (inputs, window resize, etc.)
        // You can read the io.WantCaptureMouse, io.WantCaptureKeyboard flags to tell if dear imgui wants to use your inputs.
        // - When io.WantCaptureMouse is true, do not dispatch mouse input data to your main application.
        // - When io.WantCaptureKeyboard is true, do not dispatch keyboard input data to your main application.
        // Generally you may always pass all inputs to dear imgui, and hide them from your application based on those two flags.
        SDL_Event event;
        while (SDL_PollEvent(&event))
        {
            ImGui_ImplSDL2_ProcessEvent(&event);
            if (event.type == SDL_QUIT)
                done = true;
            if (event.type == SDL_WINDOWEVENT && event.window.event == SDL_WINDOWEVENT_CLOSE && event.window.windowID == SDL_GetWindowID(window))
                done = true;
        }

        // Start the Dear ImGui frame
        ImGui_ImplOpenGL3_NewFrame();
        ImGui_ImplSDL2_NewFrame(window);
        ImGui::NewFrame();

        {
            ImVec2 mainDisplay = io.DisplaySize;
            memoryProcessesWindow("== Memory and Processes ==",
                                  ImVec2((mainDisplay.x / 2) - 20, (mainDisplay.y / 2) + 30),
                                  ImVec2((mainDisplay.x / 2) + 10, 10));
            // --------------------------------------
            systemWindow("== System ==",
                         ImVec2((mainDisplay.x / 2) - 10, (mainDisplay.y / 2) + 30),
                         ImVec2(10, 10));
            // --------------------------------------
            networkWindow("== Network ==",
                          ImVec2(mainDisplay.x - 20, (mainDisplay.y / 2) - 60),
                          ImVec2(10, (mainDisplay.y / 2) + 50));
        }

        // Rendering
        ImGui::Render();
        glViewport(0, 0, (int)io.DisplaySize.x, (int)io.DisplaySize.y);
        glClearColor(clear_color.x, clear_color.y, clear_color.z, clear_color.w);
        glClear(GL_COLOR_BUFFER_BIT);
        ImGui_ImplOpenGL3_RenderDrawData(ImGui::GetDrawData());
        SDL_GL_SwapWindow(window);
    }

    // Cleanup
    ImGui_ImplOpenGL3_Shutdown();
    ImGui_ImplSDL2_Shutdown();
    ImGui::DestroyContext();

    SDL_GL_DeleteContext(gl_context);
    SDL_DestroyWindow(window);
    SDL_Quit();

    return 0;
}
