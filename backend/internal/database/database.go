package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"doodle-clone/internal/config"
)

var Pool *pgxpool.Pool

// Connect establishes a connection to the database
func Connect() error {
	ctx := context.Background()
	connString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
	)

	var err error
	Pool, err = pgxpool.New(ctx, connString)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

// Close closes the database connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("Database connection closed")
	}
}

// RunMigrations executes all database migrations
func RunMigrations() error {
	ctx := context.Background()

	// Create tables in the correct order (respecting foreign keys)
	migrations := []string{
		createUsersTable(),
		createPollsTable(),
		createDateOptionsTable(),
		createVotesTable(),
		createCommentsTable(),
		createRefreshTokensTable(),
		createNotificationSettingsTable(),
		createNotificationsTable(),
	}

	for _, migration := range migrations {
		_, err := Pool.Exec(ctx, migration)
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Run additional migrations for existing databases
	// Add access_code column to polls if it doesn't exist
	var columnExists bool
	err := Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = 'polls'
			AND column_name = 'access_code'
		)`).Scan(&columnExists)
	if err == nil && !columnExists {
		log.Println("Adding access_code column to polls table...")
		// Add column as nullable first
		_, err = Pool.Exec(ctx, `ALTER TABLE polls ADD COLUMN IF NOT EXISTS access_code VARCHAR(20);`)
		if err != nil {
			log.Printf("Warning: failed to add access_code column: %v", err)
		} else {
			// Generate unique access codes for existing polls
			_, err = Pool.Exec(ctx, `
				UPDATE polls p1 SET access_code = (
					SELECT substr(md5(p1.id::text || random()::text), 1, 10)
					WHERE NOT EXISTS (
						SELECT 1 FROM polls p2 WHERE p2.access_code = substr(md5(p1.id::text || random()::text), 1, 10)
					)
					LIMIT 1
				);
			`)
			if err != nil {
				log.Printf("Warning: failed to generate access codes: %v", err)
			}
			// Make the column unique and not null
			_, err = Pool.Exec(ctx, `
				ALTER TABLE polls ALTER COLUMN access_code SET NOT NULL;
				ALTER TABLE polls ADD CONSTRAINT polls_access_code_key UNIQUE (access_code);
				CREATE INDEX IF NOT EXISTS idx_polls_access_code ON polls(access_code);
			`)
			if err != nil {
				log.Printf("Warning: failed to make access_code unique: %v", err)
			}
		}
	}

	// Add unique constraint to votes table if it doesn't exist
	var uniqueConstraintExists bool
	err = Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.table_constraints
			WHERE table_schema = 'public'
			AND table_name = 'votes'
			AND constraint_name = 'votes_poll_user_date_unique'
		)`).Scan(&uniqueConstraintExists)
	if err == nil && !uniqueConstraintExists {
		log.Println("Adding unique constraint to votes table...")
		_, err = Pool.Exec(ctx, `
			ALTER TABLE votes
			ADD CONSTRAINT votes_poll_user_date_unique
			UNIQUE (poll_id, date_option_id, user_id);
		`)
		if err != nil {
			log.Printf("Warning: failed to add votes unique constraint: %v", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Helper function to check if table exists
func TableExists(ctx context.Context, tableName string) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)`
	err := Pool.QueryRow(ctx, query, tableName).Scan(&exists)
	return err == nil && exists
}

// Helper function to drop all tables (for testing only)
func DropAllTables() error {
	ctx := context.Background()
	tables := []string{
		"refresh_tokens",
		"comments",
		"votes",
		"date_options",
		"polls",
		"users",
	}

	for _, table := range tables {
		_, err := Pool.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	log.Println("All tables dropped")
	return nil
}

// Helper function to get a context with timeout
func GetContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	t := 5 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return context.WithTimeout(context.Background(), t)
}
