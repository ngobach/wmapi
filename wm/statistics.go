package wm

import "time"

type DailyStatistic struct {
	Date   time.Time `json:"date"`
	Total  int       `json:"total"`
	Active int       `json:"active"`
	New    int       `json:"new"`
}

type Report struct {
	Country   *Country         `json:"country,omitempty"`
	Total     int              `json:"total,omitempty"`
	Deaths    int              `json:"deaths,omitempty"`
	Recovered int              `json:"recovered,omitempty"`
	Days      []DailyStatistic `json:"days,omitempty"`
	UpdatedAt time.Time        `json:"updated_at,omitempty"`
}
