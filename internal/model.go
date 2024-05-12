package modelID

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `json:"name"`
	Passhash string
	Email    string `json:"email"`

	Type string

	Documents []Document
	SmeID     int
}

type Sme struct {
	Rgdt time.Time `json:"rgdt"`
	gorm.Model
	GSTIN   string `json:"gstin"`
	Name    string `json:"name"`
	PAN     string `json:"pan"`
	Adr     string `json:"adr"`
	Pincode string `json:"pincode"`

	Users    []User
	Projects []Project

	Risk int `json:"risk"`
}

type Document struct {
	gorm.Model

	FileType string
	FileName string
	File     []byte

	UserID int
}

type Project struct {
	gorm.Model

	Name string
	Desc string

	SmeID int
}

type Loan struct {
	gorm.Model
}

type Bank struct {
	gorm.Model
}
