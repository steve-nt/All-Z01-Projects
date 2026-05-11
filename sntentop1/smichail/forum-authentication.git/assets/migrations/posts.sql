PRAGMA foreign_keys = ON;
-- drop table if exists activities;

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id BLOB PRIMARY KEY ,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_uuid TEXT NOT NULL,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_categories (
    post_id BLOB NOT NULL ,
    category_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uuid blob not null,
    action text not null,
    post_id BLOB ,
    comment_id BLOB  ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
    FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE
    CHECK ((action IS 'editted') OR  (action IS 'liked') OR (action IS 'disliked')OR (action IS 'commented') OR (action is 'created'))
    CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
);
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uuid blob not null,
    action text not null,
    post_id BLOB  ,
    comment_id BLOB  ,
    seen BOOL NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
    FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE
    CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
    CHECK ((action IS 'liked your post') OR (action IS 'disliked your post') OR (action IS 'disliked your comment') OR(action IS 'liked your comment') OR (action IS 'commented on your post'))
);

CREATE TABLE IF NOT EXISTS users(
uuid BLOB PRIMARY KEY ,
mail TEXT NOT NULL,
username TEXT NOT NULL,
password TEXT,
TYPE TEXT NOT NULL,
createdAt DATETIME NOT NULL,
verified BOOLEAN NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS social_users(
id INTEGER PRIMARY KEY AUTOINCREMENT,
user_uuid BLOB NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
provider TEXT NOT NULL,
provider_user_id TEXT NOT NULL,
UNIQUE(provider, provider_user_id)
);


CREATE TABLE IF NOT EXISTS sessions(
uuid BLOB NOT NULL,
cookie BLOB NOT NULL,
createdAt DATETIME NOT NULL,
expiration DATETIME NOT NULL,
absoluteExpiration DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id BLOB PRIMARY KEY ,
    content TEXT NOT NULL,
    user_uuid TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uuid TEXT NOT NULL,
    post_id blob,
    comment_id INTEGER,
    like bool,
    dislike bool,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
    CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
    CHECK ((like IS NOT TRUE AND dislike IS TRUE) OR (like IS TRUE AND dislike IS NOT TRUE))
);


INSERT INTO categories(name ) values('Tech');
INSERT INTO categories(name ) values('Politics');
INSERT INTO categories(name ) values('Economy');
INSERT INTO categories(name ) values('Ecology');
