package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DefaultEnrichmentJobType = "data-enrichment"
	// DefaultAuthorizationServerURL is the Camunda Cloud OAuth token endpoint (matches zbc.OAuthDefaultAuthzURL).
	DefaultAuthorizationServerURL = "https://login.cloud.camunda.io/oauth/token/"
)

type Config struct {
	ZeebeAddress                string
	ZeebeClientID               string
	ZeebeClientSecret           string
	AuthorizationServerURL      string
	EnrichmentJobType           string
}

// Load reads configuration from the environment. It loads a local .env file when present
// (missing file is ignored so production can rely on injected env only).
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ZeebeAddress:           strings.TrimSpace(os.Getenv("ZEEBE_ADDRESS")),
		ZeebeClientID:          strings.TrimSpace(os.Getenv("ZEEBE_CLIENT_ID")),
		ZeebeClientSecret:      strings.TrimSpace(os.Getenv("ZEEBE_CLIENT_SECRET")),
		AuthorizationServerURL: strings.TrimSpace(os.Getenv("ZEEBE_AUTHORIZATION_SERVER_URL")),
		EnrichmentJobType:      strings.TrimSpace(os.Getenv("ENRICHMENT_JOB_TYPE")),
	}

	if cfg.EnrichmentJobType == "" {
		cfg.EnrichmentJobType = DefaultEnrichmentJobType
	}

	if cfg.AuthorizationServerURL == "" {
		cfg.AuthorizationServerURL = DefaultAuthorizationServerURL
		// zbc default OAuth provider reads this from the environment.
		_ = os.Setenv("ZEEBE_AUTHORIZATION_SERVER_URL", cfg.AuthorizationServerURL)
	}

	var errs []error
	if cfg.ZeebeAddress == "" {
		errs = append(errs, errors.New("ZEEBE_ADDRESS is required"))
	}
	if cfg.ZeebeClientID == "" {
		errs = append(errs, errors.New("ZEEBE_CLIENT_ID is required"))
	}
	if cfg.ZeebeClientSecret == "" {
		errs = append(errs, errors.New("ZEEBE_CLIENT_SECRET is required"))
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("config: %w", errors.Join(errs...))
	}

	return cfg, nil
}
