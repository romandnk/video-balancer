package service

import (
	"go.uber.org/zap"
	"net/url"
	videoservice "video-balancer/internal/service/video"
)

type Video interface {
	ValidateOriginalURL(rawOriginalURL string) (*url.URL, error)
	GenerateCDNUrl(originalURL *url.URL) (string, error)
}

type Services struct {
	Video Video
}

func NewServices(CDNHost string, logger *zap.Logger) *Services {
	return &Services{
		Video: videoservice.NewVideoService(CDNHost, logger),
	}
}
