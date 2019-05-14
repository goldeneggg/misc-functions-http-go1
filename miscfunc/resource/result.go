package resource

type Result struct {
	Body       string
	StatusCode int
	Header     map[string]string
}

func NewResult(msg string, sts int) *Result {
	return &Result{
		Body:       msg,
		StatusCode: sts,
	}
}

func NewResultWithHeader(msg string, sts int, header map[string]string) *Result {
	result := NewResult(msg, sts)
	result.Header = header

	return result
}

func NewResultWithErrorAndStatus(err error, sts int) (*Result, error) {
	return NewResult(err.Error(), sts), err
}
