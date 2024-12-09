package database

import (
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Photo    []byte `json:"photo"`
}

type Message struct {
	id        int       `json:"id"`
	content   string    `json:"content"`
	photo     []byte    `json:"photo"`
	sender    int       `json:"sender"`
	receiver  int       `json:"receiver"`
	forwarded uint64    `json:"forwarded"`
	timeStamp time.Time `json:"timestamp"`
}

type Chat struct {
	id       int    `json:"id"`
	chatType string `json:"chatType"`
	name     string `json:"chatName"`
	photo    []byte `json:"chatPhoto"`
}

type Reaction struct {
	userId    int    `json:"userId"`
	messageId int    `json:"messageId"`
	emoji     string `json:"emoji"`
}
