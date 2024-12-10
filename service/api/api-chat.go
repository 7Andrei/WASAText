package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// /RIFARE
func (rt *_router) getChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chat Chat

	err := json.NewDecoder(r.Body).Decode(&chat.Id)
	if err != nil {
		fmt.Println("Error decoding chat Id Id(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpChat, err := rt.db.GetChat(chat.Id)
	chat = apiChat(tmpChat)

	if err != nil {
		fmt.Println("Error fetching chat. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("Chat found"))
	fmt.Println("Chat found", chat)
}
