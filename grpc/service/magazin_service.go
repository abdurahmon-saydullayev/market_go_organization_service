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

type MagazinService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*organization_service.UnimplementedMagazinServiceServer
}

func NewMagazinService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *MagazinService {
	return &MagazinService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *MagazinService) Create(ctx context.Context, req *organization_service.CreateMagazin) (resp *organization_service.Magazin, err error) {

	i.log.Info("---CreateMagazin----->", logger.Any("req", req))

	pKey, err := i.strg.Magazin().Create(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateMagazin>Magazin>Create--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err = i.strg.Magazin().GetByID(ctx, pKey)
	if err != nil {
		i.log.Error("!!!GetByPKeyMagazin>Magazin>Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *MagazinService) GetByID(ctx context.Context, req *organization_service.MagazinPK) (resp *organization_service.Magazin, err error) {

	i.log.Info("---GetMagazinByID------>", logger.Any("req", req))

	resp, err = i.strg.Magazin().GetByID(ctx, req)
	if err != nil {
		i.log.Error("!!!GetMagazinByID->Magazin>Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *MagazinService) GetList(ctx context.Context, req *organization_service.GetListMagazinRequest) (resp *organization_service.GetListMagazinResponse, err error) {

	i.log.Info("---GetMagazins------>", logger.Any("req", req))

	resp, err = i.strg.Magazin().GetList(ctx, req)
	if err != nil {
		i.log.Error("!!!GetMagazins->Magazin>Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *MagazinService) Update(ctx context.Context, req *organization_service.UpdateMagazin) (resp *organization_service.Magazin, err error) {

	i.log.Info("---UpdateMagazin----->", logger.Any("req", req))

	rowsAffected, err := i.strg.Magazin().Update(ctx, req)

	if err != nil {
		i.log.Error("!!!UpdateMagazin-->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Magazin().GetByID(ctx, &organization_service.MagazinPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetMagazin>Magazin>Get--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *MagazinService) UpdatePatch(ctx context.Context, req *organization_service.UpdatePatchMagazin) (resp *organization_service.Magazin, err error) {

	i.log.Info("---UpdatePatchMagazin----->", logger.Any("req", req))

	updatePatchModel := models.UpdatePatchRequest{
		Id:     req.GetId(),
		Fields: req.GetFields().AsMap(),
	}

	rowsAffected, err := i.strg.Magazin().UpdatePatch(ctx, &updatePatchModel)

	if err != nil {
		i.log.Error("!!!UpdatePatchMagazin-->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Magazin().GetByID(ctx, &organization_service.MagazinPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetMagazin>Magazin>Get--->", logger.Error(err))

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *MagazinService) Delete(ctx context.Context, req *organization_service.MagazinPK) (resp *empty.Empty, err error) {

	i.log.Info("---DeleteMagazin----->", logger.Any("req", req))

	err = i.strg.Magazin().Delete(ctx, req)
	if err != nil {
		i.log.Error("!!!DeleteMagazin>Magazin>Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}
