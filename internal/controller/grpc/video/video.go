package videogrpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync/atomic"
	"video-balancer/internal/service"
	videoservice "video-balancer/internal/service/video"
	videopb "video-balancer/proto/video/pb"
)

const redirectLimitNum uint32 = 10

type videoHandler struct {
	video                            service.Video
	requestNum                       *atomic.Uint32
	videopb.UnsafeVideoServiceServer // make sure we didn't forget to implement its methods
}

func Register(gRPCSServer *grpc.Server, video service.Video) {
	requestNum := &atomic.Uint32{}
	requestNum.Add(1)
	videopb.RegisterVideoServiceServer(gRPCSServer, &videoHandler{
		video:      video,
		requestNum: requestNum,
	})
}

func (h *videoHandler) RedirectVideo(ctx context.Context, req *videopb.RedirectVideoRequest) (*videopb.RedirectVideoResponse, error) {
	var response videopb.RedirectVideoResponse
	rawOriginalURL := req.GetVideo()
	if h.requestNum.CompareAndSwap(redirectLimitNum, 1) {
		_, err := h.video.ValidateOriginalURL(rawOriginalURL)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		response.VideoURL = rawOriginalURL
		return &response, nil
	} else {
		originalURL, err := h.video.ValidateOriginalURL(rawOriginalURL)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		cdnURL, err := h.video.GenerateCDNUrl(originalURL)
		if err != nil {
			code := codes.InvalidArgument
			if errors.Is(err, videoservice.ErrOriginalURLNil) {
				code = codes.Internal
			}
			return nil, status.Error(code, err.Error())
		}
		response.VideoURL = cdnURL
		h.requestNum.Add(1)
		return &response, nil
	}
}
