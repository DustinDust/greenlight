package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"greenlight/internal/constant"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) parseIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, nil
	}
	return id, nil
}

func (app *application) writeJSONResponse(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	j, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	j = append(j, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
	return nil
}

// Unmarshal request JSON body into object & handling error
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf(constant.JSON_SYNTAX_ERROR, syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New(constant.JSON_UNEXPECTED_EOF_ERROR)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(constant.JSON_UNMARSHAL_ERROR_FIELD, unmarshalTypeError.Field)
			}
			return fmt.Errorf(constant.JSON_UNMARSHAL_ERROR, unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New(constant.JSON_EMPTY)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}
	return nil
}
