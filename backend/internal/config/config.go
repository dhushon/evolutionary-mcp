package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds the configuration for the application.
type Config struct {
	Environment   string `mapstructure:"environment"`
	DevModeBypass bool   `mapstructure:"dev_mode_bypass"`
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
		OktaDomain      string `mapstructure:"okta_domain"`
		ClientID        string `mapstructure:"client_id"`
		ClientSecret    string `mapstructure:"client_secret"`
		SwaggerClientID string `mapstructure:"swagger_client_id"`
		RedirectURL     string `mapstructure:"redirect_url"`
	} `mapstructure:"auth"`
	TLS struct {
		Enable    bool     `mapstructure:"enable"`
		CertFile  string   `mapstructure:"cert_file"`
		KeyFile   string   `mapstructure:"key_file"`
		Hostnames []string `mapstructure:"hostnames"`
	} `mapstructure:"tls"`
}

// LoadConfig loads the configuration from a file and the environment.
func LoadConfig(envPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Also try to load .env file from project root or current dir
	viper.SetConfigType("env")
	if envPath != "" {
		viper.SetConfigFile(envPath)
		_ = viper.MergeInConfig()
	} else {
		viper.SetConfigFile(".env")
		if err := viper.MergeInConfig(); err != nil {
			viper.SetConfigFile("../.env")
			_ = viper.MergeInConfig()
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if env := viper.GetString("ENVIRONMENT"); env != "" {
		config.Environment = env
	}
	if bypass := viper.GetBool("DEV_MODE_BYPASS"); bypass {
		config.DevModeBypass = bypass
	}

	// env overrides (especially useful in containerized environments)
	if h := viper.GetString("DB_HOST"); h != "" {
		config.DB.Host = h
	}
	if p := viper.GetInt("DB_PORT"); p != 0 {
		config.DB.Port = p
	}
	if u := viper.GetString("DB_USER"); u != "" {
		config.DB.User = u
	}
	if pw := viper.GetString("DB_PASSWORD"); pw != "" {
		config.DB.Password = pw
	}
	if n := viper.GetString("DB_NAME"); n != "" {
		config.DB.Name = n
	}
	if s := viper.GetString("DB_SSLMODE"); s != "" {
		config.DB.SSLMode = s
	}
	if url := viper.GetString("ML_SIDECAR_URL"); url != "" {
		config.MLSidecar.URL = url
	}

	if d := viper.GetString("AUTH_OKTA_DOMAIN"); d != "" {
		config.Auth.OktaDomain = d
	}
	if cid := viper.GetString("AUTH_CLIENT_ID"); cid != "" {
		config.Auth.ClientID = cid
	}
	if cs := viper.GetString("AUTH_CLIENT_SECRET"); cs != "" {
		config.Auth.ClientSecret = cs
	}
	if scid := viper.GetString("AUTH_SWAGGER_CLIENT_ID"); scid != "" {
		config.Auth.SwaggerClientID = scid
	}
	if r := viper.GetString("AUTH_REDIRECT_URL"); r != "" {
		config.Auth.RedirectURL = r
	}

	// normalize OKTA issuer url (strip trailing slash if any)
	config.Auth.OktaDomain = normalizeOktaIssuer(config.Auth.OktaDomain)

	// Default Swagger Client ID to the main Client ID if not set (backward compatibility)
	if config.Auth.SwaggerClientID == "" {
		config.Auth.SwaggerClientID = config.Auth.ClientID
	}

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
