package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println("Error decoding message(sendMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message.Id, err = rt.db.SendMessage(message.Content, message.Photo, message.Sender, message.Receiver, int(message.Forwarded))
	if err != nil {
		fmt.Println("Error sending message(sendMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var message Message

	userId := r.Header.Get("Authorization")
	if userId == "" {
		fmt.Println("userId header not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Error converting string to int(forwardMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tempInt, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		fmt.Println("Error converting string to int(forwardMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message.Id = tempInt

	tempInt, err = strconv.Atoi(ps.ByName("chat_id"))
	if err != nil {
		fmt.Println("Error converting string to int(forwardMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message.Receiver = tempInt
	message.Sender = userIdInt

	err = rt.db.ForwardMessage(message.Id, message.Receiver, message.Sender)
	if err != nil {
		fmt.Println("Error forwarding message(forwardMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
