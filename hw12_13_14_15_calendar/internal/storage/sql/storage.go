package sqlstorage

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
