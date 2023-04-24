package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
)

func TestPsqlStorage_Add(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://microsoft.com"
		userID := uuid.New()

		st := &PsqlStorage{db: db}
		ctx := context.Background()

		mock.ExpectExec("INSERT").
			WithArgs(sqlmock.AnyArg(), url, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = st.Add(ctx, url, userID)
		assert.NoError(t, err)
	})

	t.Run("id already exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://github.com"
		userID := uuid.New()

		st := &PsqlStorage{db: db}
		ctx := context.Background()

		mock.ExpectExec("INSERT").
			WithArgs(sqlmock.AnyArg(), url, userID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		mock.ExpectExec("INSERT").
			WithArgs(sqlmock.AnyArg(), url, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = st.Add(ctx, url, userID)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("url already exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://google.com"
		userID := uuid.New()
		id := "linkID"

		st := &PsqlStorage{db: db}
		ctx := context.Background()

		mock.ExpectExec("INSERT").
			WithArgs(sqlmock.AnyArg(), url, userID).
			WillReturnError(&pq.Error{Code: pgerrcode.UniqueViolation})

		rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
		mock.ExpectQuery("SELECT id").
			WithArgs(url).
			WillReturnRows(rows)

		linkID, err := st.Add(ctx, url, userID)
		assert.ErrorIs(t, err, repositories.ErrURLAlreadyExists)
		assert.Equal(t, id, linkID)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("unique violation but no url", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://firefox.org"
		userID := uuid.New()

		st := &PsqlStorage{db: db}
		ctx := context.Background()

		mock.ExpectExec("INSERT").
			WithArgs(sqlmock.AnyArg(), url, userID).
			WillReturnError(&pq.Error{Code: pgerrcode.UniqueViolation})

		rows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT id").
			WithArgs(url).
			WillReturnRows(rows)

		_, err = st.Add(ctx, url, userID)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("already exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://apple.com"
		userID := uuid.New()

		st := &PsqlStorage{db: db}
		ctx := context.Background()

		mock.ExpectExec("INSERT").
			WithArgs(sqlmock.AnyArg(), url, userID).
			WillReturnError(errors.New("test"))

		_, err = st.Add(ctx, url, userID)
		assert.Error(t, err)
	})
}

func TestPsqlStorage_Get(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://imgur.com"
		id := "imgur"

		rows := sqlmock.NewRows([]string{"url", "deleted"}).AddRow(url, false)
		mock.ExpectQuery("SELECT url, deleted").
			WithArgs(id).
			WillReturnRows(rows)

		st := &PsqlStorage{db: db}
		ctx := context.Background()
		res, deleted, err := st.Get(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, url, res)
		assert.False(t, deleted)
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		id := "ajslx"

		mock.ExpectQuery("SELECT url, deleted").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		st := &PsqlStorage{db: db}
		ctx := context.Background()
		_, deleted, err := st.Get(ctx, id)

		assert.Error(t, err, sql.ErrNoRows)
		assert.False(t, deleted)
	})

	t.Run("deleted", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		url := "https://stackoverflow.com"
		id := "stack"

		rows := sqlmock.NewRows([]string{"url", "deleted"}).AddRow(url, true)
		mock.ExpectQuery("SELECT url, deleted").
			WithArgs(id).
			WillReturnRows(rows)

		st := &PsqlStorage{db: db}
		ctx := context.Background()
		longURL, deleted, err := st.Get(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, url, longURL)
		assert.True(t, deleted)
	})
}

func TestPsqlStorage_GetUserLinks(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		data := map[repositories.ID]repositories.URL{
			"smdlx": "https://impressionableracoob.com/smdlx",
			"dmxsl": "https://impressionableracoob.com/dmxsl",
			"mskls": "https://impressionableracoob.com/mskls",
		}

		rows := sqlmock.NewRows([]string{"id", "url"})
		for k, v := range data {
			rows = rows.AddRow(k, v)
		}

		mock.ExpectQuery("SELECT id, url").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(rows)

		st := &PsqlStorage{db: db}
		ctx := context.Background()
		links, err := st.GetUserLinks(ctx, uuid.New())

		assert.NoError(t, err)
		assert.Len(t, links, len(data))
		for _, link := range links {
			assert.Equal(t, data[link.ID], link.URL)
		}
	})

	t.Run("no links", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectQuery("SELECT id, url").
			WithArgs(sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		st := &PsqlStorage{db: db}
		ctx := context.Background()
		links, err := st.GetUserLinks(ctx, uuid.New())

		assert.NoError(t, err)
		assert.Len(t, links, 0)
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		rows := sqlmock.NewRows([]string{"id"}).AddRow("link2")
		mock.ExpectQuery("SELECT id, url").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(rows)

		st := &PsqlStorage{db: db}
		ctx := context.Background()
		_, err = st.GetUserLinks(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestPsqlStorage_DeleteUserLinks(t *testing.T) {
	st := &PsqlStorage{
		deleteCh: make(chan repositories.LinkData),
	}

	ctx := context.Background()
	ids := []repositories.ID{"link1", "link2", "link3"}
	userID := uuid.New()

	go func() {
		for i := 0; i < len(ids); i++ {
			link := <-st.deleteCh
			assert.Equal(t, link.ID, ids[i])
			assert.Equal(t, link.User, userID)
		}
	}()

	err := st.DeleteUserLinks(ctx, ids, userID)
	assert.NoError(t, err)
}

func TestPsqlStorage_GetStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	exp := repositories.ServiceStats{
		URLs:  10,
		Users: 5,
	}

	rows := sqlmock.NewRows([]string{"links_count", "users_count"}).AddRow(exp.URLs, exp.Users)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	st := &PsqlStorage{db: db}
	ctx := context.Background()
	stats, err := st.GetStats(ctx)

	assert.NoError(t, err)
	assert.Equal(t, exp, stats)
}

func TestPsqlStorage_Pool(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectPing()

		st := &PsqlStorage{db: db}

		ctx := context.Background()
		ok := st.Pool(ctx)

		assert.True(t, ok)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		st := &PsqlStorage{db: db}

		ctx := context.Background()
		ok := st.Pool(ctx)

		assert.False(t, ok)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestPsqlStorage_Close(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectClose()

		st := &PsqlStorage{
			db:             db,
			deleteWg:       sync.WaitGroup{},
			deleteShutdown: make(chan struct{}),
		}

		ctx := context.Background()
		err = st.Close(ctx)

		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("connection already closed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectClose().WillReturnError(sql.ErrConnDone)

		st := &PsqlStorage{
			db:             db,
			deleteWg:       sync.WaitGroup{},
			deleteShutdown: make(chan struct{}),
		}

		ctx := context.Background()
		err = st.Close(ctx)

		assert.ErrorIs(t, err, sql.ErrConnDone)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("delete wg timeout", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectClose()

		st := &PsqlStorage{
			db:             db,
			deleteWg:       sync.WaitGroup{},
			deleteShutdown: make(chan struct{}),
		}
		st.deleteWg.Add(1)

		ctx := context.Background()
		err = st.Close(ctx)

		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestPsqlStorage_doMigrate(t *testing.T) {
	t.Run("wrong dsn", func(t *testing.T) {
		st := &PsqlStorage{}
		err := st.doMigrate("wrongDSN")
		assert.Error(t, err)
	})
}
