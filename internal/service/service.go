package service

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"video-balancer/config"
	videoservice "video-balancer/internal/service/video"
)

var Module = fx.Module("services",
	fx.Provide(
		func(cfg *config.Config) string {
			return cfg.CDNHost
		},
		NewServices,
	),
)

type Video interface {
	RedirectVideo(rawVideoURL string) (string, error)
}

type Services struct {
	Video Video
}

func NewServices(CDNHost string, logger *zap.Logger) *Services {
	return &Services{
		Video: videoservice.NewVideoService(CDNHost, logger),
	}
}
