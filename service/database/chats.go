package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) GetChat(chatId int) (Chat, error) {
	var chat Chat
	err := db.c.QueryRow("SELECT id, name, photo, type FROM chats WHERE id=?", chatId).Scan(&chat.Id, &chat.Name, &chat.Photo, &chat.ChatType)
	if err == sql.ErrNoRows {
		fmt.Println("Chat not found(DB). ", err)
		return chat, err
	}
	if err != nil {
		fmt.Println("Error getting chat data(DB). ", err)
		return chat, err
	}
	return chat, nil
}

func (db *appdbimpl) CreateChat(chatName string, chatPhoto []byte, chatType string) (int, error) {

	var chatId int
	// _, err := db.c.Exec("INSERT INTO chats (name, photo, type) VALUES (?, ?, ?)", chatName, chatPhoto, chatType)
	err := db.c.QueryRow("INSERT INTO chats (name, photo, type) VALUES (?, ?, ?) RETURNING id", chatName, chatPhoto, chatType).Scan(&chatId)
	if err != nil {
		fmt.Println("Error creating chat(DB). ", err)
		return chatId, err
	}
	fmt.Println("Chat created:", chatName, chatPhoto, chatType)

	return chatId, nil
}

func (db *appdbimpl) GetAllChats() ([]Chat, error) {
	var chats []Chat
	rows, err := db.c.Query("SELECT id, name, photo, type FROM chats")
	if err != nil {
		fmt.Println("Error fetching all chats(DB). ", err)
		return chats, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.Id, &chat.Name, &chat.Photo, &chat.ChatType)
		if err != nil {
			fmt.Println("Error scanning chat data(DB). ", err)
			return chats, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (db *appdbimpl) AddParticipant(chatId int, participantId int) error {
	_, err := db.c.Exec("INSERT INTO user_chats (chatId, userId) VALUES (?, ?)", chatId, participantId)
	if err != nil {
		fmt.Println("Error adding partecipant to chat (AddPartecipant chats.go)\n", err)
		return err
	}
	return nil
}

func (db *appdbimpl) SendMessage(messageContent string, messagePhoto []byte, messageSender int, messageReceiver int) (int, error) {
	var messageId int
	err := db.c.QueryRow("INSERT INTO messages (content, photo, sender, receiver) VALUES (?, ?, ?, ?) RETURNING id", messageContent, messagePhoto, messageSender, messageReceiver).Scan(&messageId)
	if err != nil {
		fmt.Println("Error sending message(SendMessage chats.go)\n", err)
		return messageId, err
	}
	fmt.Println("Message sent:", messageContent, messagePhoto, messageSender, messageReceiver)

	return messageId, nil
}
