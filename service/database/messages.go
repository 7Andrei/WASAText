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

func (db *appdbimpl) ForwardMessage(messageId int, messageReceiver int, messageForwarded int) error {
	var oldMessage Message
	err := db.c.QueryRow("SELECT * FROM messages WHERE id = ?", messageId).Scan(&oldMessage.Id, &oldMessage.Content, &oldMessage.Photo, &oldMessage.Sender, &oldMessage.Receiver, &oldMessage.Forwarded, &oldMessage.TimeStamp)
	if err != nil {
		fmt.Println("Error forwarding message(ForwardMessage chats.go)\n", err)
		return err
	}
	oldMessage.Receiver = messageReceiver
	oldMessage.Forwarded = oldMessage.Sender
	oldMessage.Sender = messageForwarded
	_, err = db.c.Exec("INSERT INTO messages (content, photo, sender, receiver, forwarded) VALUES (?, ?, ?, ?, ?)", oldMessage.Content, oldMessage.Photo, oldMessage.Sender, oldMessage.Receiver, oldMessage.Forwarded)
	if err != nil {
		fmt.Println("Error forwarding message(ForwardMessage chats.go)\n", err)
		return err
	}
	return nil
}
