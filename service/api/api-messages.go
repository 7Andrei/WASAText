package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
