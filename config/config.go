package config

import (
	"fmt"
	"log"
	"os"

	modelID "github.com/sivaplaysmc/finIncl-backend/internal"
	"github.com/sivaplaysmc/finIncl-backend/internal/jwtgen"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"
)

type ContextKey string

const defaultDsn string = "sql:sqlatREC@tcp(0.0.0.0:3306)/gorm?parseTime=true"

type App struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	Db *gorm.DB

	Jwtgen jwtgen.JwtGenerator
}

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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db.AutoMigrate(&modelID.Sme{}, &modelID.User{}, &modelID.Document{}, &modelID.Project{})
	if err != nil {
		return nil, err
	}

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		return nil, fmt.Errorf("error : JWT_SECRET environment variable should be set ")
	}

	app := &App{
		InfoLog:  infoLog,
		ErrorLog: errrLog,

		Db: db,

		Jwtgen: jwtgen.NewJwtGenerator([]byte(jwt_secret)),
	}
	return app, nil
}
