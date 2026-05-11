CREATE TABLE IF NOT EXISTS Group_Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Group_Posts_Images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_post_id INTEGER NOT NULL,
    image_path TEXT NOT NULL,
    image_type TEXT NOT NULL CHECK (image_type IN ('JPEG', 'PNG', 'GIF')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_post_id) REFERENCES Group_Posts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Group_Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_post_id) REFERENCES Group_Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Group_Comments_Images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_comment_id INTEGER NOT NULL,
    image_path TEXT NOT NULL,
    image_type TEXT NOT NULL CHECK (image_type IN ('JPEG', 'PNG', 'GIF')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_comment_id) REFERENCES Group_Comments(id) ON DELETE CASCADE
);

