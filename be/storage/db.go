package storage

import (
	"be/constant"
	"be/logger"
	"be/model"
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

var log = logger.Logger
var db *sql.DB

// StartDB initializes the connection to a database (local or RDS) and creates the "requests" table if it doesn't exist.
func StartDB() error {
	// Get database credentials from environment variables
	dbUser := constant.APPCONFIG.DBUser
	dbPass := constant.APPCONFIG.DBPass
	dbHost := constant.APPCONFIG.DBHost
	dbName := constant.APPCONFIG.DBName

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		return fmt.Errorf("database credentials are missing in environment variables")
	}

	// URL-encode username and password
	encodedUser := url.QueryEscape(dbUser)
	encodedPass := url.QueryEscape(dbPass)

	log.Info("Connecting to database at ", dbHost)

	// Construct the DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", encodedUser, encodedPass, dbHost, dbName)

	// Open a connection to the database
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Errorf("error opening database connection: %v", err)
		return fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool configurations
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Errorf("error connecting to the database: %v", err)
		return fmt.Errorf("error connecting to the database: %w", err)
	}
	log.Info("Successfully connected to the database!")

	// Create the table if it doesn't exist
	if err := createTable(); err != nil {
		log.Error("error creating table: ", err)
		return err
	}
	log.Info("Successfully created or ensured existence of table 'requests'!")

	return nil
}

// createTable ensures the requests table exists
func createTable() error {
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS requests (
                RequestID VARCHAR(255) PRIMARY KEY,
                RequestedURL TEXT,
                IPAddress VARCHAR(45),
                RequestTime TIMESTAMP,
                AudioLink TEXT,
                LLMResponseLink TEXT,
                SpeechMarksLink TEXT,
                CleanedTextLink TEXT,
                NumberOfWords INT,
                Feedback INT DEFAULT 0,
                FeedbackText TEXT,
                TotalLatency DOUBLE PRECISION,
                SummaryLatency DOUBLE PRECISION,
                TTSLatency DOUBLE PRECISION,
                TTSPlayedDuration DOUBLE PRECISION,
                NavigatedToClean BOOLEAN
        );
        `
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	log.Info("Table 'requests' created or already exists.")
	return nil
}

// InsertRequest inserts a new request into the database
func InsertRequest(ctx context.Context, data *model.DBRow) error {
	stmt, err := db.Prepare(`
            INSERT INTO requests(
                RequestID, RequestedURL, IPAddress, RequestTime, 
                AudioLink, LLMResponseLink, SpeechMarksLink, CleanedTextLink, 
                NumberOfWords, Feedback, FeedbackText, TotalLatency, 
                SummaryLatency, TTSLatency, TTSPlayedDuration, NavigatedToClean
            )
            VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        `)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		data.RequestID,
		data.RequestedURL,
		data.IPAddress,
		data.RequestTime,
		data.AudioLink,
		data.LLMResponseLink,
		data.SpeechMarksLink,
		data.CleanedTextLink,
		data.NumberOfWords,
		data.Feedback,
		data.FeedbackText,
		data.TotalLatency,
		data.SummaryLatency,
		data.TTSLatency,
		data.TTSPlayedDuration,
		data.NavigatedToClean,
	)
	if err != nil {
		return fmt.Errorf("error executing insert statement for RequestID %s: %w", data.RequestID, err)
	}
	return nil
}

func UpdateFeedback(feedback model.Feedback) error {
	// Prepare the update statement
	stmt, err := db.Prepare(`
                UPDATE requests 
                SET Feedback = $1,
                    FeedbackText = $2,
                    TTSPlayedDuration = $3,
                    NavigatedToClean = $4
                WHERE RequestID = $5
        `)
	if err != nil {
		return fmt.Errorf("error preparing update statement: %w", err)
	}
	defer stmt.Close()

	// Execute the update
	result, err := stmt.Exec(
		feedback.Feedback,
		feedback.FeedbackText,
		feedback.TTSPlayedDuration,
		feedback.NavigatedToClean,
		feedback.RequestID,
	)
	if err != nil {
		return fmt.Errorf("error executing update statement for RequestID %s: %w", feedback.RequestID, err)
	}

	// Check if the row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no record found with RequestID: %s", feedback.RequestID)
	}
	return nil
}
