package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	// err := json.NewDecoder(r.Body).Decode(&user)
	// if err != nil {
	// 	fmt.Println("Error decoding user Id(api). ", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	var tempId string = ps.ByName("user_id")
	user_id, err := strconv.Atoi(tempId)
	if err != nil {
		fmt.Println("Error converting user id. getUser api-user.go", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Id = user_id

	tmpUser, _, erro := rt.db.GetUser(user.Id)
	user = apiUser(tmpUser)

	if erro != nil {
		fmt.Println("Error fetching user. getUser api-user.go", erro)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("User found"))

	w.Header().Set("Content-Type", "application/json")
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshalling user. getUser api-user.go", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(userJSON)
}

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	authentication := r.Header.Get("Authorization")
	headerId, err := strconv.Atoi(authentication)
	if err != nil {
		fmt.Println("Error during conversion to int setUsername api-user.go")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DBuser, available, err := rt.db.GetUser(headerId)

	if err != nil || !available {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		fmt.Println("Unauthorized. ", err)
		return
	}
	user = apiUser(DBuser)

	err = json.NewDecoder(r.Body).Decode(&user)
	// fmt.Println(user)
	if err != nil {
		fmt.Println("Error decoding username. setUsername api-user.go ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetUsername(user.Id, user.Username)
	if err != nil {
		fmt.Println("Error updating username. setUsername api-user.go", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Username updated"))
}

func (rt *_router) setPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	authentication := r.Header.Get("Authorization")
	headerId, err := strconv.Atoi(authentication)
	if err != nil {
		fmt.Println("Error during conversion to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DBuser, available, err := rt.db.GetUser(headerId)

	if err != nil || !available {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		fmt.Println("Unauthorized. ", err)
		return
	}
	user = apiUser(DBuser)

	newPhotoMulti, fileHeader, err := r.FormFile("userPhoto")
	if err != nil {
		fmt.Println("Photo not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fileHeader.Filename
	if !IsPhoto(fileName) {
		fmt.Println("Photo not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newPhoto, err := io.ReadAll(newPhotoMulti)
	if err != nil {
		fmt.Println("Error reading file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetUserPhoto(user.Id, newPhoto)
	if err != nil {
		fmt.Println("Error updating photo. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Photo updated"))
}

func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	users, err := rt.db.GetAllUsers()
	if err != nil {
		fmt.Println("Error fetching all users getAllUsers api-user.go. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshalling users. ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(usersJSON)
}
