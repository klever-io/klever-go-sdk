package http_options

import "net/http"

func DoOptions(r *http.Request, options ...IOptions) {
	for _, opt := range options {
		switch opt.GetType() {
		case Header:
			DoHeaders(r, opt.GetHeaders())
		}
	}
}

func DoHeaders(r *http.Request, headers [][2]string) {
	for _, header := range headers {
		r.Header.Set(header[0], header[1])
	}
}

func DodMarshalError(options ...IOptions) bool {
	for _, opt := range options {
		switch opt.GetType() {
		case MarshalError:
			return opt.GetMarshalError()
		}
	}
	return false
}
