package app

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"strconv"
	"video-balancer/config"
	videogrpc "video-balancer/internal/controller/grpc/video"
	"video-balancer/internal/service"
	"video-balancer/pkg/grpcserver"
	zaplogger "video-balancer/pkg/logger/zap"
)

func NewApp() fx.Option {
	return fx.Options(
		config.Module,
		ZapLoggerModule(),
		service.Module,
		GRPCServerModule(),

		CheckInitializedModules(),
	)
}

func ZapLoggerModule() fx.Option {
	return fx.Module("zap logger",
		fx.Provide(
			zaplogger.NewLogger,
			func(cfg *config.Config) zaplogger.Config {
				return cfg.ZapLogger
			},
		),
	)
}

func GRPCServerModule() fx.Option {
	return fx.Module("grpc server",
		fx.Provide(
			func(cfg *config.Config) grpcserver.Config {
				return cfg.GRPCServer
			},
			grpcserver.NewServer,
		),
		fx.Invoke(
			func(srv *grpcserver.Server, services *service.Services) {
				videogrpc.Register(srv.Srv, services.Video)
			},
			func(lc fx.Lifecycle, srv *grpcserver.Server, cfg grpcserver.Config, logger *zap.Logger, shutdowner fx.Shutdowner) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := srv.Start(); err != nil {
								logger.Error("error starting GRPC server",
									zap.Error(err),
									zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))),
								)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						srv.Stop()
						return nil
					},
				})
			}),
	)
}

func CheckInitializedModules() fx.Option {
	return fx.Module("check modules",
		fx.Invoke(
			func(cfg *config.Config) {},
			func(logger *zap.Logger) {},
			func(service *service.Services) {},
			func(srv *grpcserver.Server) {},
		),
	)
}
