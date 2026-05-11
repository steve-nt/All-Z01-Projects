// This is a network monitoring implementation for Linux systems that tracks interface statistics
// and displays them in a graphical interface using ImGui. It reads from /proc/net/dev and 
// uses system calls to gather network interface information.

// Header includes
#include "header.h"                 // Custom header file (likely contains ImGui and other project-specific declarations)
#include <chrono>                   // For time measurements (C++ Standard Library)
#include <thread>                   // For thread operations (C++ Standard Library)
#include <cstring>                  // For string operations (C Standard Library)
#include <sstream>                  // For string stream operations (C++ Standard Library)
#include <mutex>                    // For thread synchronization (C++ Standard Library)
#include <algorithm>                // For algorithms like max/min (C++ Standard Library)
#include <sys/socket.h>             // For socket operations (POSIX)
#include <netdb.h>                  // For network database operations (POSIX)
#include <iomanip>                  // For output formatting (C++ Standard Library)
#include <sys/ioctl.h>              // For device control operations (POSIX)
#include <net/if.h>                 // For network interface operations (POSIX)

using namespace std::chrono;        // Use the chrono namespace for time operations

// Structure to store network interface statistics
struct NetInterface {
    string name;                    // Interface name (e.g., eth0, wlan0)
    string ipv4;                    // IPv4 address
    string ipv6;                    // IPv6 address
    bool isUp = false;              // Interface status (up/down)
    long long prevRxBytes = 0;      // Previous received bytes count
    long long prevTxBytes = 0;      // Previous transmitted bytes count
    float speedRx = 0;              // Current receive speed in bytes/sec
    float speedTx = 0;              // Current transmit speed in bytes/sec
    vector<float> historyRx;        // History of receive speeds for graphing
    vector<float> historyTx;        // History of transmit speeds for graphing
    
    // Detailed statistics from /proc/net/dev
    long long rx_bytes = 0;         // Total received bytes
    long long tx_bytes = 0;         // Total transmitted bytes
    long long rx_packets = 0;       // Received packets count
    long long tx_packets = 0;       // Transmitted packets count
    long long rx_errors = 0;        // Receive errors
    long long tx_errors = 0;        // Transmit errors
    long long rx_drops = 0;         // Received packets dropped
    long long tx_drops = 0;         // Transmitted packets dropped
    long long rx_fifo = 0;          // FIFO buffer errors (receive)
    long long rx_frame = 0;         // Frame alignment errors
    long long rx_compressed = 0;    // Compressed packets received
    long long rx_multicast = 0;     // Multicast packets received
    long long tx_fifo = 0;          // FIFO buffer errors (transmit)
    long long tx_colls = 0;         // Collisions detected
    long long tx_carrier = 0;       // Carrier losses
    long long tx_compressed = 0;    // Compressed packets transmitted
};

// Global variables for controlling the monitoring
static bool pauseMonitoring = false;        // Pause flag for monitoring
static float speedMultiplier = 1.0f;        // Speed multiplier for simulation/testing
static map<string, NetInterface> netInterfaces; // Map of interface names to their data
static steady_clock::time_point lastUpdate; // Time point of last update
static float maxGraphValue = 1024 * 1024;   // Maximum value for graph scaling (1 MB/s default)
static const int maxHistory = 300;          // Maximum number of history points to keep
static const float minElapsed = 0.1f;       // Minimum time between updates to prevent spikes
static float smoothingFactor = 0.3f;        // Exponential smoothing factor for speed calculations

// Formats a speed value in bytes/sec to human-readable string (KB/s, MB/s, etc.)
string formatSpeed(float bytesPerSec) {
    std::stringstream ss;           // Create a string stream for formatting
    ss << std::fixed << std::setprecision(2); // Set fixed precision formatting
    
    // Convert to appropriate unit
    if (bytesPerSec >= 1024 * 1024 * 1024) {
        ss << bytesPerSec / (1024 * 1024 * 1024) << " GB/s";
    } 
    else if (bytesPerSec >= 1024 * 1024) {
        ss << bytesPerSec / (1024 * 1024) << " MB/s";
    }
    else if (bytesPerSec >= 1024) {
        ss << bytesPerSec / 1024 << " KB/s";
    }
    else {
        ss << bytesPerSec << " B/s";
    }
    return ss.str();                // Return formatted string
}

// Checks if a network interface is up using ioctl
bool checkInterfaceState(const string& ifname) {
    int sock = socket(AF_INET, SOCK_DGRAM, 0);  // Create a socket for ioctl
    if (sock < 0) {
        return false;               // Return false if socket creation fails
    }

    struct ifreq ifr;               // Interface request structure
    memset(&ifr, 0, sizeof(ifr));   // Clear the structure
    strncpy(ifr.ifr_name, ifname.c_str(), IFNAMSIZ-1); // Copy interface name

    // Get interface flags using ioctl
    if (ioctl(sock, SIOCGIFFLAGS, &ifr) < 0) {
        close(sock);                // Close socket on error
        return false;
    }

    close(sock);                    // Close the socket
    return (ifr.ifr_flags & IFF_UP) != 0; // Return true if IFF_UP flag is set
}

// Updates the IP addresses for a network interface
void updateInterfaceIPs(NetInterface& iface) {
    struct ifaddrs *ifaddr, *ifa;   // Interface address structures
    iface.ipv4.clear();             // Clear existing IPv4
    iface.ipv6.clear();             // Clear existing IPv6

    // Get linked list of network interfaces
    if (getifaddrs(&ifaddr) == -1) {
        perror("getifaddrs");       // Print error if getifaddrs fails
        return;
    }

    // Iterate through all interfaces
    for (ifa = ifaddr; ifa != NULL; ifa = ifa->ifa_next) {
        if (ifa->ifa_addr == NULL || strcmp(ifa->ifa_name, iface.name.c_str()) != 0)
            continue;               // Skip if no address or wrong interface

        // Handle IPv4 addresses
        if (ifa->ifa_addr->sa_family == AF_INET) { 
            char ip[INET_ADDRSTRLEN];                      // Buffer for IP string
            struct sockaddr_in* sa = (struct sockaddr_in*)ifa->ifa_addr; // Cast to IPv4 structure
            inet_ntop(AF_INET, &(sa->sin_addr), ip, INET_ADDRSTRLEN); // Convert binary IP to string
            iface.ipv4 = ip;                              // Store IPv4 address
        }
        // Handle IPv6 addresses
        else if (ifa->ifa_addr->sa_family == AF_INET6) { 
            char ip[INET6_ADDRSTRLEN];                    // Buffer for IP string
            struct sockaddr_in6* sa = (struct sockaddr_in6*)ifa->ifa_addr; // Cast to IPv6 structure
            inet_ntop(AF_INET6, &(sa->sin6_addr), ip, INET6_ADDRSTRLEN); // Convert binary IP to string
            iface.ipv6 = ip;                              // Store IPv6 address
        }
    }

    freeifaddrs(ifaddr);            // Free the linked list
    iface.isUp = checkInterfaceState(iface.name); // Update interface status
}

// Main function to update network statistics from /proc/net/dev
void updateNetworkStats() {
    if (pauseMonitoring) {          // Skip if monitoring is paused
        return;
    }

    ifstream netFile("/proc/net/dev"); // Open the network stats file
    if (!netFile.is_open()) {        // Check if file opened successfully
        cerr << "Failed to open /proc/net/dev" << endl;
        return;
    }

    string line;
    
    // Skip the first two header lines in /proc/net/dev
    for (int i = 0; i < 2; i++) {
        if (!getline(netFile, line)) {
            cerr << "Unexpected end of /proc/net/dev" << endl;
            netFile.close();
            return;
        }
    }

    // Calculate time since last update
    auto now = steady_clock::now();
    float elapsed = duration<float>(now - lastUpdate).count();
    lastUpdate = now;
    elapsed *= speedMultiplier;      // Apply speed multiplier

    // Ensure minimum elapsed time to prevent spikes from very small time intervals
    elapsed = max(elapsed, minElapsed);

    static float ipUpdateTimer = 0;  // Timer for periodic IP updates
    ipUpdateTimer += elapsed;        // Update the timer

    // Process each interface in /proc/net/dev
    while (getline(netFile, line)) {
        istringstream iss(line);     // Create string stream from line
        string iface;                // Interface name
        getline(iss, iface, ':');    // Get interface name (before colon)
        iface.erase(remove(iface.begin(), iface.end(), ' '), iface.end()); // Remove spaces

        if (iface.empty()) continue; // Skip empty interface names

        // Structure to hold all statistics from the line
        struct {
            long long rx_bytes, rx_packets, rx_errs, rx_drop, rx_fifo, rx_frame, rx_compressed, rx_multicast;
            long long tx_bytes, tx_packets, tx_errs, tx_drop, tx_fifo, tx_colls, tx_carrier, tx_compressed;
        } stats;

        // Parse all statistics from the line
        if (!(iss >> stats.rx_bytes >> stats.rx_packets >> stats.rx_errs >> stats.rx_drop
            >> stats.rx_fifo >> stats.rx_frame >> stats.rx_compressed >> stats.rx_multicast
            >> stats.tx_bytes >> stats.tx_packets >> stats.tx_errs >> stats.tx_drop
            >> stats.tx_fifo >> stats.tx_colls >> stats.tx_carrier >> stats.tx_compressed)) {
            cerr << "Failed to parse line for interface: " << iface << endl;
            continue;
        }

        // Get or create interface data in the map
        auto &ifaceData = netInterfaces[iface];
        ifaceData.name = iface;

        // Update IP addresses periodically (every 10 seconds)
        if (ifaceData.ipv4.empty() || ipUpdateTimer > 10.0f) {
            updateInterfaceIPs(ifaceData);
            ipUpdateTimer = 0;       // Reset timer after update
        }

        // Calculate network speeds if we have previous values
        if (ifaceData.prevRxBytes > 0 && ifaceData.prevTxBytes > 0 && elapsed > 0) {
            // Calculate raw speeds
            float newRx = float(stats.rx_bytes - ifaceData.prevRxBytes) / elapsed;
            float newTx = float(stats.tx_bytes - ifaceData.prevTxBytes) / elapsed;
            
            // Apply exponential smoothing to reduce noise
            ifaceData.speedRx = smoothingFactor * newRx + (1.0f - smoothingFactor) * ifaceData.speedRx;
            ifaceData.speedTx = smoothingFactor * newTx + (1.0f - smoothingFactor) * ifaceData.speedTx;

            // Adjust graph scale if needed
            float currentMax = std::max(ifaceData.speedRx, ifaceData.speedTx);
            if (currentMax > maxGraphValue * 0.8f) {
                maxGraphValue = currentMax * 1.2f;
            }
        }

        // Store current values as previous for next update
        ifaceData.prevRxBytes = stats.rx_bytes;
        ifaceData.prevTxBytes = stats.tx_bytes;

        // Update all interface statistics
        ifaceData.rx_bytes = stats.rx_bytes;
        ifaceData.tx_bytes = stats.tx_bytes;
        ifaceData.rx_packets = stats.rx_packets;
        ifaceData.tx_packets = stats.tx_packets;
        ifaceData.rx_errors = stats.rx_errs;
        ifaceData.tx_errors = stats.tx_errs;
        ifaceData.rx_drops = stats.rx_drop;
        ifaceData.tx_drops = stats.tx_drop;
        ifaceData.rx_fifo = stats.rx_fifo;
        ifaceData.rx_frame = stats.rx_frame;
        ifaceData.rx_compressed = stats.rx_compressed;
        ifaceData.rx_multicast = stats.rx_multicast;
        ifaceData.tx_fifo = stats.tx_fifo;
        ifaceData.tx_colls = stats.tx_colls;
        ifaceData.tx_carrier = stats.tx_carrier;
        ifaceData.tx_compressed = stats.tx_compressed;

        // Maintain history for graphing
        if (ifaceData.historyRx.size() >= maxHistory) {
            ifaceData.historyRx.erase(ifaceData.historyRx.begin());
            ifaceData.historyTx.erase(ifaceData.historyTx.begin());
        }
        ifaceData.historyRx.push_back(ifaceData.speedRx);
        ifaceData.historyTx.push_back(ifaceData.speedTx);
    }

    netFile.close();                // Close the network stats file
}

// Gets a list of all network interfaces with basic information
vector<NetworkInterface> getNetworkInterfaces() {
    vector<NetworkInterface> interfaces;
    struct ifaddrs *ifaddr, *ifa;
    
    // Get list of all network interfaces
    if (getifaddrs(&ifaddr) == -1) {
        perror("getifaddrs");
        return interfaces;
    }

    // Iterate through all interfaces
    for (ifa = ifaddr; ifa != NULL; ifa = ifa->ifa_next) {
        if (ifa->ifa_addr == NULL) continue;

        NetworkInterface ni;
        ni.name = ifa->ifa_name;

        // Add statistics if we have them
        if (netInterfaces.count(ni.name)) {
            auto &netData = netInterfaces[ni.name];
            ni.rx_bytes = netData.prevRxBytes;
            ni.tx_bytes = netData.prevTxBytes;
            ni.rx_packets = netData.rx_packets;
            ni.tx_packets = netData.tx_packets;
        }

        // Get IPv4 address if available
        if (ifa->ifa_addr->sa_family == AF_INET) {
            struct sockaddr_in *sa = (struct sockaddr_in *)ifa->ifa_addr;
            char ip[INET_ADDRSTRLEN];
            inet_ntop(AF_INET, &(sa->sin_addr), ip, INET_ADDRSTRLEN);
            ni.ipv4 = ip;
        }
        // Get IPv6 address if available
        else if (ifa->ifa_addr->sa_family == AF_INET6) {
            struct sockaddr_in6 *sa = (struct sockaddr_in6 *)ifa->ifa_addr;
            char ip[INET6_ADDRSTRLEN];
            inet_ntop(AF_INET6, &(sa->sin6_addr), ip, INET6_ADDRSTRLEN);
            ni.ipv6 = ip;
        }

        // Check if we already have this interface in our list
        bool found = false;
        for (auto& existing : interfaces) {
            if (existing.name == ni.name) {
                found = true;
                if (!ni.ipv4.empty()) existing.ipv4 = ni.ipv4;
                if (!ni.ipv6.empty()) existing.ipv6 = ni.ipv6;
                break;
            }
        }
        
        // Add new interface if not found
        if (!found) {
            interfaces.push_back(ni);
        }
    }

    freeifaddrs(ifaddr);            // Free the interface list
    return interfaces;
}

// Draws the network monitoring window using ImGui
void drawNetworkWindow(const char *id, ImVec2 size, ImVec2 position) {
    if (!pauseMonitoring) {          // Update stats if not paused
        updateNetworkStats();
    }

    ImGui::Begin(id);               // Start the ImGui window
    ImGui::SetWindowSize(size);     // Set window size
    ImGui::SetWindowPos(position);  // Set window position

    // Pause/Resume button
    if (ImGui::Button(pauseMonitoring ? "Resume Monitoring" : "Pause Monitoring")) {
        pauseMonitoring = !pauseMonitoring;
    }

    // Create tabbed interface
    if (ImGui::BeginTabBar("NetworkTabs")) {
        // Traffic tab
        if (ImGui::BeginTabItem("Traffic")) {
            // Show monitoring status
            ImGui::TextColored(pauseMonitoring ? ImVec4(1.0f, 0.5f, 0.5f, 1.0f) : ImVec4(0.5f, 1.0f, 0.5f, 1.0f),
                "%s", pauseMonitoring ? "Monitoring Paused" : "Monitoring Active");
            
            // Graph scale control
            ImGui::SliderFloat("Max Graph Value", &maxGraphValue, 1024, 1024*1024*10, "%.0f bytes/s");
            
            // Display each interface
            for (auto &[name, iface] : netInterfaces) {
                if (ImGui::CollapsingHeader((name + (iface.isUp ? " (UP)" : " (DOWN)")).c_str())) {
                    // Show IP addresses if available
                    if (!iface.ipv4.empty() || !iface.ipv6.empty()) {
                        ImGui::Text("IP Addresses:");
                        ImGui::Indent();
                        if (!iface.ipv4.empty()) {
                            ImGui::BulletText("IPv4: %s", iface.ipv4.c_str());
                        }
                        if (!iface.ipv6.empty()) {
                            ImGui::BulletText("IPv6: %s", iface.ipv6.c_str());
                        }
                        ImGui::Unindent();
                        ImGui::Separator();
                    }

                    // Show current speeds
                    ImGui::Text("Receive: %s", formatSpeed(iface.speedRx).c_str());
                    ImGui::Text("Transmit: %s", formatSpeed(iface.speedTx).c_str());
                    
                    // Detailed stats tooltip
                    if (ImGui::IsItemHovered()) {
                        ImGui::BeginTooltip();
                        ImGui::Text("Detailed Stats for %s", name.c_str());
                        ImGui::Separator();
                        ImGui::Text("RX Bytes: %lld", iface.prevRxBytes);
                        ImGui::Text("RX Packets: %lld", iface.rx_packets);
                        ImGui::Text("RX Errors: %lld", iface.rx_errors);
                        ImGui::Text("RX Drops: %lld", iface.rx_drops);
                        ImGui::Separator();
                        ImGui::Text("TX Bytes: %lld", iface.prevTxBytes);
                        ImGui::Text("TX Packets: %lld", iface.tx_packets);
                        ImGui::Text("TX Errors: %lld", iface.tx_errors);
                        ImGui::Text("TX Drops: %lld", iface.tx_drops);
                        ImGui::EndTooltip();
                    }

                    // Prepare labels for graphs
                    string rxLabel = name + " RX";
                    string txLabel = name + " TX";
                    
                    // Draw receive speed graph
                    ImGui::PlotLines(rxLabel.c_str(), iface.historyRx.data(), 
                                    iface.historyRx.size(), 0, nullptr, 
                                    0.0f, maxGraphValue, ImVec2(0, 80));
                    
                    // Draw transmit speed graph
                    ImGui::PlotLines(txLabel.c_str(), iface.historyTx.data(), 
                                    iface.historyTx.size(), 0, nullptr, 
                                    0.0f, maxGraphValue, ImVec2(0, 80));
                    
                    ImGui::Separator();
                }
            }
            ImGui::EndTabItem();
        }

        // Interfaces tab - shows all interfaces in a table
        if (ImGui::BeginTabItem("Interfaces")) {
            if (ImGui::BeginTable("NetworkInterfaces", 6, 
                ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | 
                ImGuiTableFlags_Resizable | ImGuiTableFlags_ScrollY)) {
                
                // Set up table columns
                ImGui::TableSetupColumn("Interface");
                ImGui::TableSetupColumn("Status");
                ImGui::TableSetupColumn("IPv4");
                ImGui::TableSetupColumn("IPv6");
                ImGui::TableSetupColumn("RX Speed");
                ImGui::TableSetupColumn("TX Speed");
                ImGui::TableHeadersRow();

                // Add rows for each interface
                for (auto& [name, iface] : netInterfaces) {
                    ImGui::TableNextRow();
                    
                    // Interface name
                    ImGui::TableSetColumnIndex(0);
                    ImGui::Text("%s", name.c_str());
                    
                    // Status (colored)
                    ImGui::TableSetColumnIndex(1);
                    ImGui::TextColored(iface.isUp ? ImVec4(0.0f, 1.0f, 0.0f, 1.0f) : ImVec4(1.0f, 0.0f, 0.0f, 1.0f),
                        "%s", iface.isUp ? "UP" : "DOWN");
                    
                    // IPv4 address
                    ImGui::TableSetColumnIndex(2);
                    ImGui::Text("%s", iface.ipv4.empty() ? "N/A" : iface.ipv4.c_str());
                    
                    // IPv6 address
                    ImGui::TableSetColumnIndex(3);
                    ImGui::Text("%s", iface.ipv6.empty() ? "N/A" : iface.ipv6.c_str());
                    
                    // Receive speed
                    ImGui::TableSetColumnIndex(4);
                    ImGui::Text("%s", formatSpeed(iface.speedRx).c_str());
                    
                    // Transmit speed
                    ImGui::TableSetColumnIndex(5);
                    ImGui::Text("%s", formatSpeed(iface.speedTx).c_str());
                }
                
                ImGui::EndTable();
            }
            ImGui::EndTabItem();
        }

        ImGui::EndTabBar();
    }

    ImGui::End();                   // End the ImGui window
}