package repository

import (
	"context"
	"testing"
	"time"

	pgxmock "github.com/pashagolub/pgxmock/v4"

	"github.com/lolwierd/weatherboy/be/internal/db"
	"github.com/lolwierd/weatherboy/be/internal/model"
)

func setupMock(t *testing.T) pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	db.SetDBDriver(&db.Driver{ConnPool: mock})
	return mock
}

func TestInsertBulletin(t *testing.T) {
	mock := setupMock(t)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO bulletin").
		WithArgs("vadodara", pgxmock.AnyArg(), "hi").
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))
	mock.ExpectCommit()

	b := &model.Bulletin{Location: "vadodara", IssuedAt: time.Now(), Text: "hi"}
	if err := InsertBulletin(context.Background(), b); err != nil {
		t.Fatalf("insert: %v", err)
	}

	if b.ID != 1 {
		t.Fatalf("expected ID 1 got %d", b.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestInsertBulletinRaw(t *testing.T) {
	mock := setupMock(t)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO bulletin_raw").
		WithArgs("/tmp/a.pdf", pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	br := &model.BulletinRaw{Path: "/tmp/a.pdf", FetchedAt: time.Now()}
	if err := InsertBulletinRaw(context.Background(), br); err != nil {
		t.Fatalf("insert raw: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestInsertIMDAPICall(t *testing.T) {
	mock := setupMock(t)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO imd_api_log").
		WithArgs("https://example.com", int64(123), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	call := &model.IMDAPICall{Endpoint: "https://example.com", Bytes: 123, RequestedAt: time.Now()}
	if err := InsertIMDAPICall(context.Background(), call); err != nil {
		t.Fatalf("insert api log: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
