package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Todo struct {
	Number  int    `json:"id"`
	Content string `json:"content"`
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	var (
		response []byte
	)
	if id, exists := mux.Vars(r)["id"]; exists {
		todo, err := load(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, _ = json.Marshal(todo)
	} else {
		todos, err := loadAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, _ = json.Marshal(todos)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	_, _ = w.Write(response)
}

func (todo *Todo) save() error {
	filename := fmt.Sprintf("%v.txt", todo.Number)
	return ioutil.WriteFile(filename, []byte(todo.Content), 600)
}

func load(idStr string) (*Todo, error) {
	filename := idStr + ".txt"
	content, err := ioutil.ReadFile("./todo/todos/" + filename)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	number, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	return &Todo{number, string(content)}, nil
}

func loadAll() (todos []Todo, err error) {
	files, err := ioutil.ReadDir("./todo/todos")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		numberStr := filename[:len(filename)-len(".txt")]
		todo, err := load(numberStr)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *todo)
	}
	return
}

func SetupRoutes(mux *mux.Router) {
	s := mux.PathPrefix("/todos").Subrouter()
	s.HandleFunc("", GetTodos)
	s.HandleFunc("/{id:[0-9]+}", GetTodos)
}
