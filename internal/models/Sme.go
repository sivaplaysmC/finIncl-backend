package models

import (
	"database/sql"
	"log"

	_ "github.com/go-playground/validator/v10"
	"github.com/sivaplaysmc/finIncl-backend/internal/helpers"
)

type SmeModel struct {
	db  *sql.DB
	Log *log.Logger
}

type Sme struct {
	RegisteredDate helpers.DDMMYYYY `validate:"required" json:"rgdt" `
	Address        string           `validate:"required" json:"adr"`
	GSTIN          string           `validate:"required" json:"gstin"`
	Name           string           `vaildate:"required" json:"name"`
	PAN            string           `vaildate:"required" json:"pan"`
	PinCode        string           `vaildate:"required" json:"pincode"`
	Risk           int              `vaildate:"required" json:"risk"`
}

func NewSmeModel(db *sql.DB, logger *log.Logger) *SmeModel {
	return &SmeModel{
		db:  db,
		Log: logger,
	}
}
