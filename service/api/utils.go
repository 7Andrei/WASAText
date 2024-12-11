package api

import "path/filepath"

func IsPhoto(fileName string) bool {
	fileExt := filepath.Ext(fileName)
	switch fileExt {
	case ".png", ".jpg", ".jpeg", ".webp":
		return true
	}
	return false
}

/*
func Authorized(userId int) bool {

}
*/
