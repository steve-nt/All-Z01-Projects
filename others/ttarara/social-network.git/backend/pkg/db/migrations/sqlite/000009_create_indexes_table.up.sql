-- Indexes for Users
CREATE INDEX idx_users_email ON Users(email);
CREATE INDEX idx_users_is_public ON Users(is_public);

-- Indexes for Followers
CREATE INDEX idx_followers_follower_id ON Followers(follower_id);
CREATE INDEX idx_followers_following_id ON Followers(following_id);
CREATE INDEX idx_followers_status ON Followers(status);

-- Indexes for Posts
CREATE INDEX idx_posts_user_id ON Posts(user_id);
CREATE INDEX idx_posts_created_at ON Posts(created_at DESC);
CREATE INDEX idx_posts_privacy ON Posts(privacy);

-- Indexes for Post_Visibility
CREATE INDEX idx_post_visibility_post_id ON Post_Visibility(post_id);
CREATE INDEX idx_post_visibility_user_id ON Post_Visibility(user_id);

-- Indexes for Comments
CREATE INDEX idx_comments_post_id ON Comments(post_id);
CREATE INDEX idx_comments_user_id ON Comments(user_id);

-- Indexes for Groups
CREATE INDEX idx_groups_creator_id ON Groups(creator_id);

-- Indexes for Group_Members
CREATE INDEX idx_group_members_group_id ON Group_Members(group_id);
CREATE INDEX idx_group_members_user_id ON Group_Members(user_id);

-- Indexes for Messages
CREATE INDEX idx_messages_sender_id ON Messages(sender_id);
CREATE INDEX idx_messages_recipient_id ON Messages(recipient_id);
CREATE INDEX idx_messages_group_id ON Messages(group_id);
CREATE INDEX idx_messages_created_at ON Messages(created_at DESC);

-- Indexes for Notifications
CREATE INDEX idx_notifications_user_id ON Notifications(user_id);
CREATE INDEX idx_notifications_is_read ON Notifications(is_read);
CREATE INDEX idx_notifications_created_at ON Notifications(created_at DESC);

-- Indexes for Sessions
CREATE INDEX idx_sessions_user_id ON Sessions(user_id);
CREATE INDEX idx_sessions_cookie_value ON Sessions(cookie_value);