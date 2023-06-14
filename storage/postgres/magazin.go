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

type magazinRepo struct {
	db *pgxpool.Pool
}

func NewMagazinRepo(db *pgxpool.Pool) *magazinRepo {
	return &magazinRepo{
		db: db,
	}
}

func (c *magazinRepo) Create(ctx context.Context, req *organization_service.CreateMagazin) (resp *organization_service.MagazinPK, err error) {
	id := uuid.New().String()

	query := `
		INSERT INTO "magazin" (
			id,
			name,
			filial_id,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, NOW(), NOW())
	`

	_, err = c.db.Exec(
		ctx,
		query,
		id,
		req.Name,
		req.FilialId,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &organization_service.MagazinPK{Id: id}, nil
}

func (c *magazinRepo) GetByID(ctx context.Context, req *organization_service.MagazinPK) (order *organization_service.Magazin, err error) {
	query := `
			SELECT 
		    m.id,
		    m.name,
		    f.id,
		    m.created_at,
		    m.updated_at
		FROM "magazin" AS m
		JOIN filial AS f ON f.id = m.filial_id
		WHERE m.id = $1;
	`
	var (
		id         sql.NullString
		name       sql.NullString
		filial_id  sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err = c.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&filial_id,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return order, err
	}

	order = &organization_service.Magazin{
		Id:        id.String,
		Name:      name.String,
		FilialId:  filial_id.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}

	return
}

func (c *magazinRepo) GetList(ctx context.Context, req *organization_service.GetListMagazinRequest) (resp *organization_service.GetListMagazinResponse, err error) {
	resp = &organization_service.GetListMagazinResponse{}

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
		    m.id,
		    m.name,
		    f.id,
		    m.created_at,
		    m.updated_at
		FROM "magazin" AS m
		JOIN filial AS f ON f.id = m.filial_id
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
			id         sql.NullString
			name       sql.NullString
			filial_id  sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&filial_id,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return resp, err
		}

		resp.Magazins = append(resp.Magazins, &organization_service.Magazin{
			Id:        id.String,
			Name:      name.String,
			FilialId:  filial_id.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}

	return
}

func (c *magazinRepo) Update(ctx context.Context, req *organization_service.UpdateMagazin) (resp int64, err error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"magazin"
		SET
			name = :name,
			filial_id= :filial_id,
			updated_at = now()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":        req.GetId(),
		"name":      req.GetName(),
		"filial_id": req.GetFilialId(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), nil
}

func (c *magazinRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error) {
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
			"magazin"
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

func (c *magazinRepo) Delete(ctx context.Context, req *organization_service.MagazinPK) error {
	query := `DELETE FROM "magazin" WHERE id = $1`

	_, err := c.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
