package routers

import (
	"context"
	"log"
	"net/http"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/utils"
)

const appContextKey = utils.ContextToken("App")

func appBindMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := core.GetApp()
		ctx := context.WithValue(r.Context(), appContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// adminOnly return true if user is a administrator
		isAdmin := r.Context().Value(utils.AdminContextKey).(bool)
		log.Printf("isAdmin=%v", isAdmin)
		if isAdmin != true {
			utils.JSONRespnseWithTextMessage(w, http.StatusUnauthorized, "admin only")
			return
		}
		next.ServeHTTP(w, r)
	})
}
