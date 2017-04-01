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
)

var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgWhite, color.BgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

var tlsPort = flag.String("tlsPort", "8043", "TCP port to listen for tls on")
var port = flag.String("port", "8080", "TCP port to listen for on")
var hostname = flag.String("hostname", "localhost", "FQDN of the url")

//to generate default certificates
//
// go run /usr/local/go/src/crypto/tls/generate_cert.go --host=localhost
//
// depetning on personal paths
var tls = flag.Bool("tls", false, "Use TLS (https) or not")
var cert = flag.String("cert", "crt/cert.pem", "TLS Certificate")
var key = flag.String("key", "crt/key.pem", "TLS Key")

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

	var uencMux redirecter
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, uencMux) //write this to log file later

	addr := ":" + *port
	h := &http.Server{Addr: addr, Handler: loggedRouter}
	go func() {

		log.Println(green("INFO"), "Starting HTTP server at", addr)

		if err := h.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf(red("ERROR")+" listen: %s\n", err)
		}
	}()

	addrTLS := ":" + *tlsPort
	loggedTLSRouter := handlers.CombinedLoggingHandler(os.Stdout, r) //write this to log file later
	hTLS := &http.Server{Addr: addrTLS, Handler: loggedTLSRouter}
	go func() {

		log.Println(green("INFO"), "Starting HTTPS server at", addrTLS)

		if err := hTLS.ListenAndServeTLS(*cert, *key); err != http.ErrServerClosed {
			log.Printf(red("ERROR")+" listen TLS: %s\n", err)
		}
	}()

	<-stop

	log.Println(yellow("WARN"), "Shutting down server...")

	if err := h.Shutdown(context.Background()); err != nil {
		log.Fatalf(red("ERROR ")+"could not shutdown: %v\n", err)
	}

	if err := hTLS.Shutdown(context.Background()); err != nil {
		log.Fatalf(red("ERROR ")+"could not shutdown: %v\n", err)
	}

	log.Println(red("Servers gracefully stopped"))
}

type redirecter struct{}

func (s redirecter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI())
	redirectURL := "https://" + *hostname
	if *tlsPort != "443" {
		redirectURL += ":" + *tlsPort
	}
	redirectURL += r.URL.RequestURI()
	http.Redirect(w, r, redirectURL, 301)
}
