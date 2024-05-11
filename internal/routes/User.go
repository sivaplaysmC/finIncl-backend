package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sivaplaysmc/finIncl-backend/config"
)

func GetUsersRouter(app *config.App) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", doLogin(app)).Methods(http.MethodPost)
	router.HandleFunc("/test", doTest(app))

	router.Handle("/secure", doSecure(app))

	return router
}

func doTest(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Println("Hit!!")
		fmt.Fprintln(w, "Hi There")
	}
}

func doLogin(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get login, password from json body

		loginParams := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		// check hash
		err := json.NewDecoder(r.Body).Decode(&loginParams)
		if err != nil {
			http.Error(w, "Error Decoding JSON", http.StatusUnauthorized)
			return
		}
		id, err := app.Users.ValidateUser(loginParams.Username, loginParams.Password)
		if err != nil {
			http.Error(w, "Wrong credentials", http.StatusUnauthorized)
			return
		}
		values := jwt.MapClaims{
			"id": id,
		}
		token, err := app.Jwtgen.GenToken(values)
		if err != nil {
			http.Error(w, "error generating JWT "+err.Error(), http.StatusInternalServerError)
			return
		}

		tokenCookie := &http.Cookie{
			Name:     "AuthToken",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, tokenCookie)
		w.Write([]byte(""))
	}
}

func doSecure(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hi There blud! You are logged in!")
		claims := r.Context().Value(config.ContextKey("claims")).(jwt.MapClaims)
		user, err := app.Users.GetUserByID(claims["id"].(int))
		if err != nil {
			http.Error(w, "internal server error : ", http.StatusInternalServerError)
		}
		fmt.Fprintln(w, user.Name, user.Email, user.SmeID)
	}
}

func doCreate(app *config.App) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		createUserCreds := make(map[string]interface{}, 5)
		validationRules := map[string]interface{}{
			"username": "required",
			"password": "required",
			"email":    "required,email",
			"smeName":  "required",
		}
		json.NewDecoder(r.Body).Decode(&createUserCreds)
		errs := validator.New().ValidateMap(createUserCreds, validationRules)
		if len(errs) > 0 {
			app.ErrorLog.Println(errs)
			return
		}
	}
	return handler
}
