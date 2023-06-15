package storage

import (
	"context"
	"organization_service/genproto/organization_service"
	"organization_service/models"
)

type StorageI interface {
	CloseDB()
	Filial() FilialRepoI
	Magazin() MagazinRepoI
	Staff() StaffRepoI
	Provider() ProviderRepoI
}

type FilialRepoI interface {
	Create(context.Context, *organization_service.CreateFilial) (*organization_service.FilialPK, error)
	GetByID(context.Context, *organization_service.FilialPK) (*organization_service.Filial, error)
	GetList(context.Context, *organization_service.GetListFilialRequest) (*organization_service.GetListFilialResponse, error)
	Update(context.Context, *organization_service.UpdateFilial) (int64, error)
	UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error)
	Delete(context.Context, *organization_service.FilialPK) error
}

type MagazinRepoI interface {
	Create(context.Context, *organization_service.CreateMagazin) (*organization_service.MagazinPK, error)
	GetByID(context.Context, *organization_service.MagazinPK) (*organization_service.Magazin, error)
	GetList(context.Context, *organization_service.GetListMagazinRequest) (*organization_service.GetListMagazinResponse, error)
	Update(context.Context, *organization_service.UpdateMagazin) (int64, error)
	UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error)
	Delete(context.Context, *organization_service.MagazinPK) error
}

type ProviderRepoI interface {
	Create(context.Context, *organization_service.CreateProvider) (*organization_service.ProviderPK, error)
	GetByID(context.Context, *organization_service.ProviderPK) (*organization_service.Provider, error)
	GetList(context.Context, *organization_service.GetListProviderRequest) (*organization_service.GetListProviderResponse, error)
	Update(context.Context, *organization_service.UpdateProvider) (int64, error)
	UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error)
	Delete(context.Context, *organization_service.ProviderPK) error
}

type StaffRepoI interface {
	Create(context.Context, *organization_service.CreateStaff) (*organization_service.StaffPK, error)
	GetByID(context.Context, *organization_service.StaffPK) (*organization_service.Staff, error)
	GetList(context.Context, *organization_service.GetListStaffRequest) (*organization_service.GetListStaffResponse, error)
	Update(context.Context, *organization_service.UpdateStaff) (int64, error)
	UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error)
	Delete(context.Context, *organization_service.StaffPK) error
}
