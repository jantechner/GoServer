package main

import (
	"./homepage"
	"./server"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

var (
	Addr          = os.Getenv("SERVICE_ADDR")
	RequestsLimit, _ = strconv.Atoi(os.Getenv("REQUESTS_LIMIT"))
)

func requestsCounter(cancel context.CancelFunc, logger *log.Logger, requestCounterChan chan int) {
	counter := 0
	for _ = range requestCounterChan {
		counter++
		logger.Println("requests counter: ", counter)
		if counter == RequestsLimit {
			cancel()
		}
	}
}

func handleShutdown(ctx context.Context, srv *http.Server, wg *sync.WaitGroup) {
	<-ctx.Done()
	wg.Add(1)
	ctxShutDown, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
	wg.Done()
}

func main() {
	var (
		osInterruptChan    = make(chan os.Signal, 1)
		requestCounterChan = make(chan int)
		ctx, cancel        = context.WithCancel(context.Background())
		wg                 = &sync.WaitGroup{}
		logger             = log.New(os.Stdout, "Golang server - ", log.LstdFlags)
	)

	signal.Notify(osInterruptChan, os.Interrupt)

	go func() {
		s := <-osInterruptChan
		fmt.Printf("Signal %v\n", s)
		cancel()
	}()

	go requestsCounter(cancel, logger, requestCounterChan)

	h := homepage.NewHandlers(logger, requestCounterChan)
	mux := http.NewServeMux()
	h.SetupRoutes(mux)
	srv := server.New(mux, Addr)

	go handleShutdown(ctx, srv, wg)

	logger.Println("server starting")

	if err := srv.ListenAndServe(); err != nil {
		switch err {
		case http.ErrServerClosed:
			logger.Printf("%v\n", err)
		default:
			logger.Printf("server failed to start %v\n", err)
		}
	}

	close(requestCounterChan)
	wg.Wait()
	logger.Println("Gracefully shutted down")
}
