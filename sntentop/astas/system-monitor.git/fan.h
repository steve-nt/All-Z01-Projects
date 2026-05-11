#pragma once
#include <deque>

// Struct holding fan information
struct FanInfo {
    bool active;     // Whether the fan is running or not
    int speedRPM;    // RPM of the fan
    int level;       // Optional fan level or PWM value
};

// Retrieves current fan information from sysfs
FanInfo getFanInfo();

// Renders the fan monitoring tab in ImGui
void renderFanTab();
