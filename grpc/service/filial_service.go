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

type FilialService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*organization_service.UnimplementedFilialServiceServer
}

func NewFilialService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *FilialService {
	return &FilialService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *FilialService) Create(ctx context.Context, req *organization_service.CreateFilial) (resp *organization_service.Filial, err error) {

	i.log.Info("---CreateFilial------>", logger.Any("req", req))

	pKey, err := i.strg.Filial().Create(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateFilial->Filial->Create--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err = i.strg.Filial().GetByID(ctx, pKey)
	if err != nil {
		i.log.Error("!!!GetByPKeyFilial->Filial->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *FilialService) GetByID(ctx context.Context, req *organization_service.FilialPK) (resp *organization_service.Filial, err error) {

	i.log.Info("---GetFilialByID------>", logger.Any("req", req))

	resp, err = i.strg.Filial().GetByID(ctx, req)
	if err != nil {
		i.log.Error("!!!GetFilialByID->Filial->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *FilialService) GetList(ctx context.Context, req *organization_service.GetListFilialRequest) (resp *organization_service.GetListFilialResponse, err error) {

	i.log.Info("---GetFilials------>", logger.Any("req", req))

	resp, err = i.strg.Filial().GetList(ctx, req)
	if err != nil {
		i.log.Error("!!!GetFilials->Filial->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *FilialService) Update(ctx context.Context, req *organization_service.UpdateFilial) (resp *organization_service.Filial, err error) {

	i.log.Info("---UpdateFilial------>", logger.Any("req", req))

	rowsAffected, err := i.strg.Filial().Update(ctx, req)

	if err != nil {
		i.log.Error("!!!UpdateFilial--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Filial().GetByID(ctx, &organization_service.FilialPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetFilial->Filial->Get--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *FilialService) UpdatePatch(ctx context.Context, req *organization_service.UpdatePatchFilial) (resp *organization_service.Filial, err error) {

	i.log.Info("---UpdatePatchFilial------>", logger.Any("req", req))

	updatePatchModel := models.UpdatePatchRequest{
		Id:     req.GetId(),
		Fields: req.GetFields().AsMap(),
	}

	rowsAffected, err := i.strg.Filial().UpdatePatch(ctx, &updatePatchModel)

	if err != nil {
		i.log.Error("!!!UpdatePatchFilial--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Filial().GetByID(ctx, &organization_service.FilialPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetFilial->Filial->Get--->", logger.Error(err))

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *FilialService) Delete(ctx context.Context, req *organization_service.FilialPK) (resp *empty.Empty, err error) {

	i.log.Info("---DeleteFilial------>", logger.Any("req", req))

	err = i.strg.Filial().Delete(ctx, req)
	if err != nil {
		i.log.Error("!!!DeleteFilial->Filial->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}
