package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/pikachu0310/go-backend-template/core/internal/repository"
)

// Exercise the real repository's error wrapping without a database server or photo API.
func TestGetUserErrors(t *testing.T) {
	for _, tt := range []struct {
		name, userID, body string
		dbErr              error
		status             int
		query              bool
	}{
		{name: "missing user", userID: "00000000-0000-4000-8000-000000000001", status: http.StatusNotFound, body: "{\"message\":\"Not Found\"}\n", query: true},
		{name: "invalid UUID", userID: "invalid", status: http.StatusBadRequest, body: "{\"message\":\"invalid userID\"}\n"},
		{name: "database failure", userID: "00000000-0000-4000-8000-000000000001", dbErr: errors.New("database unavailable"), status: http.StatusInternalServerError, body: "{\"message\":\"Internal Server Error\"}\n", query: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			connector := &userErrorConnector{err: tt.dbErr}
			db := sql.OpenDB(connector)
			t.Cleanup(func() { _ = db.Close() })
			h := New(repository.New(sqlx.NewDb(db, "test")))
			e := echo.New()
			e.GET("/api/v1/users/:userID", h.GetUser)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil))
			if rec.Code != tt.status || rec.Body.String() != tt.body {
				t.Fatalf("got %d %s; want %d %s", rec.Code, rec.Body.String(), tt.status, tt.body)
			}
			if connector.queried != tt.query {
				t.Fatalf("database queried = %v; want %v", connector.queried, tt.query)
			}
		})
	}
}

type userErrorConnector struct {
	err     error
	queried bool
}

func (c *userErrorConnector) Connect(context.Context) (driver.Conn, error) {
	return &userErrorConn{c}, nil
}
func (c *userErrorConnector) Driver() driver.Driver { return userErrorDriver{} }

type userErrorDriver struct{}

func (userErrorDriver) Open(string) (driver.Conn, error) { return nil, errors.New("use connector") }

type userErrorConn struct{ c *userErrorConnector }

func (*userErrorConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected Prepare")
}
func (*userErrorConn) Close() error              { return nil }
func (*userErrorConn) Begin() (driver.Tx, error) { return nil, errors.New("unexpected Begin") }
func (c *userErrorConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	c.c.queried = true
	if c.c.err != nil {
		return nil, c.c.err
	}
	return emptyUserRows{}, nil
}

type emptyUserRows struct{}

func (emptyUserRows) Columns() []string         { return []string{"id", "name", "email"} }
func (emptyUserRows) Close() error              { return nil }
func (emptyUserRows) Next([]driver.Value) error { return io.EOF }
