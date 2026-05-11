CREATE TABLE IF NOT EXISTS "notifications" (
	"id"	INTEGER NOT NULL UNIQUE,
    "user_id"	INTEGER NOT NULL,
    "actor_id" INTEGER NOT NULL,
	"type"	TEXT NOT NULL,
    "post_id" INTEGER NOT NULL,
    "comment_id"	INTEGER,
	"timestamp"	TEXT NOT NULL,
    "read" BOOLEAN NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
    FOREIGN KEY("post_id") REFERENCES "posts"("id"),
    FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
    FOREIGN KEY("user_id") REFERENCES "users"("id")
    FOREIGN KEY("actor_id") REFERENCES "users"("id")
);