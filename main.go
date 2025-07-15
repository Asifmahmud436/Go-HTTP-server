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
	"time"

	"github.com/Asifmahmud436/Go-HTTP-server/internal/database"
	"github.com/google/uuid"
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
	params.Body = strings.ReplaceAll(params.Body, "kerfuffle", "****")
	params.Body = strings.ReplaceAll(params.Body, "sharbet", "****")
	params.Body = strings.ReplaceAll(params.Body, "fornax", "****")
	json.NewEncoder(w).Encode(map[string]string{"cleaned_body": params.Body})

}

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
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

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	type UserEmail struct {
		Email string `json:"email"`
	}
	var params UserEmail
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.Email == "" {
		log.Printf("Error validating email: %s", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Couldnt create user: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Couldnt create user"})
		return
	}
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func (cfg *apiConfig) postChiprs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	type Input struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var params Input
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Error in input: %s",err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbChirp,err := cfg.DB.CreateChirp(r.Context(),database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	})
	if err!= nil{
		log.Printf("Error in creating user: %s",err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error":"Couldnt create a chirp"})
		return
	}
	type Chirp struct{
		Id uuid.UUID `json:"id"`
		CreatedAT time.Time `json:"created_at"`
		UpdatedAT time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	chirp := Chirp{
		Id: dbChirp.ID,
		CreatedAT: dbChirp.CreatedAt,
		UpdatedAT: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserId: dbChirp.UserID,
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(chirp)
	
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type","application/json")
	dbChirps, err := cfg.DB.ListChirps(r.Context())
	if(err!=nil){
		log.Printf("Couldnt get the chirps! : %s",err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error":"Error while getting all the chirps :D"})
		return
	}
	type Chirp struct{
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	chirps := make([]Chirp,len(dbChirps))
	for i, ch:= range dbChirps{
		chirps[i] = Chirp{
			Id: ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body: ch.Body,
			UserId: ch.UserID,
		}
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(chirps)
}

func main() {
	// opening the db
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Something wrong with the Database connection")
	}
	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB:             dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/assets", appAssetshandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.showHits)
	mux.HandleFunc("/admin/metrics", apiCfg.showDetailedHits)
	mux.HandleFunc("/admin/reset", apiCfg.resetHits)
	mux.HandleFunc("/api/validate_chirp", chirpValidater)
	mux.HandleFunc("POST /api/users", apiCfg.handleUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChiprs)
	mux.HandleFunc("GET /api/chirps",apiCfg.getChirps)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
