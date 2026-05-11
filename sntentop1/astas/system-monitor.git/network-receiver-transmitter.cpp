#include "header.h"
#include <imgui.h>
#include <fstream>
#include <sstream>
#include <string>
#include <map>
#include <algorithm> // For std::min

struct NetStats {
    uint64_t rx_bytes = 0, rx_packets = 0, rx_errs = 0, rx_drop = 0, rx_fifo = 0, rx_frame = 0, rx_compressed = 0, rx_multicast = 0;
    uint64_t tx_bytes = 0, tx_packets = 0, tx_errs = 0, tx_drop = 0, tx_fifo = 0, tx_colls = 0, tx_carrier = 0, tx_compressed = 0;
};

std::map<std::string, NetStats> readNetworkStats() {
    std::map<std::string, NetStats> stats;
    std::ifstream file("/proc/net/dev");
    std::string line;

    // Skip headers
    std::getline(file, line);
    std::getline(file, line);

    while (std::getline(file, line)) {
        std::istringstream iss(line);
        std::string iface;
        NetStats ns;

        std::getline(iss, iface, ':');
        iface.erase(0, iface.find_first_not_of(" ")); // Trim spaces

        iss >> ns.rx_bytes >> ns.rx_packets >> ns.rx_errs >> ns.rx_drop >> ns.rx_fifo >> ns.rx_frame >> ns.rx_compressed >> ns.rx_multicast
            >> ns.tx_bytes >> ns.tx_packets >> ns.tx_errs >> ns.tx_drop >> ns.tx_fifo >> ns.tx_colls >> ns.tx_carrier >> ns.tx_compressed;

        stats[iface] = ns;
    }

    return stats;
}

std::string formatBytes(uint64_t bytes) {
    const char* unit = "B";
    double val = static_cast<double>(bytes);

    if (val >= 1024.0 && val < 1024.0 * 1024.0) {
        val /= 1024.0;
        unit = "KB";
    } else if (val >= 1024.0 * 1024.0 && val < 1024.0 * 1024.0 * 1024.0) {
        val /= (1024.0 * 1024.0);
        unit = "MB";
    } else if (val >= 1024.0 * 1024.0 * 1024.0) {
        val /= (1024.0 * 1024.0 * 1024.0);
        unit = "GB";
    }

    char buf[64];
    snprintf(buf, sizeof(buf), "%.2f %s", val, unit);
    return std::string(buf);
}

void RenderExtraNetworkWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::SetNextWindowSize(size, ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowPos(position, ImGuiCond_FirstUseEver);

    if (!ImGui::Begin(id)) { ImGui::End(); return; }

    auto stats = readNetworkStats();

    // --- RX and TX tables ---
    if (ImGui::BeginTabBar("NetTab")) {
        if (ImGui::BeginTabItem("RX")) {
            if (ImGui::BeginTable("RXTable", 9, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg)) {
                ImGui::TableSetupColumn("Interface");
                ImGui::TableSetupColumn("Bytes");
                ImGui::TableSetupColumn("Packets");
                ImGui::TableSetupColumn("Errors");
                ImGui::TableSetupColumn("Drop");
                ImGui::TableSetupColumn("FIFO");
                ImGui::TableSetupColumn("Frame");
                ImGui::TableSetupColumn("Compressed");
                ImGui::TableSetupColumn("Multicast");
                ImGui::TableHeadersRow();

                for (const auto& [iface, ns] : stats) {
                    ImGui::TableNextRow();
                    ImGui::TableSetColumnIndex(0); ImGui::Text("%s", iface.c_str());
                    ImGui::TableSetColumnIndex(1); ImGui::Text("%llu", ns.rx_bytes);
                    ImGui::TableSetColumnIndex(2); ImGui::Text("%llu", ns.rx_packets);
                    ImGui::TableSetColumnIndex(3); ImGui::Text("%llu", ns.rx_errs);
                    ImGui::TableSetColumnIndex(4); ImGui::Text("%llu", ns.rx_drop);
                    ImGui::TableSetColumnIndex(5); ImGui::Text("%llu", ns.rx_fifo);
                    ImGui::TableSetColumnIndex(6); ImGui::Text("%llu", ns.rx_frame);
                    ImGui::TableSetColumnIndex(7); ImGui::Text("%llu", ns.rx_compressed);
                    ImGui::TableSetColumnIndex(8); ImGui::Text("%llu", ns.rx_multicast);
                }

                ImGui::EndTable();
            }
            ImGui::EndTabItem();
        }

        if (ImGui::BeginTabItem("TX")) {
            if (ImGui::BeginTable("TXTable", 9, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg)) {
                ImGui::TableSetupColumn("Interface");
                ImGui::TableSetupColumn("Bytes");
                ImGui::TableSetupColumn("Packets");
                ImGui::TableSetupColumn("Errors");
                ImGui::TableSetupColumn("Drop");
                ImGui::TableSetupColumn("FIFO");
                ImGui::TableSetupColumn("Colls");
                ImGui::TableSetupColumn("Carrier");
                ImGui::TableSetupColumn("Compressed");
                ImGui::TableHeadersRow();

                for (const auto& [iface, ns] : stats) {
                    ImGui::TableNextRow();
                    ImGui::TableSetColumnIndex(0); ImGui::Text("%s", iface.c_str());
                    ImGui::TableSetColumnIndex(1); ImGui::Text("%llu", ns.tx_bytes);
                    ImGui::TableSetColumnIndex(2); ImGui::Text("%llu", ns.tx_packets);
                    ImGui::TableSetColumnIndex(3); ImGui::Text("%llu", ns.tx_errs);
                    ImGui::TableSetColumnIndex(4); ImGui::Text("%llu", ns.tx_drop);
                    ImGui::TableSetColumnIndex(5); ImGui::Text("%llu", ns.tx_fifo);
                    ImGui::TableSetColumnIndex(6); ImGui::Text("%llu", ns.tx_colls);
                    ImGui::TableSetColumnIndex(7); ImGui::Text("%llu", ns.tx_carrier);
                    ImGui::TableSetColumnIndex(8); ImGui::Text("%llu", ns.tx_compressed);
                }

                ImGui::EndTable();
            }
            ImGui::EndTabItem();
        }

        // --- RX/TX Usage Visual Display Tabs ---
        if (ImGui::BeginTabItem("RX Usage")) {
            for (const auto& [iface, ns] : stats) {
                float gb = ns.rx_bytes / (1024.0f * 1024.0f * 1024.0f);
                float progress = std::min(gb / 2.0f, 1.0f); // Clamp to [0, 1]
                std::string label = iface + " - " + formatBytes(ns.rx_bytes) +
                                    " (" + std::to_string(ns.rx_bytes) + " bytes)";

                ImGui::Text("%s", label.c_str());
                ImGui::ProgressBar(progress, ImVec2(-1, 0));
            }
            ImGui::EndTabItem();
        }

        if (ImGui::BeginTabItem("TX Usage")) {
            for (const auto& [iface, ns] : stats) {
                float gb = ns.tx_bytes / (1024.0f * 1024.0f * 1024.0f);
                float progress = std::min(gb / 2.0f, 1.0f); // Clamp to [0, 1]
                std::string label = iface + " - " + formatBytes(ns.tx_bytes) +
                                    " (" + std::to_string(ns.tx_bytes) + " bytes)";

                ImGui::Text("%s", label.c_str());
                ImGui::ProgressBar(progress, ImVec2(-1, 0));
            }
            ImGui::EndTabItem();
        }

        ImGui::EndTabBar();
    }

    ImGui::End();
}
