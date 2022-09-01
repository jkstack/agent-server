package api

import "fmt"

type BadParam string

func (e BadParam) Error() string {
	if len(e) > 0 {
		return fmt.Sprintf("bad param: %s", string(e))
	}
	return "bad param"
}

func BadParamErr(param string) {
	panic(BadParam(param))
}

type Notfound string

func (e Notfound) Error() string {
	return fmt.Sprintf("not found: %s", string(e))
}

func NotfoundErr(what string) {
	panic(Notfound(what))
}

// MissingParam missing param error
type MissingParam string

// Error get missing param error info data, format: Missing [<data>]
func (e MissingParam) Error() string {
	return fmt.Sprintf("Missing [%s]", string(e))
}

// Timeout timeout error
type Timeout struct{}

// Error get timeout error info data, format: timeout
func (e Timeout) Error() string {
	return "timeout"
}

// InvalidType invalid agent type
type InvalidType struct {
	want string
	got  string
}

// Error get invalid type error info data, format: invalid agent type want [<type>] got [<type>]
func (e InvalidType) Error() string {
	return fmt.Sprintf("invalid agent type want [%s] got [%s]", e.want, e.got)
}
