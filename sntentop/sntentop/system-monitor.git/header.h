// This is a header file (header_H) that defines various system monitoring functions
// Header guards prevent multiple inclusions of the same file
#ifndef header_H
#define header_H

// GUI libraries for creating graphical interfaces
// ImGui: Immediate Mode GUI library (https://github.com/ocornut/imgui)
#include "imgui.h"
#include "imgui_impl_sdl.h"  // SDL2 backend for ImGui
#include "imgui_impl_opengl3.h"  // OpenGL3 backend for ImGui

// Standard C/C++ libraries
#include <stdio.h>  // Standard I/O operations (https://en.cppreference.com/w/cpp/header/cstdio)
#include <dirent.h>  // Directory operations (https://man7.org/linux/man-pages/man3/readdir.3.html)
#include <vector>  // C++ dynamic array container (https://en.cppreference.com/w/cpp/container/vector)
#include <iostream>  // Standard I/O streams (https://en.cppreference.com/w/cpp/header/iostream)
#include <cmath>  // Math functions (https://en.cppreference.com/w/cpp/header/cmath)
#include <algorithm>  // Algorithms like sort, find (https://en.cppreference.com/w/cpp/header/algorithm)
#include <fstream>  // File stream operations (https://en.cppreference.com/w/cpp/header/fstream)
#include <unistd.h>  // POSIX API for system calls (https://pubs.opengroup.org/onlinepubs/9699919799/)
#include <limits.h>  // Defines size limits of variable types (https://en.cppreference.com/w/cpp/header/climits)
#include <cpuid.h>  // CPU identification (https://man7.org/linux/man-pages/man3/__get_cpuid.3.html)
#include <sys/types.h>  // System data types (https://pubs.opengroup.org/onlinepubs/007908799/xsh/systypes.h.html)
#include <sys/sysinfo.h>  // System information (https://man7.org/linux/man-pages/man2/sysinfo.2.html)
#include <sys/statvfs.h>  // Filesystem statistics (https://man7.org/linux/man-pages/man3/statvfs.3.html)
#include <ctime>  // Time/date utilities (https://en.cppreference.com/w/cpp/header/ctime)
#include <unordered_set>  // Hash-based set container (https://en.cppreference.com/w/cpp/container/unordered_set)
#include <ifaddrs.h>  // Network interface addresses (https://man7.org/linux/man-pages/man3/getifaddrs.3.html)
#include <netinet/in.h>  // Internet address family (https://man7.org/linux/man-pages/man0/netinet_in.h.0p.html)
#include <arpa/inet.h>  // Internet address manipulation (https://man7.org/linux/man-pages/man3/inet.3.html)
#include <map>  // Associative container (https://en.cppreference.com/w/cpp/container/map)

using namespace std;  // Use standard namespace to avoid std:: prefix



// Structure to store process information from /proc/[pid]/stat
// Reference: https://man7.org/linux/man-pages/man5/proc.5.html (search for /proc/[pid]/stat)
struct Proc {
    int pid;
    string name;
    char state;
    float cpuUsage;
    long vsize;  // Virtual memory size in bytes
    long rss;    // Resident set size in bytes
    bool selected = false;  // Add selection flag
};

// Structure to store network interface information
// Reference: https://man7.org/linux/man-pages/man8/ifconfig.8.html
struct NetworkInterface {
    string name;       // Interface name (e.g., eth0, wlan0)
    string ipv4;       // IPv4 address
    string ipv6;       // IPv6 address
    long long rx_bytes;    // Received bytes
    long long tx_bytes;    // Transmitted bytes
    long long rx_packets;  // Received packets
    long long tx_packets;  // Transmitted packets
};

// Function declarations:

// Gets list of network interfaces and their statistics
// Returns vector of NetworkInterface structures
vector<NetworkInterface> getNetworkInterfaces();

// Formats network speed into human-readable string (KB/s, MB/s, etc.)
string formatSpeed(float bytesPerSec); 

// IPv4 address structure
struct IP4
{
    char *name;                   // Interface name
    char addressBuffer[INET_ADDRSTRLEN];  // Buffer for IP address (INET_ADDRSTRLEN=16)
};

// Container for network interfaces
struct Networks
{
    vector<IP4> ip4s;  // List of IPv4 interfaces
};

// Network transmission statistics
struct TX
{
    int bytes;       // Bytes transmitted
    int packets;     // Packets transmitted
    int errs;        // Transmission errors
    int drop;        // Dropped packets
    int fifo;        // FIFO buffer errors
    int frame;       // Frame errors
    int compressed;  // Compressed packets
    int multicast;   // Multicast packets
};

// Network reception statistics
struct RX
{
    int bytes;       // Bytes received
    int packets;     // Packets received
    int errs;        // Receive errors
    int drop;        // Dropped packets
    int fifo;        // FIFO buffer errors
    int colls;       // Collisions
    int carrier;     // Carrier errors
    int compressed;  // Compressed packets
};

// Disk/filesystem information structure
// Similar to output from 'df' command
struct DiskInfo {
    string filesystem;  // Device/partition name
    string size;       // Total size
    string used;       // Used space
    string available;  // Available space
    string usePercent; // Usage percentage
    string mountedOn;  // Mount point
};

// Gets disk information from system
// Returns vector of DiskInfo structures
vector<DiskInfo> getDiskInfo();

// Draws a table with disk information using ImGui
void drawDiskInfoTable(const vector<DiskInfo>& disks);

// Gets CPU information (model, cores, etc.) from /proc/cpuinfo
string CPUinfo();

// Gets operating system name
const char *getOsName();

// Gets system hostname
string getHostname();

// Gets current username
string getUsername();

// Gets current date and time
string getCurrentDateTime();

// Gets current CPU usage percentage
float getCPUUsage();

// Gets CPU temperature (from /sys/class/thermal)
float getCPUTemperature();

// Gets fan speed information (if available)
string getFanInfo(); 

// Gets memory usage information (free, used, etc.)
string getFreeMemoryInfo();

// Gets number of CPU cores
int getNumCores();

// Gets detailed memory statistics
string getMemoryStats();

// Gets list of running processes
vector<Proc> getProcesses();

// Gets total number of processes
int getTotalProcessCount();

// Draws process table using ImGui
void drawProcessTable(const vector<Proc>& processes);

// Updates network statistics
void updateNetworkStats(); 

// Draws network information window using ImGui
void drawNetworkWindow(const char* id, ImVec2 size, ImVec2 position); 

#endif  // End of header guard