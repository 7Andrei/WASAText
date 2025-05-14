package api

import (
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
	if err != nil {
		// http.Error(w, "Error during conversion to int", http.StatusBadRequest)
		return false
	}
	_, available, err := rt.db.GetUser(headerId)

	if err != nil || !available {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		return false
	}
	return true

}
