-- Add optional profile fields for settings (age is derived from date_of_birth)
ALTER TABLE Users ADD COLUMN relationship_status TEXT;
ALTER TABLE Users ADD COLUMN hobbies TEXT;
