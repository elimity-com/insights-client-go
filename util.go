package insights

import tim "time"

func parseDateTime(time tim.Time) dateTime {
	utc := time.UTC()
	day := utc.Day()
	hour := utc.Hour()
	minute := utc.Minute()
	month := utc.Month()
	mon := int(month)
	second := utc.Second()
	year := utc.Year()
	return dateTime{
		Day:    day,
		Hour:   hour,
		Minute: minute,
		Month:  mon,
		Second: second,
		Year:   year,
	}
}
