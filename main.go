package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Asifmahmud436/Go-HTTP-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func chirpValidater(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	type bodyStructs struct {
		Body string `json:"body"`
	}
	var params bodyStructs
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong!"})
		return
	}
	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Chirp is too long!"})
		return
	}
	// kerfuffle
	// sharbert
	// fornax
	w.WriteHeader(200)
	params.Body = strings.ReplaceAll(params.Body, "kerfuffle","****")
	params.Body = strings.ReplaceAll(params.Body, "sharbet","****")
	params.Body = strings.ReplaceAll(params.Body, "fornax","****")
	json.NewEncoder(w).Encode(map[string]string{"cleaned_body": params.Body})

}

type apiConfig struct {
	fileserverHits atomic.Int32
	DB *database.Queries
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) showHits(w http.ResponseWriter, r *http.Request) {
	x := cfg.fileserverHits.Load()
	result := fmt.Sprintf("Hits: %v", x)
	w.Write([]byte(result))
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

func (cfg *apiConfig) showDetailedHits(w http.ResponseWriter, r *http.Request) {
	x := cfg.fileserverHits.Load()
	w.Header().Set("Content-type", "text/html;charset=utf-8")
	fmt.Fprintf(w, `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, x)

}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-type", "text/plain;charset=utf-8")
	fmt.Fprint(w, "Hit counter has been reset to 0")
}

func main() {
	// opening the db
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres",dbURL)
	if(err!=nil){
		log.Fatal("Something wrong with the Database connection")
	}
	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/assets", appAssetshandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.showHits)
	mux.HandleFunc("/admin/metrics", apiCfg.showDetailedHits)
	mux.HandleFunc("/admin/reset", apiCfg.resetHits)
	mux.HandleFunc("/api/validate_chirp", chirpValidater)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
