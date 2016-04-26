package jsonschema

import (
	"github.com/techcampman/twitter-d-server/env"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// V1PostTweetDocument ... json schema for POST /v1/tweets
	V1PostTweetDocument *gojsonschema.Schema

	//
	// schema definitions
	//

	v1PostTweetSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"text":               text,
			"contentUrl":         contentURL,
			"inReplyToTweetId":   inReplyToTweetID,
			"inRetweetToTweetId": inRetweetToTweetID,
		},
		"additionalProperties": false,
		"required":             []interface{}{"text"},
	}
)

func initTweetSchema() {

	var err error
	V1PostTweetDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1PostTweetSchema))
	env.AssertErrForInit(err)

	return
}
