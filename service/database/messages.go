package database

func (db *appdbimpl) SendMessage(messageContent string, messagePhoto []byte, messageSender int, messageReceiver int, messageForwarded int, messageReply int) (int, error) {
	var messageId int
	err := db.c.QueryRow("INSERT INTO messages (content, photo, sender, receiver, forwarded, reply) VALUES (?, ?, ?, ?, ?, ?) RETURNING id", messageContent, messagePhoto, messageSender, messageReceiver, messageForwarded, messageReply).Scan(&messageId)
	if err != nil {
		return messageId, err
	}

	return messageId, nil
}

func (db *appdbimpl) ForwardMessage(messageId int, messageReceiver int, messageForwarded int) error {
	var oldMessage Message
	err := db.c.QueryRow("SELECT * FROM messages WHERE id = ?", messageId).Scan(&oldMessage.Id, &oldMessage.Content, &oldMessage.Photo, &oldMessage.Sender, &oldMessage.Receiver, &oldMessage.Forwarded, &oldMessage.Reply, &oldMessage.TimeStamp)
	if err != nil {
		return err
	}
	oldMessage.Receiver = messageReceiver
	oldMessage.Forwarded = oldMessage.Sender
	oldMessage.Sender = messageForwarded
	_, err = db.c.Exec("INSERT INTO messages (content, photo, sender, receiver, forwarded, reply) VALUES (?, ?, ?, ?, ?, 0)", oldMessage.Content, oldMessage.Photo, oldMessage.Sender, oldMessage.Receiver, oldMessage.Forwarded)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteMessage(messageId int) error {
	_, err := db.c.Exec("DELETE FROM messages WHERE id = ?", messageId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) AddReaction(userId int, messageId int, reaction string) error {
	_, err := db.c.Exec("INSERT INTO reactions (reaction, userId, messageId) VALUES (?, ?, ?)", reaction, userId, messageId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteReaction(userId, messageId int) error {
	_, err := db.c.Exec("DELETE FROM reactions WHERE userId=? AND messageId=? ", userId, messageId)
	if err != nil {
		return err
	}
	return nil
}
