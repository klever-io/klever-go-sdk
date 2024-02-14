package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/klever-io/klever-go-sdk/provider/utils/http_options"
)

type httpClient struct {
	http.Client
}

const defaultUserAgent = "kleversdk/1.0"

func NewHttpClient(timeout time.Duration) HttpClient {
	return &httpClient{http.Client{Timeout: timeout}}
}

func newFromError(errMessage []byte) error {
	message := string(errMessage)

	// check if error marshal
	var iErr map[string]interface{}
	err := json.Unmarshal(errMessage, &iErr)
	if err == nil {
		v, ok := iErr["error"]
		if ok {
			if message, ok = v.(string); !ok {
				message = string(errMessage)
			}
		}
	}
	return fmt.Errorf("%s", message)
}

// GetURL provides json result decode to struct
func (h *httpClient) Get(ctx context.Context, url string, target interface{}, options ...http_options.IOptions) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	http_options.DoOptions(req, options...)
	r, err := h.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode >= 400 {
		marshalError := http_options.DodMarshalError(options...)

		if marshalError {
			// try to use model
			if err := json.Unmarshal(body, &target); err == nil {
				return fmt.Errorf("request failed status %d", r.StatusCode)
			}
			return fmt.Errorf("request failed status %d - message :%v", r.StatusCode, string(body))
		}

		return newFromError(body)

	}

	return json.Unmarshal(body, target)
}

// Post provides a post using a json string
func (h *httpClient) Post(ctx context.Context, url string, body string, headers []string, target interface{}, options ...http_options.IOptions) error {
	reqBody := strings.NewReader(body)
	req, errNewReq := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if errNewReq != nil {
		return errNewReq
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Add("Content-type", "application/json; charset=UTF-8")
	for i := 0; i < len(headers); i += 2 {
		req.Header.Add(headers[i], headers[i+1])
	}

	r, err := h.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	data, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if errRead != nil {
		return errRead
	}

	if r.StatusCode >= 400 {
		marshalError := http_options.DodMarshalError(options...)
		// try to use model
		if marshalError {
			if err := json.Unmarshal(data, &target); err == nil {
				return fmt.Errorf("request failed status %d", r.StatusCode)
			}
			return fmt.Errorf("request failed status %d - message :%v", r.StatusCode, string(data))
		}

		return newFromError(data)
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}
	return nil
}
