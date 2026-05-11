#include "header.h"
#include <imgui.h>
#include <string>
#include <vector>
#include <ifaddrs.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <net/if.h>

struct NetInterface {
    std::string name;
    std::string ipv4;
};

std::vector<NetInterface> getNetworkInterfaces() {
    std::vector<NetInterface> interfaces;

    struct ifaddrs *ifaddr, *ifa;
    if (getifaddrs(&ifaddr) == -1) return interfaces;

    for (ifa = ifaddr; ifa != nullptr; ifa = ifa->ifa_next) {
        if (!ifa->ifa_addr || ifa->ifa_addr->sa_family != AF_INET) continue;

        char addr[INET_ADDRSTRLEN];
        void* in_addr = &((struct sockaddr_in*)ifa->ifa_addr)->sin_addr;

        if (inet_ntop(AF_INET, in_addr, addr, sizeof(addr))) {
            interfaces.push_back({ifa->ifa_name, addr});
        }
    }

    freeifaddrs(ifaddr);
    return interfaces;
}

void rendernetworkWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::SetNextWindowSize(size, ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowPos(position, ImGuiCond_FirstUseEver);

    if (!ImGui::Begin(id)) {
        ImGui::End();
        return;
    }

#if defined(__linux__)

    auto interfaces = getNetworkInterfaces();

    if (ImGui::BeginTable("NetTable", 2, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg)) {
        ImGui::TableSetupColumn("Interface");
        ImGui::TableSetupColumn("IPv4 Address");
        ImGui::TableHeadersRow();

        for (const auto& net : interfaces) {
            ImGui::TableNextRow();
            ImGui::TableSetColumnIndex(0); ImGui::Text("%s", net.name.c_str());
            ImGui::TableSetColumnIndex(1); ImGui::Text("%s", net.ipv4.c_str());
        }

        ImGui::EndTable();
    }

#else
    ImGui::TextColored(ImVec4(1, 0, 0, 1), "Network monitoring not supported on this OS yet.");
#endif

    ImGui::End();
}
