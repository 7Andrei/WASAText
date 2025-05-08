package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetChat(chatId int) (Chat, error) {
	var chat Chat
	err := db.c.QueryRow("SELECT id, name, photo, type FROM chats WHERE id=?", chatId).Scan(&chat.Id, &chat.Name, &chat.Photo, &chat.ChatType)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Chat not found(DB). ", err)
		return chat, err
	}
	if err != nil {
		fmt.Println("Error getting chat data(DB). ", err)
		return chat, err
	}

	rows, err := db.c.Query("SELECT u.Id, u.Name, u.Photo FROM user_chats uc JOIN users u ON uc.userId = u.Id WHERE chatId=?", chatId)
	if err != nil {
		fmt.Println("Error fetching chat participants(GetChat - chats.go). ", err)
		return chat, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			fmt.Println("Error scanning chat participants(GetChat - chats.go). ", err)
			return chat, err
		}
		chat.Participants = append(chat.Participants, user)
	}
	if rows.Err() != nil {
		fmt.Println("Error fetching chat participants(GetChat - chats.go). ", err)
		return chat, err
	}

	rows, err = db.c.Query("SELECT id, content, sender, receiver, COALESCE(forwarded, 0) AS forwarded, sentTime FROM messages WHERE receiver=?", chatId)
	if err != nil {
		fmt.Println("Error fetching chat messages(GetChat - chats.go). ", err)
		return chat, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.Content, &message.Sender, &message.Receiver, &message.Forwarded, &message.TimeStamp)
		if err != nil {
			fmt.Println("Error scanning chat messages(GetChat - chats.go). ", err)
			return chat, err
		}
		chat.Messages = append(chat.Messages, message)
	}
	if rows.Err() != nil {
		fmt.Println("Error fetching chat messages(GetChat - chats.go). ", err)
		return chat, err
	}

	for i, message := range chat.Messages {
		messageId := message.Id
		rows, err = db.c.Query("SELECT reaction, userId, messageId FROM reactions WHERE messageId=?", messageId)
		if err != nil {
			fmt.Println("Error fetching chat reactions(GetChat - chats.go). ", err)
			return chat, err
		}
		defer rows.Close()

		for rows.Next() {
			var reaction Reaction
			err := rows.Scan(&reaction.Emoji, &reaction.UserId, &reaction.MessageId)
			if err != nil {
				fmt.Println("Error scanning chat reactions(GetChat - chats.go). ", err)
				return chat, err
			}
			chat.Messages[i].Reactions = append(chat.Messages[i].Reactions, reaction)
		}
		if rows.Err() != nil {
			fmt.Println("Error fetching chat reactions(GetChat - chats.go). ", err)
			return chat, err
		}

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
	fmt.Println("Chat created:", chatName, chatType)

	return chatId, nil
}

func (db *appdbimpl) GetAllChats(userId int) ([]Chat, error) {
	var chats []Chat
	rows, err := db.c.Query("SELECT id FROM chats c JOIN user_chats uc on c.id = uc.chatId WHERE uc.userId=?", userId)
	if err != nil {
		fmt.Println("Error fetching all chats(DB). ", err)
		return chats, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var chat Chat
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("Error scanning chat data(DB). ", err)
			return chats, err
		}
		chat, err = db.GetChat(id)
		if err != nil {
			fmt.Println("Error getting chat data(DB). ", err)
			return chats, err
		}
		chats = append(chats, chat)
	}
	if rows.Err() != nil {
		fmt.Println("Error fetching chat data(DB). ", err)
		return chats, err
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

func (db *appdbimpl) SetChatName(chatId int, chatName string) error {
	_, err := db.c.Exec("UPDATE chats SET name=? WHERE id=?", chatName, chatId)
	if err != nil {
		fmt.Println("Error updating chat name. SetChatName chats.go", err)
		return err
	}
	return nil
}

func (db *appdbimpl) SetChatPhoto(chatId int, newPhoto []byte) error {
	_, err := db.c.Exec("UPDATE chats SET photo=? WHERE id=?", newPhoto, chatId)
	if err != nil {
		fmt.Println("Error updating chat photo. SetChatPhoto chats.go", err)
		return err
	}
	return nil
}

func (db *appdbimpl) LeaveChat(chatId int, userId int) error {
	_, err := db.c.Exec("DELETE FROM user_chats WHERE chatId=? AND userId=?", chatId, userId)
	if err != nil {
		fmt.Println("Error leaving chat. LeaveChat chats.go", err)
		return err
	}
	return nil
}
