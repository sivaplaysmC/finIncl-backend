package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sivaplaysmC/finIncl-backend/config"
	model "github.com/sivaplaysmC/finIncl-backend/internal"
	"github.com/sivaplaysmC/finIncl-backend/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func GetUsersRouter(app *config.App) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", doLogin(app)).Methods(http.MethodPost)
	router.HandleFunc("/test", doTest(app))

	router.Handle("/secure", middleware.AuthMiddleware(app)(doSecure(app)))
	router.Handle("/create", doCreate(app))

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

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Error Decoding JSON", http.StatusUnauthorized)
			return
		}

		user := model.User{}
		err := app.Db.Where("name = ?", username).First(&user).Error
		if err != nil {
			http.Error(w, "No such user", http.StatusNotAcceptable)
			app.ErrorLog.Println(err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(password))
		if err != nil {
			http.Error(w, "Wrong credentials", http.StatusUnauthorized)
			app.ErrorLog.Println(err)
			return
		}

		values := jwt.MapClaims{
			"id": user.ID,
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
			MaxAge:   3600,
		}

		http.SetCookie(w, tokenCookie)
		w.Write([]byte(""))
	}
}

func doSecure(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hi There blud! You are logged in!")
		claims := r.Context().Value(config.ContextKey("claims")).(jwt.MapClaims)
		user := model.User{}
		err := app.Db.First(&user, claims["id"].(int)).Error
		if err != nil {
			http.Error(w, "internal server error : ", http.StatusInternalServerError)
			return
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
			"gstin":    "required",
			"type":     "required",
		}
		json.NewDecoder(r.Body).Decode(&createUserCreds)
		errs := validator.New().ValidateMap(createUserCreds, validationRules)
		if len(errs) > 0 {
			http.Error(w, "invalid input", http.StatusNotAcceptable)
			return
		}

		sme := model.Sme{}
		err := app.Db.
			Where("Name like ?", createUserCreds["smeName"].(string)).
			First(&sme).Error
		// smeId, err := app.Smes.GetSmeID(createUserCreds["smeName"].(string))

		client := &http.Client{}
		request, err := http.NewRequest("GET", fmt.Sprintf("https://gst-return-status.p.rapidapi.com/free/gstin/%v", sme.GSTIN), nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-RapidAPI-Key", "2048709673msha68caf7a53b895cp1cdce5jsn716f33d84e2e")
		request.Header.Add("X-RapidAPI-Host", "gst-return-status.p.rapidapi.com")

		resp, err := client.Do(request)

		respJson := make(map[string]interface{}, 10)
		json.NewDecoder(resp.Body).Decode(&respJson)
		success := respJson["success"].(bool)

		if !success {
			http.Error(w, "invalid GSTIN error", http.StatusNotAcceptable)
			return
		}
		// data := respJson["data"].(map[string]interface{})

		if err != nil {
			http.Error(w, "Sme not found", http.StatusNotAcceptable)
			app.ErrorLog.Println(err)
			return
		}

		passHash, _ := bcrypt.GenerateFromPassword([]byte(createUserCreds["password"].(string)), 12)
		user := model.User{
			Name:     createUserCreds["username"].(string),
			Email:    createUserCreds["email"].(string),
			SmeID:    int(sme.ID),
			Passhash: string(passHash),
		}
		err = app.Db.Create(&user).Error
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			app.ErrorLog.Println(err)
			return
		}

		w.WriteHeader(200)
	}
	return handler
}
