#include "header.h"
#include <fstream>
#include <sstream>
#include <map>
#include <ifaddrs.h>
#include <netinet/in.h>
#include <arpa/inet.h>

std::map<std::string, NetStats> getNetworkStats() {
    std::map<std::string, NetStats> stats;
    std::ifstream file("/proc/net/dev");
    std::string line;
    int lineNum = 0;
    while (std::getline(file, line)) {
        lineNum++;
        if (lineNum <= 2) continue; // skip headers
        std::istringstream iss(line);
        std::string iface;
        std::getline(iss, iface, ':');
        // trim spaces
        iface.erase(0, iface.find_first_not_of(" \t"));
        iface.erase(iface.find_last_not_of(" \t") + 1);
        if (iface.empty()) continue;
        NetStats ns = {};
        iss >> ns.rx.bytes >> ns.rx.packets >> ns.rx.errs >> ns.rx.drop >> ns.rx.fifo >> ns.rx.frame >> ns.rx.compressed >> ns.rx.multicast
            >> ns.tx.bytes >> ns.tx.packets >> ns.tx.errs >> ns.tx.drop >> ns.tx.fifo >> ns.tx.colls >> ns.tx.carrier >> ns.tx.compressed;
        stats[iface] = ns;
    }
    return stats;
}

std::map<std::string, std::string> getIPv4Addresses() {
    std::map<std::string, std::string> result;
    struct ifaddrs *ifaddr, *ifa;
    if (getifaddrs(&ifaddr) == -1) return result;
    for (ifa = ifaddr; ifa != nullptr; ifa = ifa->ifa_next) {
        if (!ifa->ifa_addr) continue;
        if (ifa->ifa_addr->sa_family == AF_INET) {
            char addr[INET_ADDRSTRLEN];
            void* in_addr = &((struct sockaddr_in*)ifa->ifa_addr)->sin_addr;
            inet_ntop(AF_INET, in_addr, addr, INET_ADDRSTRLEN);
            result[ifa->ifa_name] = addr;
        }
    }
    freeifaddrs(ifaddr);
    return result;
}
