package api

import (
	"time"

	"github.com/7Andrei/WASAText/service/database"
)

type User struct {
	Id       int    `json:"userId"`
	Username string `json:"userName"`
	Photo    []byte `json:"userPhoto"`
}

func apiUser(user database.User) User {
	return User{
		Id:       user.Id,
		Username: user.Username,
		Photo:    user.Photo,
	}
}

type Chat struct {
	Id           int       `json:"id"`
	Name         string    `json:"chatName"`
	Photo        []byte    `json:"chatPhoto"`
	ChatType     string    `json:"chatType"`
	Participants []User    `json:"chatParticipants"`
	Messages     []Message `json:"chatMessages"`
}

func apiChat(chat database.Chat) Chat {
	participants := make([]User, len(chat.Participants))
	messages := make([]Message, len(chat.Messages))
	for i, user := range chat.Participants {
		participants[i] = apiUser(user)
	}
	for i, message := range chat.Messages {
		messages[i] = apiMessage(message)
	}

	return Chat{
		Id:           chat.Id,
		Name:         chat.Name,
		Photo:        chat.Photo,
		ChatType:     chat.ChatType,
		Participants: participants,
		Messages:     messages,
	}
}

type Message struct {
	Id        int        `json:"id"`
	Content   string     `json:"text"`
	Photo     []byte     `json:"photo"`
	Sender    int        `json:"sender"`
	Receiver  int        `json:"receiver"`
	Forwarded int        `json:"forwarded"`
	TimeStamp time.Time  `json:"dateTime"`
	Reactions []Reaction `json:"reactions"`
	Reply     int        `json:"reply"`
}

func apiMessage(message database.Message) Message {
	reactions := make([]Reaction, len(message.Reactions))
	for i, reaction := range message.Reactions {
		reactions[i] = Reaction{
			UserId:    reaction.UserId,
			MessageId: reaction.MessageId,
			Emoji:     reaction.Emoji,
		}
	}

	return Message{
		Id:        message.Id,
		Content:   message.Content,
		Photo:     message.Photo,
		Sender:    message.Sender,
		Receiver:  message.Receiver,
		Forwarded: message.Forwarded,
		TimeStamp: message.TimeStamp,
		Reactions: reactions,
		Reply:     message.Reply,
	}
}

type Reaction struct {
	UserId    int    `json:"userId"`
	MessageId int    `json:"messageId"`
	Emoji     string `json:"reaction"`
}
