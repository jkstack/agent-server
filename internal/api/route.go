package api

type Route struct {
	Method string
	Uri    string
}

func MakeRoute(method, uri string) Route {
	return Route{Method: method, Uri: uri}
}
