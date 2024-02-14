package utils

import (
	"context"

	"github.com/klever-io/klever-go-sdk/provider/utils/http_options"
)

type HttpClient interface {
	Get(ctx context.Context, url string, target interface{}, options ...http_options.IOptions) error
	Post(ctx context.Context, url string, body string, headers []string, target interface{}, options ...http_options.IOptions) error
}
