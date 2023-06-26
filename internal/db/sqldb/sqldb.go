package sqldb

import (
	"database/sql"

	"vault/internal/db/queries"
	"vault/internal/item"
)

// DB is sql-like database.
type SQLStore struct {
	*sql.DB
}

// Save saves service to chat.
func (db SQLStore) Save(chatID int64, service string, cred item.Credentials) error {
	prep, err := queries.GetPreparedStatement(queries.AddService)
	if err != nil {
		return err
	}
	_, err = prep.Exec(service, cred.Login, cred.Password, chatID, cred.Login, cred.Password, chatID)
	return err
}

// Get gets service from chat.
func (db SQLStore) Get(chatID int64, service string) (item.Credentials, error) {
	prep, err := queries.GetPreparedStatement(queries.GetService)
	if err != nil {
		return item.Credentials{}, err
	}

	var cred item.Credentials
	err = prep.QueryRow(service, chatID).Scan(&cred.Login, &cred.Password)
	return cred, err
}

// Delete deletes service from chat.
func (db SQLStore) Delete(chatID int64, serviceName string) error {
	prep, err := queries.GetPreparedStatement(queries.DeleteService)
	if err != nil {
		return err
	}

	r, err := prep.Exec(serviceName, chatID)
	if err != nil {
		return err
	}
	a, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if a == 0 {
		return queries.ErrNotFound
	}
	return nil
}

// GetLang gets language for chat.
func (db SQLStore) GetLang(chatID int64) (string, error) {
	prep, err := queries.GetPreparedStatement(queries.GetLang)
	if err != nil {
		return "", err
	}

	var lang string
	err = prep.QueryRow(chatID).Scan(&lang)
	return lang, err
}

// SetLang sets language for chat.
func (db SQLStore) SetLang(chatID int64, lang string) error {
	prep, err := queries.GetPreparedStatement(queries.AddOrUpdateChatLang)
	if err != nil {
		return err
	}
	_, err = prep.Exec(chatID, lang, lang)
	return err
}
