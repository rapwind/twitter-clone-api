package jsonschema

import (
	"github.com/techcampman/twitter-d-server/env"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// V1CreateInstallationDocument ... json schema for POST /v1/apps
	V1CreateInstallationDocument *gojsonschema.Schema
	// V1UpdateInstallationDocument ... json schema for POST /v1/apps
	V1UpdateInstallationDocument *gojsonschema.Schema

	//
	// schema definitions
	//

	v1CreateInstallationSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"clientType":  clientType,
			"deviceToken": deviceToken,
		},
		"additionalProperties": false,
		"required":             []interface{}{"clientType"},
	}

	v1UpdateInstallationSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"deviceToken": deviceToken,
			"arnEndpoint": arnEndpoint,
		},
		"additionalProperties": false,
	}
)

func initInstallationSchema() {

	var err error
	V1CreateInstallationDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1CreateInstallationSchema))
	env.AssertErrForInit(err)

	V1UpdateInstallationDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1UpdateInstallationSchema))
	env.AssertErrForInit(err)

	return
}
