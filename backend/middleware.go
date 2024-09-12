package main

import (
	"github.com/c00rni/Swiss-financial-events/internal/database"
	"log"
	"net/http"
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
