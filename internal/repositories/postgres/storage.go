// Package postgres содержит хранилище интерфейса Storager для взаимодействия с базой данных Postgres.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres init for golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // file init for golang-migrate
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

const (
	deleteBufferSize    = 100
	deleteBufferTimeout = time.Second
	shutdownTimeout     = 15 * time.Second
)

// PsqlStorage - структура для хранилища Postgres.
type PsqlStorage struct {
	db             *sql.DB
	deleteCh       chan repositories.LinkData
	deleteWg       sync.WaitGroup
	deleteShutdown chan struct{}
}

// NewPsqlStorage - конструктор для PsqlStorage.
func NewPsqlStorage(dsn string) (*PsqlStorage, error) {
	st := &PsqlStorage{
		deleteCh:       make(chan repositories.LinkData),
		deleteShutdown: make(chan struct{}),
	}

	var err error
	st.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = st.doMigrate(dsn)
	if err != nil {
		return nil, err
	}

	st.deleteWg.Add(1)
	go st.deleteUserLinksWorker(context.Background(), deleteBufferSize, deleteBufferTimeout)

	return st, nil
}

// Add - сократить ссылку.
func (st *PsqlStorage) Add(
	ctx context.Context,
	url repositories.URL,
	userID repositories.User,
) (id repositories.ID, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	for {
		id, err = utils.GenRandomID()
		if err != nil {
			log.Printf("generate id failed: %v", err)
			return "", err
		}

		var res sql.Result
		res, err = st.db.ExecContext(
			ctx,
			`INSERT INTO links (id, url, user_id) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`,
			id, url, userID,
		)

		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := st.db.QueryRowContext(ctx, `SELECT id FROM links WHERE url = $1`, url)
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

		var aff int64
		aff, err = res.RowsAffected()
		if err != nil {
			return id, err
		}

		if aff == 1 {
			break
		}
	}

	return id, nil
}

// Get - получить оригинальную ссылку по ID.
func (st *PsqlStorage) Get(ctx context.Context, id repositories.ID) (url repositories.URL, deleted bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	row := st.db.QueryRowContext(
		ctx,
		`SELECT url, deleted FROM links WHERE id = $1`,
		id,
	)

	err = row.Scan(&url, &deleted)
	if err != nil {
		log.Printf("query failed: %v", err)
		return "", false, err
	}

	return url, deleted, nil
}

// GetUserLinks - получить все ссылки пользователя.
func (st *PsqlStorage) GetUserLinks(
	ctx context.Context,
	user repositories.User,
) (data []repositories.LinkData, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	data = make([]repositories.LinkData, 0)

	rows, err := st.db.QueryContext(
		ctx,
		`SELECT id, url FROM links WHERE user_id = $1 AND deleted = FALSE`,
		user,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return data, nil
	}
	if err != nil {
		log.Printf("query failed: %v", err)
		return nil, err
	}
	if rows.Err() != nil {
		log.Printf("rows failed: %v", err)
		return nil, err
	}

	for rows.Next() {
		link := repositories.LinkData{
			User:    user,
			Deleted: false,
		}
		err = rows.Scan(&link.ID, &link.URL)
		if err != nil {
			log.Printf("row scan failed: %v", err)
			return nil, err
		}
		data = append(data, link)
	}

	return data, nil
}

// DeleteUserLinks - удалить ссылки пользователя.
func (st *PsqlStorage) DeleteUserLinks(_ context.Context, ids []repositories.ID, user repositories.User) error {
	for _, id := range ids {
		st.deleteCh <- repositories.LinkData{ID: id, User: user}
	}

	return nil
}

// GetStats - получить статистику сервиса.
func (st *PsqlStorage) GetStats(ctx context.Context) (repositories.ServiceStats, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stats := repositories.ServiceStats{}

	row := st.db.QueryRowContext(ctx,
		`SELECT COUNT(id) AS links_count, COUNT(DISTINCT user_id) AS users_count FROM links`,
	)

	err := row.Scan(&stats.URLs, &stats.Users)
	if err != nil {
		log.Printf("query failed: %v", err)
		return stats, err
	}

	return stats, nil
}

// Pool - проверить соединение с базой данных.
func (st *PsqlStorage) Pool(ctx context.Context) (ok bool) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return st.db.PingContext(ctx) == nil
}

// Close - мягко завершить работу хранилища.
func (st *PsqlStorage) Close(_ context.Context) error {
	close(st.deleteShutdown)

	c := make(chan struct{})
	go func() {
		defer close(c)
		st.deleteWg.Wait()
	}()
	select {
	case <-c:
	case <-time.After(shutdownTimeout):
		log.Print("storage close timeout exceed")
	}

	return st.db.Close()
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

func (st *PsqlStorage) deleteUserLinksWorker(ctx context.Context, bufferSize int, bufferTimeout time.Duration) {
	ids := make([]repositories.ID, 0, bufferSize)
	users := make([]repositories.User, 0, bufferSize)

worker:
	for {
		select {
		case <-st.deleteShutdown:
			break worker
		case <-ctx.Done():
			break worker
		default:
		}

		ids = ids[:0]
		users = users[:0]

		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, bufferTimeout)

	loop:
		for {
			select {
			case v := <-st.deleteCh:
				ids = append(ids, v.ID)
				users = append(users, v.User)
				if len(ids) == bufferSize {
					timeoutCancel()
					break loop
				}
			case <-timeoutCtx.Done():
				timeoutCancel()
				break loop
			case <-st.deleteShutdown:
				timeoutCancel()
				break loop
			}
		}

		if len(ids) == 0 {
			continue
		}

		ctxLocal, cancelLocal := context.WithTimeout(ctx, time.Second*10)

		_, err := st.db.ExecContext(
			ctxLocal,
			`UPDATE links SET deleted = TRUE 
             FROM (select unnest($1::text[]) as id, unnest($2::uuid[]) as user) as data_table
             WHERE links.id = data_table.id AND user_id = data_table.user`,
			pq.Array(ids), pq.Array(users),
		)
		if err != nil {
			log.Printf("update failed: %v", err)
		}

		cancelLocal()
	}

	st.deleteWg.Done()
}
