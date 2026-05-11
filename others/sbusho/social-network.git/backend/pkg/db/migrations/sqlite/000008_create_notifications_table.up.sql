CREATE TABLE IF NOT EXISTS Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,   -- the user who will receive the notification
    type TEXT NOT NULL, -- follow_request, group_invitation, group_join_request, group_event_response, message
    related_user_id INTEGER, -- the user who is related to the notification
    related_group_id INTEGER, -- the group who is related to the notification
    related_post_id INTEGER, -- the post who is related to the notification
    related_comment_id INTEGER, -- the comment who is related to the notification
    related_event_id INTEGER, -- the event who is related to the notification
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (related_user_id) REFERENCES Users(user_id) ON DELETE SET NULL,
    FOREIGN KEY (related_group_id) REFERENCES Groups(id) ON DELETE SET NULL,
    FOREIGN KEY (related_post_id) REFERENCES Posts(post_id) ON DELETE SET NULL,
    FOREIGN KEY (related_comment_id) REFERENCES Comments(id) ON DELETE SET NULL,
    FOREIGN KEY (related_event_id) REFERENCES Group_Events(id) ON DELETE SET NULL
);