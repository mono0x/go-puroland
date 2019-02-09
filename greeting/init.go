package greeting

import "time"

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}
