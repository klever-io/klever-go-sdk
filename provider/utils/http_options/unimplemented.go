package http_options

var _ IOptions = (*HTTPOptions)(nil)

type HTTPOptions struct {
	Type OptionType `json:"type"`
}

func (s *HTTPOptions) GetType() OptionType {
	return s.Type
}

func (s *HTTPOptions) GetHeaders() [][2]string {
	return make([][2]string, 0)
}

func (s *HTTPOptions) GetMarshalError() bool {
	return false
}
