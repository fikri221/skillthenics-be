package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return NewError(http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
	}

	if err := validate.Struct(v); err != nil {
		return NewError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
	}

	return nil
}
