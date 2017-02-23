package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"

	"github.com/fatih/color"
)

var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgWhite, color.BgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func init() {
	version, err := strconv.Atoi(runtime.Version()[4:])
	if err != nil {
		log.Fatalln(red("ERROR"), "Couln't reliably find golang version [comment this code block out if you are using more than go1.8]", err)
	}
	if version < 8 {
		log.Fatalln(red("ERROR"), "Minimum version go1.8 required for this project")
	}
}

func main() {
	stop := make(chan os.Signal)

	signal.Notify(stop, os.Interrupt)

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":2017"
	}

	h := &http.Server{Addr: addr, Handler: &server{}}

	go func() {
		log.Println(green("INFO"), "Starting HTTP server at", addr)

		if err := h.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf(red("ERROR")+" listen: %s\n", err)
		}
	}()

	<-stop

	log.Println(yellow("WARN"), "Shutting down server...")

	if err := h.Shutdown(context.Background()); err != nil {
		log.Fatalf(red("ERROR ")+"could not shutdown: %v\n", err)
	}

	log.Println(red("Server gracefully stopped"))
}

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
