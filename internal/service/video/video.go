package videoservice

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync/atomic"
)

const (
	redirectVideoFunc string = "RedirectVideo"
	redirectLimitNum  uint32 = 10
)

var (
	ErrParsingURL = errors.New("error parsing url")
)

type VideoService struct {
	CDNHost    string
	requestNum *atomic.Uint32
	logger     *zap.Logger
}

func NewVideoService(CDNHost string, logger *zap.Logger) *VideoService {
	requestNum := &atomic.Uint32{}
	requestNum.Add(1) // to improve readability of redirectLimitNum constant
	return &VideoService{
		CDNHost:    CDNHost,
		requestNum: requestNum,
		logger:     logger,
	}
}

func (s *VideoService) RedirectVideo(rawVideoURL string) (string, error) {
	log := s.logger.With(zap.String("func", redirectVideoFunc))

	var redirectURL string

	rawVideoURL = strings.TrimSpace(rawVideoURL)
	videoURL, err := url.ParseRequestURI(rawVideoURL)
	if err != nil {
		log.Error("url.ParseRequestURI", zap.Error(err))
		return redirectURL, fmt.Errorf("%w - %w", ErrParsingURL, err)
	}

	if s.requestNum.CompareAndSwap(redirectLimitNum, 0) {
		// TODO: cdn
	} else {
		s.requestNum.Add(1)
		// TODO: another
	}

	videoURL.Host = s.CDNHost

	return videoURL.String(), nil
}
