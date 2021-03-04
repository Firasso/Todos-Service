package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	types "github.com/Firasso/DemoGoApi/types"
	_assert "github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetTodos(t *testing.T) {
	t.Run("returns the response body and 200 status code", func(t *testing.T) {
		assert := _assert.New(t)
		db, mock := NewMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"}).
			AddRow(1, "314fb35b-f8ca-45d7-96ff-ed586a5ee8ae", "my first todo in db", false, "123").
			AddRow(2, "3fa29283-15c9-49d6-b847-c030bc85c614", "my second todo in db", true, "12345")

		mock.ExpectQuery("^SELECT (.+) FROM todos").
			WillReturnRows(rows)

		ctx := context.WithValue(context.Background(), types.DBContext, db)

		req, err := http.NewRequestWithContext(ctx, "GET", "localhost:8000/todos", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rec := httptest.NewRecorder()

		GetTodos(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("could not read response body : %v", err)
		}

		expected := `[{"uuid":"314fb35b-f8ca-45d7-96ff-ed586a5ee8ae","text":"my first todo in db","completed":false,"created_on":"123"},{"uuid":"3fa29283-15c9-49d6-b847-c030bc85c614","text":"my second todo in db","completed":true,"created_on":"12345"}]`
		assert.Equal(expected, string(b))
		assert.Equal(http.StatusOK, res.StatusCode)
	})

	t.Run("returns error status code 500 when DB error", func(t *testing.T) {
		assert := _assert.New(t)
		db, mock := NewMock()
		defer db.Close()

		expected := "db internal error"

		mock.ExpectQuery("^SELECT (.+) FROM todos").
			WillReturnError(fmt.Errorf(expected))

		ctx := context.WithValue(context.Background(), types.DBContext, db)

		req, err := http.NewRequestWithContext(ctx, "GET", "localhost:8000/todos", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rec := httptest.NewRecorder()

		GetTodos(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("could not read response body : %v", err)
		}

		assert.Equal(http.StatusInternalServerError, res.StatusCode)
		assert.Equal(expected, strings.TrimSpace(string(b)))
	})

	t.Run("returns error status code 500 when DB not in context", func(t *testing.T) {
		assert := _assert.New(t)

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, "GET", "localhost:8000/todos", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rec := httptest.NewRecorder()

		GetTodos(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("could not read response body : %v", err)
		}

		expected := "db context missing"
		assert.Equal(http.StatusInternalServerError, res.StatusCode)
		assert.Equal(expected, strings.TrimSpace(string(b)))
	})
}
