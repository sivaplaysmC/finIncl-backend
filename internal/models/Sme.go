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

func (smes *SmeModel) CreateSme(sme *Sme) error {
	_, err := smes.db.Query(
		`INSERT into SME
     (adr, GSTIN, name ,PAN, Pincode, rgdt, risk) 
     values (? ,? ,? ,? ,? ,? ,?)`,
		sme.Address,
		sme.GSTIN,
		sme.Name,
		sme.PAN,
		sme.PinCode,
		sme.RegisteredDate.Format("2006-01-02"),
		sme.Risk,
	)
	return err
}

func (smes *SmeModel) UpdateSmeID(sme *Sme) error {
	_, err := smes.db.Exec(
		`Update SME set risk=? where name like ?`,
		sme.Risk, sme.Name,
	)
	return err
}

func (smes *SmeModel) GetSmeID(name string) (int, error) {
	smeid := 0
	err := smes.db.QueryRow("select id from sme where name like ? ", name).Scan(&smeid)
	return smeid, err
}
