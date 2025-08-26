package configuration

import "github.com/Netflix/go-env"

type Http struct {
	Addr string `env:"HTTP_ADDR"`
	Port string `env:"HTTP_PORT"`
}

type PostgreSQL struct {
	Host    string `env:"DB_HOST"`
	Port    string `env:"DB_PORT"`
	User    string `env:"DB_USER"`
	Pass    string `env:"DB_PASS"`
	Name    string `env:"DB_NAME"`
	SSLMode string `env:"DB_SSL_MODE"`
}

type S3 struct {
	Endpoint        string `env:"S3_ENDPOINT"`
	AccessKeyID     string `env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY"`
	BucketName      string `env:"S3_BUCKET_NAME"`
	UseSSL          bool   `env:"S3_USE_SSL"`
}

type Mail struct {
	Host     string `env:"MAIL_SMTP_HOST"`
	Port     int    `env:"MAIL_SMTP_PORT"`
	Username string `env:"MAIL_SMTP_USERNAME"`
	Password string `env:"MAIL_SMTP_PASSWORD"`
	FromName string `env:"MAIL_FROM_NAME"`
}

type Logger struct {
	OutputPath      string `env:"LOGGER_OUTPUT_PATH"`
	ErrorOutputPath string `env:"LOGGER_ERROR_OUTPUT_PATH"`
	ServiceName     string `env:"LOGGER_SERVICE_NAME"`
}

type Config struct {
	Http       Http
	PostgreSQL PostgreSQL
	S3         S3
	Mail       Mail
	Logger     Logger
}

func Load() (*Config, error) {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
