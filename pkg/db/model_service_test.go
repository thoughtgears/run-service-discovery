package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_ValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid Cloud Run URL", "https://example-service-12345-uc.a.run.app", false},
		{"Valid URL format", "https://example.com", false},
		{"Empty URL", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{URL: tt.url}
			err := s.ValidateURL()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_ValidateEnv(t *testing.T) {
	tests := []struct {
		name        string
		environment Environment
		wantErr     bool
	}{
		{"Valid prod environment", Prod, false},
		{"Valid staging environment", Staging, false},
		{"Valid dev environment", Dev, false},
		{"Invalid environment", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{Environment: tt.environment}
			err := s.ValidateEnv()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			"Valid service",
			Service{
				Name:        "example-service",
				URL:         "https://example-service-12345-uc.a.run.app",
				Environment: Prod,
			},
			false,
		},
		{
			"Valid URL format",
			Service{
				Name:        "example-service",
				URL:         "https://example.com",
				Environment: Prod,
			},
			false,
		},
		{
			"Invalid environment",
			Service{
				Name:        "example-service",
				URL:         "https://example-service-12345-uc.a.run.app",
				Environment: "invalid",
			},
			true,
		},
		{
			"Empty name",
			Service{
				Name:        "",
				URL:         "https://example-service-12345-uc.a.run.app",
				Environment: Prod,
			},
			true,
		},
		{
			"Empty URL",
			Service{
				Name:        "example-service",
				URL:         "",
				Environment: Prod,
			},
			true,
		},
		{
			"Empty environment",
			Service{
				Name:        "example-service",
				URL:         "https://example-service-12345-uc.a.run.app",
				Environment: "",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
