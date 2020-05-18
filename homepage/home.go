package homepage

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const message = "Hello World!"

type Handlers struct {
	logger *log.Logger
	notifier chan<- int
}

func (h *Handlers) Home(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	time.Sleep(1 * time.Millisecond)
	_, _ = writer.Write([]byte(message))
}

func (h *Handlers) Counter(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h.notifier <- 1
		next(writer, request)
	}
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		next(writer, request)
		h.logger.Printf("request processed in %s\n", time.Since(startTime))
	}
}

func (h *Handlers) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", h.Logger(h.Counter(h.Home)))
}

func NewHandlers(logger *log.Logger, notifier chan<- int) *Handlers {
	return &Handlers{
		logger: logger,
		notifier: notifier,
	}
}
