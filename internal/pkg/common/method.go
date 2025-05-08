package common

// MethodValue ...
type MethodValue struct {
	Value string
}

// MethodEnum ...
type MethodEnum struct {
	ANY     *MethodValue
	GET     *MethodValue
	QUERY   *MethodValue
	POST    *MethodValue
	PUT     *MethodValue
	DELETE  *MethodValue
	OPTIONS *MethodValue
}

// APIMethod Published enum
var APIMethod = MethodEnum{
	ANY:     &MethodValue{Value: "ANY"},
	GET:     &MethodValue{Value: "GET"},
	QUERY:   &MethodValue{Value: "QUERY"},
	POST:    &MethodValue{Value: "POST"},
	PUT:     &MethodValue{Value: "PUT"},
	DELETE:  &MethodValue{Value: "DELETE"},
	OPTIONS: &MethodValue{Value: "OPTIONS"},
}
