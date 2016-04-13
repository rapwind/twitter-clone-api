package jsonschema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/errors"
)

var subMatchReg = regexp.MustCompile("\"(.*)\"")

// JSONSchema validates a value by json schema
func JSONSchema(b []byte, s *gojsonschema.Schema) (err error) {

	if b == nil {
		err = fmt.Errorf("target bytes is nil pointer")
		return
	}

	if s == nil {
		err = fmt.Errorf("json schema document is nil pointer")
		return
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		logger.Error(err)
		err = errors.InvalidSyntax()
		return
	}

	if result, _ := s.Validate(gojsonschema.NewGoLoader(m)); !result.Valid() {
		err = convertToResponseErr(result.Errors()[0], m)
	}

	return
}

func convertToResponseErr(e gojsonschema.ResultError, input map[string]interface{}) *errors.ResponseError {

	s := strings.Split(e.Context().String(), ".")
	if len(s) > 1 {
		return errors.BadParams(s[1], fmt.Sprint(input[s[1]]))
	}

	r := subMatchReg.FindStringSubmatch(e.Description())
	if len(r) > 1 {
		return errors.BadParams(r[1], fmt.Sprint(input[r[1]]))
	}
	return errors.BadParams(e.Context().String(), e.Description())
}
