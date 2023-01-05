package repositories

import (
	"context"
	"time"

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
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	st := &PsqlStorage{
		db: db,
	}

	exists, err := st.checkIsTablesExists()
	if err != nil {
		return nil, err
	}

	if !exists {
		err = st.createTables()
		if err != nil {
			return nil, err
		}
	}

	return st, nil
}

func (st *PsqlStorage) checkIsTablesExists() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := st.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'links')`)

	var result bool

	err := row.Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (st *PsqlStorage) createTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := st.db.Exec(ctx,
		`CREATE TABLE links (
			id varchar(255) NOT NULL UNIQUE,
			url varchar(255) NOT NULL,
			user_id uuid NOT NULL)`)
	return err
}

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
