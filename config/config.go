package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/sivaplaysmc/finIncl-backend/internal/models"
)

type App struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	Users *models.UserModel
	Smes  *models.SmeModel
}

type ContextKey string

const defaultDsn string = "sql:sqlatREC@tcp(0.0.0.0:3306)/finIncl"

func GetApp() (*App, error) {
	infoLog := log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errrLog := log.New(os.Stdout, "[ERRR] ", log.LstdFlags|log.Lshortfile)

	argc := len(os.Args)
	dsn := ""
	if argc >= 2 {
		dsn = os.Args[1]
	} else {
		infoLog.Println("data source name not supplied, defaulting to", defaultDsn)
		dsn = defaultDsn
	}

	db, err := getDBPool(dsn)
	if err != nil {
		return nil, err
	}

	app := &App{
		Users:    models.NewUserModel(db, infoLog),
		Smes:     models.NewSmeModel(db, infoLog),
		InfoLog:  infoLog,
		ErrorLog: errrLog,
	}
	return app, nil
}

func getDBPool(dataSoruceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSoruceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
