package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

// Email represents an email message
type Email struct {
	ID         string
	Subject    string
	BodyText   string
	BodyHTML   string
	FromEmail  string
	ToEmail    string
	CCEmail    string
	Date       int64
	Labels     string
	AISummary  string
	Status     string
	ThreadID   string
	AICategory string
}

// EmailDB represents the SQLite database connection
type EmailDB struct {
	db *sql.DB
}

// getHomeDir returns the home directory path for the current platform
func getHomeDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		home := os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
		if home == "" {
			return "", fmt.Errorf("unable to determine home directory")
		}
		return home, nil
	default:
		home := os.Getenv("HOME")
		if home == "" {
			return "", fmt.Errorf("unable to determine home directory")
		}
		return home, nil
	}
}

// NewEmailDB creates a new EmailDB instance
func NewEmailDB() (*EmailDB, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(homeDir, ".terminalmail", "terminalmail.db")

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open database with FTS5 support
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Set pragmas for performance
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=1000000",
		"PRAGMA foreign_keys=ON",
		"PRAGMA temp_store=MEMORY",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, err
		}
	}

	// Verify FTS5 support
	var fts5Support bool
	err = db.QueryRow("SELECT sqlite_compileoption_used('ENABLE_FTS5')").Scan(&fts5Support)
	if err != nil {
		return nil, err
	}

	if !fts5Support {
		return nil, fmt.Errorf("SQLite FTS5 support is not enabled. Build with -tags fts5")
	}
	// Create schema
	schemaPath := "internal/database/schema.sql"
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}
	return &EmailDB{db: db}, nil
}

// SaveEmail saves an email to the database
func (e *EmailDB) SaveEmail(email Email) error {
	// Validate required fields
	if email.ID == "" {
		return fmt.Errorf("email ID is required")
	}

	// Start transaction
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update email
	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO emails (
			id, subject, body_text, body_html, from_email, to_email, cc_email, 
			date, labels, ai_summary, status, thread_id, ai_category
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		email.ID, email.Subject, email.BodyText, email.BodyHTML, email.FromEmail,
		email.ToEmail, email.CCEmail, email.Date, email.Labels, email.AISummary,
		email.Status, email.ThreadID, email.AICategory,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// GetEmails retrieves emails with pagination
func (e *EmailDB) GetEmails(limit, offset int) ([]Email, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT id, subject, body_text, body_html, from_email, to_email, cc_email, 
		       date, labels, ai_summary, status, thread_id, ai_category
		FROM emails
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := e.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		err := rows.Scan(
			&email.ID, &email.Subject, &email.BodyText, &email.BodyHTML,
			&email.FromEmail, &email.ToEmail, &email.CCEmail, &email.Date,
			&email.Labels, &email.AISummary, &email.Status, &email.ThreadID,
			&email.AICategory,
		)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

// SearchEmails performs full-text search on emails
func (e *EmailDB) SearchEmails(query string, limit, offset int) ([]Email, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	sqlQuery := `
		SELECT e.id, e.subject, e.body_text, e.body_html, e.from_email, e.to_email, e.cc_email, 
		       e.date, e.labels, e.ai_summary, e.status, e.thread_id, e.ai_category
		FROM emails e
		JOIN email_search es ON e.id = es.id
		WHERE email_search MATCH ?
		ORDER BY rank
		LIMIT ? OFFSET ?
	`

	rows, err := e.db.Query(sqlQuery, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		err := rows.Scan(
			&email.ID, &email.Subject, &email.BodyText, &email.BodyHTML,
			&email.FromEmail, &email.ToEmail, &email.CCEmail, &email.Date,
			&email.Labels, &email.AISummary, &email.Status, &email.ThreadID,
			&email.AICategory,
		)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

// Close closes the database connection
func (e *EmailDB) Close() error {
	return e.db.Close()
}
