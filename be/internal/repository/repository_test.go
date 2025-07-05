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



func TestInsertBulletinRaw(t *testing.T) {
	mock := setupMock(t)
	defer mock.Close()

	mock.ExpectQuery("INSERT INTO bulletin_raw").
		WithArgs("/tmp/a.pdf", pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

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

func TestInsertNowcastRaw(t *testing.T) {
	mock := setupMock(t)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO nowcast_raw").
		WithArgs("vadodara", pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	nr := &model.NowcastRaw{Location: "vadodara", Data: []byte("{}"), FetchedAt: time.Now()}
	if err := InsertNowcastRaw(context.Background(), nr); err != nil {
		t.Fatalf("insert raw: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestInsertNowcastCategory(t *testing.T) {
	mock := setupMock(t)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO nowcast_category").
		WithArgs(1, 2, int16(3)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	c := &model.NowcastCategory{NowcastID: 1, Category: 2, Value: 3};
	if err := InsertNowcastCategory(context.Background(), c); err != nil {
		t.Fatalf("insert cat: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}