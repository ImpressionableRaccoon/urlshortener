package postgres

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
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

	err = st.doMigrate(dsn)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (st *PsqlStorage) doMigrate(dsn string) error {
	m, err := migrate.New("file://migrations/postgres", dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (st *PsqlStorage) Add(ctx context.Context, url repositories.URL, userID repositories.User) (id repositories.ID, err error) {
	ctxLocal, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var res pgconn.CommandTag

	for {
		id, err = utils.GenRandomID()
		if err != nil {
			log.Printf("generate id failed: %v", err)
			return "", err
		}

		res, err = st.db.Exec(ctxLocal,
			"INSERT INTO links (id, url, user_id) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
			id, url, userID)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := st.db.QueryRow(ctxLocal, `SELECT id FROM links WHERE url = $1`, url)
			err = row.Scan(&id)
			if err != nil {
				log.Printf("query failed: %v", err)
				return "", err
			}
			return id, repositories.ErrURLAlreadyExists
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

func (st *PsqlStorage) Get(ctx context.Context, id repositories.ID) (url repositories.URL, deleted bool, err error) {
	ctxLocal, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	row := st.db.QueryRow(ctxLocal, `SELECT url, deleted FROM links WHERE id = $1`, id)
	err = row.Scan(&url, &deleted)
	if err != nil {
		log.Printf("query failed: %v", err)
	}

	return
}

func (st *PsqlStorage) GetUserLinks(ctx context.Context, user repositories.User) (data []repositories.UserLink, err error) {
	ctxLocal, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := st.db.Query(ctxLocal, `SELECT id, url FROM links WHERE user_id = $1 AND deleted = FALSE`, user)
	if err != nil {
		log.Printf("query failed: %v", err)
		return nil, err
	}

	data = make([]repositories.UserLink, 0)

	for rows.Next() {
		link := repositories.UserLink{}
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
	ctxLocal, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return st.db.Ping(ctxLocal) == nil
}

func (st *PsqlStorage) DeleteUserLinks(ctx context.Context, ids []repositories.ID, user repositories.User) error {
	ctxLocal, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := st.db.Exec(ctxLocal, `UPDATE links SET deleted = TRUE WHERE id = ANY($1) AND user_id = $2`, pq.Array(ids), user)
	return err
}
