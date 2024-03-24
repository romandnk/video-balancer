package videogrpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"sync"
	"video-balancer/internal/service"
	videoservice "video-balancer/internal/service/video"
	videopb "video-balancer/proto/video/pb"
)

const redirectLimitNum uint64 = 10

type videoHandler struct {
	video                            service.Video
	mu                               sync.Mutex
	cdnRequestCount                  map[string]uint64
	videopb.UnsafeVideoServiceServer // make sure we didn't forget to implement its methods
}

func Register(gRPCSServer *grpc.Server, video service.Video) {
	videopb.RegisterVideoServiceServer(gRPCSServer, &videoHandler{
		video:           video,
		mu:              sync.Mutex{},
		cdnRequestCount: make(map[string]uint64),
	})
}

func (h *videoHandler) RedirectVideo(ctx context.Context, req *videopb.RedirectVideoRequest) (*videopb.RedirectVideoResponse, error) {
	var response videopb.RedirectVideoResponse
	rawOriginalURL := req.GetVideo()
	originalURL, err := h.video.ValidateOriginalURL(rawOriginalURL)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	domains := strings.Split(originalURL.Hostname(), ".")
	clusterName := domains[0]
	if clusterName == "" {
		return nil, status.Error(codes.InvalidArgument, videoservice.ErrEmptyClusterName.Error())
	}

	h.mu.Lock()
	requestNum, ok := h.cdnRequestCount[clusterName]
	if !ok {
		h.cdnRequestCount[clusterName] = 1
	} else {
		h.cdnRequestCount[clusterName]++
	}
	h.mu.Unlock()

	if (requestNum+1)%redirectLimitNum == 0 {
		response.VideoURL = rawOriginalURL
	} else {
		cdnURL, err := h.video.GenerateCDNUrl(originalURL, clusterName)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		response.VideoURL = cdnURL
	}

	return &response, nil
}
