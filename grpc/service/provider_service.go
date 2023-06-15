package service

import (
	"context"
	"organization_service/config"
	"organization_service/genproto/organization_service"
	"organization_service/grpc/client"
	"organization_service/models"
	"organization_service/pkg/logger"
	"organization_service/storage"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProviderService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*organization_service.UnimplementedProviderServiceServer
}

func NewProviderService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *ProviderService {
	return &ProviderService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *ProviderService) Create(ctx context.Context, req *organization_service.CreateProvider) (resp *organization_service.Provider, err error) {

	i.log.Info("---CreateProvider------>", logger.Any("req", req))

	pKey, err := i.strg.Provider().Create(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateProvider->Provider->Create--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err = i.strg.Provider().GetByID(ctx, pKey)
	if err != nil {
		i.log.Error("!!!GetByPKeyProvider->Provider->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *ProviderService) GetByID(ctx context.Context, req *organization_service.ProviderPK) (resp *organization_service.Provider, err error) {

	i.log.Info("---GetProviderByID------>", logger.Any("req", req))

	resp, err = i.strg.Provider().GetByID(ctx, req)
	if err != nil {
		i.log.Error("!!!GetProviderByID->Provider->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *ProviderService) GetList(ctx context.Context, req *organization_service.GetListProviderRequest) (resp *organization_service.GetListProviderResponse, err error) {

	i.log.Info("---GetProviders------>", logger.Any("req", req))

	resp, err = i.strg.Provider().GetList(ctx, req)
	if err != nil {
		i.log.Error("!!!GetProviders->Provider->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *ProviderService) Update(ctx context.Context, req *organization_service.UpdateProvider) (resp *organization_service.Provider, err error) {

	i.log.Info("---UpdateProvider------>", logger.Any("req", req))

	rowsAffected, err := i.strg.Provider().Update(ctx, req)

	if err != nil {
		i.log.Error("!!!UpdateProvider--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Provider().GetByID(ctx, &organization_service.ProviderPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetProvider->Provider->Get--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *ProviderService) UpdatePatch(ctx context.Context, req *organization_service.UpdatePatchProvider) (resp *organization_service.Provider, err error) {

	i.log.Info("---UpdatePatchProvider------>", logger.Any("req", req))

	updatePatchModel := models.UpdatePatchRequest{
		Id:     req.GetId(),
		Fields: req.GetFields().AsMap(),
	}

	rowsAffected, err := i.strg.Provider().UpdatePatch(ctx, &updatePatchModel)

	if err != nil {
		i.log.Error("!!!UpdatePatchProvider--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Provider().GetByID(ctx, &organization_service.ProviderPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetProvider->Provider->Get--->", logger.Error(err))

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *ProviderService) Delete(ctx context.Context, req *organization_service.ProviderPK) (resp *empty.Empty, err error) {

	i.log.Info("---DeleteProvider------>", logger.Any("req", req))

	err = i.strg.Provider().Delete(ctx, req)
	if err != nil {
		i.log.Error("!!!DeleteProvider->Provider->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}
