package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"vault-bot/internal/database/postgres"
	"vault-bot/internal/database/queries"
	"vault-bot/internal/database/sqlite"
	"vault-bot/internal/secret"
)

// Storage is an interface that allows to use different databases.
type Storage interface {
	Save(chatID int64, service string, secret secret.Credentials) error
	Get(chatID int64, service string) (secret.Credentials, error)
	Delete(chatID int64, service string) error
	GetLang(chatID int64) (string, error)
	SetLang(chatID int64, lang string) error
}

// DB is a struct that contains all methods for working with user services.
type DB struct {
	ramStorage  *sync.Map
	storage     Storage
	langStorage *sync.Map
}

// ErrServiceNotFound is returned when user service is not found.
var ErrServiceNotFound = errors.New("not found")

// New DB constructor.
func New(dbType, dataSrcName string) (*DB, error) {
	var rs Storage

	switch dbType {
	case "postgres":
		db, err := sql.Open("postgres", dataSrcName)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}

		rs, err = postgres.New(db, "file://schema/postgres")
		if err != nil {
			return nil, fmt.Errorf("new postgres: %w", err)
		}

		err = queries.Prepare(db, "postgres")
		if err != nil {
			return nil, fmt.Errorf("prepare db: %w", err)
		}
	case "sqlite", "test":
		db, err := sql.Open("sqlite", dataSrcName)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}

		if dbType == "test" {
			rs, err = sqlite.New(db, "file://../../schema/sqlite")
			if err != nil {
				return nil, fmt.Errorf("new sqlite: %w", err)
			}
		} else {
			rs, err = sqlite.New(db, "file://schema/sqlite")
			if err != nil {
				return nil, fmt.Errorf("new sqlite: %w", err)
			}
		}

		err = queries.Prepare(db, "sqlite")
		if err != nil {
			return nil, fmt.Errorf("prepare db: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown db type: %s", dbType)
	}

	return &DB{
		ramStorage:  &sync.Map{},
		langStorage: &sync.Map{},
		storage:     rs,
	}, nil
}

// Save saves user service
func (s *DB) Save(chatID int64, service string, secret secret.Credentials) error {
	us, err := s.getUserStorage(chatID)
	if err != nil && !errors.Is(err, ErrServiceNotFound) {
		return err
	}

	us.Store(service, secret)
	return s.storage.Save(chatID, service, secret)
}

func (s *DB) getUserStorage(chatID int64) (*sync.Map, error) {
	us, _ := s.ramStorage.LoadOrStore(chatID, &sync.Map{})

	db, ok := us.(*sync.Map)
	if !ok {
		log.Println("db is not *sync.Map")
		return nil, ErrServiceNotFound
	}

	return db, nil
}

// Get gets user service
func (s *DB) Get(chatID int64, service string) (secret.Credentials, error) {
	us, err := s.getUserStorage(chatID)
	if err != nil {
		if !errors.Is(err, ErrServiceNotFound) {
			return secret.Credentials{}, err
		}

		p, err := s.storage.Get(chatID, service)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return secret.Credentials{}, ErrServiceNotFound
			}
			return secret.Credentials{}, err
		}
		return p, nil
	}

	value, ok := us.Load(service)
	if !ok {
		p, err := s.storage.Get(chatID, service)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return secret.Credentials{}, ErrServiceNotFound
			}
			return secret.Credentials{}, err
		}
		return p, nil
	}

	cred, ok := value.(secret.Credentials)
	if !ok {
		return secret.Credentials{}, ErrServiceNotFound
	}

	return cred, nil
}

// Delete deletes user service
func (s *DB) Delete(chatID int64, serviceName string) error {
	us, err := s.getUserStorage(chatID)
	if err != nil {
		return err
	}

	us.Delete(serviceName)
	err = s.storage.Delete(chatID, serviceName)
	if err != nil {
		if errors.Is(err, ErrServiceNotFound) {
			return ErrServiceNotFound
		}
		return fmt.Errorf("storage delete: %w", err)
	}
	return nil
}

// GetLang gets user language
func (s *DB) GetLang(chatID int64) (string, error) {
	lang, loaded := s.langStorage.LoadOrStore(chatID, "en")
	if !loaded {
		lang, err := s.storage.GetLang(chatID)
		if err != nil {
			return "", fmt.Errorf("get lang: %w", err)
		}
		s.langStorage.Store(chatID, lang)
		return lang, nil
	}

	l, ok := lang.(string)
	if !ok {
		return "", ErrServiceNotFound
	}

	return l, nil
}

// SetLang sets user language
func (s *DB) SetLang(chatID int64, lang string) error {
	s.langStorage.Store(chatID, lang)
	err := s.storage.SetLang(chatID, lang)
	if err != nil {
		return fmt.Errorf("set lang: %w", err)
	}
	return nil
}
