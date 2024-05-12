package routes

import (
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sivaplaysmc/finIncl-backend/config"
	modelID "github.com/sivaplaysmc/finIncl-backend/internal"
	"github.com/sivaplaysmc/finIncl-backend/internal/middleware"
)

func GetContentRoutes(app *config.App) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/", middleware.AuthMiddleware(app)(doUpload(app)))
	return router
}

func doUpload(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 32)
		val := r.Context().Value(config.ContextKey("claims")).(jwt.MapClaims)

		file, header, err := r.FormFile("upload")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		bytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		doc := modelID.Document{}
		doc.File = bytes
		Mime, err := mimetype.DetectFile(header.Filename)
		if err == nil {
			doc.FileType = Mime.String()
		} else {
			doc.FileType = "application/json"
		}
		doc.FileName = header.Filename
		doc.UserID = int(val["id"].(float64))

		app.Db.Save(&doc)
	}
}
