package videoservice

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"path"
	"strings"
)

const (
	generateCDNUrlFunc      string = "GenerateCDNUrl"
	validateOriginalURLFunc string = "ValidateOriginalURL"
)

const httpScheme string = "http"

var (
	ErrParsingURL       = errors.New("error parsing url")
	ErrEmptyFirstDomain = errors.New("empty first domain")
	ErrOriginalURLNil   = errors.New("original url is nil")
)

type VideoService struct {
	CDNHost string
	logger  *zap.Logger
}

func NewVideoService(CDNHost string, logger *zap.Logger) *VideoService {
	return &VideoService{
		CDNHost: CDNHost,
		logger:  logger,
	}
}

func (s *VideoService) ValidateOriginalURL(rawOriginalURL string) (*url.URL, error) {
	log := s.logger.With(zap.String("func", validateOriginalURLFunc))

	rawOriginalURL = strings.TrimSpace(rawOriginalURL)
	originalURL, err := url.ParseRequestURI(rawOriginalURL)
	if err != nil {
		log.Error("url.ParseRequestURI", zap.Error(err))
		return originalURL, fmt.Errorf("%w - %w", ErrParsingURL, err)
	}

	return originalURL, nil
}

func (s *VideoService) GenerateCDNUrl(originalURL *url.URL) (string, error) {
	log := s.logger.With(zap.String("func", generateCDNUrlFunc))
	if originalURL == nil {
		return "", ErrOriginalURLNil
	}
	log.Info("starting generating cdn url", zap.String("original url", originalURL.String()))
	var cdnURL url.URL
	domains := strings.Split(originalURL.Hostname(), ".")
	firstDomain := domains[0]
	if firstDomain == "" {
		log.Error("empty first domain", zap.Error(ErrEmptyFirstDomain))
		return "", ErrEmptyFirstDomain
	}

	cdnURL.Scheme = httpScheme
	cdnURL.Host = s.CDNHost
	cdnURL.Path = firstDomain
	cdnURL.Path = path.Join(cdnURL.Path, originalURL.Path)
	return cdnURL.String(), nil
}
