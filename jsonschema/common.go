package jsonschema

var (
	clientType = map[string]interface{}{
		"type":      "string",
		"minLength": 1.0,
		"maxLength": 15.0,
	}

	deviceToken = map[string]interface{}{
		"type":      "string",
		"minLength": 1.0,
		"maxLength": 64.0,
	}

	screenName = map[string]interface{}{
		"type":      "string",
		"minLength": 1.0,
		"maxLength": 15.0,
	}

	passwordHash = map[string]interface{}{
		"type":      "string",
		"minLength": 1.0,
		"maxLength": 64.0,
	}
)
