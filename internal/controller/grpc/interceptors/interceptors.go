package interceptors

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LoggingInterceptor(logger *zap.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		logInfos := []zap.Field{zap.String("method", info.FullMethod), zap.String("processing time", duration.String())}
		if err != nil {
			logInfos = append(logInfos, zap.String("errors", err.Error()))
		}

		logger.Info("Request info", logInfos...)

		return resp, err
	}
}
