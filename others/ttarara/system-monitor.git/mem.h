#ifndef MEM_H
#define MEM_H

struct MemoryStats {
    float totalRamMB;
    float usedRamMB;
    float totalSwapMB;
    float usedSwapMB;
};

MemoryStats getMemoryStats();

#endif // MEM_H
