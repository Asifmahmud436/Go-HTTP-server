package main

import (
	"fmt"
	"log"
	"net/http"
)

func apphandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html> Welcome to Chirpy</html>"))
}

func appAssetshandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<pre><a href="logo.png">logo.png</a>Expecting body to contain: </pre>`))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) 
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/app/page/", apphandler)
	mux.HandleFunc("/assets/", appAssetshandler)
	mux.HandleFunc("/healthz", healthzHandler) 

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	server := http.Server{Handler: mux, Addr: ":8080"}

	fmt.Println("Starting the server!!!")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
