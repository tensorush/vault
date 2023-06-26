package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"vault/internal/db/postgres"
	"vault/internal/db/queries"
	"vault/internal/item"
)

// Store is an interface that allows to use different databases.
type Store interface {
	Save(chatID int64, service string, secret item.Credentials) error
	Get(chatID int64, service string) (item.Credentials, error)
	Delete(chatID int64, service string) error
	GetLang(chatID int64) (string, error)
	SetLang(chatID int64, lang string) error
}

// DB is a struct that contains all methods for working with user services.
type DB struct {
	ramStore  *sync.Map
	store     Store
	langStore *sync.Map
}

// ErrServiceNotFound is returned when user service is not found.
var ErrServiceNotFound = errors.New("not found")

// New DB constructor.
func New(dataSrcName string) (*DB, error) {
	var rs Store

	db, err := sql.Open("postgres", dataSrcName)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	rs, err = postgres.New(db, "file://internal/db/migrations")
	if err != nil {
		return nil, fmt.Errorf("new postgres: %w", err)
	}

	err = queries.Prepare(db, "postgres")
	if err != nil {
		return nil, fmt.Errorf("prepare db: %w", err)
	}

	return &DB{
		ramStore:  &sync.Map{},
		langStore: &sync.Map{},
		store:     rs,
	}, nil
}

// Save saves user service
func (s *DB) Save(chatID int64, service string, secret item.Credentials) error {
	us, err := s.getUserStore(chatID)
	if err != nil && !errors.Is(err, ErrServiceNotFound) {
		return err
	}

	us.Store(service, secret)
	return s.store.Save(chatID, service, secret)
}

func (s *DB) getUserStore(chatID int64) (*sync.Map, error) {
	us, _ := s.ramStore.LoadOrStore(chatID, &sync.Map{})

	db, ok := us.(*sync.Map)
	if !ok {
		log.Println("db is not *sync.Map")
		return nil, ErrServiceNotFound
	}

	return db, nil
}

// Get gets user service
func (s *DB) Get(chatID int64, service string) (item.Credentials, error) {
	us, err := s.getUserStore(chatID)
	if err != nil {
		if !errors.Is(err, ErrServiceNotFound) {
			return item.Credentials{}, err
		}

		p, err := s.store.Get(chatID, service)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return item.Credentials{}, ErrServiceNotFound
			}
			return item.Credentials{}, err
		}
		return p, nil
	}

	value, ok := us.Load(service)
	if !ok {
		p, err := s.store.Get(chatID, service)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return item.Credentials{}, ErrServiceNotFound
			}
			return item.Credentials{}, err
		}
		return p, nil
	}

	cred, ok := value.(item.Credentials)
	if !ok {
		return item.Credentials{}, ErrServiceNotFound
	}

	return cred, nil
}

// Delete deletes user service
func (s *DB) Delete(chatID int64, serviceName string) error {
	us, err := s.getUserStore(chatID)
	if err != nil {
		return err
	}

	us.Delete(serviceName)
	err = s.store.Delete(chatID, serviceName)
	if err != nil {
		if errors.Is(err, ErrServiceNotFound) {
			return ErrServiceNotFound
		}
		return fmt.Errorf("store delete: %w", err)
	}
	return nil
}

// GetLang gets user language
func (s *DB) GetLang(chatID int64) (string, error) {
	lang, loaded := s.langStore.LoadOrStore(chatID, "en")
	if !loaded {
		lang, err := s.store.GetLang(chatID)
		if err != nil {
			return "", fmt.Errorf("get lang: %w", err)
		}
		s.langStore.Store(chatID, lang)
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
	s.langStore.Store(chatID, lang)
	err := s.store.SetLang(chatID, lang)
	if err != nil {
		return fmt.Errorf("set lang: %w", err)
	}
	return nil
}
