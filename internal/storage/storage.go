package storage

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	logger *logrus.Logger
	dbConn *sql.DB
}

func NewStorage(logger *logrus.Logger, dbConn *sql.DB) *Storage {
	return &Storage{
		logger: logger,
		dbConn: dbConn,
	}
}
