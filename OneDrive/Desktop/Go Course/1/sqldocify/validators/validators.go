package validators

import (
	"errors"
	"regexp"
)

type ServerValidator struct{}

func (v *ServerValidator) ValidateConnectionPath(config string) (bool, error) {
	// Define a regex pattern for the connection path format
	pattern := `^[a-zA-Z0-9]+:[^@]+@tcp\((\d+\.){3}\d+:\d+\)/[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, config)
	if err != nil {
		return false, err
	}

	if !matched {
		return false, errors.New("invalid connection path format")
	}
	return true, nil
}
