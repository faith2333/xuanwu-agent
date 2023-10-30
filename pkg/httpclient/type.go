package httpclient

import "strings"

type Method string

const (
	MethodGET    Method = "GET"
	MethodPOST   Method = "POST"
	MethodDELETE Method = "DELETE"
	MethodPUT    Method = "PUT"
)

var AllMethods = []Method{
	MethodGET, MethodPOST, MethodDELETE, MethodPUT,
}

func (m Method) String() string {
	return string(m)
}

func (m Method) Upper() Method {
	return Method(strings.ToUpper(m.String()))
}

func (m Method) IsSupported() bool {
	for _, aMethod := range AllMethods {
		if m.Upper() == aMethod {
			return true
		}
	}

	return false
}
