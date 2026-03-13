package boundary

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/msskobelina/fit-profi/internal/delivery/controller"
)

type jsonIO struct {
	validator *validator.Validate
}

func New() controller.IO {
	return &jsonIO{
		validator: validator.New(),
	}
}

func (c *jsonIO) Read(request interface{}, reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(request); err != nil {
		return errors.New("invalid request body")
	}

	if err := c.validator.Struct(request); err != nil {
		return err
	}

	return nil
}

func (c *jsonIO) Error(err error, r *http.Request, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func (c *jsonIO) Fatal(err error, r *http.Request, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "internal server error",
	})
}

func (c *jsonIO) Result(response interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
