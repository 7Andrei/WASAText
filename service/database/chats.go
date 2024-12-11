package database

import (
	"fmt"
)

func (db *appdbimpl) GetChat(chatId int) (Chat, error) {
	var chat Chat
	err := db.c.QueryRow("SELECT id, name, photo, type FROM chats WHERE id=?", chatId).Scan(&chat.Id, &chat.Name, &chat.Photo, &chat.ChatType)
	if err != nil {
		fmt.Println("Error getting user data. ", err)
		return chat, err
	}
	return chat, nil
}

/*
func (db *appdbimpl) CreateChat(chat Chat) (int, error) {

	_, err := db.c.Exec("INSERT INTO chats (name) VALUES (?)", userName)
	if err != nil {
		fmt.Println("Error creating 1 user. ", err)
		return userId, err
	}
	err = db.c.QueryRow("SELECT id FROM users WHERE name=?", userName).Scan(&userId)
	if err != nil {
		fmt.Println("Error creating 2 user. ", err)
		return userId, err
	}
	//fmt.Println("User ID:", userId)
	return userId, nil
}
*/
