package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	EnvKeyPort         int    `env:"PORT" env-default:"4000"`
	EnvKeyEnvironment  string `env:"ENVIRONMENT" env-default:"local"`
	EnvKeyAppUrl       string `env:"APP_URL" env-default:"http://localhost:3000/"`
	EnvKeySignKey      string `env:"SIGN_KEY" env-default:"8kzz3S4mVmx3BodlQiba"`
	EnvKeyContactEmail string `env:"CONTACT_EMAIL" env-default:"nicolas.dev.py@gmail.com"`
	EnvKeySendgridKey  string `env:"SENDGRID_KEY" env-default:"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"`

	EnvKeyPermissionHost  string `env:"PERMISSION_HOST" env-default:"localhost:50051"`
	EnvKeyPermissionToken string `env:"PERMISSION_TOKEN" env-default:"afd9ec25-4bd7-4884-8106-3822d7e721b4"`

	EnvKeyDbHost     string `env:"DB_HOST" env-default:"dpg-cps6o2aj1k6c738in0p0-a"`
	EnvKeyDbName     string `env:"DB_NAME" env-default:"ghhapi"`
	EnvKeyDbUserName string `env:"DB_USER" env-default:"postgre"`
	EnvKeyDbPort     string `env:"DB_PORT" env-default:"5432"`
	EnvKeyDbPassword string `env:"DB_PASS" env-default:"yp7qQzci9yxxzuePek8e3nY1onOvY6AK"`
	EnvKeyDbSSL      string `env:"DB_SSL" env-default:"disable"`
	EnvKeyDbSSLCa    string `env:"DB_SSL_CA" env-default:"disable"`
	EnvKeyDbSSLCert  string `env:"DB_SSL_CERT" env-default:"disable"`
	EnvKeyDbSSLKey   string `env:"DB_SSL_KEY" env-default:"disable"`
	EnvKeyDbTimeZone string `env:"DB_TIMEZONE" env-default:"America/Santiago"`

	EnvKeyWsChileUser string `env:"WS_USER" env-default:"AutentiaManager"`
	EnvKeyWsChilePass string `env:"WS_PASS" env-default:"4gj=U4A%F5"`
	EnvKeyWsChileOper string `env:"WS_OPER" env-default:"0000005555-7"`
	EnvKeyWsChileUrl  string `env:"WS_URL" env-default:"http://cap.autentia.cl/"`

	EnvKeyWsLatamUser string `env:"WS_USER_LATAM" env-default:"AutentiaManager"`
	EnvKeyWsLatamPass string `env:"WS_PASS_LATAM" env-default:"4gj=U4A%F5"`
	EnvKeyWsLatamOper string `env:"WS_OPER_LATAM" env-default:"0000005555-7"`
	EnvKeyWsLatamUrl  string `env:"WS_URL_LATAM" env-default:"http://stg.autentia.io/"`
}

func New() *Config {
	var cfg = Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		fmt.Println("Can't read the config")
		os.Exit(1)
	}

	return &cfg
}

func (c *Config) Environment() string {
	return c.EnvKeyEnvironment
}

func (c *Config) Port() string {
	return fmt.Sprintf(":%d", c.EnvKeyPort)
}

func (c *Config) AppUrl() string {
	return c.EnvKeyAppUrl
}

func (c *Config) PermissionHost() string {
	return c.EnvKeyPermissionHost
}

func (c *Config) PermissionToken() string {
	return c.EnvKeyPermissionToken
}

func (c *Config) SendgridKey() string {
	return c.EnvKeySendgridKey
}

func (c *Config) ContactEmail() string {
	return c.EnvKeyContactEmail
}

func (c *Config) SignKey() string {
	return c.EnvKeySignKey
}

func (c *Config) DbHost() string {
	return c.EnvKeyDbHost
}

func (c *Config) DbName() string {
	return c.EnvKeyDbName
}

func (c *Config) DbUserName() string {
	return c.EnvKeyDbUserName
}

func (c *Config) DbPort() string {
	return c.EnvKeyDbPort
}

func (c *Config) DbPassword() string {
	return c.EnvKeyDbPassword
}

func (c *Config) DbSSL() string {
	return c.EnvKeyDbSSL
}

func (c *Config) DbSSLCa() string {
	return c.EnvKeyDbSSLCa
}

func (c *Config) DbSSLCert() string {
	return c.EnvKeyDbSSLCert
}

func (c *Config) DbSSLKey() string {
	return c.EnvKeyDbSSLKey
}

func (c *Config) DbTimeZone() string {
	return c.EnvKeyDbTimeZone
}

func (c *Config) WsUser(country string) string {
	if country == CHILE {
		return c.EnvKeyWsChileUser
	}
	return c.EnvKeyWsLatamUser
}

func (c *Config) WsPass(country string) string {
	if country == CHILE {
		return c.EnvKeyWsChilePass
	}
	return c.EnvKeyWsLatamPass
}

func (c *Config) WsOper(country string) string {
	if country == CHILE {
		return c.EnvKeyWsChileOper
	}
	return c.EnvKeyWsLatamOper
}

func (c *Config) WsUrl(country string) string {
	if country == CHILE {
		return c.EnvKeyWsChileUrl
	}
	return c.EnvKeyWsLatamUrl
}

func (c *Config) Path(country string) string {
	if country == CHILE {
		return ""
	}
	return ""
}
