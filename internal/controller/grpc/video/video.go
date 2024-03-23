package videogrpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"video-balancer/internal/service"
	videopb "video-balancer/proto/video/pb"
)

type videoHandler struct {
	video                            service.Video
	videopb.UnsafeVideoServiceServer // make sure we didn't forget to implement its methods
}

func Register(gRPCSServer *grpc.Server, video service.Video) {
	videopb.RegisterVideoServiceServer(gRPCSServer, &videoHandler{
		video: video,
	})
}

func (h *videoHandler) RedirectVideo(ctx context.Context, req *videopb.RedirectVideoRequest) (*videopb.RedirectVideoResponse, error) {
	videoURL, err := h.video.RedirectVideo(req.GetVideo())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &videopb.RedirectVideoResponse{
		VideoURL: videoURL,
	}, nil
}
