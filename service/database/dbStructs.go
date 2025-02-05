package database

import (
	"time"
)

type User struct {
	Id       int    `json:"userId"`
	Username string `json:"userName"`
	Photo    []byte `json:"userPhoto"`
}

type Message struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	Photo     []byte    `json:"photo"`
	Sender    int       `json:"sender"`
	Receiver  int       `json:"receiver"`
	Forwarded uint64    `json:"forwarded"`
	TimeStamp time.Time `json:"timestamp"`
}

type Chat struct {
	Id           int       `json:"id"`
	ChatType     string    `json:"chatType"`
	Name         string    `json:"chatName"`
	Photo        []byte    `json:"chatPhoto"`
	Participants []User    `json:"chatParticipants"`
	Messages     []Message `json:"chatMessages"`
}

type Reaction struct {
	UserId    int    `json:"userId"`
	MessageId int    `json:"messageId"`
	Emoji     string `json:"emoji"`
}
