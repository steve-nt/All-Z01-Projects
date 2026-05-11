- Created tables which will be used for all next parts 
Notes:

- FOR 
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
ON Followers Table
If following_id user has is_public = TRUE: auto-set status = 'accepted' on insert
If following_id user has is_public = FALSE: set status = 'pending' on insert

- CASCADE for most relationships: when a user/group/post is deleted, related data is cleaned up automatically.


ALL REQUIRED TABLES:

1. **Authentication**
- Users table: email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me
- Sessions table: for session/cookie management
2. **Followers** 
- Followers table: follower_id, following_id, status (pending/accepted/declined)
- Supports follow requests and auto-accept for public profiles
3. **Profile** 
- Users table: all profile fields
- is_public field: public/private profile toggle
- Profile data queryable from Users + Posts + Followers tables
4. **Posts** 
- Posts table: content, privacy (public/almost_private/private)
- Posts_Images table: images/GIFs for posts
- Post_Visibility table: for private posts (selected followers)
- Comments table: comments on posts
- Comments_Images table: images/GIFs in comments
5. **Groups** 
- Groups table: title (group_name), description, creator
- Group_Members table: membership tracking
- Group_Invitations table: invitations with accept/decline
- Group_Join_Requests table: join requests with accept/decline
- Group_Posts table: posts within groups
- Group_Posts_Images table: images for group posts
- Group_Comments table: comments on group posts
- Group_Comments_Images table: images in group comments
- Group_Events table: events with title, description, event_datetime
- Group_Event_Responses table: going/not_going responses
6. **Chat** 
- Messages table: supports both private (recipient_id) and group (group_id) messages
- CHECK constraint ensures either recipient_id OR group_id is set
- TEXT content supports emojis
7. **Notifications** 
- Notifications table: type field for different notification types
- related_* fields: links to related entities (user, group, post, comment, event)
- Supports all required notification types

**20 tables covering**:
*Core*: Users, Sessions
*Social*: Followers
*Content*: Posts, Posts_Images, Post_Visibility, Comments, Comments_Images
*Groups*: Groups, Group_Members, Group_Invitations, Group_Join_Requests, Group_Posts, Group_Posts_Images, Group_Comments, Group_Comments_Images, Group_Events, Group_Event_Responses
*Communication*: Messages
*System*: Notifications