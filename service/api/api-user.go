package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tempId = ps.ByName("user_id")
	user_id, err := strconv.Atoi(tempId)
	if err != nil {
		http.Error(w, "Error converting user id", http.StatusBadRequest)
		return
	}
	user.Id = user_id

	tmpUser, _, erro := rt.db.GetUser(user.Id)
	user = apiUser(tmpUser)

	if erro != nil {
		http.Error(w, "Error fetching user", http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("User found"))

	w.Header().Set("Content-Type", "application/json")
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshalling user", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)
}

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authentication := r.Header.Get("Authorization")
	headerId, err := strconv.Atoi(authentication)
	if err != nil {
		http.Error(w, "Error during conversion to int", http.StatusBadRequest)
		return
	}
	DBuser, available, err := rt.db.GetUser(headerId)

	if err != nil || !available {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user = apiUser(DBuser)

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error decoding username", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	err = rt.db.SetUsername(user.Id, user.Username)
	if err != nil {
		http.Error(w, "Error updating username", http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Username updated"))
}

func (rt *_router) setPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authentication := r.Header.Get("Authorization")
	headerId, err := strconv.Atoi(authentication)
	if err != nil {
		http.Error(w, "Error during conversion to int", http.StatusBadRequest)
		return
	}
	DBuser, available, err := rt.db.GetUser(headerId)

	if err != nil || !available {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user = apiUser(DBuser)

	newPhotoMulti, fileHeader, err := r.FormFile("userPhoto")
	if err != nil {
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}

	fileName := fileHeader.Filename
	if !IsPhoto(fileName) {
		http.Error(w, "File is not a photo", http.StatusBadRequest)
		return
	}

	newPhoto, err := io.ReadAll(newPhotoMulti)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	err = rt.db.SetUserPhoto(user.Id, newPhoto)
	if err != nil {
		http.Error(w, "Error updating photo", http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Photo updated"))
}

func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if !Authorized(r, rt) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := rt.db.GetAllUsers()
	if err != nil {
		http.Error(w, "Error fetching all users", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error marshalling users", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(usersJSON)
}
