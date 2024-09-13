package main

import (
	"github.com/c00rni/Swiss-financial-events/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareSession(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session_id")
		email, ok := session.Values["email"]
		log.Println(email, ok)
		if !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		user, err := cfg.DB.GetUserByEmail(r.Context(), email.(string))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get the user form thr database.")
			log.Println(err)
			return
		}

		handler(w, r, user)
	})
}

func (cfg *apiConfig) middlewareApiToken(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := extractAuthorization(r, "Bearer ")
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Databse error")
			log.Println(err)
			return
		}

		_, err = cfg.DB.CreateRequest(r.Context(), database.CreateRequestParams{
			ID:     uuid.NewString(),
			UserID: user.ID,
			Date:   time.Now(),
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database error")
			log.Println(err)
			return
		}
		handler(w, r, user)
	})
}
