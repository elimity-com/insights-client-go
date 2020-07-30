package insights

// Type represents one of the data types supported by all Elimity Insights servers.
type Type struct {
	data typeData
}

// NewBooleanType creates a new boolean type.
func NewBooleanType() Type {
	data := booleanTypeData{}
	return Type{data: data}
}

// NewDateType creates a new date type.
func NewDateType() Type {
	data := dateTypeData{}
	return Type{data: data}
}

// NewDateTimeType creates a new date-time type.
func NewDateTimeType() Type {
	data := dateTimeTypeData{}
	return Type{data: data}
}

// NewNumberType creates a new number type.
func NewNumberType() Type {
	data := numberTypeData{}
	return Type{data: data}
}

// NewStringType creates a new string type.
func NewStringType() Type {
	data := stringTypeData{}
	return Type{data: data}
}

// NewTimeType creates a new time type.
func NewTimeType() Type {
	data := timeTypeData{}
	return Type{data: data}
}

type booleanTypeData struct{}

type dateTypeData struct{}

type dateTimeTypeData struct{}

type numberTypeData struct{}

type typeData = interface{}

type stringTypeData struct{}

type timeTypeData struct{}
