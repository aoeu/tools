package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	args := struct {
		port string
		URI  string
	}{}
	flag.StringVar(&args.URI, "uri", "/", "The URI to serve and dump HTTP requests that have been made to it.")
	flag.StringVar(&args.port, "port", ":8080", "The port to serve the URI on and dump HTTP requests from.")
	flag.Parse()
	http.HandleFunc(args.URI, func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not dump client request: %v\n", err)
		}
		fmt.Println(string(b))
	})
	log.Fatal(http.ListenAndServe(args.port, nil))
}
