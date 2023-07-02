package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"greenlight/internal/constant"
	"greenlight/internal/validator"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// get string from url query values
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// read a string into array of strings, seperated by ","
func (app *application) readCSV(qs url.Values, key string, seperator string, defaultValues []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValues
	}
	return strings.Split(csv, seperator)
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
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
	// 1Mb of maximum bytes
	maxSize := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
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

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxSize)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}
	return nil
}
