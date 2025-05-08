package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) Login(userName string) (int, error) {
	var userId int
	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := db.c.Exec("INSERT INTO users (name) VALUES (?)", userName)
		if err != nil {
			fmt.Println("Error creating user. Login users.go", err)
			return userId, err
		}
		err = db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
		if err != nil {
			fmt.Println("Error fetching user. Login users.go ", err)
			return userId, err
		}
	}
	return userId, nil
}

func (db *appdbimpl) GetUser(userId int) (User, bool, error) {
	var user User
	err := db.c.QueryRow("SELECT id, name, photo FROM users WHERE id=?", userId).Scan(&user.Id, &user.Username, &user.Photo)
	if err != nil {
		fmt.Println("Error getting user data. GetUser users.go", err)
		return user, false, err
	}
	return user, true, nil
}

func (db *appdbimpl) SetUsername(userId int, newName string) error {
	_, err := db.c.Exec("UPDATE users SET name=? WHERE id=?", newName, userId)
	if err != nil {
		fmt.Println("Error updating username. SetUsername users.go", err)
		return err
	}
	return nil
}

func (db *appdbimpl) SetUserPhoto(userId int, newPhoto []byte) error {
	_, err := db.c.Exec("UPDATE users SET photo=? WHERE id=?", newPhoto, userId)
	if err != nil {
		fmt.Println("Error updating user photo. SetUserPhoto users.go", err)
		return err
	}
	return nil
}

func (db *appdbimpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT id, name, photo FROM users")
	if err != nil {
		fmt.Println("Error getting all users. GetAllUsers users.go. ", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			fmt.Println("Error getting user data. GetAllUsers users.go. ", err)
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		fmt.Println("Error fetching user data. GetAllUsers users.go. ", err)
		return nil, err
	}
	return users, nil
}
