package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"organization_service/genproto/organization_service"
	"organization_service/models"
	"organization_service/pkg/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type providerRepo struct {
	db *pgxpool.Pool
}

func NewProviderRepo(db *pgxpool.Pool) *providerRepo {
	return &providerRepo{
		db: db,
	}
}

func (c *providerRepo) Create(ctx context.Context, req *organization_service.CreateProvider) (resp *organization_service.ProviderPK, err error) {
	id := uuid.New().String()

	query := `
		INSERT INTO "provider" (
			id,
			name,
			phone,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, NOW(), NOW())
	`

	_, err = c.db.Exec(
		ctx,
		query,
		id,
		req.Name,
		req.Phone,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &organization_service.ProviderPK{Id: id}, nil
}

func (c *providerRepo) GetByID(ctx context.Context, req *organization_service.ProviderPK) (Provider *organization_service.Provider, err error) {
	query := `
		SELECT 
		    id,
			name,
			phone,
			status,
			created_at,
			updated_at
		FROM "provider"
		WHERE id = $1;
	`
	var (
		id         sql.NullString
		name       sql.NullString
		phone      sql.NullString
		status     sql.NullInt32
		created_at sql.NullString
		updated_at sql.NullString
	)

	err = c.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&phone,
		&status,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, err
	}

	Provider = &organization_service.Provider{
		Id:        id.String,
		Name:      name.String,
		Phone:     phone.String,
		Status:    status.Int32,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}

	return
}

func (c *providerRepo) GetList(ctx context.Context, req *organization_service.GetListProviderRequest) (resp *organization_service.GetListProviderResponse, err error) {
	resp = &organization_service.GetListProviderResponse{}

	var (
		query  string
		limit  = ""
		offset = " OFFSET 0 "
		params = make(map[string]interface{})
		filter = " WHERE TRUE "
		sort   = " ORDER BY created_at DESC"
	)

	query = `
	   SELECT 
	   		COUNT(*) OVER(),
			   id,
			   name,
			   phone,
			   status,
			   created_at,
			   updated_at
		FROM "provider" 
	`
	if len(req.GetSearch()) > 0 {
		filter += " AND name ILIKE '%' || '" + req.Search + "' || '%' "
	}
	if req.GetLimit() > 0 {
		limit = " LIMIT :limit"
		params["limit"] = req.Limit
	}
	if req.GetOffset() > 0 {
		offset = " OFFSET :offset"
		params["offset"] = req.Offset
	}
	query += filter + sort + offset + limit

	query, args := helper.ReplaceQueryParams(query, params)

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			phone      sql.NullString
			status     sql.NullInt32
			created_at sql.NullString
			updated_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&phone,
			&status,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return resp, err
		}

		resp.Providers = append(resp.Providers, &organization_service.Provider{
			Id:        id.String,
			Name:      name.String,
			Phone:     phone.String,
			Status:    status.Int32,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}

	return
}

func (c *providerRepo) Update(ctx context.Context, req *organization_service.UpdateProvider) (resp int64, err error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"provider"
		SET
			name = :name,
			phone= :phone,
			status = :status,
			updated_at = now()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":     req.GetId(),
		"name":   req.GetName(),
		"phone":  req.GetPhone(),
		"status": req.GetStatus(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), nil
}

func (c *providerRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error) {
	var (
		set   = " SET "
		ind   = 0
		query string
	)

	if len(req.Fields) == 0 {
		err = errors.New("no updates provided")
		return
	}

	req.Fields["id"] = req.Id

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s ", key, key)
		if ind != len(req.Fields)-1 {
			set += ", "
		}
		ind++
	}

	query = `
		UPDATE
			"provider"
	` + set + ` , updated_at = now()
		WHERE
			id = :id
	`

	query, args := helper.ReplaceQueryParams(query, req.Fields)

	result, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), err
}

func (c *providerRepo) Delete(ctx context.Context, req *organization_service.ProviderPK) error {
	query := `DELETE FROM "provider" WHERE id = $1`

	_, err := c.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
