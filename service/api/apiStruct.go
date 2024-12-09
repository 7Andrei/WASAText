package api

import (
	"github.com/7Andrei/WASAText/service/database"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Photo    []byte `json:"photo"`
}

func apiUser(user database.User) User {
	return User{
		Id:       user.Id,
		Username: user.Username,
		Photo:    user.Photo,
	}
}
