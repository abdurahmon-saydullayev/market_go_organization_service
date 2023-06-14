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

type filialRepo struct {
	db *pgxpool.Pool
}

func NewFilialRepo(db *pgxpool.Pool) *filialRepo {
	return &filialRepo{
		db: db,
	}
}

func (c *filialRepo) Create(ctx context.Context, req *organization_service.CreateFilial) (resp *organization_service.FilialPK, err error) {
	id := uuid.New().String()

	filial_code := helper.CombineFirstLetters(req.Name, "")

	query := `
		INSERT INTO "filial" (
			id,
			filial_code,
			name,
			address,
			phone,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err = c.db.Exec(
		ctx,
		query,
		id,
		filial_code,
		req.Name,
		req.Address,
		req.Phone,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &organization_service.FilialPK{Id: id}, nil
}

func (c *filialRepo) GetByID(ctx context.Context, req *organization_service.FilialPK) (order *organization_service.Filial, err error) {
	query := `
		SELECT 
		id,
		filial_code,
		name,
		address,
		phone,
		created_at,
		updated_at
		FROM "filial"
		WHERE id = $1
	`

	var (
		id          sql.NullString
		filial_code sql.NullString
		name        sql.NullString
		address     sql.NullString
		phone       sql.NullString
		created_at  sql.NullString
		updated_at  sql.NullString
	)

	err = c.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&filial_code,
		&name,
		&address,
		&phone,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return order, err
	}

	order = &organization_service.Filial{
		Id:         id.String,
		FilialCode: filial_code.String,
		Name:       name.String,
		Address:    address.String,
		Phone:      phone.String,
		CreatedAt:  created_at.String,
		UpdatedAt:  updated_at.String,
	}

	return
}

func (c *filialRepo) GetList(ctx context.Context, req *organization_service.GetListFilialRequest) (resp *organization_service.GetListFilialResponse, err error) {
	resp = &organization_service.GetListFilialResponse{}

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
		filial_code,
		name,
		address,
		phone,
		TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS'),
		TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	FROM "filial"
	`
	if len(req.GetSearch()) > 0 {
		filter += " AND filial_code ILIKE '%' || '" + req.Search + "' || '%' "
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
			id          sql.NullString
			filial_code sql.NullString
			name        sql.NullString
			address     sql.NullString
			phone       sql.NullString
			created_at  sql.NullString
			updated_at  sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&filial_code,
			&name,
			&address,
			&phone,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return resp, err
		}

		resp.Filials = append(resp.Filials, &organization_service.Filial{
			Id:         id.String,
			FilialCode: filial_code.String,
			Name:       name.String,
			Address:    address.String,
			Phone:      phone.String,
			CreatedAt:  created_at.String,
			UpdatedAt:  updated_at.String,
		})
	}

	return
}

func (c *filialRepo) Update(ctx context.Context, req *organization_service.UpdateFilial) (resp int64, err error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"filial"
		SET
			filial_code = :filial_code,
			name = :name,
			address = :address,
			phone = :phone,
			updated_at = now()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":          req.GetId(),
		"filial_code": req.GetFilialCode(),
		"name":        req.GetName(),
		"address":     req.GetAddress(),
		"phone":       req.GetPhone(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), nil
}

func (c *filialRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error) {
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
			"filial"
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

func (c *filialRepo) Delete(ctx context.Context, req *organization_service.FilialPK) error {
	query := `DELETE FROM "filial" WHERE id = $1`

	_, err := c.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
