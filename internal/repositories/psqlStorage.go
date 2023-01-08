package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/golang-migrate/migrate/v4"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type PsqlStorage struct {
	db *pgxpool.Pool
}

func NewPsqlStorage(dsn string) (*PsqlStorage, error) {
	st := &PsqlStorage{}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	st.db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	err = st.createTables()
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (st *PsqlStorage) createTables() error {
	m, err := migrate.New("file://migrations/postgres", configs.DatabaseDSN)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// TODO: ctx from request r.Context()

func (st *PsqlStorage) Add(url URL, userID User) (id ID, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var res pgconn.CommandTag

	for {
		id, err = utils.GetRandomID()
		if err != nil {
			return "", err
		}

		res, err = st.db.Exec(ctx,
			"INSERT INTO links (id, url, user_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			id, url, userID)
		if err != nil {
			return "", err
		}

		if res.RowsAffected() == 1 {
			break
		}
	}

	return id, err
}

func (st *PsqlStorage) Get(id ID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var url URL

	row := st.db.QueryRow(ctx, `SELECT url FROM links WHERE id = $1`, id)
	err := row.Scan(&url)
	return url, err
}

func (st *PsqlStorage) GetUserLinks(user User) (data []UserLink, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data = make([]UserLink, 0)

	rows, err := st.db.Query(ctx, `SELECT id, url FROM links WHERE user_id = $1`, user)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := UserLink{}
		err = rows.Scan(&link.ID, &link.URL)
		if err != nil {
			return nil, err
		}
		data = append(data, link)
	}

	return data, nil
}

func (st *PsqlStorage) Pool() bool {
	return st.db.Ping(context.Background()) == nil
}
