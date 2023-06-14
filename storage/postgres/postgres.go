package postgres

import (
	"context"
	"fmt"
	"organization_service/config"
	"organization_service/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db     *pgxpool.Pool
	filial storage.FilialRepoI
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:     pool,
		filial: NewFilialRepo(pool),
	}, nil
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) Filial() storage.FilialRepoI {
	if s.filial == nil {
		s.filial = NewFilialRepo(s.db)
	}
	return s.filial
}
