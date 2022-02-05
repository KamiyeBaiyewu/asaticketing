package storage

import (
	"io"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/cache"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/storage/database"
)

// Storage is a closer interface
type Storage interface {
	//database.Database
	//  cache.Cache
	io.Closer
}

type storage struct {
	db    database.Database
	Cache cache.Cache
}

// NewStorage - creates an new instance of storage
func NewStorage() (Storage, error) {
	return &storage{}, nil

	/* 	// Connecct to the Database
	   	db, err := database.New()
	   	if err != nil {

	   		return nil, errors.Wrap(err, "Error Connecting to Database")
	   	}

	   	 redisCache, err := cache.New()
	   	if err != nil {
	   		return nil, errors.Wrap(err, "Error Connecting to Cache")
	   	}

	   	return &storage{db: db}, nil */
}

func (s *storage) DB() database.Database {
	return s.db
}
func (s *storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	/* if err := s.Cache.Close(); err != nil {
		return err
	} */

	return nil
}
