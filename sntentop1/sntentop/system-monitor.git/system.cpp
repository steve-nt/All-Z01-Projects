// This header likely contains custom type definitions (like 'string') and function prototypes
// Sources: 
// - For general header usage: https://en.cppreference.com/w/cpp/preprocessor/include
#include "header.h"

// unistd.h provides access to POSIX operating system API (e.g., gethostname, getlogin_r)
// Sources:
// - POSIX API reference: https://pubs.opengroup.org/onlinepubs/9699919799/
#include <unistd.h>       

// time.h provides time/date manipulation functions (e.g., time, localtime, strftime)
// Sources:
// - C time reference: https://en.cppreference.com/w/c/chrono
#include <time.h>         

// stdio.h provides standard input/output functions (e.g., FILE operations, printf)
// Sources:
// - C I/O reference: https://en.cppreference.com/w/c/io
#include <stdio.h>        

// string.h provides string manipulation functions (e.g., memset, memcpy)
// Sources:
// - C string reference: https://en.cppreference.com/w/c/string/byte
#include <string.h>       

// CPUinfo() - Retrieves the CPU brand string using CPUID instruction
// Returns: string containing full CPU name
// Sources:
// - CPUID instruction: https://en.wikipedia.org/wiki/CPUID
// - Intel CPUID reference: https://software.intel.com/content/www/us/en/develop/articles/intel-sdm.html
string CPUinfo()
{
    // Buffer to store CPU brand string (64 bytes)
    char CPUBrandString[0x40];
    // Array to store CPUID results (EAX, EBX, ECX, EDX registers)
    unsigned int CPUInfo[4] = {0, 0, 0, 0};

    // Execute CPUID with 0x80000000 to get highest extended function supported
    __cpuid(0x80000000, CPUInfo[0], CPUInfo[1], CPUInfo[2], CPUInfo[3]);
    // Store highest extended function ID
    unsigned int nExIds = CPUInfo[0];

    // Initialize brand string buffer with zeros
    memset(CPUBrandString, 0, sizeof(CPUBrandString));

    // Loop through extended CPUID functions to get brand string parts
    for (unsigned int i = 0x80000000; i <= nExIds; ++i)
    {
        // Execute CPUID for current function
        __cpuid(i, CPUInfo[0], CPUInfo[1], CPUInfo[2], CPUInfo[3]);

        // Copy results to appropriate positions in brand string
        if (i == 0x80000002)
            memcpy(CPUBrandString, CPUInfo, sizeof(CPUInfo));
        else if (i == 0x80000003)
            memcpy(CPUBrandString + 16, CPUInfo, sizeof(CPUInfo));
        else if (i == 0x80000004)
            memcpy(CPUBrandString + 32, CPUInfo, sizeof(CPUInfo));
    }
    // Convert char array to string object and return
    string str(CPUBrandString);
    return str;
}

// getOsName() - Detects the operating system at compile time
// Returns: const char* with OS name
// Sources:
// - Predefined macros: https://sourceforge.net/p/predef/wiki/OperatingSystems/
const char *getOsName()
{
// Check for Windows 32-bit
#ifdef _WIN32
    return "Windows 32-bit";
// Check for Windows 64-bit
#elif _WIN64
    return "Windows 64-bit";
// Check for Mac OS
#elif __APPLE__ || __MACH__
    return "Mac OSX";
// Check for Linux
#elif __linux__
    return "Linux";
// Check for FreeBSD
#elif __FreeBSD__
    return "FreeBSD";
// Check for Unix
#elif __unix || __unix__
    return "Unix";
// Unknown OS
#else
    return "Other";
#endif
}

// getCurrentDateTime() - Gets current date/time as formatted string
// Returns: string with current date/time in YYYY-MM-DD HH:MM:SS format
// Sources:
// - time() function: https://en.cppreference.com/w/c/chrono/time
// - strftime() function: https://en.cppreference.com/w/c/chrono/strftime
string getCurrentDateTime() {
    // Get current time as time_t (seconds since epoch)
    time_t now = time(0);
    // Buffer for formatted time string
    char buf[80];
    // Format time into buffer using local time
    strftime(buf, sizeof(buf), "%Y-%m-%d %X", localtime(&now));
    // Convert to string and return
    return string(buf);
}

// getHostname() - Gets system hostname
// Returns: string with hostname
// Sources:
// - gethostname(): https://linux.die.net/man/2/gethostname
string getHostname() {
    // Buffer for hostname (max length defined by HOST_NAME_MAX)
    char hostname[HOST_NAME_MAX];
    // Get hostname from system
    gethostname(hostname, HOST_NAME_MAX);
    // Convert to string and return
    return string(hostname);
}

// getUsername() - Gets current username
// Returns: string with username
// Sources:
// - getlogin_r(): https://linux.die.net/man/3/getlogin_r
string getUsername() {
    // Buffer for username (max length defined by LOGIN_NAME_MAX)
    char username[LOGIN_NAME_MAX];
    // Get username from system
    getlogin_r(username, LOGIN_NAME_MAX);
    // Convert to string and return
    return string(username);
}


// CPUStats struct to store CPU usage metrics
struct CPUStats 
{
    long long user, nice, system, idle, iowait, irq, softirq, steal, guest, guestNice;
};

// getCPUUsage() - Calculates CPU usage as percentage using /proc/stat
// Returns: float representing CPU usage percentage (0-100)
// Sources:
// - /proc/stat format: https://www.kernel.org/doc/html/latest/filesystems/proc.html
// - CPU usage calculation: https://stackoverflow.com/a/23376195
float getCPUUsage() {
    // Static variables to maintain state between calls
    static CPUStats prevStats;
    static float movingAverage = 0.0f;
    static const int sampleCount = 5;
    static float samples[sampleCount] = {0};
    static int currentSample = 0;
    
    CPUStats currStats;
    
    // Open /proc/stat to read CPU metrics
    FILE* file = fopen("/proc/stat", "r");
    if (!file) return movingAverage;
    
    // Read CPU metrics from first line of /proc/stat
    fscanf(file, "cpu %lld %lld %lld %lld %lld %lld %lld %lld %lld %lld",
           &currStats.user, &currStats.nice, &currStats.system, &currStats.idle,
           &currStats.iowait, &currStats.irq, &currStats.softirq, &currStats.steal,
           &currStats.guest, &currStats.guestNice);
    fclose(file);
    
    // Calculate idle and non-idle times
    const long long prevIdle = prevStats.idle + prevStats.iowait;
    const long long currIdle = currStats.idle + currStats.iowait;
    
    const long long prevNonIdle = prevStats.user + prevStats.nice + prevStats.system + 
                                 prevStats.irq + prevStats.softirq + prevStats.steal;
    const long long currNonIdle = currStats.user + currStats.nice + currStats.system + 
                                  currStats.irq + currStats.softirq + currStats.steal;
    
    // Calculate total times
    const long long prevTotal = prevIdle + prevNonIdle;
    const long long currTotal = currIdle + currNonIdle;
    
    // Calculate differences since last measurement
    const long long totald = currTotal - prevTotal;
    const long long idled = currIdle - prevIdle;
    
    // Save current stats for next call
    prevStats = currStats;
    
    // Avoid division by zero
    if (totald == 0) return movingAverage;
    
    // Calculate current CPU usage percentage
    float currentUsage = (float)(totald - idled) / totald * 100.0f;
    
    // Update moving average buffer
    samples[currentSample] = currentUsage;
    currentSample = (currentSample + 1) % sampleCount;
    
    // Calculate new moving average
    float sum = 0.0f;
    for (int i = 0; i < sampleCount; i++) {
        sum += samples[i];
    }
    movingAverage = sum / sampleCount;
    
    return movingAverage;
}

// getCPUTemperature() - Reads CPU temperature from thermal zone
// Returns: float with temperature in Celsius
// Sources:
// - Linux thermal sysfs: https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt
float getCPUTemperature() {
    // Open thermal zone temperature file
    std::ifstream tempFile("/sys/class/thermal/thermal_zone0/temp");
    float temperature = 0.0f;

    if (tempFile.is_open()) {
        // Read temperature (typically in millidegrees Celsius)
        int tempMilliC;
        tempFile >> tempMilliC;
        // Convert to degrees Celsius
        temperature = tempMilliC / 1000.0f; 
        tempFile.close();
    }

    return temperature;
}

// getFanInfo() - Attempts to read fan speed from hardware monitoring interface
// Returns: string with fan speed or "Not available"
// Sources:
// - Linux hwmon: https://www.kernel.org/doc/html/latest/hwmon/hwmon-kernel-api.html
string getFanInfo() {
    string fanInfo;
    char path[256];

    // Check common hwmon paths for fan speed
    for (int i = 0; i < 5; ++i) {
        // Try different hwmon indices
        snprintf(path, sizeof(path), "/sys/class/hwmon/hwmon%d/fan1_input", i);
        ifstream file(path);
        if (file) {
            // Read fan speed if file exists
            int speed;
            file >> speed;
            file.close();
            fanInfo = "Fan Speed: " + to_string(speed) + " RPM";
            return fanInfo;
        }
    }
    return "Fan Speed: Not available";
}