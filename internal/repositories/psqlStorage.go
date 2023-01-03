package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PsqlStorage struct {
	DB *pgxpool.Pool
}

func NewPsqlStorage(dsn string) (*PsqlStorage, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	st := &PsqlStorage{
		DB: db,
	}

	return st, nil
}

func (st *PsqlStorage) Add(url URL, userID User) (id ID, err error) {
	return "", nil
}

func (st *PsqlStorage) Get(id ID) (string, error) {
	return "", nil
}

func (st *PsqlStorage) GetUserLinks(user User) (data []UserLink) {
	return nil
}

func (st *PsqlStorage) IsUserExists(userID User) bool {
	return true
}

func (st *PsqlStorage) Pool() bool {
	return st.DB.Ping(context.Background()) == nil
}
