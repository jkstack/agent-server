package api

type Route struct {
	Method     string
	Uri        string
	MetricName string
}

func MakeRoute(method, uri, name string) Route {
	return Route{
		Method:     method,
		Uri:        uri,
		MetricName: name,
	}
}
