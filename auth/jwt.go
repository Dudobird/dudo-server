package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zhangmingkai4315/dudo-server/models"
	"github.com/zhangmingkai4315/dudo-server/utils"
)

var (
	guestURL = []string{
		"/api/user/new",
		"/api/user/login",
	}
)

// JWTAuthentication is a middleware for all request
// it will stop the request when jwt authencticate is fail
func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		for _, url := range guestURL {
			if url == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			utils.JSONRespnseWithTextMessage(w, http.StatusForbidden, "missing auth token")
			return
		}
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "malformed token")
			return
		}
		tokenFromHeader := splitted[1]

		userToken := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenFromHeader, userToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil {
			utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "token process fail")
			return
		}

		if !token.Valid {
			utils.JSONRespnseWithTextMessage(w, http.StatusForbidden, "token valid fail")
			return
		}

		ctx := context.WithValue(r.Context(), "user", userToken.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
