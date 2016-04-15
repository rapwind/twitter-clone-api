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

	name = map[string]interface{}{
		"type":      "string",
		"minLength": 1.0,
		"maxLength": 20.0,
	}

	phoneNumber = map[string]interface{}{
		"type":      "string",
		"minLength": 10.0,
		"maxLength": 20.0,
	}

	email = map[string]interface{}{
		"type":      "string",
		"format":    "email",
		"minLength": 4.0,
		"maxLength": 256.0,
	}

	screenName = map[string]interface{}{
		"type":      "string",
		"minLength": 4.0,
		"maxLength": 15.0,
	}

	accountName = map[string]interface{}{
		"type":      "string",
		"minLength": 4.0,
		"maxLength": 256.0,
	}

	passwordHash = map[string]interface{}{
		"type":      "string",
		"minLength": 1.0,
		"maxLength": 64.0,
	}
)
