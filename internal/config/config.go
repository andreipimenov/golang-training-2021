package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

var c = config()

type Config struct {
	Port             int           `json:"port" env:"PORT" envDefault:"8080"`
	ExternalAPIToken string        `json:"external_api_token" env:"EXTERNAL_API_TOKEN,notEmpty,unset,file"`
	DBConnString     string        `json:"db_conn_string" env:"DB_CONN_STRING,notEmpty,unset,file"`
	DbDriverName     string        `json:"db_driver_name"`
	DbMigrations     string        `json:"db_migrations"`
	LogLevel         int8          `json:"log_level"`
	RepositoryType   string        `json:"repository_type"`
	MongoSettings    MongoSettings `json:"mongo_settings,omitempty"`
}

type MongoSettings struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type Profiles struct {
	Active   string            `json:"active"`
	Profiles map[string]Config `json:"profiles"`
}

func config() *Config {
	cfg, err := yamlConfig()
	if err == nil {
		return cfg
	}
	log.Err(err).Msg("Yaml configuration error")
	cfg, err = envConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Env configuration error")
	}
	return cfg
}

func yamlConfig() (*Config, error) {
	p, err := ioutil.ReadFile("profiles.yml")
	if err != nil {
		return nil, err
	}
	cfg := &Profiles{}
	err = yaml.Unmarshal(p, cfg)
	prf := cfg.Profiles[cfg.Active]
	return &prf, err
}

func envConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}

func Get() Config {
	return *c
}
