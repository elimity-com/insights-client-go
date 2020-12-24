package insights

import (
	"encoding/json"
	tim "time"
)

// Value represents a value that can be assigned to an entity for an attribute type.
type Value struct {
	value interface{}
}

// NewBooleanValue creates a new boolean value with the given value.
func NewBooleanValue(value bool) Value {
	val := booleanValue{value: value}
	return Value{value: val}
}

// NewDateTimeValue creates a new date-time value with the given value.
func NewDateTimeValue(value tim.Time) Value {
	val := dateTimeValue{value: value}
	return Value{value: val}
}

// NewDateValue creates a new date value with the given value.
func NewDateValue(value tim.Time) Value {
	val := dateValue{value: value}
	return Value{value: val}
}

// NewNumberValue creates a new number value with the given value.
func NewNumberValue(value float64) Value {
	val := numberValue{value: value}
	return Value{value: val}
}

// NewStringValue creates a new string value with the given value.
func NewStringValue(value string) Value {
	val := stringValue{value: value}
	return Value{value: val}
}

// NewTimeValue creates a new time value with the given value.
func NewTimeValue(value tim.Time) Value {
	val := timeValue{value: value}
	return Value{value: val}
}

func (v Value) model() valueModel {
	switch value := v.value.(type) {
	case booleanValue:
		return makeValueModel("boolean", value.value)

	case dateValue:
		time := value.value
		utc := time.UTC()
		day := utc.Day()
		month := utc.Month()
		mon := int(month)
		year := utc.Year()
		date := date{
			Day:   day,
			Month: mon,
			Year:  year,
		}
		return makeValueModel("date", date)

	case dateTimeValue:
		time := parseDateTime(value.value)
		return makeValueModel("dateTime", time)

	case numberValue:
		return makeValueModel("number", value.value)

	case stringValue:
		return makeValueModel("string", value.value)

	case timeValue:
		t := value.value
		utc := t.UTC()
		hour := utc.Hour()
		minute := utc.Minute()
		second := utc.Second()
		time := time{
			Hour:   hour,
			Minute: minute,
			Second: second,
		}
		return makeValueModel("time", time)

	default:
		panic("unreachable")
	}
}

type booleanValue struct {
	value bool
}

type date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type dateTime struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Month  int `json:"month"`
	Second int `json:"second"`
	Year   int `json:"year"`
}

type dateTimeValue struct {
	value tim.Time
}

type dateValue struct {
	value tim.Time
}

type numberValue struct {
	value float64
}

type stringValue struct {
	value string
}

type time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

type timeValue struct {
	value tim.Time
}

type valueModel struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func makeValueModel(typ string, value interface{}) valueModel {
	val, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return valueModel{
		Type:  typ,
		Value: val,
	}
}
