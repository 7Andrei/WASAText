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

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// tmpReactionId := ps.ByName("reaction_id")
	// fmt.Println("tmpReactionId:", tmpReactionId)
	// reactionId, err := strconv.Atoi(tmpReactionId)
	// if err != nil {
	// 	fmt.Println("Error converting string to int(deleteReaction api-message.go)\n", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	err = rt.db.DeleteReaction(userIdInt, messageId)
	if err != nil {
		fmt.Println("Error deleting reaction(deleteReaction api-message.go)\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
