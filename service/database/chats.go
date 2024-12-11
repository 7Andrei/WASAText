package database

import (
	"fmt"
)

func (db *appdbimpl) GetChat(chatId int) (Chat, error) {
	var chat Chat
	err := db.c.QueryRow("SELECT id, name, photo, type FROM chats WHERE id=?", chatId).Scan(&chat.Id, &chat.Name, &chat.Photo, &chat.ChatType)
	if err != nil {
		fmt.Println("Error getting chat data. ", err)
		return chat, err
	}
	return chat, nil
}

func (db *appdbimpl) CreateChat(chatName string, chatPhoto []byte, chatType string) error {

	_, err := db.c.Exec("INSERT INTO chats (name, photo, type) VALUES (?, ?, ?)", chatName, chatPhoto, chatType)
	if err != nil {
		fmt.Println("Error creating chat. ", err)
		return err
	}

	fmt.Println("Chat created:", chatName, chatPhoto, chatType)
	return nil
}
