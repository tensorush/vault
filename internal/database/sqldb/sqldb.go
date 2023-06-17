package sqldb

import (
	"database/sql"

	"vault-bot/internal/database/queries"
	"vault-bot/internal/secret"
)

// DB is sql-like database.
type SQLStorage struct {
	*sql.DB
}

// Save saves service to chat.
func (db SQLStorage) Save(chatID int64, service string, cred secret.Credentials) error {
	prep, err := queries.GetPreparedStatement(queries.AddService)
	if err != nil {
		return err
	}
	_, err = prep.Exec(service, cred.Login, cred.Password, chatID, cred.Login, cred.Password, chatID)
	return err
}

// Get gets service from chat.
func (db SQLStorage) Get(chatID int64, service string) (secret.Credentials, error) {
	prep, err := queries.GetPreparedStatement(queries.GetService)
	if err != nil {
		return secret.Credentials{}, err
	}

	var cred secret.Credentials
	err = prep.QueryRow(service, chatID).Scan(&cred.Login, &cred.Password)
	return cred, err
}

// Delete deletes service from chat.
func (db SQLStorage) Delete(chatID int64, serviceName string) error {
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
func (db SQLStorage) GetLang(chatID int64) (string, error) {
	prep, err := queries.GetPreparedStatement(queries.GetLang)
	if err != nil {
		return "", err
	}

	var lang string
	err = prep.QueryRow(chatID).Scan(&lang)
	return lang, err
}

// SetLang sets language for chat.
func (db SQLStorage) SetLang(chatID int64, lang string) error {
	prep, err := queries.GetPreparedStatement(queries.AddOrUpdateChatLang)
	if err != nil {
		return err
	}
	_, err = prep.Exec(chatID, lang, lang)
	return err
}
