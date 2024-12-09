package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
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

	tmpUser, erro := rt.db.GetUser(user.Id)
	user = apiUser(tmpUser)

	if erro != nil {
		fmt.Println("Error fetching user. ", erro)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("User found"))
	fmt.Println("User found", user)
}
