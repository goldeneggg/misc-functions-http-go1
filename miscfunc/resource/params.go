package resource

type Params struct {
	Path        string
	Method      string
	QueryParams map[string]string
	PathParams  map[string]string
	Header      map[string]string
	Body        string
	Stage       string
}
