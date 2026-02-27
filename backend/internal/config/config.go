package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds the configuration for the application.
type Config struct {
	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
	MLSidecar struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"ml_sidecar"`
	Auth struct {
		OktaDomain   string `mapstructure:"okta_domain"`
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	} `mapstructure:"auth"`
	TLS struct {
		Enable    bool   `mapstructure:"enable"`
		CertFile  string `mapstructure:"cert_file"`
		KeyFile   string `mapstructure:"key_file"`
		Hostnames []string `mapstructure:"hostnames"`
	} `mapstructure:"tls"`
}

// LoadConfig loads the configuration from a file and the environment.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// normalize OKTA issuer url (strip trailing slash if any)
	config.Auth.OktaDomain = normalizeOktaIssuer(config.Auth.OktaDomain)

	return &config, nil
}

// normalizeOktaIssuer ensures the provided Okta issuer string is in a
// predictable form. It removes any trailing slash and leaves the scheme and
// path intact. This allows users to paste the full URL from the Okta admin
// console without worrying about double prefixes.
func normalizeOktaIssuer(input string) string {
	iss := strings.TrimSpace(input)
	if strings.HasSuffix(iss, "/") {
		iss = strings.TrimRight(iss, "/")
	}
	return iss
}
