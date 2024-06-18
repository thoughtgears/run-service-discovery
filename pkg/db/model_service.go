package db

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const cloudRunURLRegex = `^https://(?:[-a-zA-Z0-9]+\.)*[-a-zA-Z0-9]+\.[a-zA-Z]{2,}(:[0-9]+)?(/.*)?$`

var cloudRunURLPattern = regexp.MustCompile(cloudRunURLRegex)

type Environment string

const (
	Prod    Environment = "prod"
	Staging Environment = "staging"
	Dev     Environment = "dev"
)

// Service represents a service in the database
// ID is a sha256 hash of the service name
// Name is the name of the service
// URL is the URL of the service
type Service struct {
	ID          string      `json:"-" firestore:"id"`
	Name        string      `json:"name" firestore:"name" validate:"required,min=1"`
	URL         string      `json:"url" firestore:"url" validate:"required,min=1"`
	Environment Environment `json:"environment" firestore:"environment" validate:"required,min=1"`
}

// Validate the service model ensuring URL and Environment are valid
// It also uses the go-playground/validator to validate the struct
func (s *Service) Validate() error {
	validate := validator.New()

	if err := s.ValidateURL(); err != nil {
		return err
	}

	if err := s.ValidateEnv(); err != nil {
		return err
	}

	return validate.Struct(s)
}

// ValidateURL checks if the service URL matches the Cloud Run URL pattern.
func (s *Service) ValidateURL() error {
	if s.URL == "" {
		return errors.New("empty url not allowed")
	}
	if !cloudRunURLPattern.MatchString(s.URL) {
		return errors.New("invalid Cloud Run URL format")
	}
	return nil
}

// ValidateEnv checks if the service environment is valid.
func (s *Service) ValidateEnv() error {
	switch s.Environment {
	case Prod, Staging, Dev:
		return nil
	default:
		return errors.New("invalid environment, allowed values: prod, staging, dev")
	}
}
