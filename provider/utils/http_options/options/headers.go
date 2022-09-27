package options

import "github.com/klever-io/klever-go-sdk/provider/utils/http_options"

type Headers struct {
	*http_options.HTTPOptions
	Headers [][2]string `json:"headers"`
}

func NewHeaders(headers [][2]string) http_options.IOptions {
	opts := Headers{
		HTTPOptions: &http_options.HTTPOptions{Type: http_options.Header},
		Headers:     headers,
	}

	return &opts
}

func (s *Headers) GetHeaders() [][2]string {
	return s.Headers
}
