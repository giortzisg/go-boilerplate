package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	InvalidContentTypeError = "Invalid content type"
)

func Decoder[RequestType any](r *http.Request) (*RequestType, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("%s: %s", InvalidContentTypeError, r.Header.Get("Content-Type"))
	}

	var input RequestType
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			return nil, err
		}
	}

	return &input, nil
}

func Encoder[ResponseType any](w http.ResponseWriter, data *ResponseType, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}
