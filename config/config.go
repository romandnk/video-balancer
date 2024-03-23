package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/fx"
	"video-balancer/pkg/grpcserver"
	zaplogger "video-balancer/pkg/logger/zap"
)

var Module = fx.Module("config", fx.Provide(NewConfig))

const configPath string = "./config/config.yaml"

type Config struct {
	ZapLogger  zaplogger.Config  `yaml:"zap_logger"`
	GRPCServer grpcserver.Config `yaml:"grpc_server"`
	CDNHost    string            `env:"CDN_HOST" env-required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return &cfg, fmt.Errorf("NewConfig - cleanenv.ReadConfig - %w", err)
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("NewConfig - cleanenv.UpdateEnv - %w", err)
	}
	fmt.Printf("%+v\n", cfg.GRPCServer)
	return &cfg, nil
}
