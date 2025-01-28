package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chat Chat

	// err := json.NewDecoder(r.Body).Decode(&chat)
	// if err != nil {
	// 	fmt.Println("Error decoding chat Id(api). ", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	var tempId string = ps.ByName("chat_id")
	chat_Id, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting chat id(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat.Id = chat_Id

	tmpChat, err := rt.db.GetChat(chat.Id)
	if err == sql.ErrNoRows {
		fmt.Println("Chat not found(api). ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Println("Error fetching chat(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chat = apiChat(tmpChat)
	// _, _ = w.Write([]byte("Chat found"))
	fmt.Println("Chat found", chat)

	w.Header().Set("Content-Type", "application/json")
	chatJSON, err := json.Marshal(chat)
	if err != nil {
		fmt.Println("Error marshalling chat(api). ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(chatJSON)
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
