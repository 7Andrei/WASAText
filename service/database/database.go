/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	InsertUser() error
	GetUserName() (string, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		userStmt := `CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY, name VARCHAR NOT NULL , photo BLOB, identifier VARCHAR NOT NULL);`
		_, err = db.Exec(userStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure users: %w", err)
		}

		messagesStmt := `CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY,
												content TEXT, photo BLOB,
												sender INT NOT NULL,
												receiver INT NOT NULL,
												forwarded INT,
												sentTime TEXT DEFAULT CURRENT_TIMESTAMP,
												FOREIGN KEY(sender) references users(id),
												FOREIGN KEY(receiver) references chats(id),
												FOREIGN KEY(forwarded) references chats(id)
												);`
		_, err = db.Exec(messagesStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messages: %w", err)
		}

		chatsStmt := `CREATE TABLE IF NOT EXISTS chats (id INTEGER NOT NULL PRIMARY KEY, name VARCHAR NOT NULL , photo BLOB, type VARCHAR NOT NULL);`
		_, err = db.Exec(chatsStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure chats: %w", err)
		}

		reactionsStmt := `CREATE TABLE IF NOT EXISTS reactions (id INTEGER NOT NULL PRIMARY KEY, reaction VARCHAR NOT NULL, userId INT NOT NULL, messageId INT NOT NULL,
						  FOREIGN KEY(userId) references users(id), FOREIGN KEY(messageId) references messages(id));`
		_, err = db.Exec(reactionsStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure reactions: %w", err)
		}

		user_messagesStmt := `CREATE TABLE IF NOT EXISTS user_messages (userId INTEGER NOT NULL, messageId INT NOT NULL,
							  FOREIGN KEY (userId) references users(id), FOREIGN KEY(messageId) references messages(id));`
		_, err = db.Exec(user_messagesStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user_messages: %w", err)
		}

		user_chatsStmt := `CREATE TABLE IF NOT EXISTS user_chats (userId INTEGER NOT NULL, chatId INT NOT NULL,
							  FOREIGN KEY (userId) references users(id), FOREIGN KEY(chatId) references chats(id));`
		_, err = db.Exec(user_chatsStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user_chats: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
