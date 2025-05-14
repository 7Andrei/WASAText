package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var message Message
	var newPhoto []byte
	var checkPhoto = true
	var checkContent = true
	newPhoto = nil

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newPhotoMulti, fileHeader, err := r.FormFile("photo")
	if err != nil {
		checkPhoto = false
	}

	if fileHeader != nil {

		fileName := fileHeader.Filename
		if !IsPhoto(fileName) {
			checkPhoto = false
		}

		newPhoto, err = io.ReadAll(newPhotoMulti)
		if err != nil {
			checkPhoto = false
		}
		message.Photo = newPhoto
	} else {
		checkPhoto = false
		message.Photo = nil
	}

	authorization := r.Header.Get("Authorization")
	message.Sender, err = strconv.Atoi(authorization)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	message.Receiver, err = strconv.Atoi(ps.ByName("chat_id"))
	if err != nil {
		http.Error(w, "Error fetching receiver id", http.StatusBadRequest)
		return
	}

	message.Content = r.FormValue("text")
	if message.Content == "" {
		checkContent = false
	}

	if checkPhoto || checkContent {
		message.Id, err = rt.db.SendMessage(message.Content, message.Photo, message.Sender, message.Receiver, int(message.Forwarded))
		if err != nil {
			http.Error(w, "Error sending message", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Message content and photo not found", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var message Message

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := r.Header.Get("Authorization")
	if userId == "" {
		http.Error(w, "userId header not found", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Error converting string to int(forwardMessage api-message.go)", http.StatusBadRequest)
		return
	}

	tempInt, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		http.Error(w, "Error converting string to int(forwardMessage api-message.go)", http.StatusBadRequest)
		return
	}
	message.Id = tempInt

	tempInt, err = strconv.Atoi(ps.ByName("chat_id"))
	if err != nil {
		http.Error(w, "Error converting string to int(forwardMessage api-message.go)", http.StatusBadRequest)
		return
	}
	message.Receiver = tempInt
	message.Sender = userIdInt

	err = rt.db.ForwardMessage(message.Id, message.Receiver, message.Sender)
	if err != nil {
		http.Error(w, "Error forwarding message", http.StatusBadRequest)
		return
	}
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := r.Header.Get("Authorization")
	if userId == "" {
		http.Error(w, "userId header not found", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Error converting string to int", http.StatusBadRequest)
		return
	}
	fmt.Println("userIdInt:", userIdInt)

	tempInt, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		http.Error(w, "Error converting string to int", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(tempInt)
	if err != nil {
		http.Error(w, "Error deleting message", http.StatusBadRequest)
		return
	}
}

func (rt *_router) addReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var reaction Reaction

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := r.Header.Get("Authorization")
	if userId == "" {
		http.Error(w, "userId header not found", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Error converting string to int", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		http.Error(w, "Error decoding reaction", http.StatusBadRequest)
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		http.Error(w, "Error converting string to int", http.StatusBadRequest)
		return
	}

	err = rt.db.AddReaction(userIdInt, messageId, reaction.Emoji)
	if err != nil {
		http.Error(w, "Error adding reaction", http.StatusBadRequest)
		return
	}
}

func (rt *_router) deleteReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := r.Header.Get("Authorization")
	if userId == "" {
		http.Error(w, "userId header not found", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Error converting string to int", http.StatusBadRequest)
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		http.Error(w, "Error converting string to int", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteReaction(userIdInt, messageId)
	if err != nil {
		http.Error(w, "Error deleting reaction", http.StatusBadRequest)
		return
	}
}
