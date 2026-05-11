# Learning Guide: Database Migrations with SQLite

## What Are Migrations?

**Migrations** are a way to manage changes to your database schema (structure) over time. Think of them like version control (Git) for your database.

### Why Use Migrations?

1. **Version Control**: Track all changes to your database structure
2. **Reproducibility**: Anyone can set up the same database structure
3. **Rollback**: Undo changes if something goes wrong
4. **Team Collaboration**: Multiple developers can work on the same database without conflicts
5. **Production Safety**: Apply changes in a controlled, tested way

## How Migrations Work

### Migration Files

Each migration has two files:
- **`.up.sql`**: Applies the migration (creates tables, adds columns, etc.)
- **`.down.sql`**: Rolls back the migration (drops tables, removes columns, etc.)

### Example Migration

**000001_create_users_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS Users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    -- ... more columns
);
```

**000001_create_users_table.down.sql:**
```sql
DROP TABLE IF EXISTS Users;
```

### Migration Naming Convention

Migrations are numbered sequentially:
- `000001_` - First migration
- `000002_` - Second migration
- `000003_` - Third migration
- etc.

The number determines the order in which migrations are applied.

## How golang-migrate Works

1. **Migration Tracking**: Creates a special table (`schema_migrations`) to track which migrations have been applied
2. **Version Checking**: Compares the current database version with available migration files
3. **Selective Application**: Only applies migrations that haven't been run yet
4. **Atomic Operations**: Each migration runs as a single transaction (all or nothing)

## Understanding the Code

### 1. Database Connection (`sql.Open`)

```go
db, err := sql.Open("sqlite3", "file:path?params")
```

- Opens a connection to SQLite database
- The connection string includes parameters:
  - `_foreign_keys=1`: Enables foreign key constraints
  - `_journal_mode=WAL`: Uses Write-Ahead Logging for better performance

### 2. Migration Driver (`sqlite3.WithInstance`)

```go
driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
```

- Creates a driver that golang-migrate uses to interact with SQLite
- Wraps the database connection with migration-specific functionality

### 3. Migration Source (`file://`)

```go
migrationsURL := "file:///absolute/path/to/migrations"
```

- Tells golang-migrate where to find migration files
- Uses file:// protocol (like http:// but for local files)

### 4. Applying Migrations (`m.Up()`)

```go
err = m.Up()
```

- Applies all pending migrations
- Returns `migrate.ErrNoChange` if all migrations are already applied
- Creates/updates tables based on `.up.sql` files

## Common Migration Scenarios

### Scenario 1: First Time Setup
- Database doesn't exist
- All migrations are applied in order
- All tables are created

### Scenario 2: Adding a New Migration
- Database already exists with some tables
- Only new migrations are applied
- Existing data is preserved

### Scenario 3: Migration Failure
- If a migration fails, the database is marked as "dirty"
- You need to manually fix the issue
- Then you can continue with migrations

## Best Practices

1. **Never Edit Old Migrations**: Once a migration is applied in production, don't change it. Create a new migration instead.

2. **Test Migrations**: Always test both up and down migrations before deploying.

3. **Backup Before Migrating**: In production, backup your database before running migrations.

4. **One Change Per Migration**: Keep migrations focused on one logical change.

5. **Use Transactions**: Migrations should be atomic (all or nothing).

## Troubleshooting

### Error: "no change"
- This means all migrations are already applied. This is normal and not an error.

### Error: "dirty database"
- A previous migration failed partway through
- You need to manually fix the database state
- Check the `schema_migrations` table

### Error: "file not found"
- Check that the migrations path is correct
- Use absolute paths, not relative paths
- Ensure the migrations directory exists

## Next Steps

1. **Run the code**: Initialize the database and see migrations in action
2. **Check the database**: Use a SQLite browser to see the created tables
3. **Add a new migration**: Create a new migration file and see it applied
4. **Experiment**: Try rolling back a migration (advanced)

## Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Go database/sql Tutorial](https://go.dev/doc/database/sql-intro)

