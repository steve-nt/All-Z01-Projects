-- SQLite 3.35.0+ supports DROP COLUMN; older versions may need manual migration
ALTER TABLE Users DROP COLUMN relationship_status;
ALTER TABLE Users DROP COLUMN hobbies;
