package validator

import (
	"errors"
	"net/url"
)

func ValidateURL(rawUrl string) error {
	parsed, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" &&
		parsed.Scheme != "https" {
		return errors.New("invalid scheme")
	}

	if parsed.Host == "" {
		return errors.New("missing host")
	}

	return nil
}
