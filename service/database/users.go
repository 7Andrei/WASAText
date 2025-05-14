package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) Login(userName string) (int, error) {
	var userId int
	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := db.c.Exec("INSERT INTO users (name) VALUES (?)", userName)
		if err != nil {
			return userId, err
		}
		err = db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
		if err != nil {
			return userId, err
		}
	}
	return userId, nil
}

func (db *appdbimpl) GetUser(userId int) (User, bool, error) {
	var user User
	err := db.c.QueryRow("SELECT id, name, photo FROM users WHERE id=?", userId).Scan(&user.Id, &user.Username, &user.Photo)
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

func (db *appdbimpl) SetUsername(userId int, newName string) error {
	_, err := db.c.Exec("UPDATE users SET name=? WHERE id=?", newName, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetUserPhoto(userId int, newPhoto []byte) error {
	_, err := db.c.Exec("UPDATE users SET photo=? WHERE id=?", newPhoto, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT id, name, photo FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return users, nil
}
