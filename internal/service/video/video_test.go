package videoservice

import (
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/url"
	"testing"
)

func TestVideoService_ValidateOriginalURL(t *testing.T) {
	testCases := []struct {
		name           string
		originalURL    string
		expectedResult url.URL
		expectedError  error
	}{
		{
			name:        "OK",
			originalURL: "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8",
			expectedResult: url.URL{
				Scheme: "http",
				Host:   "s1.origin-cluster",
				Path:   "/video/123/xcg2djHckad.m3u8",
			},
			expectedError: nil,
		},
		{
			name:           "Empty scheme",
			originalURL:    "://s1.origin-cluster/video/123/xcg2djHckad.m3u8",
			expectedResult: url.URL{},
			expectedError:  errors.New("error parsing url - parse \"://s1.origin-cluster/video/123/xcg2djHckad.m3u8\": missing protocol scheme"),
		},
		{
			name:           "Space in host name",
			originalURL:    "http:// /video/123/xcg2djHckad.m3u8",
			expectedResult: url.URL{},
			expectedError:  errors.New("error parsing url - parse \"http:// /video/123/xcg2djHckad.m3u8\": invalid character \" \" in host name"),
		},
	}

	testCDNHost := "cdn.ru"
	logger, err := zap.NewProduction()
	require.NoError(t, err)
	videoService := NewVideoService(testCDNHost, logger)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualURL, err := videoService.ValidateOriginalURL(tc.originalURL)
			if tc.expectedError == nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}
			require.Equal(t, tc.expectedResult, actualURL)
		})
	}
}
