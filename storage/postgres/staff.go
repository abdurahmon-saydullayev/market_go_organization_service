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

type staffRepo struct {
	db *pgxpool.Pool
}

func NewStaffRepo(db *pgxpool.Pool) *staffRepo {
	return &staffRepo{
		db: db,
	}
}

func (c *staffRepo) Create(ctx context.Context, req *organization_service.CreateStaff) (resp *organization_service.StaffPK, err error) {
	id := uuid.New().String()

	query := `
		INSERT INTO "staff" (
			id,
			first_name,
			last_name,
			phone,
			login,
			password,
			staff_type,
			magazin_id,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`

	_, err = c.db.Exec(
		ctx,
		query,
		id,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Login,
		req.Password,
		req.StaffType,
		req.MagazinId,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &organization_service.StaffPK{Id: id}, nil
}

func (c *staffRepo) GetByID(ctx context.Context, req *organization_service.StaffPK) (staff *organization_service.Staff, err error) {
	query := `
			SELECT 
		    s.id,
		    s.first_name,
			s.last_name,
			s.phone,
			s.login,
			s.password,
			s.staff_type,
		    m.id,
		    s.created_at,
		    s.updated_at
		FROM "staff" AS s
		JOIN "magazin" AS m ON m.id = s.magazin_id
		WHERE s.id = $1;
	`
	var (
		id         sql.NullString
		first_name sql.NullString
		last_name  sql.NullString
		phone      sql.NullString
		login      sql.NullString
		password   sql.NullString
		staff_type sql.NullString
		magazin_id sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err = c.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&first_name,
		&last_name,
		&phone,
		&login,
		&password,
		&staff_type,
		&magazin_id,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return staff, err
	}

	staff = &organization_service.Staff{
		Id:        id.String,
		FirstName: first_name.String,
		LastName:  last_name.String,
		Phone:     phone.String,
		Login:     login.String,
		Password:  password.String,
		StaffType: staff_type.String,
		MagazinId: magazin_id.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}

	return
}

func (c *staffRepo) GetList(ctx context.Context, req *organization_service.GetListStaffRequest) (resp *organization_service.GetListStaffResponse, err error) {
	resp = &organization_service.GetListStaffResponse{}

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
			s.id,
			s.first_name,
			s.last_name,
			s.phone,
			s.login,
			s.password,
			s.staff_type,
			m.id,
			s.created_at,
			s.updated_at
		FROM "staff" AS m
		JOIN filial AS f ON f.id = m.filial_id
	`
	if len(req.GetSearch()) > 0 {
		filter += " AND first_name ILIKE '%' || '" + req.Search + "' || '%' "
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
			first_name sql.NullString
			last_name  sql.NullString
			phone      sql.NullString
			login      sql.NullString
			password   sql.NullString
			staff_type sql.NullString
			magazin_id sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&last_name,
			&phone,
			&login,
			&password,
			&staff_type,
			&magazin_id,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return resp, err
		}

		resp.Staffs = append(resp.Staffs, &organization_service.Staff{
			Id:        id.String,
			FirstName: first_name.String,
			LastName:  last_name.String,
			Phone:     phone.String,
			Login:     login.String,
			Password:  password.String,
			StaffType: staff_type.String,
			MagazinId: magazin_id.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}

	return
}

func (c *staffRepo) Update(ctx context.Context, req *organization_service.UpdateStaff) (resp int64, err error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"staff"
		SET
			first_name = :first_name,
			last_name= :last_name,
			phone = :phone,
			login = :login,
			password = : password,
			staff_type = :staff_type,
			magazin_id = :magazin_id,
			updated_at = now()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":        req.GetId(),
		"first_name":      req.GetFirstName(),
		"last_name": req.GetLastName(),
		"phone": req.GetPhone(),
		"login": req.GetLogin(),
		"password": req.GetPassword(),
		"staff_type": req.GetStaffType(),
		"magazin_id": req.GetMagazinId(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), nil
}

func (c *staffRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (resp int64, err error) {
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
			"staff"
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

func (c *staffRepo) Delete(ctx context.Context, req *organization_service.StaffPK) error {
	query := `DELETE FROM "staff" WHERE id = $1`

	_, err := c.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
