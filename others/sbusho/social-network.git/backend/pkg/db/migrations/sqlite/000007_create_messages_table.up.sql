CREATE TABLE IF NOT EXISTS Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER,
    group_id INTEGER,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (sender_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES Users(user_id) ON DELETE SET NULL,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE SET NULL,
    CHECK (
        (recipient_id IS NOT NULL AND group_id IS NULL) OR 
        (recipient_id IS NULL AND group_id IS NOT NULL)
    )
);
