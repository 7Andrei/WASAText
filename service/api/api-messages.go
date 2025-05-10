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
		// http.Error(w, "Photo not found", http.StatusBadRequest)
		// return
	}

	if fileHeader != nil {

		fileName := fileHeader.Filename
		if !IsPhoto(fileName) {
			checkPhoto = false
			// http.Error(w, "File is not a photo", http.StatusBadRequest)
			// return
		}

		newPhoto, err = io.ReadAll(newPhotoMulti)
		if err != nil {
			checkPhoto = false
			// http.Error(w, "Error reading file", http.StatusBadRequest)
			// return
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
		// http.Error(w, "Message content not found", http.StatusBadRequest)
		// return
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

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := r.Header.Get("Authorization")
	if userId == "" {
		fmt.Println("userId header not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Error converting string to int(deleteMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("userIdInt:", userIdInt)

	tempInt, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		fmt.Println("Error converting string to int(deleteMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(tempInt)
	if err != nil {
		fmt.Println("Error deleting message(deleteMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
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
		fmt.Println("userId header not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Error converting string to int(deleteMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		fmt.Println("Error decoding reaction(addReaction api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		fmt.Println("Error converting string to int(addReaction api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.AddReaction(userIdInt, messageId, reaction.Emoji)
	if err != nil {
		fmt.Println("Error adding reaction(addReaction api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
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
		fmt.Println("userId header not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Error converting string to int(deleteMessage api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		fmt.Println("Error converting string to int(deleteReaction api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteReaction(userIdInt, messageId)
	if err != nil {
		fmt.Println("Error deleting reaction(deleteReaction api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
