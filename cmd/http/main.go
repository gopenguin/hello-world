package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	flagAddr string
	flagHelp bool
)

func init() {
	flag.StringVar(&flagAddr, "addr", "127.0.0.1:8080", "The address and port to listen to")
	flag.BoolVar(&flagHelp, "help", false, "Print this message")
}

var (
	result []int64
)

func main() {
	flag.Parse()

	if flagHelp {
		printUsage()
		return
	}

	workQueue := make(chan int64, 5)

	createWorker(workQueue)

	server := newHTTPServer(workQueue)

	log.Printf("Listening on %s\n", flagAddr)
	log.Fatal(server.ListenAndServe())
}

func newHTTPServer(workQueue chan<- int64) *http.Server {
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/submit":
			handleSubmit(w, r, workQueue)
		case "/result":
			handleResult(w, r)
		default:
			fmt.Fprint(w, "Results: /result\nSubmit: /submit")
		}
	}

	server := &http.Server{
		Addr:         flagAddr,
		Handler:      requestLogger(httpHandler),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	return server
}

func handleSubmit(w http.ResponseWriter, r *http.Request, workQueue chan<- int64) {

	headerValue := r.URL.Query().Get("value")
	intVal, err := strconv.ParseInt(headerValue, 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s can't be parsed as a 64 bit integer: %v", headerValue, err)
		return
	}

	select {
	case workQueue <- intVal:
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "created job %d\n result will be availabe at /result", intVal)
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "The workqueue is currently full. Please try again later")
	}
}

func handleResult(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "The result are:\n%v", result)
}

func requestLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Printf("%v %v", r.Method, r.URL)
	}
}

func createWorker(workQueue <-chan int64) {
	go func() {
		for work := range workQueue {
			time.Sleep(time.Duration(work) * time.Millisecond)
			log.Printf("Done processing %d", work)
			result = append(result, work)
		}
	}()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
