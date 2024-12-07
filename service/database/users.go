package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetUserName(userId string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM users WHERE id=?", userId).Scan(&name)
	return name, err
}
