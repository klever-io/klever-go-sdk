package utils

import "github.com/klever-io/klever-go-sdk/provider/utils/http_options"

type HttpClient interface {
	Get(url string, target interface{}, options ...http_options.IOptions) error
	Post(url string, body string, headers []string, target interface{}, options ...http_options.IOptions) error
}
