Backend Division (4 Parts)
**Sofia** Part 1: Database & Authentication Foundation
Responsibility: Database setup, user authentication, and session management

- Set up SQLite database connection and structure
- Create and implement migration system (all migration files)
- Build user registration system (handling all required and optional fields)
- Implement login/logout functionality
- Create session and cookie management system
- Password encryption with bcrypt
- Image upload handling for avatars

Key deliverables:

Database schema and migrations
User registration and login endpoints
Session middleware
Image storage system


**Georgia** Part 2: User Profiles & Following System
Responsibility: User profiles, follower relationships, and privacy settings

Create user profile endpoints (public/private profiles)
Implement follow request system (send, accept, decline)
Build follower/following management
Handle profile privacy toggles (public/private)
Display user activity and information
Manage follow request notifications

Key deliverables:

Profile API endpoints
Following/follower logic
Privacy system implementation
Profile data retrieval


**Charoula** Part 3: Posts & Groups
Responsibility: Post creation, groups, and group events

Create post system (with image/GIF support)
Implement post privacy levels (public, almost private, private)
Build commenting system on posts
Create group management (create, invite, request to join)
Implement group posts and comments
Build event creation system within groups
Handle event responses (going/not going)

Key deliverables:

Posts API endpoints
Groups API endpoints
Events system
Privacy filtering logic


**Andy** Part 4: Real-time Communication (WebSocket & Chat)
Responsibility: Private messaging, group chat, and notifications

Set up WebSocket connections
Implement private messaging system
Build group chat rooms
Create notification system for all types:

Follow requests
Group invitations
Group join requests
Event notifications


Handle emoji support in messages
Real-time message delivery

Key deliverables:

WebSocket server setup
Chat endpoints and logic
Notification system
Real-time connection management


Suggested Work Order

Start with Part 1 - Everyone needs the database and authentication working
Parts 2 and 3 can work in parallel once Part 1 is done
Part 4 starts last but can begin planning/setup earlier

Integration Points to Coordinate

Part 1 creates the user tables that Parts 2, 3, 4 will use
Part 2's following system affects Part 3's post privacy and Part 4's messaging permissions
Part 4 needs user data from Part 1, following data from Part 2, and group data from Part 3

Each person should create clear API documentation for their endpoints so others can integrate easily. Consider having regular sync meetings to ensure the parts work together smoothly.