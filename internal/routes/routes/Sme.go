package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sivaplaysmC/finIncl-backend/config"
	model "github.com/sivaplaysmC/finIncl-backend/internal"
)

func GetSmeRoutes(app *config.App) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/create", doCreateSme(app))
	router.HandleFunc("/setRisk", doSetRisk(app))

	return router
}

// smeCreation controller
func doCreateSme(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sme := &model.Sme{}
		err := json.NewDecoder(r.Body).Decode(&sme)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Sprintf("https://gst-return-status.p.rapidapi.com/free/gstin/%v", sme.GSTIN)
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
		data := respJson["data"].(map[string]interface{})
		sme.Adr = data["adr"].(string)
		ddmmyy, err := time.Parse("02/01/2006", data["rgdt"].(string))
		sme.Rgdt = ddmmyy
		sme.PAN = data["pan"].(string)
		sme.Pincode = data["pincode"].(string)

		if err != nil {
			http.Error(w, "invalid JSON recieved", http.StatusNonAuthoritativeInfo)
			app.ErrorLog.Println(err)
			return
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(sme)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		err = app.Db.Create(&sme).Error
		fmt.Println(sme.Pincode)
		if err != nil {
			app.ErrorLog.Println("Fuk you")
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
	}
}

func doSetRisk(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sme := &model.Sme{}
		nameAndRisk := make(map[string]interface{}, 2)
		json.NewDecoder(r.Body).Decode(&nameAndRisk)
		err := app.Db.Where("name = ?", nameAndRisk["name"].(string)).First(&sme).Error
		if err != nil {
			http.Error(w, "internal sserver error", http.StatusInternalServerError)
			app.ErrorLog.Println(err)
			return
		}
		if shit := nameAndRisk["risk"]; shit != nil {
			sme.Risk = shit.(int)
			app.Db.Save(&sme)
		}
	}
}

func doGetSmes(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := make([]model.User, 10)
		err := app.Db.Select("name").Find(&users).Error
		if err != nil {
			http.Error(w, "shit happens", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}
