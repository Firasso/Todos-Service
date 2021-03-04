package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	types "github.com/Firasso/DemoGoApi/types"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-chi/chi"
)

var todos = types.Todos{}

var getReqSuccessCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "get_todos_success",
		Help: "Number of total successful GET /todos requests.",
	})

func init() {
	prometheus.MustRegister(getReqSuccessCounter)
}

// A todosResponse.
// swagger:response todosResponse
type todosResponse struct {
	// All todos in db
	// in: body
	Body []types.Todo
}

//	swagger:route GET /todos todos list
// 		Returns a list of todos
// 		Responses:
//			200: todosResponse
func GetTodos(w http.ResponseWriter, req *http.Request) {
	db, ok := req.Context().Value(types.DBContext).(*sql.DB)
	if !ok {
		http.Error(w, "db context missing", http.StatusInternalServerError)
		return
	}

	rows, err := todos.GetAll(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	getReqSuccessCounter.Inc()
}

func AddTodo(w http.ResponseWriter, req *http.Request) {
	db, ok := req.Context().Value(types.DBContext).(*sql.DB)
	if !ok {
		http.Error(w, "db context missing", http.StatusInternalServerError)
		return
	}

	var t types.NewTodo
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err = todos.AddTodo(t, db)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetSingleTodo(w http.ResponseWriter, req *http.Request) {
	db, ok := req.Context().Value(types.DBContext).(*sql.DB)
	if !ok {
		http.Error(w, "db context missing", http.StatusInternalServerError)
		return
	}

	var todo *types.Todo
	var err error

	if todoID := chi.URLParam(req, "todoID"); todoID != "" {
		todo, err = todos.GetTodoByID(todoID, db)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateTodo(w http.ResponseWriter, req *http.Request) {
	db, ok := req.Context().Value(types.DBContext).(*sql.DB)
	if !ok {
		http.Error(w, "db context missing", http.StatusInternalServerError)
		return
	}

	if todoID := chi.URLParam(req, "todoID"); todoID != "" {
		var t types.NewTodo
		t.UUID = todoID

		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ok, err = todos.UpdateTodoByID(t, db)

		if !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func DeleteTodo(w http.ResponseWriter, req *http.Request) {
	db, ok := req.Context().Value(types.DBContext).(*sql.DB)
	if !ok {
		http.Error(w, "db context missing", http.StatusInternalServerError)
		return
	}

	if todoID := chi.URLParam(req, "todoID"); todoID != "" {
		ok, err := todos.RemoveTodoByID(todoID, db)
		if !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
