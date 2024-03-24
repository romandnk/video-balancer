package videogrpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"net/url"
	"testing"
	mockservice "video-balancer/internal/service/mock"
	videopb "video-balancer/proto/video/pb"
)

func startGRPCServer() (*grpc.Server, *bufconn.Listener) {
	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)

	srv := grpc.NewServer()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	return srv, listener
}

func getDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestVideoHandler_RedirectVideo(t *testing.T) {
	type args struct {
		rawOriginalUrl        string
		validatedUrl          url.URL
		cdnUrl                string
		expectedErrorValidate error
		expectedErrorCDN      error
	}

	type mockBehaviour func(m *mockservice.MockVideo, args args)

	testCases := []struct {
		name          string
		args          args
		mock          mockBehaviour
		expectedVideo string
		expectedError error
	}{
		{
			name: "OK - return cdn url",
			args: args{
				rawOriginalUrl: "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8",
				validatedUrl: url.URL{
					Scheme: "http",
					Host:   "s1.origin-cluster",
					Path:   "/video/123/xcg2djHckad.m3u8",
				},
				cdnUrl: "http://cdn.ru/s1/video/123/xcg2djHckad.m3u8",
			},
			mock: func(m *mockservice.MockVideo, args args) {
				m.EXPECT().ValidateOriginalURL(args.rawOriginalUrl).Return(args.validatedUrl, args.expectedErrorValidate)
				m.EXPECT().GenerateCDNUrl(args.validatedUrl).Return(args.cdnUrl, args.expectedErrorCDN)
			},
			expectedVideo: "http://cdn.ru/s1/video/123/xcg2djHckad.m3u8",
		},
		{
			name: "Invalid original url - scheme is empty",
			args: args{
				rawOriginalUrl:        "://s1.origin-cluster/video/123/xcg2djHckad.m3u8",
				expectedErrorValidate: errors.New("error parsing url - parse \"://s1.origin-cluster/video/123/xcg2djHckad.m3u8\": missing protocol scheme"),
			},
			mock: func(m *mockservice.MockVideo, args args) {
				m.EXPECT().ValidateOriginalURL(args.rawOriginalUrl).Return(args.validatedUrl, args.expectedErrorValidate)
			},
			expectedError: errors.New("rpc error: code = InvalidArgument desc = error parsing url - parse \"://s1.origin-cluster/video/123/xcg2djHckad.m3u8\": missing protocol scheme"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv, lis := startGRPCServer()
			defer srv.Stop()
			defer lis.Close()

			videoService := mockservice.NewMockVideo(ctrl)
			Register(srv, videoService)

			ctx := context.Background()

			conn, err := grpc.DialContext(ctx, "",
				grpc.WithContextDialer(getDialer(lis)),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()

			client := videopb.NewVideoServiceClient(conn)

			tc.mock(videoService, tc.args)

			res, err := client.RedirectVideo(ctx, &videopb.RedirectVideoRequest{
				Video: tc.args.rawOriginalUrl,
			})
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedVideo, res.GetVideoURL())
		})
	}
}

func TestVideoHandler_RedirectVideo_DependsOnRequestCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, lis := startGRPCServer()
	defer srv.Stop()
	defer lis.Close()

	videoService := mockservice.NewMockVideo(ctrl)
	Register(srv, videoService)

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := videopb.NewVideoServiceClient(conn)

	rawOriginalUrl := "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8"
	validatedUrl := url.URL{
		Scheme: "http",
		Host:   "s1.origin-cluster",
		Path:   "/video/123/xcg2djHckad.m3u8",
	}
	expectedCDNUrl := "http://cdn.ru/s1/video/123/xcg2djHckad.m3u8"

	const requestCount int = 20
	for requestNum := 1; requestNum <= requestCount; requestNum++ {
		if requestNum%10 == 0 {
			videoService.EXPECT().ValidateOriginalURL(rawOriginalUrl).Return(validatedUrl, nil)
			res, err := client.RedirectVideo(ctx, &videopb.RedirectVideoRequest{
				Video: rawOriginalUrl,
			})

			require.NoError(t, err)
			require.Equal(t, rawOriginalUrl, res.GetVideoURL())
		} else {
			videoService.EXPECT().ValidateOriginalURL(rawOriginalUrl).Return(validatedUrl, nil)
			videoService.EXPECT().GenerateCDNUrl(validatedUrl).Return(expectedCDNUrl, nil)

			res, err := client.RedirectVideo(ctx, &videopb.RedirectVideoRequest{
				Video: rawOriginalUrl,
			})
			require.NoError(t, err)
			require.Equal(t, expectedCDNUrl, res.GetVideoURL())
		}
	}
}
