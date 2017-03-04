package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgWhite, color.BgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

var tlsPort = flag.String("tlsPort", "8043", "TCP port to listen for tls on")
var port = flag.String("port", "8080", "TCP port to listen for on")

//to generate default certificates
//
// go run /usr/local/go/src/crypto/tls/generate_cert.go --host=localhost
//
// depetning on personal paths
var tls = flag.Bool("tls", false, "Use TLS (https) or not")
var cert = flag.String("cert", "crt/cert.pem", "TLS Certificate")
var key = flag.String("key", "crt/key.pem", "TLS Key")

var r = mux.NewRouter()

func init() {
	version, err := strconv.Atoi(runtime.Version()[4:])
	if err != nil {
		log.Fatalln(red("ERROR"), "Couln't reliably find golang version [comment this code block out if you are using more than go1.8]", err)
	}
	if version < 8 {
		log.Fatalln(red("ERROR"), "Minimum version go1.8 required for this project")
	}

	flag.Parse()

}

func main() {
	stop := make(chan os.Signal)

	signal.Notify(stop, os.Interrupt)
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, r)

	addr := ":" + *port
	h := &http.Server{Addr: addr, Handler: loggedRouter}
	go func() {

		log.Println(green("INFO"), "Starting HTTP server at", addr)

		if err := h.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf(red("ERROR")+" listen: %s\n", err)
		}
	}()

	addrTLS := ":" + *tlsPort
	hTLS := &http.Server{Addr: addrTLS, Handler: loggedRouter, Certificate: *cert, Key: *key}
	go func() {

		log.Println(green("INFO"), "Starting HTTPS server at", addr)

		if err := hTLS.ListenAndServeTLS(); err != http.ErrServerClosed {
			log.Printf(red("ERROR")+" listen TLS: %s\n", err)
		}
	}()

	<-stop

	log.Println(yellow("WARN"), "Shutting down server...")

	if err := h.Shutdown(context.Background()); err != nil {
		log.Fatalf(red("ERROR ")+"could not shutdown: %v\n", err)
	}

	log.Println(red("Server gracefully stopped"))
}
