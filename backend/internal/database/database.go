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
	}

	for _, migration := range migrations {
		_, err := Pool.Exec(ctx, migration)
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
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
