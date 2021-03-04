package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"

	_assert "github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var todos Todos

func TestGetAll(t *testing.T) {
	t.Run("returns all todos", func(t *testing.T) {
		assert := _assert.New(t)

		db, mock := NewMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"}).
			AddRow(1, "314fb35b-f8ca-45d7-96ff-ed586a5ee8ae", "my first todo in db", false, "123").
			AddRow(2, "3fa29283-15c9-49d6-b847-c030bc85c614", "my second todo in db", true, "12345")

		mock.ExpectQuery("^SELECT (.+) FROM todos").
			WillReturnRows(rows)

		got, err := todos.GetAll(db)

		var d []Todo
		if e := json.Unmarshal(got, &d); e != nil {
			t.Fatalf("Unable to parse return value of todos.GetAll %q into slice of Todo, '%v'", got, e)
		}

		var expected = []Todo{
			{0, uuid.FromStringOrNil("314fb35b-f8ca-45d7-96ff-ed586a5ee8ae"), "my first todo in db", false, "123"},
			{0, uuid.FromStringOrNil("3fa29283-15c9-49d6-b847-c030bc85c614"), "my second todo in db", true, "12345"},
		}

		assert.Nil(err)
		assert.Equal(expected, d)
	})

	t.Run("returns nil if no rows in the DB", func(t *testing.T) {
		assert := _assert.New(t)
		db, mock := NewMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"})

		mock.ExpectQuery("^SELECT (.+) FROM todos").
			WillReturnRows(rows)

		got, err := todos.GetAll(db)
		var d []Todo
		if e := json.Unmarshal(got, &d); e != nil {
			t.Fatalf("Unable to parse return value of todos.GetAll %q into slice of Todo, '%v'", got, e)
		}

		assert.Nil(err)
		assert.Nil(d)
	})

	t.Run("returns error if no DB", func(t *testing.T) {
		assert := _assert.New(t)
		db, mock := NewMock()
		defer db.Close()

		sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"})

		mock.ExpectQuery("^SELECT (.+) FROM todos").
			WillReturnError(fmt.Errorf("db internal error"))

		got, err := todos.GetAll(db)

		assert.Nil(got)
		assert.NotNil(err)
	})
}

func TestAddTodo(t *testing.T) {
	t.Run("adds todo to DB", func(t *testing.T) {
		var todo = NewTodo{"random-string", "new todo", true}
		assert := _assert.New(t)

		db, mock := NewMock()
		defer db.Close()

		sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"})

		mock.ExpectExec("INSERT INTO todos (.*)").
			WithArgs(todo.Text, sqlmock.AnyArg(), todo.Completed).
			WillReturnResult(sqlmock.NewResult(1, 1))

		success, err := todos.AddTodo(todo, db)

		assert.Nil(err)
		assert.True(success)
	})

	t.Run("return error if cannot add todo to DB", func(t *testing.T) {
		var todo = NewTodo{"random-string", "new todo", true}
		assert := _assert.New(t)

		db, mock := NewMock()
		defer db.Close()

		sqlmock.NewRows([]string{"Id", "uuid", "text", "completed", "created_on"})

		mock.ExpectExec("INSERT INTO todos (.*)").
			WithArgs(todo.Text, sqlmock.AnyArg(), todo.Completed).
			WillReturnError(fmt.Errorf("db internal error"))

		success, err := todos.AddTodo(todo, db)

		assert.NotNil(err)
		assert.False(success)
	})
}
