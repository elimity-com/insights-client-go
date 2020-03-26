package insights

import (
	"encoding/json"
	"strconv"
	"time"
)

// Value represents a value that can be assigned to an entity for an attribute type.
type Value struct {
	data valueData
}

// NewBooleanValue creates a new boolean value with the given value.
func NewBooleanValue(value bool) Value {
	data := booleanValueData{value: value}
	return Value{data: data}
}

// NewDateValue creates a new date value with the given value.
func NewDateValue(value time.Time) Value {
	data := dateValueData{value: value}
	return Value{data: data}
}

// NewDateTimeValue creates a new date-time value with the given value.
func NewDateTimeValue(value time.Time) Value {
	data := dateTimeValueData{value: value}
	return Value{data: data}
}

// NewNumberValue creates a new number value with the given value.
func NewNumberValue(value float64) Value {
	data := numberValueData{value: value}
	return Value{data: data}
}

// NewStringValue creates a new string value with the given value.
func NewStringValue(value string) Value {
	data := stringValueData{value: value}
	return Value{data: data}
}

// NewTimeValue creates a new time value with the given value.
func NewTimeValue(value time.Time) Value {
	data := timeValueData{value: value}
	return Value{data: data}
}

func (v Value) model() valueModel {
	switch data := v.data.(type) {
	case booleanValueData:
		value := strconv.FormatBool(data.value)
		return valueModel{
			Type:  "boolean",
			Value: value,
		}

	case dateValueData:
		value := data.value.Format("2006-01-02")
		return valueModel{
			Type:  "date",
			Value: value,
		}

	case dateTimeValueData:
		value := data.value.Format(time.RFC3339)
		return valueModel{
			Type:  "dateTime",
			Value: value,
		}

	case numberValueData:
		bs, err := json.Marshal(data.value)
		if err != nil {
			panic(err)
		}
		value := string(bs)
		return valueModel{
			Type:  "number",
			Value: value,
		}

	case stringValueData:
		return valueModel{
			Type:  "string",
			Value: data.value,
		}

	case timeValueData:
		value := data.value.Format("15:04:05Z07:00")
		return valueModel{
			Type:  "time",
			Value: value,
		}

	default:
		panic("unreachable")
	}
}

type booleanValueData struct {
	value bool
}

type dateValueData struct {
	value time.Time
}

type dateTimeValueData struct {
	value time.Time
}

type numberValueData struct {
	value float64
}

type stringValueData struct {
	value string
}

type timeValueData struct {
	value time.Time
}

type valueData = interface{}

type valueModel struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
