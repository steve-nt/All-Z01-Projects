Phase 1: Foundation (start here)
WebSocket server setup + connection hub
Why first: Required for everything else
What: Upgrade HTTP to WebSocket, create hub to manage connections, authenticate connections using existing session system
Benefit: Establishes the infrastructure

Phase 2: Quick win
Real-time notifications
Why second: High value, low complexity
What: Push notifications via WebSocket when they’re created
Benefit: Immediate value; notifications are already created in your codebase, just need to push them
Integration: Hook into existing notification creation in groups/handlers.go and profile/profileHandler.go

Phase 3: Core messaging
Private messaging
Why third: Simpler than group chat (1-to-1)
What: Send/receive private messages, message history API
Benefit: Core feature, validates the foundation

Phase 4: Extended features
Group chat
Why fourth: Builds on private messaging patterns
What: Broadcast to group members, group message history
Benefit: Extends messaging to groups

Phase 5: Polish
Emoji support
Why last: Can be added anytime
What: Ensure Unicode/emoji handling (likely already works)
Benefit: Nice-to-have enhancement
Why this order
Incremental: Each step builds on the previous
Testable: You can test after each phase
Early value: Notifications provide immediate benefit
Progressive complexity: Start simple, add complexity
Low risk: Foundation first reduces rework