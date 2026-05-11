// Including necessary headers:
// header.h - Likely contains custom project headers and forward declarations
// algorithm - Provides STL algorithms like sorting, searching (https://en.cppreference.com/w/cpp/header/algorithm)
// chrono - For time-related functions (https://en.cppreference.com/w/cpp/header/chrono)
// fstream - For file input/output operations (https://en.cppreference.com/w/cpp/header/fstream)
// sstream - For string stream operations (https://en.cppreference.com/w/cpp/header/sstream)
// dirent.h - For directory operations (POSIX) (https://man7.org/linux/man-pages/man3/readdir.3.html)
// sys/stat.h - For file status information (https://man7.org/linux/man-pages/man2/stat.2.html)
// unistd.h - POSIX API for system calls (https://man7.org/linux/man-pages/man0/unistd.h.0p.html)
// string.h - C string operations (https://en.cppreference.com/w/c/string/byte)
// array - STL array container (https://en.cppreference.com/w/cpp/container/array)
// memory - For smart pointers (https://en.cppreference.com/w/cpp/header/memory)
// stdexcept - Standard exception classes (https://en.cppreference.com/w/cpp/header/stdexcept)
#include "header.h"
#include <algorithm>
#include <chrono>  
#include <fstream>
#include <sstream>
#include <dirent.h>
#include <sys/stat.h>
#include <unistd.h>
#include <string.h>
#include <array>
#include <memory>
#include <stdexcept>

// Function to get number of CPU cores
// Uses sysconf() with _SC_NPROCESSORS_ONLN parameter to get online processors
// Source: https://man7.org/linux/man-pages/man3/sysconf.3.html
int getNumCores() {
    return sysconf(_SC_NPROCESSORS_ONLN);
}

// Function to get memory statistics as a formatted string
// Uses sysinfo() system call to get memory information
// Source: https://man7.org/linux/man-pages/man2/sysinfo.2.html
string getMemoryStats() {
    struct sysinfo memInfo; 
    sysinfo(&memInfo);  // Populate memInfo struct with system memory data
    
    // Calculate total physical memory (RAM)
    // mem_unit contains size of memory unit in bytes
    long long totalPhysMem = memInfo.totalram;
    totalPhysMem *= memInfo.mem_unit;
    
    // Calculate used physical memory
    long long physMemUsed = memInfo.totalram - memInfo.freeram;
    physMemUsed *= memInfo.mem_unit;
    
    // Format the output string with memory usage statistics
    std::stringstream ss;
    ss << "Memory Usage: " << (physMemUsed / (1024 * 1024)) << " MB / " 
       << (totalPhysMem / (1024 * 1024)) << " MB (" 
       << (physMemUsed * 100 / totalPhysMem) << "%)";
    
    return ss.str();
}

// Function to get free memory information using 'free -h' command
// Uses popen() to execute shell command and read its output
// Source: https://man7.org/linux/man-pages/man3/popen.3.html
string getFreeMemoryInfo() {
    string result;
    FILE* pipe = popen("free -h", "r");  // Open pipe to free command
    if (!pipe) return "Error executing free command";
    
    // Read command output line by line
    char buffer[128];
    while (fgets(buffer, sizeof(buffer), pipe) != nullptr) {
        result += buffer;
    }
    
    pclose(pipe);  // Close the pipe
    return result;
}

// Function to execute a shell command and return its output
// Uses popen() with RAII wrapper (unique_ptr) for automatic cleanup
// Source: https://man7.org/linux/man-pages/man3/popen.3.html
string exec_command(const char* cmd) {
    array<char, 128> buffer;  // Fixed-size buffer for reading output
    string result;
    // Define pclose_func type for custom deleter in unique_ptr
    using pclose_func = int(*)(FILE*);
    // Create unique_ptr with custom deleter to ensure pipe is closed
    unique_ptr<FILE, pclose_func> pipe(popen(cmd, "r"), pclose);
    if (!pipe) {
        throw runtime_error("popen() failed!");
    }
    // Read output line by line
    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }
    return result;
}

// Function to get list of running processes
// Tries multiple methods: top, ps, and direct /proc reading
// Returns vector of Proc structs containing process information
vector<Proc> getProcesses() {
    vector<Proc> processes;
    unordered_set<int> seen_pids;  // Track seen PIDs to avoid duplicates

    // First attempt: use top command for process information
    try {
        // Execute top command in batch mode (-b), 1 iteration (-n 1), sorted by CPU (-o %CPU)
        string top_output = exec_command("top -b -n 1 -o %CPU 2>/dev/null");
        istringstream iss(top_output);
        string line;
        
        // Skip header lines (first 6 lines in top output)
        for (int i = 0; i < 6; i++) {
            if (!getline(iss, line)) break;
        }
        
        // Process each line of top output
        while (getline(iss, line)) {
            if (line.empty()) continue;
            
            istringstream line_stream(line);
            vector<string> columns;
            string column;
            
            // Split line into columns
            while (line_stream >> column) {
                columns.push_back(column);
            }
            
            // Need at least 9 columns (PID through COMMAND)
            if (columns.size() < 9) continue;
            
            try {
                Proc p;
                p.pid = stoi(columns[0]);  // First column is PID
                p.cpuUsage = stof(columns[8]); // CPU% is typically 9th column
                
                // Validate PID and CPU values
                if (p.pid <= 0 || p.cpuUsage < 0) continue;
                
                // Process name might be last column
                p.name = columns[11];
                for (size_t i = 12; i < columns.size(); i++) {
                    p.name += " " + columns[i];
                }
                
                // Skip if we've already seen this PID
                if (!seen_pids.insert(p.pid).second) continue;
                
                // Get additional info from /proc filesystem
                char path[256];
                snprintf(path, sizeof(path), "/proc/%d/stat", p.pid);
                ifstream stat_file(path);
                if (stat_file) {
                    string stat_line;
                    getline(stat_file, stat_line);
                    istringstream stat_iss(stat_line);
                    
                    // Read through stat file to get needed fields
                    string dummy;
                    for (int i = 0; i < 2; i++) stat_iss >> dummy; // Skip pid and comm
                    stat_iss >> p.state;  // Process state
                    for (int i = 0; i < 19; i++) stat_iss >> dummy; // Skip to vsize
                    stat_iss >> p.vsize >> p.rss;  // Virtual and resident memory size
                    
                    // Convert pages to bytes (typically 4KB pages)
                    p.rss *= sysconf(_SC_PAGESIZE);
                }
                
                processes.push_back(p);
            } catch (...) {
                continue;  // Skip problematic entries
            }
        }
        
        if (!processes.empty()) return processes;
    } catch (...) {
        // Fall through to ps attempt if top fails
    }

    // Second attempt: use ps command if top failed
    try {
        // Execute ps command with specific columns: pid, cpu%, command, state, vsize, rss
        string ps_output = exec_command("ps -eo pid,pcpu,comm,state,vsize,rss --no-headers --sort -pcpu 2>/dev/null");
        istringstream iss(ps_output);
        string line;
        
        while (getline(iss, line)) {
            if (line.empty()) continue;
            
            Proc p;
            istringstream line_stream(line);
            string token;
            
            // Parse columns: PID (1), CPU% (2), COMMAND (3), STATE (4), VSZ (5), RSS (6)
            if (!(line_stream >> p.pid >> p.cpuUsage >> p.name >> token)) continue;
            
            // Validate PID and CPU
            if (p.pid <= 0 || p.cpuUsage < 0) continue;
            
            p.state = token.empty() ? '?' : token[0];
            
            if (line_stream >> p.vsize >> p.rss) {
                if (seen_pids.insert(p.pid).second) { // Check for duplicates
                    processes.push_back(p);
                }
            }
        }
        
        if (!processes.empty()) return processes;
    } catch (...) {
        // Fall through to /proc attempt if ps fails
    }
    
    // Final fallback: read from /proc directly
    DIR *dir = opendir("/proc");  // Open /proc directory
    if (dir) {
        struct dirent *entry;
        while ((entry = readdir(dir)) != nullptr) {
            // Check if directory name is a number (PID)
            if (entry->d_type == DT_DIR && isdigit(entry->d_name[0])) {
                try {
                    int pid = stoi(entry->d_name);
                    if (pid <= 0 || seen_pids.count(pid)) continue;
                    
                    // Build path to process stat file
                    char path[256];
                    snprintf(path, sizeof(path), "/proc/%d/stat", pid);
                    ifstream stat_file(path);
                    if (!stat_file) continue;
                    
                    Proc p;
                    p.pid = pid;
                    
                    string stat_line;
                    getline(stat_file, stat_line);
                    istringstream stat_iss(stat_line);
                    
                    // Read stat file
                    string comm;
                    stat_iss >> p.pid >> comm >> p.state;
                    
                    // Clean up command name (remove parentheses)
                    if (comm.size() > 2) {
                        p.name = comm.substr(1, comm.size() - 2);
                    } else {
                        p.name = comm;
                    }
                    
                    // Skip fields to get to memory info
                    string dummy;
                    for (int i = 0; i < 19; i++) stat_iss >> dummy;
                    stat_iss >> p.vsize >> p.rss;
                    p.rss *= sysconf(_SC_PAGESIZE);
                    
                    // Try to get CPU usage from /proc/pid/stat
                    long utime, stime;
                    long start_time;
                    stat_iss.seekg(0);
                    for (int i = 0; i < 13; i++) stat_iss >> dummy;
                    stat_iss >> utime >> stime;
                    for (int i = 0; i < 4; i++) stat_iss >> dummy;
                    stat_iss >> start_time;
                    
                    // Calculate CPU usage (simplified)
                    long total_time = utime + stime;
                    p.cpuUsage = total_time / 100.0f; // Rough approximation
                    
                    processes.push_back(p);
                    seen_pids.insert(pid);
                } catch (...) {
                    continue;  // Skip problematic entries
                }
            }
        }
        closedir(dir);
    }
    
    return processes;
}

// Function to get total process count
// Tries multiple methods: top, ps, and /proc directory scanning
int getTotalProcessCount() {
    try {
        // First try top command
        string top_output = exec_command("top -b -n 1 -o %CPU 2>/dev/null | wc -l");
        istringstream iss(top_output);
        int count;
        if (iss >> count) {
            // Subtract header lines (typically 7 lines)
            return count - 7;
        }
    } catch (...) {
        // Fall through to ps attempt
    }

    try {
        // Then try ps command
        string ps_output = exec_command("ps -e --no-headers 2>/dev/null");
        istringstream iss(ps_output);
        return count(istreambuf_iterator<char>(iss), 
                   istreambuf_iterator<char>(), '\n');
    } catch (...) {
        // Fall back to /proc directory scanning
        int count = 0;
        DIR *dir = opendir("/proc");
        if (!dir) return 0;

        struct dirent *entry;
        while ((entry = readdir(dir)) != nullptr) {
            if (entry->d_type == DT_DIR && isdigit(entry->d_name[0])) {
                count++;
            }
        }
        closedir(dir);
        return count;
    }
}

// Function to draw process table using ImGui
// Creates a table with columns: PID, Name, State, CPU, Memory
// Source: https://github.com/ocornut/imgui
// Add these helper functions for selection management
void selectAllProcesses(vector<Proc>& processes, bool select = true) {
    for (auto& p : processes) {
        p.selected = select;
    }
}

void clearProcessSelection(vector<Proc>& processes) {
    selectAllProcesses(processes, false);
}

int countSelectedProcesses(const vector<Proc>& processes) {
    return count_if(processes.begin(), processes.end(), 
        [](const Proc& p) { return p.selected; });
}

// Modified drawProcessTable function with enhanced selection features
void drawProcessTable(vector<Proc>& processes) {
    // Display selection controls above the table
    int selectedCount = countSelectedProcesses(processes);
    ImGui::Text("Selected: %d", selectedCount);
    ImGui::SameLine();
    
    if (ImGui::SmallButton("Select All")) {
        selectAllProcesses(processes);
    }
    ImGui::SameLine();
    
    if (ImGui::SmallButton("Clear")) {
        clearProcessSelection(processes);
    }
    ImGui::SameLine();
    
    // Optional: Select first N processes button
    if (ImGui::SmallButton("Select 3")) {
        clearProcessSelection(processes);
        for (size_t i = 0; i < min(processes.size(), 3ul); i++) {
            processes[i].selected = true;
        }
    }
    ImGui::SameLine();
    
    // Optional: Select processes matching criteria (e.g., high CPU)
    if (ImGui::SmallButton("Select High CPU")) {
        for (auto& p : processes) {
            p.selected = (p.cpuUsage > 10.0f);  // Select processes using >10% CPU
        }
    }

    // Begin the process table
    if (ImGui::BeginTable("Processes", 6, 
        ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | 
        ImGuiTableFlags_ScrollY | ImGuiTableFlags_Resizable)) {
        
        // Set up columns
        ImGui::TableSetupColumn("Select", ImGuiTableColumnFlags_WidthFixed, 30);
        ImGui::TableSetupColumn("PID", ImGuiTableColumnFlags_WidthFixed, 60);
        ImGui::TableSetupColumn("Name", ImGuiTableColumnFlags_WidthStretch);
        ImGui::TableSetupColumn("State", ImGuiTableColumnFlags_WidthFixed, 40);
        ImGui::TableSetupColumn("CPU%", ImGuiTableColumnFlags_WidthFixed, 60);
        ImGui::TableSetupColumn("Memory", ImGuiTableColumnFlags_WidthFixed, 80);
        ImGui::TableHeadersRow();
        
        // Populate table rows
        for (auto& p : processes) {
            ImGui::TableNextRow();
            
            // Selection checkbox
            ImGui::TableSetColumnIndex(0);
            ImGui::PushID(p.pid);
            ImGui::Checkbox("", &p.selected);
            ImGui::PopID();
            
            // PID
            ImGui::TableSetColumnIndex(1);
            ImGui::Text("%d", p.pid);
            
            // Name
            ImGui::TableSetColumnIndex(2);
            ImGui::Text("%s", p.name.c_str());
            
            // State
            ImGui::TableSetColumnIndex(3);
            ImGui::Text("%c", p.state);
            
            // CPU Usage (colored by usage)
            ImGui::TableSetColumnIndex(4);
            if (p.cpuUsage > 50.0f) {
                ImGui::PushStyleColor(ImGuiCol_Text, ImVec4(1.0f, 0.0f, 0.0f, 1.0f)); // Red
            } else if (p.cpuUsage > 20.0f) {
                ImGui::PushStyleColor(ImGuiCol_Text, ImVec4(1.0f, 1.0f, 0.0f, 1.0f)); // Yellow
            }
            ImGui::Text("%.1f%%", p.cpuUsage);
            if (p.cpuUsage > 20.0f) {
                ImGui::PopStyleColor();
            }
            
            // Memory
            ImGui::TableSetColumnIndex(5);
            ImGui::Text("%.1f MB", p.vsize / (1024.0 * 1024.0));
        }
        ImGui::EndTable();
    }

        // Display action buttons when 3 or more processes are selected
    if (selectedCount >= 3) {
        ImGui::Separator();
        ImGui::Text("Actions for %d selected processes:", selectedCount);
        
        if (ImGui::Button("Kill Selected")) {
            // Implement kill functionality
            for (const auto& p : processes) {
                if (p.selected) {
                    // WARNING: Actual kill implementation would go here
                    // string cmd = "kill -9 " + to_string(p.pid);
                    // system(cmd.c_str());
                }
            }
            // Clear selection after action
            clearProcessSelection(processes);
        }
        
        ImGui::SameLine();
        
        if (ImGui::Button("Priority Up")) {
            // Implement priority change
            for (const auto& p : processes) {
                if (p.selected) {
                    // string cmd = "renice -n -5 -p " + to_string(p.pid);
                    // system(cmd.c_str());
                }
            }
        }
        
        ImGui::SameLine();
        
        if (ImGui::Button("Save to File")) {
            // Save selected processes info to file
            ofstream out("selected_processes.txt");
            if (out) {
                for (const auto& p : processes) {
                    if (p.selected) {
                        out << p.pid << "\t" << p.name << "\t" 
                            << p.cpuUsage << "%\t" 
                            << (p.vsize / (1024 * 1024)) << "MB\n";
                    }
                }
            }
        }
    }
}

// Function to get disk information using df command
// Returns vector of DiskInfo structs
vector<DiskInfo> getDiskInfo() {
    vector<DiskInfo> disks;
    FILE* pipe = popen("df -h", "r");  // Execute df -h command
    if (!pipe) return disks;

    char buffer[256];
    // Skip the header line
    fgets(buffer, sizeof(buffer), pipe);
    
    // Process each line of df output
    while (fgets(buffer, sizeof(buffer), pipe) != nullptr) {
        DiskInfo disk;
        istringstream iss(buffer);
        // Parse columns: filesystem, size, used, available, use%, mounted on
        iss >> disk.filesystem 
            >> disk.size 
            >> disk.used 
            >> disk.available 
            >> disk.usePercent;
        
        // The rest of the line is the mount point (might contain spaces)
        string mountedOn;
        getline(iss, mountedOn);
        // Trim leading whitespace
        disk.mountedOn = mountedOn.substr(mountedOn.find_first_not_of(" \t"));
        
        disks.push_back(disk);
    }
    
    pclose(pipe);
    return disks;
}

// Function to draw disk information table using ImGui
// Creates a table with columns: Filesystem, Size, Used, Available, Use%, Mounted On
void drawDiskInfoTable(const vector<DiskInfo>& disks) {
    if (ImGui::BeginTable("Disks", 6, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_ScrollY)) {
        // Set up table columns
        ImGui::TableSetupColumn("Filesystem");
        ImGui::TableSetupColumn("Size");
        ImGui::TableSetupColumn("Used");
        ImGui::TableSetupColumn("Available");
        ImGui::TableSetupColumn("Use%");
        ImGui::TableSetupColumn("Mounted On");
        ImGui::TableHeadersRow();
        
        // Populate table rows with disk data
        for (const auto& disk : disks) {
            ImGui::TableNextRow();
            
            ImGui::TableSetColumnIndex(0);
            ImGui::Text("%s", disk.filesystem.c_str());
            
            ImGui::TableSetColumnIndex(1);
            ImGui::Text("%s", disk.size.c_str());
            
            ImGui::TableSetColumnIndex(2);
            ImGui::Text("%s", disk.used.c_str());
            
            ImGui::TableSetColumnIndex(3);
            ImGui::Text("%s", disk.available.c_str());
            
            ImGui::TableSetColumnIndex(4);
            // Color the percentage based on usage
            int percent = atoi(disk.usePercent.c_str());
            if (percent > 90) {
                ImGui::PushStyleColor(ImGuiCol_Text, ImVec4(1.0f, 0.0f, 0.0f, 1.0f));  // Red
            } else if (percent > 70) {
                ImGui::PushStyleColor(ImGuiCol_Text, ImVec4(1.0f, 1.0f, 0.0f, 1.0f)); // Yellow
            }
            ImGui::Text("%s", disk.usePercent.c_str());
            if (percent > 70) {
                ImGui::PopStyleColor();
            }
            
            ImGui::TableSetColumnIndex(5);
            ImGui::Text("%s", disk.mountedOn.c_str());
        }
        ImGui::EndTable();
    }
}