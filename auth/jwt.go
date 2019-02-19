package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Dudobird/dudo-server/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

var (
	guestURL = []string{
		"/api/auth/signup",
		"/api/auth/signin",
	}
)

// ContextToken package exported type
type ContextToken string

// TokenContextKey a simple wrapper for ContextToken
const TokenContextKey = ContextToken("MyAppToken")

// JWTAuthentication is a middleware for all request
// it will stop the request when jwt authencticate is fail
func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := config.GetConfig()
		tokenSecrect := config.Application.Token
		requestPath := r.URL.Path
		for _, url := range guestURL {
			if url == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			utils.JSONRespnseWithTextMessage(w, http.StatusUnauthorized, "missing auth token")
			return
		}
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			utils.JSONRespnseWithTextMessage(w, http.StatusUnauthorized, "malformed token")
			return
		}
		tokenFromHeader := splitted[1]

		userToken := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenFromHeader, userToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecrect), nil
		})

		if err != nil {
			utils.JSONRespnseWithTextMessage(w, http.StatusUnauthorized, "token process fail")
			return
		}

		if !token.Valid {
			utils.JSONRespnseWithTextMessage(w, http.StatusUnauthorized, "token valid fail")
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, userToken.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
