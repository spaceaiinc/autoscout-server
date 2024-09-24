package utility

import "time"

var (
	Tokyo *time.Location
)

func init() {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	Tokyo = tz
}

// time.Timeでcreated_atとupdated_at以外の初期値入れる時に使用
func EarliestTime() (t time.Time) {
	t = time.Date(1, time.January, 1, 1, 0, 0, 0, time.UTC)
	return t
}
 