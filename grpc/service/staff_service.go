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

type StaffService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*organization_service.UnimplementedStaffServiceServer
}

func NewStaffService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *StaffService {
	return &StaffService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *StaffService) Create(ctx context.Context, req *organization_service.CreateStaff) (resp *organization_service.Staff, err error) {

	i.log.Info("---CreateStaff------>", logger.Any("req", req))

	pKey, err := i.strg.Staff().Create(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateStaff->Staff->Create--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err = i.strg.Staff().GetByID(ctx, pKey)
	if err != nil {
		i.log.Error("!!!GetByPKeyStaff->Staff->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *StaffService) GetByID(ctx context.Context, req *organization_service.StaffPK) (resp *organization_service.Staff, err error) {

	i.log.Info("---GetStaffByID------>", logger.Any("req", req))

	resp, err = i.strg.Staff().GetByID(ctx, req)
	if err != nil {
		i.log.Error("!!!GetStaffByID->Staff->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *StaffService) GetList(ctx context.Context, req *organization_service.GetListStaffRequest) (resp *organization_service.GetListStaffResponse, err error) {

	i.log.Info("---GetStaffs------>", logger.Any("req", req))

	resp, err = i.strg.Staff().GetList(ctx, req)
	if err != nil {
		i.log.Error("!!!GetStaffs->Staff->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *StaffService) Update(ctx context.Context, req *organization_service.UpdateStaff) (resp *organization_service.Staff, err error) {

	i.log.Info("---UpdateStaff------>", logger.Any("req", req))

	rowsAffected, err := i.strg.Staff().Update(ctx, req)

	if err != nil {
		i.log.Error("!!!UpdateStaff--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Staff().GetByID(ctx, &organization_service.StaffPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetStaff->Staff->Get--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *StaffService) UpdatePatch(ctx context.Context, req *organization_service.UpdatePatchStaff) (resp *organization_service.Staff, err error) {

	i.log.Info("---UpdatePatchStaff------>", logger.Any("req", req))

	updatePatchModel := models.UpdatePatchRequest{
		Id:     req.GetId(),
		Fields: req.GetFields().AsMap(),
	}

	rowsAffected, err := i.strg.Staff().UpdatePatch(ctx, &updatePatchModel)

	if err != nil {
		i.log.Error("!!!UpdatePatchStaff--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.Staff().GetByID(ctx, &organization_service.StaffPK{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetStaff->Staff->Get--->", logger.Error(err))

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *StaffService) Delete(ctx context.Context, req *organization_service.StaffPK) (resp *empty.Empty, err error) {

	i.log.Info("---DeleteStaff------>", logger.Any("req", req))

	err = i.strg.Staff().Delete(ctx, req)
	if err != nil {
		i.log.Error("!!!DeleteStaff->Staff->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}
