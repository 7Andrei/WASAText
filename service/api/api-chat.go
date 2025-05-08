package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chat Chat

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tempId string = ps.ByName("chat_id")
	chat_Id, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, "Error converting chat id", http.StatusBadRequest)
		return
	}
	chat.Id = chat_Id

	tmpChat, err := rt.db.GetChat(chat.Id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error fetching chat", http.StatusBadRequest)
		return
	}

	chat = apiChat(tmpChat)

	w.Header().Set("Content-Type", "application/json")
	chatJSON, err := json.Marshal(chat)
	if err != nil {
		http.Error(w, "Error marshalling chat", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(chatJSON)
}

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chat Chat

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIdHeader := r.Header.Get("Authorization")
	if userIdHeader == "" {
		http.Error(w, "userId header not found", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdHeader)
	if err != nil {
		http.Error(w, "Error converting userId header to int", http.StatusBadRequest)
		return
	}

	chat.ChatType = r.FormValue("chatType")
	if chat.ChatType == "group" {
		newPhotoMulti, fileHeader, err := r.FormFile("chatPhoto")
		if err != nil {
			http.Error(w, "Photo not found", http.StatusBadRequest)
			return
		}

		fileName := fileHeader.Filename
		if !IsPhoto(fileName) {
			http.Error(w, "Photo not found", http.StatusBadRequest)
			return
		}

		newPhoto, err := io.ReadAll(newPhotoMulti)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}

		chat.Photo = newPhoto
	}

	chat.Name = r.FormValue("chatName")
	participants := r.FormValue("chatParticipants")
	fmt.Println("utenti ", participants)
	err = json.Unmarshal([]byte(participants), &chat.Participants)
	if err != nil {
		http.Error(w, "Error decoding participants", http.StatusBadRequest)
		return
	}
	fmt.Println(chat.Participants)

	if !(chat.ChatType == "private" || chat.ChatType == "group") {
		http.Error(w, "Chat type can only be private or group", http.StatusBadRequest)
		return
	}
	if (len(chat.Participants) > 1) && (chat.ChatType == "private") {
		http.Error(w, "Chat type is private but more than one participant", http.StatusBadRequest)
		return
	}

	chat.Id, err = rt.db.CreateChat(chat.Name, chat.Photo, chat.ChatType)
	if err != nil {
		http.Error(w, "Error creating chat", http.StatusBadRequest)
		return
	}

	for _, participant := range chat.Participants {
		err := rt.db.AddParticipant(chat.Id, participant.Id)
		if err != nil {
			http.Error(w, "Error adding participant", http.StatusBadRequest)
			return
		}
	}
	err = rt.db.AddParticipant(chat.Id, userId)
	if err != nil {
		http.Error(w, "Error adding participant", http.StatusBadRequest)
		return
	}

	w.Write([]byte(strconv.Itoa(chat.Id)))
}

func (rt *_router) getAllChats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chats []Chat

	if !Authorized(r, rt) {
		fmt.Println("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userIdHeader := r.Header.Get("Authorization")
	if userIdHeader == "" {
		fmt.Println("userId header not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Println("userId header:", userIdHeader)

	userId, err := strconv.Atoi(userIdHeader)
	if err != nil {
		fmt.Println("Error converting userId header to int getAllChats api-chat.go", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatList, err := rt.db.GetAllChats(userId)
	if err != nil {
		fmt.Println("Error fetching chat list(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, chat := range chatList {
		chats = append(chats, apiChat(chat))
	}

	w.Header().Set("Content-Type", "application/json")
	chatsJSON, err := json.Marshal(chats)
	if err != nil {
		fmt.Println("Error marshalling chat list(api). ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(chatsJSON)
}

func (rt *_router) setChatName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chat Chat
	var chatId int

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tempId string = ps.ByName("chat_id")
	chatId, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting chat id setChatName api-chat.go. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		fmt.Println("Error decoding chat name. setChatName api-chat.go ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if chat.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetChatName(chatId, chat.Name)
	if err != nil {
		fmt.Println("Error updating chat name. setChatName api-chat.go", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Chat Name updated"))
}

func (rt *_router) setChatPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chatId int

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tempId string = ps.ByName("chat_id")
	chatId, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting chat id setChatPhoto api-chat.go. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newPhotoMulti, fileHeader, err := r.FormFile("chatPhoto")
	if err != nil {
		fmt.Println("Photo not found setChatPhoto api-chat.go")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fileHeader.Filename
	if !IsPhoto(fileName) {
		fmt.Println("Photo not found setChatPhoto api-chat.go")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newPhoto, err := io.ReadAll(newPhotoMulti)
	if err != nil {
		fmt.Println("Error reading file setChatPhoto api-chat.go")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetChatPhoto(chatId, newPhoto)
	if err != nil {
		fmt.Println("Error updating photo. setChatPhoto api-chat.go ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Photo updated"))
}

func (rt *_router) addUserToChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var chatId int
	var chat Chat

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tempId string = ps.ByName("chat_id")
	chatId, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting chat id addUserToChat api-chat.go. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		fmt.Println("Error decoding participant id. addUserToChat api-chat.go ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, participant := range chat.Participants {
		err = rt.db.AddParticipant(chatId, participant.Id)
		if err != nil {
			fmt.Println("Error adding partecipant (AddParticipants api-chat)\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	_, _ = w.Write([]byte("Participant added"))
}

func (rt *_router) leaveChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var chatId int
	var userId int

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tempId string = ps.ByName("chat_id")
	chatId, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting chat id leaveChat api-chat.go. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authentication := r.Header.Get("Authorization")
	userId, err = strconv.Atoi(authentication)
	if err != nil {
		fmt.Println("Error during conversion to int leaveChat api-chat.go")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.LeaveChat(chatId, userId)
	if err != nil {
		fmt.Println("Error leaving chat. leaveChat api-chat.go ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Chat left"))
}
