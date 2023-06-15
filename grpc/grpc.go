package grpc

import (
	"organization_service/config"
	"organization_service/genproto/organization_service"
	"organization_service/grpc/client"
	"organization_service/grpc/service"
	"organization_service/pkg/logger"
	"organization_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	organization_service.RegisterFilialServiceServer(grpcServer, service.NewFilialService(cfg, log, strg, srvc))
	organization_service.RegisterMagazinServiceServer(grpcServer, service.NewMagazinService(cfg, log, strg, srvc))
	organization_service.RegisterProviderServiceServer(grpcServer, service.NewProviderService(cfg, log, strg, srvc))
	organization_service.RegisterStaffServiceServer(grpcServer, service.NewStaffService(cfg, log, strg, srvc))

	reflection.Register(grpcServer)
	return
}
