package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"

	"vault/internal/db"
	"vault/internal/item"
)

const defaultLanguage = "en"

// Vault is the main struct for the application logic.
type Vault struct {
	db     *db.DB
	cipher cipher.Block
	logger *zap.Logger
}

// New creates a new Vault.
func New(db *db.DB, key string, logger *zap.Logger) (*Vault, error) {
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &Vault{
		db:     db,
		cipher: cipher,
		logger: logger,
	}, nil
}

// Get returns the secret from the database.
func (v *Vault) Get(chatID int64, service string) (item.Credentials, error) {
	service, err := v.Hash(service)
	if err != nil {
		err = fmt.Errorf("vault.Hash: %w", err)
		v.logger.Warn(err.Error())
		return item.Credentials{}, err
	}

	cred, err := v.db.Get(chatID, service)
	if err != nil {
		err = fmt.Errorf("vault.Get: %w", err)
		v.logger.Warn(err.Error())
		return item.Credentials{}, err
	}
	cred.Login, err = v.Decrypt(cred.Login)
	if err != nil {
		err = fmt.Errorf("vault.Decrypt: %w", err)
		v.logger.Warn(err.Error())
		return item.Credentials{}, err
	}

	cred.Password, err = v.Decrypt(cred.Password)
	if err != nil {
		err = fmt.Errorf("vault.Decrypt: %w", err)
		v.logger.Warn(err.Error())
		return item.Credentials{}, err
	}

	return cred, nil
}

// Save saves the secret to the database.
func (v *Vault) Save(chatID int64, service, login, password string) (err error) {
	login, err = v.Encrypt(login)
	if err != nil {
		err = fmt.Errorf("vault.Encrypt: %w", err)
		v.logger.Warn(err.Error())
		return err
	}

	password, err = v.Encrypt(password)
	if err != nil {
		err = fmt.Errorf("vault.Encrypt: %w", err)
		v.logger.Warn(err.Error())
		return err
	}

	service, err = v.Hash(service)
	if err != nil {
		err = fmt.Errorf("vault.Hash: %w", err)
		v.logger.Warn(err.Error())
		return err
	}

	if err := v.db.Save(chatID, service, item.Credentials{Login: login, Password: password}); err != nil {
		err = fmt.Errorf("vault.Save: %w", err)
		v.logger.Warn(err.Error())
		return err
	}

	return nil
}

// Delete deletes the secret from the database.
func (v *Vault) Delete(chatID int64, service string) (err error) {
	service, err = v.Hash(service)
	if err != nil {
		err = fmt.Errorf("vault.Hash: %w", err)
		v.logger.Warn(err.Error())
		return err
	}
	if err := v.db.Delete(chatID, service); err != nil {
		err = fmt.Errorf("vault.Delete: %w", err)
		v.logger.Warn(err.Error())
		return err
	}
	return nil
}

// GetLang returns the language of the user.
func (v *Vault) GetLang(chatID int64) string {
	l, err := v.db.GetLang(chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			v.SetLang(chatID, defaultLanguage)
			return defaultLanguage
		}
		err = fmt.Errorf("vault.GetLang: %w", err)
		v.logger.Warn(err.Error())
		return defaultLanguage
	}
	return l
}

// SetLang sets the language of the user.
func (v *Vault) SetLang(chatID int64, lang string) {
	err := v.db.SetLang(chatID, lang)
	if err != nil {
		err = fmt.Errorf("vault.SetLang: %w", err)
		v.logger.Warn(err.Error())
	}
}

// Encrypt encrypts the text.
func (v *Vault) Encrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	} else if len(text) < aes.BlockSize {
		text += strings.Repeat(" ", aes.BlockSize-len(text))
	}

	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		err = fmt.Errorf("io.ReadFull: %w", err)
		v.logger.Warn(err.Error())
		return "", err
	}

	stream := cipher.NewCFBEncrypter(v.cipher, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts the text.
func (v *Vault) Decrypt(text string) (string, error) {
	if text == "" || len(text) < aes.BlockSize {
		return text, nil
	}

	ciphertext, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		err = fmt.Errorf("base64.RawStdEncoding.DecodeString: %w", err)
		v.logger.Warn(err.Error())
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(v.cipher, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return strings.Trim(string(ciphertext), " "), nil
}

// Hash hashes the text.
func (v *Vault) Hash(text string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		err = fmt.Errorf("hash.Write: %w", err)
		v.logger.Warn(err.Error())
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil)), nil
}
