package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chat Chat

	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		fmt.Println("Error decoding chat Id(api). ", err)
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

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//var user User
	var chat Chat

	if !Authorized(r, rt) {
		fmt.Println("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
	}

	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		fmt.Println("Error decoding chat Id(api). ", err, "chat", chat)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !(chat.ChatType == "private" || chat.ChatType == "group") {
		fmt.Println("Chat type can only be private or group")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Dati chat", chat.Name, chat.ChatType)

	err = rt.db.CreateChat(chat.Name, chat.Photo, chat.ChatType)
	if err != nil {
		fmt.Println("Error creating chat. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("chat created"))
	fmt.Println("chat created", chat)
}
