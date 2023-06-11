package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d minutes", r)
	quotedJSONVaue := strconv.Quote(jsonValue)
	return []byte(quotedJSONVaue), nil
}

// Custom unmarshalling runtime fields.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// check the format ( "__ minutes" )
	unquotedValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return errors.New("invalid runtime format")
	}

	parts := strings.Split(unquotedValue, " ")
	if len(parts) != 2 || parts[1] != "minutes" {
		return errors.New("invalid runtime format")
	}

	// fail to convert the value back to int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return errors.New("invalid runtime format")
	}

	// set the runtime value to the unmarshal on through pointer dereference
	*r = Runtime(i)
	return nil
}
