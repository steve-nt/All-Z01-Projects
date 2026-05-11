#include "header.h"
#include "cpu.h"
#include <SDL2/SDL.h>
#include <chrono>
#include "temp.h"
#include "mem.h"
#include "process.h"
#include "fan.h"
#include <set>


// === OpenGL Loader Definitions ===
#if defined(IMGUI_IMPL_OPENGL_LOADER_GL3W)
#include <GL/gl3w.h>
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLEW)
#include <GL/glew.h>
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD)
#include <glad/glad.h>
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLAD2)
#include <glad/gl.h>
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING2)
#define GLFW_INCLUDE_NONE
#include <glbinding/Binding.h>
#include <glbinding/gl/gl.h>
using namespace gl;
#elif defined(IMGUI_IMPL_OPENGL_LOADER_GLBINDING3)
#define GLFW_INCLUDE_NONE
#include <glbinding/glbinding.h>
#include <glbinding/gl/gl.h>
using namespace gl;
#else
#include IMGUI_IMPL_OPENGL_LOADER_CUSTOM
#endif

// === System Info Window ===
void systemWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::Begin(id);
    ImGui::SetWindowSize(id, size);
    ImGui::SetWindowPos(id, position);

    ImGui::Text("[U] User: %s", getUserName().c_str());
    ImGui::Text("[H] Host: %s", getHostName().c_str());
    ImGui::Text("[O] OS: %s", getOSName().c_str());
    ImGui::Text("[C] CPU: %s", getCPUModel().c_str());
    ImGui::Separator();

    if (ImGui::BeginTabBar("SystemTabs")) {

        if (ImGui::BeginTabItem("CPU")) {
            static std::vector<float> cpuHistory;
            const int graphSize = 100;
            static double lastUpdateTime = 0;
            static float currentUsage = 0.0f;

            static bool pauseGraph = false;
            static float refreshRate = 0.3f;
            static float yAxisMax = 100.0f;

            ImGui::Checkbox("Animate", &pauseGraph);
            ImGui::SliderFloat("Update Rate (sec)", &refreshRate, 0.2f, 1.0f);
            ImGui::SameLine();
            ImGui::Text("- FPS: %.1f", 1.0f / refreshRate);
            ImGui::SliderFloat("Y-Axis Scale", &yAxisMax, 10.0f, 200.0f);

            double currentTime = ImGui::GetTime();
            if (!pauseGraph && (currentTime - lastUpdateTime >= refreshRate)) {
                currentUsage = getCPUUsage();
                lastUpdateTime = currentTime;

                if (cpuHistory.size() >= graphSize)
                    cpuHistory.erase(cpuHistory.begin());
                cpuHistory.push_back(currentUsage);
            }

            ImVec2 graphSizeVec = ImVec2(ImGui::GetContentRegionAvail().x, 100);
            ImGui::PlotLines("##CPU", cpuHistory.data(), cpuHistory.size(), 0, nullptr, 0.0f, yAxisMax, graphSizeVec);

            ImGui::SetCursorPosX(ImGui::GetCursorPosX() - graphSizeVec.x);
            ImGui::SetCursorPosY(ImGui::GetCursorPosY() - graphSizeVec.y - ImGui::GetTextLineHeight());
            ImGui::Text("Current CPU: %.1f%%", currentUsage);

            ImGui::EndTabItem();
        }

        if (ImGui::BeginTabItem("Fan")) {
            static std::vector<float> fanHistory;
            static double lastFanUpdate = 0.0;
            static float refreshRate = 1.0f;
            static float yMax = 6000.0f;
            static bool animate = true;
            static double lastFanCheck = 0.0;
            static int fanSpeed = 0;

            double now = ImGui::GetTime();
            if (now - lastFanCheck > 0.5) {
                fanSpeed = getFanSpeed();
                lastFanCheck = now;
            }

            bool active = fanSpeed > 0;
            std::string level = fanSpeed > 2500 ? "Manual" : "Auto";

            ImGui::Text("Fan Status: %s", active ? "Enabled" : "Disabled");
            ImGui::Text("Fan Speed: %d RPM", fanSpeed);
            ImGui::Text("Level: %s", level.c_str());

            ImGui::Separator();
            ImGui::Checkbox("Animate", &animate);
            ImGui::SliderFloat("FPS", &refreshRate, 0.1f, 2.0f);
            ImGui::SliderFloat("Max Scale", &yMax, 1000.0f, 10000.0f);

            double currentTime = ImGui::GetTime();
            if (animate && (currentTime - lastFanUpdate >= (1.0f / refreshRate))) {
                lastFanUpdate = currentTime;

                if (fanHistory.size() >= 100)
                    fanHistory.erase(fanHistory.begin());

                fanHistory.push_back((float)fanSpeed);
            }

            ImVec2 graphSize = ImVec2(ImGui::GetContentRegionAvail().x, 100);
            ImGui::PlotLines("##FanGraph", fanHistory.data(), fanHistory.size(), 0, nullptr, 0.0f, yMax, graphSize);

            ImGui::SetCursorPosX(ImGui::GetCursorPosX() - graphSize.x);
            ImGui::SetCursorPosY(ImGui::GetCursorPosY() - graphSize.y - ImGui::GetTextLineHeight());
            ImGui::Text("Current Fan Speed: %d RPM", fanSpeed);

            ImGui::EndTabItem();
        }

        if (ImGui::BeginTabItem("Thermal")) {
    static std::vector<float> tempHistory;
    static double lastTempUpdate = 0.0;
    static float refreshRate = 1.0f;
    static float yMax = 100.0f;
    static bool animate = true;
    static float currentTemp = 0.0f;

    double now = ImGui::GetTime();

    if (animate && (now - lastTempUpdate >= (1.0f / refreshRate))) {
        currentTemp = getCPUTemperature();
        lastTempUpdate = now;

        if (tempHistory.size() >= 100)
            tempHistory.erase(tempHistory.begin());
        tempHistory.push_back(currentTemp);
    }

    // === Label ===
    ImGui::Text("temperature = %.1f °C", currentTemp);

    // === Controls ===
    ImGui::Checkbox("Animate", &animate);
    ImGui::SliderFloat("FPS", &refreshRate, 0.2f, 2.0f);
    ImGui::SliderFloat("Max scale", &yMax, 30.0f, 120.0f);

    // === Graph ===
    ImVec2 graphSize = ImVec2(ImGui::GetContentRegionAvail().x, 100);
    ImGui::PlotLines("##ThermalGraph", tempHistory.data(), tempHistory.size(), 0,
                     nullptr, 20.0f, yMax, graphSize);

    // === Overlay: temp = XX ===
    ImGui::SetCursorPosX(ImGui::GetCursorPosX() - graphSize.x);
    ImGui::SetCursorPosY(ImGui::GetCursorPosY() - graphSize.y - ImGui::GetTextLineHeight());
    ImGui::Text("temp = %.2f", currentTemp);

    ImGui::EndTabItem();
}

        ImGui::EndTabBar();
    }

    // Process state breakdown
    int running = 0, sleeping = 0, stopped = 0, zombie = 0, other = 0;
    for (const auto& proc : getProcesses()) {
        char state = proc.state.empty() ? '?' : proc.state[0];
        switch (state) {
            case 'R': running++; break;
            case 'S': sleeping++; break;
            case 'T': stopped++; break;
            case 'Z': zombie++; break;
            default: other++; break;
        }
    }
    int total = running + sleeping + stopped + zombie + other;
    ImGui::Text("Total number of processes: %d", total);
    ImGui::Text("Processes: Running: %d, Sleeping: %d, Stopped: %d, Zombie: %d, Other: %d", running, sleeping, stopped, zombie, other);

    ImGui::End();
}

// === Memory and Processes Window ===
void memoryProcessesWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::Begin(id);
    ImGui::SetWindowSize(id, size);
    ImGui::SetWindowPos(id, position);

    // === Στατιστικά μνήμης από το /proc/meminfo
    MemoryStats stats = getMemoryStats();

    // === Υπολόγισε ποσοστά
    float ramUsageRatio = stats.usedRamMB / stats.totalRamMB;
    float swapUsageRatio = stats.usedSwapMB / stats.totalSwapMB;

    // === Εμφάνιση RAM
    ImGui::Text("RAM Usage: %.2f MB / %.2f MB", stats.usedRamMB, stats.totalRamMB);
    ImGui::ProgressBar(ramUsageRatio, ImVec2(-1.0f, 20.0f));

    // === Εμφάνιση SWAP
    ImGui::Text("SWAP Usage: %.2f MB / %.2f MB", stats.usedSwapMB, stats.totalSwapMB);
    ImGui::ProgressBar(swapUsageRatio, ImVec2(-1.0f, 20.0f));

    // === Εμφάνιση Δίσκου ===
    DiskStats disk = getDiskStats();
    float diskUsageRatio = (disk.totalGB > 0.0f) ? (disk.usedGB / disk.totalGB) : 0.0f;
    ImGui::Text("Disk Usage: %.2f GB / %.2f GB", disk.usedGB, disk.totalGB);
    ImGui::ProgressBar(diskUsageRatio, ImVec2(-1.0f, 20.0f));

    ImGui::Separator();
    if (ImGui::BeginTabBar("MemProcTabs")) {
        if (ImGui::BeginTabItem("Processes")) {
            static char filter[64] = "";
            ImGui::InputText("Filter", filter, sizeof(filter));
            static std::set<int> selectedPIDs;
            // === Πίνακας διεργασιών ===
            if (ImGui::BeginTable("Process Table", 5, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_Resizable)) {
                ImGui::TableSetupColumn("PID");
                ImGui::TableSetupColumn("Name");
                ImGui::TableSetupColumn("State");
                ImGui::TableSetupColumn("CPU (%)");
                ImGui::TableSetupColumn("Memory (%)");
                ImGui::TableHeadersRow();
                for (const auto &proc : getProcesses()) {
                    std::string pidStr = std::to_string(proc.pid);
                    if (filter[0] && proc.name.find(filter) == std::string::npos && pidStr.find(filter) == std::string::npos) continue;
                    ImGui::TableNextRow();
                    ImGui::TableSetColumnIndex(0);
                    bool selected = selectedPIDs.count(proc.pid) > 0;
                    if (ImGui::Selectable((pidStr + "##row").c_str(), selected, ImGuiSelectableFlags_SpanAllColumns)) {
                        if (ImGui::GetIO().KeyCtrl) {
                            if (selected) selectedPIDs.erase(proc.pid);
                            else selectedPIDs.insert(proc.pid);
                        } else {
                            selectedPIDs.clear();
                            selectedPIDs.insert(proc.pid);
                        }
                    }
                    ImGui::TableSetColumnIndex(1); ImGui::Text("%s", proc.name.c_str());
                    ImGui::TableSetColumnIndex(2); ImGui::Text("%s", proc.state.c_str());
                    ImGui::TableSetColumnIndex(3); ImGui::Text("%.2f", proc.cpuUsage);
                    ImGui::TableSetColumnIndex(4); ImGui::Text("%.2f", proc.memUsage);
                }
                ImGui::EndTable();
            }
            ImGui::Text("Selected processes: %d", (int)selectedPIDs.size());
            ImGui::EndTabItem();
        }
        // Future tabs can be added here
        ImGui::EndTabBar();
    }

    ImGui::End();
}



// === Network Window ===
void networkWindow(const char *id, ImVec2 size, ImVec2 position)
{
    ImGui::Begin(id);
    ImGui::SetWindowSize(id, size);
    ImGui::SetWindowPos(id, position);

    auto stats = getNetworkStats();
    auto ipv4s = getIPv4Addresses();
    if (ImGui::BeginTabBar("NetworkTabs")) {
        if (ImGui::BeginTabItem("RX (Network Receiver)")) {
            if (ImGui::BeginTable("RXTable", 10, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_Resizable)) {
                ImGui::TableSetupColumn("Interface");
                ImGui::TableSetupColumn("IPv4");
                ImGui::TableSetupColumn("Bytes");
                ImGui::TableSetupColumn("Packets");
                ImGui::TableSetupColumn("Errs");
                ImGui::TableSetupColumn("Drop");
                ImGui::TableSetupColumn("Fifo");
                ImGui::TableSetupColumn("Frame");
                ImGui::TableSetupColumn("Compressed");
                ImGui::TableSetupColumn("Multicast");
                ImGui::TableHeadersRow();
                for (const auto& [iface, ns] : stats) {
                    ImGui::TableNextRow();
                    ImGui::TableSetColumnIndex(0); ImGui::Text("%s", iface.c_str());
                    ImGui::TableSetColumnIndex(1); ImGui::Text("%s", ipv4s.count(iface) ? ipv4s[iface].c_str() : "-");
                    ImGui::TableSetColumnIndex(2); ImGui::Text("%lld", (long long)ns.rx.bytes);
                    ImGui::TableSetColumnIndex(3); ImGui::Text("%d", ns.rx.packets);
                    ImGui::TableSetColumnIndex(4); ImGui::Text("%d", ns.rx.errs);
                    ImGui::TableSetColumnIndex(5); ImGui::Text("%d", ns.rx.drop);
                    ImGui::TableSetColumnIndex(6); ImGui::Text("%d", ns.rx.fifo);
                    ImGui::TableSetColumnIndex(7); ImGui::Text("%d", ns.rx.frame);
                    ImGui::TableSetColumnIndex(8); ImGui::Text("%d", ns.rx.compressed);
                    ImGui::TableSetColumnIndex(9); ImGui::Text("%d", ns.rx.multicast);
                }
                ImGui::EndTable();
            }
            ImGui::Separator();
            ImGui::Text("RX Usage (0-2GB):");
            for (const auto& [iface, ns] : stats) {
                float rxGB = ns.rx.bytes / (1024.0f * 1024.0f * 1024.0f);
                float rxMB = ns.rx.bytes / (1024.0f * 1024.0f);
                float rxKB = ns.rx.bytes / 1024.0f;
                std::string rxStr;
                if (rxGB >= 1.0f) rxStr = std::to_string(rxGB) + " GB";
                else if (rxMB >= 1.0f) rxStr = std::to_string(rxMB) + " MB";
                else rxStr = std::to_string(rxKB) + " KB";
                float rxRatio = std::min((float)(ns.rx.bytes / (2.0f * 1024 * 1024 * 1024)), 1.0f);
                ImGui::Text("%s: %s", iface.c_str(), rxStr.c_str());
                ImGui::ProgressBar(rxRatio, ImVec2(-1.0f, 18.0f));
            }
            ImGui::EndTabItem();
        }
        if (ImGui::BeginTabItem("TX (Network Transmitter)")) {
            if (ImGui::BeginTable("TXTable", 10, ImGuiTableFlags_Borders | ImGuiTableFlags_RowBg | ImGuiTableFlags_Resizable)) {
                ImGui::TableSetupColumn("Interface");
                ImGui::TableSetupColumn("IPv4");
                ImGui::TableSetupColumn("Bytes");
                ImGui::TableSetupColumn("Packets");
                ImGui::TableSetupColumn("Errs");
                ImGui::TableSetupColumn("Drop");
                ImGui::TableSetupColumn("Fifo");
                ImGui::TableSetupColumn("Colls");
                ImGui::TableSetupColumn("Carrier");
                ImGui::TableSetupColumn("Compressed");
                ImGui::TableHeadersRow();
                for (const auto& [iface, ns] : stats) {
                    ImGui::TableNextRow();
                    ImGui::TableSetColumnIndex(0); ImGui::Text("%s", iface.c_str());
                    ImGui::TableSetColumnIndex(1); ImGui::Text("%s", ipv4s.count(iface) ? ipv4s[iface].c_str() : "-");
                    ImGui::TableSetColumnIndex(2); ImGui::Text("%lld", (long long)ns.tx.bytes);
                    ImGui::TableSetColumnIndex(3); ImGui::Text("%d", ns.tx.packets);
                    ImGui::TableSetColumnIndex(4); ImGui::Text("%d", ns.tx.errs);
                    ImGui::TableSetColumnIndex(5); ImGui::Text("%d", ns.tx.drop);
                    ImGui::TableSetColumnIndex(6); ImGui::Text("%d", ns.tx.fifo);
                    ImGui::TableSetColumnIndex(7); ImGui::Text("%d", ns.tx.colls);
                    ImGui::TableSetColumnIndex(8); ImGui::Text("%d", ns.tx.carrier);
                    ImGui::TableSetColumnIndex(9); ImGui::Text("%d", ns.tx.compressed);
                }
                ImGui::EndTable();
            }
            ImGui::Separator();
            ImGui::Text("TX Usage (0-2GB):");
            for (const auto& [iface, ns] : stats) {
                float txGB = ns.tx.bytes / (1024.0f * 1024.0f * 1024.0f);
                float txMB = ns.tx.bytes / (1024.0f * 1024.0f);
                float txKB = ns.tx.bytes / 1024.0f;
                std::string txStr;
                if (txGB >= 1.0f) txStr = std::to_string(txGB) + " GB";
                else if (txMB >= 1.0f) txStr = std::to_string(txMB) + " MB";
                else txStr = std::to_string(txKB) + " KB";
                float txRatio = std::min((float)(ns.tx.bytes / (2.0f * 1024 * 1024 * 1024)), 1.0f);
                ImGui::Text("%s: %s", iface.c_str(), txStr.c_str());
                ImGui::ProgressBar(txRatio, ImVec2(-1.0f, 18.0f));
            }
            ImGui::EndTabItem();
        }
        ImGui::EndTabBar();
    }
    ImGui::End();
}

// === Main ===
int main(int, char **)
{
    if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_TIMER | SDL_INIT_GAMECONTROLLER) != 0)
    {
        printf("Error: %s\n", SDL_GetError());
        return -1;
    }

    const char *glsl_version = "#version 130";
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_FLAGS, 0);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_PROFILE_MASK, SDL_GL_CONTEXT_PROFILE_CORE);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MAJOR_VERSION, 3);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MINOR_VERSION, 0);
    SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER, 1);
    SDL_GL_SetAttribute(SDL_GL_DEPTH_SIZE, 24);
    SDL_GL_SetAttribute(SDL_GL_STENCIL_SIZE, 8);

    SDL_WindowFlags window_flags = (SDL_WindowFlags)(SDL_WINDOW_OPENGL | SDL_WINDOW_RESIZABLE | SDL_WINDOW_ALLOW_HIGHDPI);
    SDL_Window *window = SDL_CreateWindow("System Monitor", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, 1280, 720, window_flags);
    SDL_GLContext gl_context = SDL_GL_CreateContext(window);
    SDL_GL_MakeCurrent(window, gl_context);
    SDL_GL_SetSwapInterval(1);

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

    IMGUI_CHECKVERSION();
    ImGui::CreateContext();
    ImGuiIO &io = ImGui::GetIO();
    ImGui::StyleColorsDark();
    ImGui_ImplSDL2_InitForOpenGL(window, gl_context);
    ImGui_ImplOpenGL3_Init(glsl_version);

    ImVec4 clear_color = ImVec4(0.0f, 0.0f, 0.0f, 0.0f);

    bool done = false;
    while (!done)
    {
        SDL_Event event;
        while (SDL_PollEvent(&event))
        {
            ImGui_ImplSDL2_ProcessEvent(&event);
            if (event.type == SDL_QUIT)
                done = true;
            if (event.type == SDL_WINDOWEVENT && event.window.event == SDL_WINDOWEVENT_CLOSE &&
                event.window.windowID == SDL_GetWindowID(window))
                done = true;
        }

        ImGui_ImplOpenGL3_NewFrame();
        ImGui_ImplSDL2_NewFrame(window);
        ImGui::NewFrame();

        ImVec2 mainDisplay = io.DisplaySize;
        memoryProcessesWindow("== Memory and Processes ==",
                              ImVec2((mainDisplay.x / 2) - 20, (mainDisplay.y / 2) + 30),
                              ImVec2((mainDisplay.x / 2) + 10, 10));

        systemWindow("== System Monitor ==",
                     ImVec2((mainDisplay.x / 2) - 10, (mainDisplay.y / 2) + 30),
                     ImVec2(10, 10));

        networkWindow("== Network ==",
                      ImVec2(mainDisplay.x - 20, (mainDisplay.y / 2) - 60),
                      ImVec2(10, (mainDisplay.y / 2) + 50));

        ImGui::Render();
        glViewport(0, 0, (int)io.DisplaySize.x, (int)io.DisplaySize.y);
        glClearColor(clear_color.x, clear_color.y, clear_color.z, clear_color.w);
        glClear(GL_COLOR_BUFFER_BIT);
        ImGui_ImplOpenGL3_RenderDrawData(ImGui::GetDrawData());
        SDL_GL_SwapWindow(window);
    }

    ImGui_ImplOpenGL3_Shutdown();
    ImGui_ImplSDL2_Shutdown();
    ImGui::DestroyContext();
    SDL_GL_DeleteContext(gl_context);
    SDL_DestroyWindow(window);
    SDL_Quit();
    return 0;
}