package api

import "fmt"

type BadParam string

func (e BadParam) Error() string {
	return fmt.Sprintf("bad param: %s", string(e))
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
