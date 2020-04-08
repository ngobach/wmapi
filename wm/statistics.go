package wm

import "time"

type DailyStatistic struct {
	Date   time.Time
	Total  int
	Active int
	New    int
}

type Report struct {
	Country   *Country
	Total     int
	Active    int
	Deaths    int
	Recovered int
	Days      []DailyStatistic
	UpdatedAt time.Time
}
