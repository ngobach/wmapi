package wm

import "time"

type DailyStatistic struct {
	Date   time.Time `json:"date"`
	Total  int       `json:"total"`
	Active int       `json:"active"`
	New    int       `json:"new"`
}

type Report struct {
	Country   *Country         `json:"country"`
	Total     int              `json:"total"`
	Deaths    int              `json:"deaths"`
	Recovered int              `json:"recovered"`
	Days      []DailyStatistic `json:"days"`
	UpdatedAt time.Time        `json:"updated_at"`
}
