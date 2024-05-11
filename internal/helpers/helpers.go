package helpers

import (
	"fmt"
	"strings"
	"time"
)

type DDMMYYYY struct {
	time.Time
}

func (d DDMMYYYY) MarshalJSON() ([]byte, error) {
	year, month, date := d.Date()
	return []byte(fmt.Sprintf("%2d-%2d-%4d", date, month, year)), nil
}

func (d *DDMMYYYY) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	tim, err := time.Parse("02-01-2006", str)
	if err != nil {
		return err
	}
	*d = DDMMYYYY{Time: tim}
	return nil
}
