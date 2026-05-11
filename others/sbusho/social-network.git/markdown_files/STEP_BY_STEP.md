# Step-by-Step Learning Guide: Setting Up Your First Database with Migrations

Follow these steps to understand and test your database setup.

## Prerequisites

Make sure you have Go installed (version 1.21 or later).

## Step 1: Install Dependencies

Open a terminal in the `backend` directory and run:

```bash
go mod tidy
```

This will download all the required packages:
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/golang-migrate/migrate/v4` - Migration tool

**What this does**: Downloads the packages specified in `go.mod` so you can use them in your code.

## Step 2: Understand the Migration Files

Look at your migration files in `backend/pkg/db/migrations/sqlite/`:

1. **Open** `000001_create_users_table.up.sql`
   - This creates the Users and Sessions tables
   - Notice the SQL syntax: `CREATE TABLE IF NOT EXISTS`
   - The `IF NOT EXISTS` prevents errors if the table already exists

2. **Open** `000001_create_users_table.down.sql`
   - This drops (deletes) the tables
   - Used to rollback the migration if needed

3. **Look at** other migration files
   - Each one builds on the previous ones
   - They're numbered sequentially (000001, 000002, etc.)

**Learning Point**: Migrations are applied in numerical order. The number determines the sequence.

## Step 3: Read the Code (Understanding Phase)

Open `backend/pkg/db/sqlite/sqlite.go` and read through it:

1. **Read the comments** - They explain what each section does
2. **Follow the flow**:
   - `InitDB()` is the main function
   - It opens a database connection
   - Then calls `runMigrations()` to apply migrations
   - Finally stores the connection in a global variable

3. **Key concepts to understand**:
   - `sql.Open()` - Opens database connection
   - `sqlite3.WithInstance()` - Creates migration driver
   - `migrate.NewWithDatabaseInstance()` - Sets up migration system
   - `m.Up()` - Applies migrations

**Don't worry if you don't understand everything yet!** The next steps will help.

## Step 4: Test the Setup

Create a simple test file to see it in action:

1. **Run the example**:
   ```bash
   cd backend
   go run example_usage.go
   ```

2. **What should happen**:
   - Creates `data/social_network.db` file
   - Applies all 9 migrations
   - Creates all tables
   - Prints success messages

3. **If you see errors**:
   - Check that you're in the `backend` directory
   - Make sure `go mod tidy` ran successfully
   - Check the error message for clues

## Step 5: Inspect the Database

After running the example, you can inspect what was created:

1. **Install a SQLite browser** (optional but helpful):
   - [DB Browser for SQLite](https://sqlitebrowser.org/) (free, GUI)
   - Or use command line: `sqlite3 data/social_network.db`

2. **Open** `data/social_network.db` in the browser

3. **What you'll see**:
   - All the tables from your migrations
   - A special table called `schema_migrations` (created by golang-migrate)
   - The `schema_migrations` table tracks which migrations have been applied

4. **Explore the tables**:
   - Click on each table to see its structure
   - Notice the columns match what's in your migration files
   - The foreign keys and constraints are in place

## Step 6: Understand What Happened

After running the code, think about what happened:

1. **Database file created**: `data/social_network.db` didn't exist, so it was created
2. **Migrations applied**: All 9 migration files were read and executed
3. **Tables created**: All tables defined in migrations now exist
4. **Version tracked**: The `schema_migrations` table records that all migrations were applied

## Step 7: Experiment (Optional)

Try these experiments to deepen your understanding:

### Experiment 1: Run it again
```bash
go run example_usage.go
```
- Notice it says "no change" - this is normal!
- Migrations only run if they haven't been applied yet

### Experiment 2: Delete the database
```bash
rm data/social_network.db
go run example_usage.go
```
- The database is recreated
- All migrations are applied again
- This shows migrations are repeatable

### Experiment 3: Check migration version
Look at the `schema_migrations` table:
```sql
SELECT * FROM schema_migrations;
```
- You'll see the version number (should be 9, since you have 9 migrations)
- This is how golang-migrate tracks what's been applied

## Step 8: Integrate into Your Server

When you're ready to use this in your actual server:

1. **In `server.go`**, add this at the start of `main()`:
   ```go
   import "social-network/backend/pkg/db/sqlite"
   
   func main() {
       // Initialize database
       if err := sqlite.InitDB("data/social_network.db"); err != nil {
           log.Fatalf("Failed to initialize database: %v", err)
       }
       defer sqlite.CloseDB() // Close when server shuts down
       
       // ... rest of your server code
   }
   ```

2. **Use the database** in other parts of your code:
   ```go
   db, err := sqlite.GetDB()
   if err != nil {
       // handle error
   }
   // Now use db to run queries
   ```

## Common Questions

### Q: What if I need to add a new table later?
**A**: Create a new migration file (e.g., `000010_create_new_table.up.sql`) and it will be applied automatically.

### Q: Can I modify an existing migration?
**A**: Only if it hasn't been applied yet. Once applied, create a new migration instead.

### Q: What if a migration fails?
**A**: The database is marked as "dirty". You'll need to manually fix the issue, then continue.

### Q: How do I rollback a migration?
**A**: Use `m.Down()` instead of `m.Up()`, but this is advanced and rarely needed.

## Next Steps

Now that you understand migrations:
1. ✅ You've set up the database connection
2. ✅ You've implemented the migration system
3. ✅ You understand how it works

**You've completed Item 5!** 🎉

You can now move on to:
- Item 7: Build user registration system
- Item 8: Implement login/logout functionality
- Item 9: Create session and cookie management system

All of these will use the database connection you just set up!

