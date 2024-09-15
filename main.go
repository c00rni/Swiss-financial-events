package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/c00rni/Swiss-financial-events/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"time"
)

var randomStr, _ = generateApiKey()
var store = sessions.NewCookieStore([]byte(randomStr))

type apiConfig struct {
	port string
	DB   *database.Queries
	auth *oauth2.Config
}

type Goauth struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func handleReadyness(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return
}

func (cfg *apiConfig) handleGetEvents(w http.ResponseWriter, r *http.Request, user database.User) {

	location := r.URL.Query().Get("location")
	if location == "" {
		location = "%"
	}
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%"
	}
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		topic = "%"
	}
	params := database.GetFilteredEventsParams{
		Location: location,
		Name:     category,
		Name_2:   topic,
	}

	events, err := cfg.DB.GetFilteredEvents(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, events)
	return
}

func (cfg *apiConfig) handleGetCategories(w http.ResponseWriter, r *http.Request, user database.User) {
	names, err := cfg.DB.GetCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, names)
	return
}

func (cfg *apiConfig) handleGetTopics(w http.ResponseWriter, r *http.Request, user database.User) {
	names, err := cfg.DB.GetTopics(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, names)
	return
}

func (cfg *apiConfig) handleOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := cfg.auth.Exchange(context.Background(), code)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	client := cfg.auth.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "Can't reach the Oauth provider.")
		return
	}

	decoder := json.NewDecoder(resp.Body)
	jsonResp := &Goauth{}
	err = decoder.Decode(jsonResp)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	apiKey, err := generateApiKey()
	if err != nil {
		log.Println("Failed to create a API key:", err)
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), jsonResp.Email)
	if err != nil {
		newUser := database.CreateUserParams{
			ID:            uuid.NewString(),
			Email:         jsonResp.Email,
			VerifiedEmail: jsonResp.VerifiedEmail,
			Name:          jsonResp.Name,
			GivenName:     jsonResp.GivenName,
			FamilyName:    jsonResp.FamilyName,
			Picture:       jsonResp.Picture,
			Token:         token.RefreshToken,
			ApiKey:        apiKey,
		}
		user, err = cfg.DB.CreateUser(r.Context(), newUser)
		if err != nil {
			log.Println("Failed to create an user:", err)
			respondWithError(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}
	today := time.Now()
	aMonthEarlier := today.AddDate(0, -1, 0)
	requests, err := cfg.DB.GetUserRequests(r.Context(), database.GetUserRequestsParams{
		UserID: user.ID,
		Date:   aMonthEarlier,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	nbRequests := len(requests)

	cookie := &http.Cookie{
		Name:     "user",
		Value:    fmt.Sprintf("apiKey=%v,name=%v,picture=%v,nbRequests=%v", user.ApiKey, user.Name, user.Picture, nbRequests),
		MaxAge:   60 * 60 * 24,
		Secure:   true,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	session, _ := store.Get(r, "session_id")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 3,
		HttpOnly: false,
	}
	session.Values["email"] = user.Email
	err = session.Save(r, w)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error.")
		return
	}

	http.Redirect(w, r, "/dashboard/", http.StatusTemporaryRedirect)
}

func (cfg *apiConfig) handleOauth(w http.ResponseWriter, r *http.Request) {
	url := cfg.auth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (cfg *apiConfig) handleUserInfo(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		Email      string `json:"email"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture_url"`
		ApiKey     string `json:"api_key"`
	}

	resp := response{
		Email:      user.Email,
		Name:       user.Name,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Picture:    user.Picture,
		ApiKey:     user.ApiKey,
	}
	respondWithJSON(w, http.StatusOK, resp)
	return
}

func (cfg *apiConfig) handleRenewApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	apiKey, err := generateApiKey()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		log.Println(err)
		return
	}

	user, err = cfg.DB.UpdateUserApiKey(r.Context(), database.UpdateUserApiKeyParams{
		ApiKey: apiKey,
		Email:  user.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		log.Println(err)
		return
	}

	cookie := &http.Cookie{
		Name:     "user",
		Value:    fmt.Sprintf("apiKey=%v,name=%v,picture=%v", user.ApiKey, user.Name, user.Picture),
		MaxAge:   60 * 60 * 24,
		Secure:   true,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	type response struct {
		ApiKey string `json:"api_key"`
	}
	respondWithJSON(w, http.StatusOK, response{ApiKey: user.ApiKey})
	return
}

func (cfg *apiConfig) handleDashboard(w http.ResponseWriter, r *http.Request, user database.User) {
	staticHandler := http.StripPrefix("/dashboard", http.FileServer(http.Dir("./_dist/")))
	staticHandler.ServeHTTP(w, r)
	return
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	staticHandler := http.FileServer(http.Dir("./_dist"))
	staticHandler.ServeHTTP(w, r)
	return
}

func (cfg *apiConfig) handleLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	session, _ := store.Get(r, "session_id")
	_, ok := session.Values["email"]
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	session.Options.MaxAge = -1

	cookie := &http.Cookie{
		Name:     "user",
		Value:    "",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	err := session.Save(r, w)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error.")
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.AllowContentType("application/json", "text/plain", "text/css"))
	router.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	router.Use(middleware.CleanPath)
	router.Use(middleware.Timeout(time.Second * 60))
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("'PORT' environment variable is not set.")
		return
	}

	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		log.Fatalln("'DATABASE_URL' environment variable is not set.")
		return
	}

	authConf := &oauth2.Config{
		ClientID:     os.Getenv("GCP_CLIENT_ID"),
		ClientSecret: os.Getenv("GCP_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_CALLBACK_URI"),
		Scopes: []string{
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}

	apiCfg := apiConfig{
		port: port,
		DB:   nil,
		auth: authConf,
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("libsql", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		dbQueries := database.New(db)
		apiCfg.DB = dbQueries
		log.Println("Connected to database!")
	}

	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)

	// Require Api key
	v1Router := chi.NewRouter()
	v1Router.Get("/events", apiCfg.middlewareApiToken(apiCfg.handleGetEvents))
	v1Router.Get("/categories", apiCfg.middlewareApiToken(apiCfg.handleGetCategories))
	v1Router.Get("/topics", apiCfg.middlewareApiToken(apiCfg.handleGetTopics))
	v1Router.Get("/healthz", apiCfg.middlewareApiToken(handleReadyness))
	apiRouter.Mount("/v1", v1Router)

	// Require Authentificaiton
	router.Group(func(r chi.Router) {
		router.Get("/dashboard", apiCfg.middlewareSession(apiCfg.handleDashboard))
		router.Get("/user", apiCfg.middlewareSession(apiCfg.handleUserInfo))
		router.Get("/logout", apiCfg.middlewareSession(apiCfg.handleLogout))
		router.Get("/renew", apiCfg.middlewareSession(apiCfg.handleRenewApiKey))
	})

	// Public routes
	router.Get("/*", apiCfg.handleLogin)
	router.Get("/auth/oauth", apiCfg.handleOauth)
	router.Get("/auth/callback", apiCfg.handleOauthCallback)

	go apiCfg.scrapCfasociety(time.Hour * 6)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + apiCfg.port,
	}

	log.Printf("Serving on port: %s\n", apiCfg.port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}
