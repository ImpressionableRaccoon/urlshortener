package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgerrcode"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	err = st.doMigrate()
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (st *PsqlStorage) doMigrate() error {
	m, err := migrate.New("file://migrations/postgres", configs.DatabaseDSN)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (st *PsqlStorage) Add(ctx context.Context, url URL, userID User) (id ID, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var res pgconn.CommandTag

	for {
		id, err = utils.GenRandomID()
		if err != nil {
			log.Printf("generate id failed: %v", err)
			return "", err
		}

		res, err = st.db.Exec(ctx,
			"INSERT INTO links (id, url, user_id) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
			id, url, userID)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := st.db.QueryRow(ctx, `SELECT id FROM links WHERE url = $1`, url)
			err = row.Scan(&id)
			if err != nil {
				log.Printf("query failed: %v", err)
				return "", err
			}
			return id, ErrURLAlreadyExists
		}
		if err != nil {
			log.Printf("exec failed: %v", err)
			return "", err
		}

		if res.RowsAffected() == 1 {
			break
		}
	}

	return id, err
}

func (st *PsqlStorage) Get(ctx context.Context, id ID) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var url URL
	row := st.db.QueryRow(ctx, `SELECT url FROM links WHERE id = $1`, id)
	err := row.Scan(&url)
	if err != nil {
		log.Printf("query failed: %v", err)
	}

	return url, err
}

func (st *PsqlStorage) GetUserLinks(ctx context.Context, user User) (data []UserLink, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := st.db.Query(ctx, `SELECT id, url FROM links WHERE user_id = $1`, user)
	if err != nil {
		log.Printf("query failed: %v", err)
		return nil, err
	}

	data = make([]UserLink, 0)

	for rows.Next() {
		link := UserLink{}
		err = rows.Scan(&link.ID, &link.URL)
		if err != nil {
			log.Printf("row scan failed: %v", err)
			return nil, err
		}
		data = append(data, link)
	}

	return data, nil
}

func (st *PsqlStorage) Pool(ctx context.Context) bool {
	return st.db.Ping(ctx) == nil
}
