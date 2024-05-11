package middleware

import (
	"context"
	"net/http"

	"github.com/sivaplaysmc/finIncl-backend/config"
)

type (
	middleware func(http.Handler) http.Handler
	a          string
)

func AuthMiddleware(app *config.App) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("AuthToken")
			if err != nil {
				http.Error(w, "Unauthorized, please login", http.StatusUnauthorized)
				return
			}
			tokenString := cookie.Value
			claims, verificationError := app.Jwtgen.VerifyToken(tokenString)
			if verificationError != nil {
				app.ErrorLog.Println("token verification error")
				http.Error(w, "bad authrization", http.StatusNetworkAuthenticationRequired)
				return
			}
			newRequest := r.WithContext(context.WithValue(r.Context(), config.ContextKey("claims"), claims))
			next.ServeHTTP(w, newRequest)
		})
	}
}
