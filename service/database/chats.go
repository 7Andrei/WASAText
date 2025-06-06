package database

import (
	"database/sql"
	"errors"
	"time"
)

func (db *appdbimpl) GetChat(chatId int, flagSingle bool, userId int) (Chat, error) {
	var chat Chat
	err := db.c.QueryRow("SELECT id, name, photo, type FROM chats WHERE id=?", chatId).Scan(&chat.Id, &chat.Name, &chat.Photo, &chat.ChatType)
	if errors.Is(err, sql.ErrNoRows) {
		return chat, err
	}
	if err != nil {
		return chat, err
	}

	rows, err := db.c.Query("SELECT u.Id, u.Name, u.Photo FROM user_chats uc JOIN users u ON uc.userId = u.Id WHERE chatId=?", chatId)
	if err != nil {
		return chat, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			return chat, err
		}
		chat.Participants = append(chat.Participants, user)
	}
	if rows.Err() != nil {
		return chat, err
	}

	rows, err = db.c.Query("SELECT id, content, sender, receiver, COALESCE(forwarded, 0) AS forwarded, sentTime, photo, reply FROM messages WHERE receiver=?", chatId)
	if err != nil {
		return chat, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.Content, &message.Sender, &message.Receiver, &message.Forwarded, &message.TimeStamp, &message.Photo, &message.Reply)
		if err != nil {
			return chat, err
		}
		chat.Messages = append(chat.Messages, message)
	}
	if rows.Err() != nil {
		return chat, err
	}

	for i, message := range chat.Messages {
		messageId := message.Id
		rows, err = db.c.Query("SELECT reaction, userId, messageId FROM reactions WHERE messageId=?", messageId)
		if err != nil {
			return chat, err
		}
		defer rows.Close()

		for rows.Next() {
			var reaction Reaction
			err := rows.Scan(&reaction.Emoji, &reaction.UserId, &reaction.MessageId)
			if err != nil {
				return chat, err
			}
			chat.Messages[i].Reactions = append(chat.Messages[i].Reactions, reaction)
		}
		if rows.Err() != nil {
			return chat, err
		}

	}

	if flagSingle {
		_, err = db.c.Exec("UPDATE user_chats SET lastAccess = CURRENT_TIMESTAMP WHERE chatId = ? AND userId = ?", chatId, userId)
		if err != nil {
			return chat, err
		}
	}

	return chat, nil
}

func (db *appdbimpl) CreateChat(chatName string, chatPhoto []byte, chatType string) (int, error) {

	var chatId int
	err := db.c.QueryRow("INSERT INTO chats (name, photo, type) VALUES (?, ?, ?) RETURNING id", chatName, chatPhoto, chatType).Scan(&chatId)
	if err != nil {
		return chatId, err
	}

	return chatId, nil
}

func (db *appdbimpl) GetAllChats(userId int) ([]Chat, error) {
	var chats []Chat
	rows, err := db.c.Query("SELECT id FROM chats c JOIN user_chats uc on c.id = uc.chatId WHERE uc.userId=?", userId)
	if err != nil {
		return chats, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var chat Chat
		err := rows.Scan(&id)
		if err != nil {
			return chats, err
		}
		chat, err = db.GetChat(id, false, 0)
		if err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}
	if rows.Err() != nil {
		return chats, err
	}
	return chats, nil
}

func (db *appdbimpl) AddParticipant(chatId int, participantId int) error {
	_, err := db.c.Exec("INSERT INTO user_chats (chatId, userId) VALUES (?, ?)", chatId, participantId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetChatName(chatId int, chatName string) error {
	_, err := db.c.Exec("UPDATE chats SET name=? WHERE id=?", chatName, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetChatPhoto(chatId int, newPhoto []byte) error {
	_, err := db.c.Exec("UPDATE chats SET photo=? WHERE id=?", newPhoto, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) LeaveChat(chatId int, userId int) error {
	_, err := db.c.Exec("DELETE FROM user_chats WHERE chatId=? AND userId=?", chatId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) MessageSeen(chatId int, userId int) ([]time.Time, error) {
	var lastAccesses []time.Time
	rows, err := db.c.Query("SELECT lastAccess FROM user_chats WHERE chatId=? AND userId != ?", chatId, userId)
	if err != nil {
		return lastAccesses, err
	}
	defer rows.Close()
	for rows.Next() {
		var lastAccess time.Time
		err := rows.Scan(&lastAccess)
		if err != nil {
			return lastAccesses, err
		}
		lastAccesses = append(lastAccesses, lastAccess)
	}
	if rows.Err() != nil {
		return lastAccesses, err
	}
	return lastAccesses, nil
}
