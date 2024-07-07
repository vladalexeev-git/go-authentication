package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
	oauth2google "golang.org/x/oauth2/google"
	"os"
	"time"
)

type (
	Config struct {
		HTTP        `yaml:"http"`
		Logger      `yaml:"logger"`
		Postgres    `yaml:"postgres"`
		AccessToken `yaml:"access_token"`
		Session     `yaml:"session"`
		MongoDB     `yaml:"mongodb"`
		CSRFToken   `yaml:"csrf-token"`
		Redis       `yaml:"redis"`
		SocialAuth  `yaml:"social_auth"`
	}

	HTTP struct {
		Port             string `yaml:"port"`
		CorsAllowOrigins string `yaml:"cors_allow_origins"`
	}

	Logger struct {
		Env string `yaml:"env"`
	}

	Postgres struct {
		PoolMax int    `yaml:"pool_max"`
		URL     string `env:"PG_URL"`
	}

	Session struct {
		TTL          time.Duration `yaml:"ttl"`
		CookieKey    string        `yaml:"cookie_key"`
		CookieDomain string        `yaml:"cookie_domain"`
		//CookiePath     string        `yaml:"cookie_path"`
		CookieSecure   bool `yaml:"cookie_secure"`
		CookieHttpOnly bool `yaml:"cookie_httponly"`
	}

	CSRFToken struct {
		CSRFttl       time.Duration `yaml:"ttl"`
		CSRFCookieKey string        `yaml:"cookie_key"`
		CSRFHeaderKey string        `yaml:"header_key"`
	}

	AccessToken struct {
		TTL        time.Duration `yaml:"ttl"`
		SigningKey string        `yaml:"signing_key"`
	}

	MongoDB struct {
		DbName   string `yaml:"db_name"`
		URI      string `yaml:"uri" env:"MONGO_URI"`
		Username string `yaml:"username" env:"MONGO_USER"`
		Password string `yaml:"password" env:"MONGO_PASS"`
	}

	Redis struct {
		Addr     string `env-required:"true" env:"REDIS_ADDR"`
		Password string `env-required:"true" env:"REDIS_PASSWORD"`
	}
)

type SocialAuth struct {
	GitHubClientID     string `yaml:"github_client_id" env-required:"true" env:"GH_CLIENT_ID"`
	GitHubClientSecret string `env-required:"true" env:"GH_CLIENT_SECRET"`
	GitHubScope        string `yaml:"github_scope" env-required:"true" env:"GH_SCOPE"`

	GoogleClientID     string `yaml:"google_client_id" env-required:"true" env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env-required:"true" env:"GOOGLE_CLIENT_SECRET"`
	GoogleScope        string `yaml:"google_scope" env-required:"true" env:"GOOGLE_SCOPE"`
}

func (sa *SocialAuth) Endpoints() map[string]oauth2.Endpoint {
	return map[string]oauth2.Endpoint{
		"github": oauth2github.Endpoint,
		"google": oauth2google.Endpoint,
	}
}

func (sa *SocialAuth) Scopes() map[string]string {
	return map[string]string{
		"google": sa.GoogleScope,
	}
}

func (sa *SocialAuth) ClientIDs() map[string]string {
	return map[string]string{
		"google": sa.GoogleClientID,
	}
}

func (sa *SocialAuth) ClientSecrets() map[string]string {
	return map[string]string{
		"google": sa.GoogleClientSecret,
	}
}

func MustLoad() *Config {
	var cfg Config

	configPath := fetchConfigPath()

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
