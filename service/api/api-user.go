package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	err = json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	if err != nil {
		fmt.Println("Error decoding username(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.SetUsername(user.Id, user.Username)
	if err != nil {
		fmt.Println("Error updating username. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Username updated"))
}

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Error decoding user Id(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpUser, _, erro := rt.db.GetUser(user.Id)
	//RISOLVERE
	// available = true
	// if available {
	// 	fmt.Println("ok")
	// }
	user = apiUser(tmpUser)

	if erro != nil {
		fmt.Println("Error fetching user. ", erro)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("User found"))
	fmt.Println("User found", user)
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

	newPhotoMulti, fileHeader, err := r.FormFile("photo")
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
