package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/fx"
	"strings"
	"video-balancer/pkg/grpcserver"
	zaplogger "video-balancer/pkg/logger/zap"
)

var Module = fx.Module("config", fx.Provide(NewConfig))

const configPath string = "./config/config.yaml"

var (
	ErrEmptyCDNHost = errors.New("CDN_HOST env cannot be empty")
)

type Config struct {
	ZapLogger  zaplogger.Config  `yaml:"zap_logger"`
	GRPCServer grpcserver.Config `yaml:"grpc_server"`
	CDNHost    string            `env:"CDN_HOST" env-required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("NewConfig - cleanenv.ReadConfig - %w", err)
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("NewConfig - cleanenv.UpdateEnv - %w", err)
	}

	CDNHost, err := validateCDNHost(cfg.CDNHost)
	if err != nil {
		return nil, fmt.Errorf("NewConfig - validateCDNHost - %w", err)
	}
	cfg.CDNHost = CDNHost

	return &cfg, nil
}

func validateCDNHost(CDNHost string) (string, error) {
	CDNHost = strings.TrimSpace(CDNHost)
	if CDNHost == "" {
		return CDNHost, ErrEmptyCDNHost
	}
	return CDNHost, nil
}
