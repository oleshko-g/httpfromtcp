package request

var versionsSupported = map[string]struct{}{
	"1.1": {},
}

func VersionSupported(s string) bool {
	_, ok := versionsSupported[s]
	return ok
}

var methodsSupported = map[string]struct{}{
	"GET":  {},
	"POST": {},
}

func MethodSupported(s string) bool {
	_, ok := methodsSupported[s]
	return ok
}
