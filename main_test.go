package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_assert "github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestRouting(t *testing.T) {
	t.Run("GET /todos -> body 200", func(t *testing.T) {
		assert := _assert.New(t)

		db, mock := NewMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"}).
			AddRow(1, "314fb35b-f8ca-45d7-96ff-ed586a5ee8ae", "my first todo in db", false, "123")

		mock.ExpectQuery("^SELECT (.+) FROM todos").WillReturnRows(rows)

		SetupMiddleware(db)
		SetupRoutes()

		srv := httptest.NewServer(Router)
		defer srv.Close()

		res, err := http.Get(srv.URL + "/todos")
		if err != nil {
			t.Fatalf("error in response: %v", err)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("could not read response body : %v", err)
		}

		expected := `[{"uuid":"314fb35b-f8ca-45d7-96ff-ed586a5ee8ae","text":"my first todo in db","completed":false,"created_on":"123"}]`

		assert.Equal(expected, string(b))
	})
}
