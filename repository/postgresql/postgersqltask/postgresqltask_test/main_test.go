package postgresqltask_test

import (
	"taskmaneger/repository/postgresql"
	"testing"
)

func setupTestDB(t *testing.T) *postgresql.DB {
	connectionString := "host=localhost port=5436 user=admin dbname=task_manager_test_db password=password123 sslmode=disable"

	db, err := postgresql.NewDB(connectionString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		t.Fatalf("Failed to auto-migrate database: %v", err)
	}

	return db
}

func teardownTestDB(t *testing.T, db *postgresql.DB) {
	if err := db.Conn.Exec("DELETE FROM tasks").Error; err != nil {
		t.Fatalf("Failed to clean up tasks table: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close database connection: %v", err)
	}
}
