// Including a custom header file likely containing project-specific definitions
// Header files typically contain function declarations, macros, and type definitions
// Source: https://en.cppreference.com/w/cpp/preprocessor/include
#include "header.h"

// Including the SDL (Simple DirectMedia Layer) library header
// SDL is a cross-platform development library for audio, input, and graphics
// Source: https://wiki.libsdl.org/SDL2/FrontPage
#include <SDL.h>

// Including the unordered_set container from the C++ Standard Library
// unordered_set is a hash-based set implementation for fast lookups
// Source: https://en.cppreference.com/w/cpp/container/unordered_set
#include <unordered_set>

// Declaring a static vector to store CPU usage history
// 'static' makes this variable local to this translation unit
// vector is a dynamic array container from C++ STL
// Source: https://en.cppreference.com/w/cpp/container/vector
static std::vector<float> cpuHistory;

// Constant integer defining how many data points to store in history
// Used to limit the size of cpuHistory vector
static const int historySize = 200;

// Maximum Y-value for the CPU graph's vertical axis
// Helps in scaling the graph visualization
static float graphMaxY = 100.0f;  

// Interval (in seconds) at which the CPU graph should update (30 FPS)
// 1.0f / 30.0f means updating 30 times per second
static float updateInterval = 1.0f / 30.0f; 

// Accumulator tracking time since last CPU graph update
// Used with delta time in game loops
static float timeSinceLastUpdate = 0.0f;

// Boolean flag to pause CPU graph updates when true
static bool pauseGraph = false;

// Similar to cpuHistory but for temperature data
static std::vector<float> tempHistory;

// Maximum Y-value for temperature graph
static float tempGraphMaxY = 100.0f;

// Update interval for temperature graph (also targeting 30 FPS)
static float tempUpdateInterval = 1.0f / 30.0f;

// Accumulator for temperature graph update timing
static float tempTimeSinceLastUpdate = 0.0f;

// Flag to pause temperature graph updates
static bool pauseTempGraph = false;

// Similar to cpuHistory but for fan speed data
static std::vector<float> fanHistory;

// Maximum Y-value for fan speed graph (higher because RPM values are larger)
static float fanGraphMaxY = 6000.0f; 

// Update interval for fan speed graph (30 FPS)
static float fanUpdateInterval = 1.0f / 30.0f; 

// Accumulator for fan speed graph update timing
static float fanTimeSinceLastUpdate = 0.0f;

// Flag to pause fan speed graph updates
static bool pauseFanGraph = false;

// This is a preprocessor conditional block that selects which OpenGL loader to use based on compile-time definitions
// The OpenGL loader is responsible for loading OpenGL function pointers at runtime
#if defined(IMGUI_IMPL_OPENGL_LOADER_GL3W)
// If IMGUI_IMPL_OPENGL_LOADER_GL3W is defined, use gl3w loader
// gl3w is a minimal OpenGL core profile loader (https://github.com/skaslev/gl3w)
#include <GL/gl3w.h> 
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLEW)
// If IMGUI_IMPL_OPENGL_LOADER_GLEW is defined, use GLEW loader
// GLEW (OpenGL Extension Wrangler Library) is a cross-platform loader (http://glew.sourceforge.net/)
#include <GL/glew.h> 
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD)
// If IMGUI_IMPL_OPENGL_LOADER_GLAD is defined, use GLAD loader
// GLAD is a multi-language loader generator (https://github.com/Dav1dde/glad)
#include <glad/glad.h> 
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD2)
// If IMGUI_IMPL_OPENGL_LOADER_GLAD2 is defined, use GLAD2 loader
// GLAD2 is the newer version of GLAD with a different API (https://github.com/Dav1dde/glad)
#include <glad/gl.h> 
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING2)
// If IMGUI_IMPL_OPENGL_LOADER_GLBINDING2 is defined, use glbinding2
// glbinding is a C++ binding for OpenGL (https://github.com/cginternals/glbinding)
#define GLFW_INCLUDE_NONE      // Tell GLFW not to include its own OpenGL headers
#include <glbinding/Binding.h> // Main glbinding header
#include <glbinding/gl/gl.h>   // OpenGL API definitions
using namespace gl;            // Use gl namespace for OpenGL functions
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING3)
// If IMGUI_IMPL_OPENGL_LOADER_GLBINDING3 is defined, use glbinding3
// This is a newer version of glbinding with some API changes
#define GLFW_INCLUDE_NONE        // Tell GLFW not to include its own OpenGL headers
#include <glbinding/glbinding.h> // Main glbinding header
#include <glbinding/gl/gl.h>     // OpenGL API definitions
using namespace gl;              // Use gl namespace for OpenGL functions
#else
// If none of the above loaders are specified, use a custom loader
// The custom loader header path is defined by IMGUI_IMPL_OPENGL_LOADER_CUSTOM macro
#include IMGUI_IMPL_OPENGL_LOADER_CUSTOM
#endif 


// This function creates a system information window using Dear ImGui
// ImGui is an immediate-mode GUI library for C++ (https://github.com/ocornut/imgui)
void systemWindow(const char *id, ImVec2 size, ImVec2 position)
{
    // Begin a new ImGui window with the given ID
    // ImGui windows are containers for UI elements (https://github.com/ocornut/imgui/wiki)
    ImGui::Begin(id);
    
    // Set the window size using the provided ImVec2 (2D vector) parameter
    // ImVec2 contains x (width) and y (height) components (https://github.com/ocornut/imgui/wiki/Getting-Started)
    ImGui::SetWindowSize(size);
    
    // Set the window position using the provided ImVec2 parameter
    // Positions are relative to the top-left corner of the application window
    ImGui::SetWindowPos(position);

    // Add a text label "System Information" to the window
    // ImGui::Text displays non-interactive text (https://github.com/ocornut/imgui/wiki/Widgets)
    ImGui::Text("System Information");
    
    // Add a horizontal separator line for visual organization
    // Separators help group related UI elements (https://github.com/ocornut/imgui/wiki/Widgets#separator)
    ImGui::Separator();

    // Begin a tab bar container named "SystemTabs"
    // Tab bars allow organizing content into multiple pages (https://github.com/ocornut/imgui/wiki/Tabs)
    if (ImGui::BeginTabBar("SystemTabs")) {
        
        // Begin a tab item labeled "CPU" - this creates one tab page
        // Returns true if the tab is currently selected/visible
        if (ImGui::BeginTabItem("CPU")) {
            // Display hostname information
            // getHostname() returns a string containing the system's hostname
            ImGui::Text("Hostname: %s", getHostname().c_str());
            
            // Display username information
            // getUsername() returns a string containing the current username
            ImGui::Text("Username: %s", getUsername().c_str());
            
            // Display operating system name
            // getOsName() returns a string with OS information
            ImGui::Text("Operating System: %s", getOsName());
            
            // Display CPU information
            // CPUinfo() returns a string with CPU model/details
            ImGui::Text("CPU: %s", CPUinfo().c_str());
            
            // Static variables to track CPU usage display
            // These retain their values between function calls
            static float lastCpuUpdateTime = 0.0f;  // Last time CPU usage was updated
            static float displayedCpuUsage = 0.0f;  // Current displayed CPU usage value
            
            // Get current time from ImGui's internal clock
            // ImGui::GetTime() returns seconds since application start (https://github.com/ocornut/imgui/wiki/Internal-Functions)
            float currentTime = ImGui::GetTime();
            
            // Update CPU usage display every 0.5 seconds (throttle updates)
            if (currentTime - lastCpuUpdateTime > 0.5f) {
                // getCPUUsage() returns current CPU usage percentage
                displayedCpuUsage = getCPUUsage();
                lastCpuUpdateTime = currentTime;
            }
            
            // Display the CPU usage percentage with 2 decimal places
            ImGui::Text("CPU Usage: %.2f%%", displayedCpuUsage);
            
            // Add another separator for visual organization
            ImGui::Separator();
            
            // Label for the CPU performance graph section
            ImGui::Text("CPU Performance Graph");

            // Slider to control the maximum Y-axis value of the graph
            // graphMaxY is a float variable storing the current max value
            ImGui::SliderFloat("Max Y (%)", &graphMaxY, 10.0f, 100.0f);
            
            // Slider to control graph update frequency (frames per second)
            // updateInterval is a float storing time between updates
            ImGui::SliderFloat("FPS", &updateInterval, 1.0f / 144.0f, 1.0f); 
            
            // Button to toggle graph pausing
            // pauseGraph is a boolean tracking pause state
            // Button label changes based on current state
            if (ImGui::Button(pauseGraph ? "Resume Graph" : "Pause Graph"))
                pauseGraph = !pauseGraph;

            // Track time since last graph update
            static float lastTime = ImGui::GetTime();
            float deltaTime = currentTime - lastTime;
            lastTime = currentTime;

            // Accumulate time since last graph update
            timeSinceLastUpdate += deltaTime;

            // Update graph data if not paused and enough time has passed
            if (!pauseGraph && timeSinceLastUpdate >= updateInterval) {
                // Remove oldest data point if history is full
                if (cpuHistory.size() >= historySize)
                    cpuHistory.erase(cpuHistory.begin());
                
                // Add current CPU usage to history
                cpuHistory.push_back(displayedCpuUsage);  
                
                // Reset update timer
                timeSinceLastUpdate = 0.0f;
            }

            // Display the CPU usage graph if we have data
            if (!cpuHistory.empty()) {
                // ImGui::PlotLines creates a line graph from data
                // Parameters: label, data array, array size, text overlay,
                //             scale min, scale max, display size
                ImGui::PlotLines("CPU Usage (%)", cpuHistory.data(), cpuHistory.size(), 0,
                                 nullptr, 0.0f, graphMaxY, ImVec2(0, 100));
            }

            // End the CPU tab item
            ImGui::EndTabItem();
        }

        
// This code is part of a graphical user interface (GUI) using Dear ImGui (ImGui), a lightweight immediate-mode GUI library for C++
// Source: https://github.com/ocornut/imgui

// Begins a tab item labeled "Fan" in a tab bar. Returns true if the tab is selected/active.
// Source: https://github.com/ocornut/imgui/wiki/FAQ#how-can-i-have-multiple-widgets-with-the-same-label
if (ImGui::BeginTabItem("Fan")) {
    
    // Gets fan information as a string from some external function (implementation not shown here)
    // The string likely contains information like fan speed in RPM (Revolutions Per Minute)
    string fanStr = getFanInfo();
    
    // Displays the fan information text in the GUI
    // %s is a format specifier for string, c_str() converts C++ string to C-style string
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L2308
    ImGui::Text("%s", fanStr.c_str());

    // Adds a horizontal separator line in the GUI to visually separate sections
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L2342
    ImGui::Separator();
    
    // Displays a label for the fan speed graph section
    ImGui::Text("Fan Speed Graph");
    
    // Creates a slider to adjust the maximum Y-axis value (RPM) for the graph
    // The slider ranges from 1000 to 8000 RPM, and modifies the fanGraphMaxY variable
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1955
    ImGui::SliderFloat("Max Y (RPM)", &fanGraphMaxY, 1000.0f, 8000.0f);
    
    // Creates a slider to adjust the update frequency (frames per second) of the fan graph
    // The slider ranges from ~0.007 (1/144) to 1 FPS, and modifies fanUpdateInterval
    // "%.3f" formats the number to show 3 decimal places
    // "##Fan" adds a unique identifier to distinguish this slider from others with same label
    ImGui::SliderFloat("FPS##Fan", &fanUpdateInterval, 1.0f / 144.0f, 1.0f, "%.3f");
    
    // Creates a button that toggles between "Pause" and "Resume" states for the graph
    // The button text changes based on the pauseFanGraph boolean value
    // Clicking the button toggles the pauseFanGraph variable
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1838
    if (ImGui::Button(pauseFanGraph ? "Resume Fan Graph" : "Pause Fan Graph"))
        pauseFanGraph = !pauseFanGraph;

    // Gets the current time from ImGui's internal timer (seconds since program start)
    // Used to calculate time between updates
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L2624
    static float lastFanTime = ImGui::GetTime();
    float currentFanTime = ImGui::GetTime();
    
    // Calculates time elapsed since last update
    float deltaFanTime = currentFanTime - lastFanTime;
    
    // Updates the last time variable for next frame
    lastFanTime = currentFanTime;
    
    // Accumulates time since last graph update
    fanTimeSinceLastUpdate += deltaFanTime;

    // Extracts the RPM value from the fan information string using sscanf
    // Parses the string looking for "Fan Speed: %d RPM" format and stores the number in fanRPM
    // Source: https://en.cppreference.com/w/cpp/io/c/fscanf
    int fanRPM = 0;
    sscanf(fanStr.c_str(), "Fan Speed: %d RPM", &fanRPM);

    // Updates the fan history data if not paused and enough time has passed
    if (!pauseFanGraph && fanTimeSinceLastUpdate >= fanUpdateInterval) {
        // If history buffer is full, remove oldest entry
        if (fanHistory.size() >= historySize)
            fanHistory.erase(fanHistory.begin());
        
        // Add current RPM reading to history (converted to float)
        fanHistory.push_back(static_cast<float>(fanRPM));
        
        // Reset update timer
        fanTimeSinceLastUpdate = 0.0f;
    }

    // Draws the fan RPM graph if there is data available
    if (!fanHistory.empty()) {
        // Creates a line plot showing RPM history
        // Parameters: label, data array, data count, array offset (0 for all),
        // overlay text (nullptr for none), scale minimum (0), scale maximum (fanGraphMaxY),
        // graph dimensions (0 = auto width, 100 pixels height)
        // Source: https://github.com/ocornut/imgui/wiki/Plotting
        ImGui::PlotLines("Fan RPM", fanHistory.data(), fanHistory.size(), 0,
                         nullptr, 0.0f, fanGraphMaxY, ImVec2(0, 100));
    }

    // Ends the "Fan" tab item (must match BeginTabItem)
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L2776
    ImGui::EndTabItem();
}

// This code snippet is part of a graphical user interface (GUI) application using Dear ImGui
// (https://github.com/ocornut/imgui), a popular immediate mode GUI library for C++.

// Begins a tab item labeled "Thermal" in a tab bar
// ImGui tabs allow organizing content into separate panels
// Source: https://github.com/ocornut/imgui/wiki/FAQ#how-can-i-have-multiple-widgets-with-the-same-label
if (ImGui::BeginTabItem("Thermal")) {
    // Gets the current CPU temperature from a system-specific function
    // The implementation of getCPUTemperature() would be platform-dependent
    float temperature = getCPUTemperature();
    
    // Displays the current CPU temperature as text with 2 decimal places
    // ImGui::Text formats and displays text similar to printf
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L223
    ImGui::Text("CPU Temperature: %.2f°C", temperature);

    // Adds a horizontal separator line for visual organization
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L269
    ImGui::Separator();
    
    // Displays a label for the temperature graph section
    ImGui::Text("Temperature Graph");
    
    // Creates a slider to control the maximum Y-axis value for the temperature graph
    // The &tempGraphMaxY parameter is the variable being modified
    // Range is from 30.0°C to 120.0°C
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L422
    ImGui::SliderFloat("Max Y (°C)", &tempGraphMaxY, 30.0f, 120.0f);
    
    // Creates a slider to control the update frequency of the temperature graph
    // The "##Thermal" suffix ensures unique ID for this slider
    // Range is from 1/144th of a second to 1 second
    ImGui::SliderFloat("FPS##Thermal", &tempUpdateInterval, 1.0f / 144.0f, 1.0f);
    
    // Creates a button that toggles between "Pause" and "Resume" states
    // The label changes based on the pauseTempGraph boolean state
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L372
    if (ImGui::Button(pauseTempGraph ? "Resume Temp Graph" : "Pause Temp Graph"))
        pauseTempGraph = !pauseTempGraph;

    // Stores the last time temperature was updated (static maintains value between calls)
    static float lastTempTime = ImGui::GetTime();
    
    // Gets the current time from ImGui's internal clock
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1697
    float currentTempTime = ImGui::GetTime();
    
    // Calculates time elapsed since last temperature update
    float deltaTempTime = currentTempTime - lastTempTime;
    
    // Updates the last update time to current time
    lastTempTime = currentTempTime;

    // Accumulates time since last graph update
    tempTimeSinceLastUpdate += deltaTempTime;

    // Updates the temperature history if not paused and enough time has passed
    if (!pauseTempGraph && tempTimeSinceLastUpdate >= tempUpdateInterval) {
        // Maintains history size by removing oldest entry if at capacity
        if (tempHistory.size() >= historySize)
            tempHistory.erase(tempHistory.begin());
        
        // Adds current temperature to history
        tempHistory.push_back(temperature);
        
        // Resets the update timer
        tempTimeSinceLastUpdate = 0.0f;
    }

    // Only draw the graph if there's data available
    if (!tempHistory.empty()) {
        // Creates a line plot of temperature history
        // Parameters: label, data array, data count, value offset, overlay text
        //             scale min, scale max, and graph dimensions
        // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1777
        ImGui::PlotLines("Temperature (°C)", tempHistory.data(), tempHistory.size(), 0,
                         nullptr, 0.0f, tempGraphMaxY, ImVec2(0, 100));
    }

    // Ends the "Thermal" tab item
    // Must be called for each BeginTabItem()
    // Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1474
    ImGui::EndTabItem();
}

// Ends the tab bar that contains all tab items
// Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1471
ImGui::EndTabBar();

// Adds a horizontal separator between sections
ImGui::Separator();

// Opens the Linux system file that contains uptime information
// /proc/uptime contains system uptime and idle time in seconds
// Source: https://man7.org/linux/man-pages/man5/proc.5.html
FILE* uptimeFile = fopen("/proc/uptime", "r");

// If the file opened successfully...
if (uptimeFile) {
    // Variables to store uptime and idle time values
    float uptime, idle;
    
    // Reads the two float values from the file
    fscanf(uptimeFile, "%f %f", &uptime, &idle);
    
    // Closes the file as we're done with it
    fclose(uptimeFile);

    // Converts total seconds into days, hours, minutes, seconds
    int days = uptime / 86400;           // 86400 seconds in a day
    int hours = ((int)uptime % 86400) / 3600;  // 3600 seconds in an hour
    int minutes = ((int)uptime % 3600) / 60;    // 60 seconds in a minute
    int seconds = (int)uptime % 60;             // Remainder is seconds

    // Displays formatted uptime information
    ImGui::Text("System Uptime: %d days, %d hours, %d minutes, %d seconds",
                days, hours, minutes, seconds);
}

// Displays current date/time from a custom getCurrentDateTime() function
// The implementation would return a string with formatted date/time
ImGui::Text("Current Date/Time: %s", getCurrentDateTime().c_str());

// Ends the current ImGui window
// Must be called for each Begin() call
// Source: https://github.com/ocornut/imgui/blob/master/imgui.h#L1318
ImGui::End();
    }
}

// This function creates a window displaying memory and process information using Dear ImGui
// Reference: https://github.com/ocornut/imgui
void memoryProcessesWindow(const char *id, ImVec2 size, ImVec2 position)
{
    // Begin a new ImGui window with the given ID
    // Reference: https://github.com/ocornut/imgui/wiki/Getting-Started
    ImGui::Begin(id);
    
    // Set the window size using the provided dimensions
    // Reference: https://github.com/ocornut/imgui/wiki/Window-Sizing
    ImGui::SetWindowSize(id, size);
    
    // Set the window position using the provided coordinates
    // Reference: https://github.com/ocornut/imgui/wiki/Window-Sizing#window-position
    ImGui::SetWindowPos(id, position);

    // Create a tab bar to organize memory and process information
    // Reference: https://github.com/ocornut/imgui/wiki/Tabs
    if (ImGui::BeginTabBar("MemoryProcessTabs"))
    {
        // First tab: Memory information
        // Reference: https://github.com/ocornut/imgui/wiki/Tabs
        if (ImGui::BeginTabItem("Memory"))
        {
            // Get system memory information using sysinfo structure
            // Reference: http://man7.org/linux/man-pages/man2/sysinfo.2.html
            struct sysinfo memInfo;
            sysinfo(&memInfo);

            // Calculate total physical memory (RAM) in bytes
            // mem_unit is the memory unit size in bytes
            // Reference: http://man7.org/linux/man-pages/man2/sysinfo.2.html
            long long totalPhysMem = memInfo.totalram * memInfo.mem_unit;
            
            // Calculate used physical memory (total - free)
            long long physMemUsed = (memInfo.totalram - memInfo.freeram) * memInfo.mem_unit;
            
            // Calculate memory usage percentage
            float memoryUsagePercent = (float)physMemUsed / totalPhysMem * 100.0f;

            // Display memory usage information in MB
            ImGui::Text("Memory Usage: %.1f MB / %.1f MB (%.1f%%)", 
                (float)physMemUsed / (1024 * 1024), 
                (float)totalPhysMem / (1024 * 1024),
                memoryUsagePercent);

            // Create a progress bar showing memory usage
            // Push/PopStyleColor temporarily changes the color
            // Reference: https://github.com/ocornut/imgui/wiki/Progress-Bars
            ImGui::PushStyleColor(ImGuiCol_PlotHistogram, ImVec4(0.2f, 0.6f, 1.0f, 1.0f));
            ImGui::ProgressBar(memoryUsagePercent / 100.0f, ImVec2(-1, 20), "");
            ImGui::PopStyleColor();

            // Variables for memory usage history graph
            static std::vector<float> memoryHistory;          // Stores historical memory usage values
            static const int memoryHistorySize = 200;         // Maximum number of points to store
            static float memoryGraphMaxY = 100.0f;           // Maximum Y-axis value for the graph
            static float memoryUpdateInterval = 1.0f / 30.0f; // Update interval in seconds (30 FPS)
            static float memoryTimeSinceLastUpdate = 0.0f;    // Time since last update
            static bool pauseMemoryGraph = false;             // Pause graph updates when true

            // Calculate time delta for smooth graph updates
            static float lastMemoryTime = ImGui::GetTime();
            float currentMemoryTime = ImGui::GetTime();
            float deltaMemoryTime = currentMemoryTime - lastMemoryTime;
            lastMemoryTime = currentMemoryTime;
            memoryTimeSinceLastUpdate += deltaMemoryTime;

            // Update memory history if not paused and enough time has passed
            if (!pauseMemoryGraph && memoryTimeSinceLastUpdate >= memoryUpdateInterval) {
                // Maintain fixed-size history by removing oldest point if needed
                if (memoryHistory.size() >= memoryHistorySize)
                    memoryHistory.erase(memoryHistory.begin());
                // Add current memory usage to history
                memoryHistory.push_back(memoryUsagePercent);
                memoryTimeSinceLastUpdate = 0.0f;
            }

            // UI controls for the memory graph
            ImGui::Separator();
            ImGui::Text("Memory Usage History");
            ImGui::SliderFloat("Max Y (%)##Memory", &memoryGraphMaxY, 10.0f, 100.0f);
            ImGui::SliderFloat("FPS##Memory", &memoryUpdateInterval, 1.0f / 144.0f, 1.0f);
            if (ImGui::Button(pauseMemoryGraph ? "Resume Memory Graph" : "Pause Memory Graph"))
                pauseMemoryGraph = !pauseMemoryGraph;

            // Display the memory usage history graph if we have data
            if (!memoryHistory.empty()) {
                // PlotLines creates a line graph of the data
                // Reference: https://github.com/ocornut/imgui/wiki/Plotting
                ImGui::PlotLines("##MemoryUsage", memoryHistory.data(), memoryHistory.size(), 0,
                                 nullptr, 0.0f, memoryGraphMaxY, ImVec2(0, 100));
            }

            // Display swap memory information if available
            if (memInfo.totalswap > 0) {
                ImGui::Separator();
                // Calculate swap memory usage
                long long swapTotal = memInfo.totalswap * memInfo.mem_unit;
                long long swapFree = memInfo.freeswap * memInfo.mem_unit;
                long long swapUsed = swapTotal - swapFree;
                float swapUsagePercent = (float)swapUsed / swapTotal * 100.0f;

                // Display swap usage information
                ImGui::Text("Swap Usage: %.1f MB / %.1f MB (%.1f%%)", 
                    (float)swapUsed / (1024 * 1024), 
                    (float)swapTotal / (1024 * 1024),
                    swapUsagePercent);

                // Create a progress bar for swap usage with different color
                ImGui::PushStyleColor(ImGuiCol_PlotHistogram, ImVec4(0.8f, 0.4f, 0.6f, 1.0f));
                ImGui::ProgressBar(swapUsagePercent / 100.0f, ImVec2(-1, 20), "");
                ImGui::PopStyleColor();
            }

            // Display additional memory information from 'free -h' command
            ImGui::Separator();
            ImGui::Text("Memory Summary (free -h):");
            string freeInfo = getFreeMemoryInfo();
            ImGui::TextUnformatted(freeInfo.c_str());

            // End the Memory tab
            ImGui::EndTabItem();
        }

        // Second tab: Process information
        if (ImGui::BeginTabItem("Processes"))
        {
            // Get total number of running processes
            int totalProcesses = getTotalProcessCount();
            ImGui::Text("Total Running Processes: %d", totalProcesses);

            // Process filtering UI
            static char processNameFilter[256] = "";  // Filter by process name
            static char processPidFilter[64] = "";    // Filter by process ID
            ImGui::Separator();
            ImGui::Text("Filters:");
            // Input text fields for filtering
            // Reference: https://github.com/ocornut/imgui/wiki/Widgets#input-text
            ImGui::InputText("Name##Filter", processNameFilter, IM_ARRAYSIZE(processNameFilter));
            ImGui::InputText("PID##Filter", processPidFilter, IM_ARRAYSIZE(processPidFilter));

            // Get all processes and apply filters
            std::vector<Proc> allProcesses = getProcesses();
            std::vector<Proc> filteredProcesses;

            // Filter processes based on name and PID
            for (const auto& proc : allProcesses) {
                // Check if process name matches filter (if any)
                bool nameMatch = strlen(processNameFilter) == 0 || 
                                 strstr(proc.name.c_str(), processNameFilter);
                
                // Check if process PID matches filter (if any)
                bool pidMatch = true;
                if (strlen(processPidFilter) > 0) {
                    try {
                        // Convert PID filter to integer and compare
                        int pidFilter = std::stoi(processPidFilter);
                        pidMatch = proc.pid == pidFilter;
                    } catch (...) {
                        // Handle invalid PID filter (non-numeric input)
                        pidMatch = false;
                    }
                }

                // Add to filtered list if both filters match
                if (nameMatch && pidMatch) {
                    filteredProcesses.push_back(proc);
                }
            }

            // Show filtered process count
            ImGui::Text("Filtered Processes: %d/%d", (int)filteredProcesses.size(), totalProcesses);

            // Process selection functionality
            static std::unordered_set<int> selectedPids;  // Stores selected process PIDs
            
            // Process list display
            ImGui::Separator();
            ImGui::Text("Running Processes:");
            // Create scrollable child window for process list
            // Reference: https://github.com/ocornut/imgui/wiki/Child-Windows
            ImGui::BeginChild("ProcessList", ImVec2(0, 0), true);

            // Create a table to display process information
            // Reference: https://github.com/ocornut/imgui/wiki/Tables
            if (ImGui::BeginTable("Processes", 5, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_Resizable)) {
                // Set up table columns
                ImGui::TableSetupColumn("PID");
                ImGui::TableSetupColumn("Name");
                ImGui::TableSetupColumn("State");
                ImGui::TableSetupColumn("CPU");
                ImGui::TableSetupColumn("Memory (MB)");
                ImGui::TableHeadersRow();

                // Populate table with process information
                for (const auto& p : filteredProcesses) {
                    ImGui::TableNextRow();
                    ImGui::TableSetColumnIndex(0);

                    // Process selection functionality
                    bool isSelected = selectedPids.count(p.pid) > 0;
                    std::string label = std::to_string(p.pid) + "##" + p.name;

                    // Create selectable first column that spans all columns
                    if (ImGui::Selectable(label.c_str(), isSelected, ImGuiSelectableFlags_SpanAllColumns)) {
                        // Toggle selection status
                        if (isSelected)
                            selectedPids.erase(p.pid);
                        else
                            selectedPids.insert(p.pid);
                    }

                    // Fill remaining columns with process information
                    ImGui::TableSetColumnIndex(1);
                    ImGui::Text("%s", p.name.c_str());

                    ImGui::TableSetColumnIndex(2);
                    ImGui::Text("%c", p.state);

                    ImGui::TableSetColumnIndex(3);
                    float cpuUsage = p.cpuUsage;
                    ImGui::Text("%.1f%%", cpuUsage);

                    ImGui::TableSetColumnIndex(4);
                    ImGui::Text("%.1f MB", p.vsize / (1024.0 * 1024.0));
                }

                ImGui::EndTable();
            }

            ImGui::EndChild();

            // Display selected PIDs if any processes are selected
            if (!selectedPids.empty()) {
                ImGui::Separator();
                ImGui::Text("Selected PIDs:");
                for (int pid : selectedPids) {
                    ImGui::SameLine();
                    ImGui::Text("%d", pid);
                }
            }

            // End the Processes tab
            ImGui::EndTabItem();
        }

   /*
This is a C++ program using Dear ImGui (Immediate Mode GUI) to create a system monitoring dashboard.
It displays disk, memory, process, and network information in separate windows using SDL2 for window
creation and OpenGL3 for rendering.

Main components:
1. SDL2 - Handles window creation and input events
2. OpenGL - Renders the graphics
3. Dear ImGui - Creates the GUI interface
4. Custom functions - Get system information and display it

Sources:
- ImGui: https://github.com/ocornut/imgui
- SDL2: https://www.libsdl.org/
- OpenGL: https://www.opengl.org/
*/

// Display disk information in a tabbed interface
if (ImGui::BeginTabItem("Disk")) {
    // Get disk information and store it in a static vector to persist between frames
    static vector<DiskInfo> disks = getDiskInfo();
    
    // Loop through all disks to find the root filesystem (mounted on "/")
    for (const auto& disk : disks) {
        if (disk.mountedOn == "/") {
            // Display root filesystem information
            ImGui::Text("Root Filesystem: %s", disk.filesystem.c_str());
            // Display disk space information (total, used, available)
            ImGui::Text("Total: %s, Used: %s (%s), Available: %s", 
                disk.size.c_str(), 
                disk.used.c_str(), 
                disk.usePercent.c_str(),
                disk.available.c_str());
            
            // Convert usage percentage to integer for the progress bar
            int percent = atoi(disk.usePercent.c_str());
            // Set progress bar color based on usage level
            ImGui::PushStyleColor(ImGuiCol_PlotHistogram, 
                percent > 90 ? ImVec4(1.0f, 0.0f, 0.0f, 1.0f) : // Red if >90%
                percent > 70 ? ImVec4(1.0f, 1.0f, 0.0f, 1.0f) : // Yellow if >70%
                ImVec4(0.2f, 0.8f, 0.2f, 1.0f)); // Green otherwise
            // Display progress bar showing disk usage
            ImGui::ProgressBar(percent / 100.0f, ImVec2(-1, 20));
            ImGui::PopStyleColor();
            break;
        }
    }
    
    // Add a separator line
    ImGui::Separator();
    // Display heading for all filesystems
    ImGui::Text("All Filesystems:");
    // Draw a table with all disk information
    drawDiskInfoTable(disks);
    
    // Button to refresh disk information
    if (ImGui::Button("Refresh Disk Info")) {
        disks = getDiskInfo();
    }
    
    // End the disk tab
    ImGui::EndTabItem();
}

// End the tab bar that contains all tabs
ImGui::EndTabBar();

// End the main ImGui window
ImGui::End();
    }
}

/*
Network window function - displays network statistics
Parameters:
- id: Window identifier
- size: Window dimensions
- position: Window position on screen
*/
void networkWindow(const char *id, ImVec2 size, ImVec2 position)
{
    // Begin a new ImGui window
    ImGui::Begin(id);
    // Set window size
    ImGui::SetWindowSize(id, size);
    // Set window position
    ImGui::SetWindowPos(id, position);
    // Update network statistics
    updateNetworkStats();                    
    // Draw the network information
    drawNetworkWindow(id, size, position);   
    
    // End the network window
    ImGui::End();
}

/*
Main function - entry point of the program
Initializes all systems and runs the main loop
*/
int main(int, char **)
{
    // Initialize SDL2 with video, timer and game controller subsystems
    // SDL provides cross-platform access to system hardware
    if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_TIMER | SDL_INIT_GAMECONTROLLER) != 0)
    {
        printf("Error: %s\n", SDL_GetError());
        return -1;
    }

    // Set OpenGL version to 3.0
    const char *glsl_version = "#version 130";
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_FLAGS, 0);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_PROFILE_MASK, SDL_GL_CONTEXT_PROFILE_CORE);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MAJOR_VERSION, 3);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MINOR_VERSION, 0);

    // Configure OpenGL buffer settings
    SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER, 1); // Enable double buffering
    SDL_GL_SetAttribute(SDL_GL_DEPTH_SIZE, 24);   // 24-bit depth buffer
    SDL_GL_SetAttribute(SDL_GL_STENCIL_SIZE, 8);  // 8-bit stencil buffer
    
    // Create window with OpenGL context
    SDL_WindowFlags window_flags = (SDL_WindowFlags)(SDL_WINDOW_OPENGL | SDL_WINDOW_RESIZABLE | SDL_WINDOW_ALLOW_HIGHDPI);
    SDL_Window *window = SDL_CreateWindow("Dear ImGui SDL2+OpenGL3 example", 
        SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, 1280, 720, window_flags);
    
    // Create OpenGL context
    SDL_GLContext gl_context = SDL_GL_CreateContext(window);
    SDL_GL_MakeCurrent(window, gl_context);
    SDL_GL_SetSwapInterval(1); // Enable vsync

    // Initialize OpenGL loader (different options for different loaders)
#if defined(IMGUI_IMPL_OPENGL_LOADER_GL3W)
    bool err = gl3wInit() != 0;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLEW)
    bool err = glewInit() != GLEW_OK;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD)
    bool err = gladLoadGL() == 0;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD2)
    bool err = gladLoadGL((GLADloadfunc)SDL_GL_GetProcAddress) == 0; 
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING2)
    bool err = false;
    glbinding::Binding::initialize();
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING3)
    bool err = false;
    glbinding::initialize([](const char *name) { return (glbinding::ProcAddress)SDL_GL_GetProcAddress(name); });
#else
    bool err = false; 
#endif
    if (err)
    {
        fprintf(stderr, "Failed to initialize OpenGL loader!\n");
        return 1;
    }

    // Setup Dear ImGui context
    IMGUI_CHECKVERSION();
    ImGui::CreateContext();
    
    // Get reference to ImGui IO (input/output) configuration
    ImGuiIO &io = ImGui::GetIO();

    // Setup ImGui style (dark theme)
    ImGui::StyleColorsDark();

    // Initialize ImGui backends for SDL2 and OpenGL3
    ImGui_ImplSDL2_InitForOpenGL(window, gl_context);
    ImGui_ImplOpenGL3_Init(glsl_version);

    // Set clear color (transparent black)
    ImVec4 clear_color = ImVec4(0.0f, 0.0f, 0.0f, 0.0f);

    // Main application loop
    bool done = false;
    while (!done)
    {
        // Process SDL events (input, window events, etc.)
        SDL_Event event;
        while (SDL_PollEvent(&event))
        {
            // Let ImGui process the event
            ImGui_ImplSDL2_ProcessEvent(&event);
            // Quit if window is closed
            if (event.type == SDL_QUIT)
                done = true;
            if (event.type == SDL_WINDOWEVENT && event.window.event == SDL_WINDOWEVENT_CLOSE && event.window.windowID == SDL_GetWindowID(window))
                done = true;
        }

        // Start new ImGui frame
        ImGui_ImplOpenGL3_NewFrame();
        ImGui_ImplSDL2_NewFrame(window);
        ImGui::NewFrame();

        // Create application windows
        {
            // Get main display size
            ImVec2 mainDisplay = io.DisplaySize;
            
            // Create memory and processes window (top-left)
            memoryProcessesWindow("== Memory and Processes ==",
                                  ImVec2((mainDisplay.x / 2) - 20, (mainDisplay.y / 2) + 30),
                                  ImVec2((mainDisplay.x / 2) + 10, 10));
            
            // Create system info window (top-right)
            systemWindow("== System ==",
                         ImVec2((mainDisplay.x / 2) - 10, (mainDisplay.y / 2) + 30),
                         ImVec2(10, 10));
            
            // Create network window (bottom, full width)
            networkWindow("== Network ==",
                          ImVec2(mainDisplay.x - 20, (mainDisplay.y / 2) - 60),
                          ImVec2(10, (mainDisplay.y / 2) + 50));
        }

        // Render ImGui to screen
        ImGui::Render();
        // Set OpenGL viewport to cover entire window
        glViewport(0, 0, (int)io.DisplaySize.x, (int)io.DisplaySize.y);
        // Clear screen with specified color
        glClearColor(clear_color.x, clear_color.y, clear_color.z, clear_color.w);
        glClear(GL_COLOR_BUFFER_BIT);
        // Render ImGui draw data using OpenGL
        ImGui_ImplOpenGL3_RenderDrawData(ImGui::GetDrawData());
        // Swap front and back buffers (double buffering)
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

