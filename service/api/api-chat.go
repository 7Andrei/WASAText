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

	// err := json.NewDecoder(r.Body).Decode(&chat)
	// if err != nil {
	// 	fmt.Println("Error decoding chat Id(api). ", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	if !Authorized(r, rt) {
		fmt.Println("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tempId string = ps.ByName("chat_id")
	chat_Id, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting chat id(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat.Id = chat_Id

	tmpChat, err := rt.db.GetChat(chat.Id)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Chat not found(api). ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Println("Error fetching chat(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Println("Temp Chat:", tmpChat.Messages)
	chat = apiChat(tmpChat)
	// fmt.Println("Chat nuova:", chat.Messages)

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
		w.Write([]byte("Unauthorized"))
		return
	}

	userIdHeader := r.Header.Get("Authorization")
	if userIdHeader == "" {
		fmt.Println("userId header not found")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("userId header not found"))
		return
	}

	userId, err := strconv.Atoi(userIdHeader)
	if err != nil {
		fmt.Println("Error converting userId header to int createChat api-chat.go", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error converting userId header to int"))
		return
	}

	chat.ChatType = r.FormValue("chatType")
	if chat.ChatType == "group" {
		newPhotoMulti, fileHeader, err := r.FormFile("chatPhoto")
		if err != nil {
			fmt.Println("Photo not found", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Photo not found"))
			return
		}

		fileName := fileHeader.Filename
		// fmt.Println("File name:")
		// fmt.Println("File name:", fileName)
		if !IsPhoto(fileName) {
			fmt.Println("Photo not found")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Photo not found"))
			return
		}

		newPhoto, err := io.ReadAll(newPhotoMulti)
		if err != nil {
			fmt.Println("Error reading file")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error reading file"))
			return
		}

		chat.Photo = newPhoto
	}

	chat.Name = r.FormValue("chatName")
	participants := r.FormValue("chatParticipants")
	fmt.Println("utenti ", participants)
	err = json.Unmarshal([]byte(participants), &chat.Participants)
	if err != nil {
		fmt.Println("Error decoding participants:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding participants"))
		return
	}
	fmt.Println(chat.Participants)

	if !(chat.ChatType == "private" || chat.ChatType == "group") {
		fmt.Println("Chat type can only be private or group")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Chat type can only be private or group"))
		return
	}
	if (len(chat.Participants) > 1) && (chat.ChatType == "private") {
		fmt.Println("Chat type is private but more than one participant")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Chat type is private but more than one participant was selected"))
		return
	}

	// fmt.Println("Dati chat:", chat.Name, chat.ChatType)
	// fmt.Println("partecipanti:", chat.Participants)

	chat.Id, err = rt.db.CreateChat(chat.Name, chat.Photo, chat.ChatType)
	if err != nil {
		fmt.Println("Error creating chat. ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating chat"))
		return
	}

	// fmt.Println("chat created", chat)

	for _, participant := range chat.Participants {
		err := rt.db.AddParticipant(chat.Id, participant.Id)
		if err != nil {
			fmt.Println("Error adding partecipant (AddParticipants api-chat)\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error adding partecipant"))
			return
		}
	}
	err = rt.db.AddParticipant(chat.Id, userId)
	if err != nil {
		fmt.Println("Error adding partecipant (AddParticipants api-chat)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error adding partecipant"))
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

	// authentication := r.Header.Get("Authorization")
	// headerId, err := strconv.Atoi(authentication)
	// if err != nil {
	// 	fmt.Println("Error during conversion to int setChatName api-chat.go")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

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

	// authentication := r.Header.Get("Authorization")
	// headerId, err := strconv.Atoi(authentication)
	// if err != nil {
	// 	fmt.Println("Error during conversion to int")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// DBuser, available, err := rt.db.GetUser(headerId)

	// if err != nil || !available {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	fmt.Println("Unauthorized. ", err)
	// 	return
	// }
	// user = apiUser(DBuser)

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
