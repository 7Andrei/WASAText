package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

func IsPhoto(fileName string) bool {
	fileExt := filepath.Ext(fileName)
	switch fileExt {
	case ".png", ".jpg", ".jpeg", ".webp":
		return true
	}
	return false
}

func Authorized(r *http.Request, rt *_router) bool {

	authentication := r.Header.Get("Authorization")
	headerId, err := strconv.Atoi(authentication)
	fmt.Println(authentication)
	if err != nil {
		fmt.Println("Error during conversion to int (Authorized utils.go)\n", err)
		// w.WriteHeader(http.StatusBadRequest)
		return false
	}
	_, available, err := rt.db.GetUser(headerId)

	if err != nil || !available {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		fmt.Println("Unauthorized. ", err)
		return false
	}
	return true

}
