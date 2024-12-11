package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) Login(userName string) (int, error) {
	var userId int
	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
	if err == sql.ErrNoRows {
		_, err := db.c.Exec("INSERT INTO users (name) VALUES (?)", userName)
		if err != nil {
			fmt.Println("Error creating 1 user. ", err)
			return userId, err
		}
		err = db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
		if err != nil {
			fmt.Println("Error creating 2 user. ", err)
			return userId, err
		}
	}
	//fmt.Println("User ID:", userId)
	return userId, nil
}

func (db *appdbimpl) GetUser(userId int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, name, photo FROM users WHERE id=?", userId).Scan(&user.Id, &user.Username, &user.Photo)
	if err != nil {
		fmt.Println("Error getting user data. ", err)
		return user, err
	}
	return user, nil
}

func (db *appdbimpl) SetUsername(userId int, newName string) error {
	_, err := db.c.Exec("UPDATE users SET name=? WHERE id=?", newName, userId)
	if err != nil {
		fmt.Println("Error updating username. ", err)
		return err
	}
	return nil
}

func (db *appdbimpl) SetUserPhoto(userId int, newPhoto []byte) error {
	_, err := db.c.Exec("UPDATE users SET photo=? WHERE id=?", newPhoto, userId)
	if err != nil {
		fmt.Println("Error updating user photo. ", err)
		return err
	}
	return nil
}
