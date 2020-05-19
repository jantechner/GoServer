package todo

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) GetTodos(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) PostTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todo.Number = Counter.GetNewId()

	err = todo.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Printf("New todo created %v\n", todo)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		next(writer, request)
		h.logger.Printf("%s %s processed in %s\n", request.Method, request.URL.Path, time.Since(startTime))
	}
}

func (h *Handlers) SetupRoutes(mux *mux.Router) {
	s := mux.PathPrefix("/todos").Subrouter()
	s.HandleFunc("", h.Logger(h.GetTodos)).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", h.Logger(h.GetTodos))
	s.HandleFunc("", h.Logger(h.PostTodo)).Methods("POST")
}

func NewHandlers(logger *log.Logger) *Handlers {
	if err := Counter.Init(); err != nil {
		logger.Fatalln("todos counter init error", err.Error())
	}
	return &Handlers{logger}
}



