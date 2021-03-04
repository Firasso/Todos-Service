package types

import (
	"database/sql"
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Todo struct {
	id        int
	UUID      uuid.UUID `json:"uuid"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedOn string    `json:"created_on"`
}

type NewTodo struct {
	UUID      string
	Text      string
	Completed bool
}

type Todos struct {
}

func (t Todos) GetAll(db *sql.DB) ([]byte, error) {
	rows, err := db.Query("SELECT * FROM todos;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo

		err = rows.Scan(&todo.id, &todo.UUID, &todo.Text, &todo.Completed, &todo.CreatedOn)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(todos)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func (t *Todos) AddTodo(a NewTodo, db *sql.DB) (bool, error) {
	id := uuid.NewV4()
	_, err := db.Exec("INSERT INTO todos (text,uuid,completed) VALUES ($1,$2,$3)", a.Text, id, a.Completed)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t Todos) GetTodoByID(id string, db *sql.DB) (*Todo, error) {
	var todo Todo

	row := db.QueryRow("SELECT * FROM todos WHERE uuid=$1", id)

	switch err := row.Scan(&todo.id, &todo.UUID, &todo.Text, &todo.Completed, &todo.CreatedOn); err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		return &todo, nil
	default:
		panic(err)
	}
}

func (t Todos) UpdateTodoByID(n NewTodo, db *sql.DB) (bool, error) {
	_, err := db.Exec("UPDATE todos SET text=$1, completed=$2 WHERE uuid=$3;", n.Text, n.Completed, n.UUID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *Todos) RemoveTodoByID(id string, db *sql.DB) (bool, error) {
	_, err := db.Exec("DELETE FROM todos WHERE uuid=$1;", id)
	if err != nil {
		return false, err
	}

	return true, nil
}
