package database

import (
	"fmt"
)

func (db *appdbimpl) SendMessage(messageContent string, messagePhoto []byte, messageSender int, messageReceiver int, messageForwarded int) (int, error) {
	var messageId int
	err := db.c.QueryRow("INSERT INTO messages (content, photo, sender, receiver, forwarded) VALUES (?, ?, ?, ?, ?) RETURNING id", messageContent, messagePhoto, messageSender, messageReceiver, messageForwarded).Scan(&messageId)
	if err != nil {
		fmt.Println("Error sending message(SendMessage chats.go)\n", err)
		return messageId, err
	}
	fmt.Println("Message sent:", messageContent, messagePhoto, messageSender, messageReceiver)

	return messageId, nil
}
