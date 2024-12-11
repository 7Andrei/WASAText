package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) loginUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	//fmt.Println(user)
	if err != nil {
		fmt.Println("Error decoding user(api). ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Id, err = rt.db.Login(user.Username)
	if err != nil {
		fmt.Println("Error loggin in. ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(user)
	JSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(JSON)
}
