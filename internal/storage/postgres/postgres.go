package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/risetyll/finuslugi/internal/entities/material"
	"github.com/risetyll/finuslugi/internal/entities/requests"
	lu "github.com/risetyll/finuslugi/internal/logger/utils"
	"github.com/risetyll/finuslugi/internal/storage"
)

type Connector interface {
	GetProvider() string
	GetConnect() string
}

type Postgres struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(connector Connector, logger *slog.Logger) (*Postgres, error) {
	const op = "storage.New"

	db, err := sql.Open(connector.GetProvider(), connector.GetConnect())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{db, logger}, nil
}

func (p *Postgres) Init() error {
	const op = "storage.Init"

	p.logger.Debug("starting init transaction")
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			p.logger.Error("rolling back transaction due to error")
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				p.logger.Error("failed to rollback transaction: %v", lu.Error(rollbackErr))
			}
		}
	}()

	p.logger.Debug("creating table")
	_, err = tx.Exec(storage.TableSchema)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Debug("creating material index")
	_, err = tx.Exec(storage.MaterialIndexSchema)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Debug("creating date index")
	_, err = tx.Exec(storage.DateIndexSchema)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Debug("creating composite index")
	_, err = tx.Exec(storage.CompositeIndexSchema)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Debug("completing transaction")
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Postgres) CreateMaterial(req *requests.CreateMaterialRequest) error {
	const op = "storage.CreateMaterial"

	s.logger.Debug("building query")
	query, args, err := squirrel.
		Insert("materials").
		Columns("material_type", "publication_status", "title", "content").
		Values(req.Type, req.Status, req.Title, req.Content).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Debug("query: %s", slog.AnyValue(query).Any())

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Postgres) GetMaterialById(req *requests.GetMaterialByIdRequest) (*material.Material, error) {
	const op = "storage.GetMaterialById"

	s.logger.Debug("building query")
	queryBuilder := squirrel.Select("*").
		From("materials").
		Where(squirrel.Eq{"uuid": req.UUID})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Debug("query: %s", slog.AnyValue(query).Any())

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	if !rows.Next() {
		s.logger.Debug("no rows found for UUID: %d", slog.AnyValue(req.UUID).Any())
		return nil, nil
	}

	var material material.Material
	if err := rows.Scan(
		&material.UUID,
		&material.Type,
		&material.Status,
		&material.Title,
		&material.Content,
		&material.CreatedAt,
		&material.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &material, nil
}

func (s *Postgres) UpdateMaterial(req *requests.UpdateMaterialRequest) error {
	const op = "storage.UpdateMaterial"

	s.logger.Debug("building query")
	queryBuilder := squirrel.
		Update("materials").
		Where(squirrel.Eq{"uuid": req.UUID})

	if req.Status != "" {
		queryBuilder = queryBuilder.Set("publication_status", req.Status)
	}
	if req.Title != "" {
		queryBuilder = queryBuilder.Set("title", req.Title)
	}
	if req.Content != "" {
		queryBuilder = queryBuilder.Set("content", req.Content)
	}

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Debug("query: %s", slog.AnyValue(query).Any())

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Postgres) GetMaterials(req *requests.GetMaterialsRequest) ([]*material.Material, error) {
	const op = "storage.GetMaterials"

	s.logger.Debug("building query")
	queryBuilder := squirrel.Select("*").From("materials")

	if req.Type != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"material_type": req.Type})
	}
	if !req.CreatedFrom.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"creation_date": req.CreatedFrom})
	}
	if !req.CreatedTo.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"creation_date": req.CreatedTo})
	}

	queryBuilder = queryBuilder.
		Limit(uint64(req.PageSize)).
		Offset(uint64((req.Page - 1) * req.PageSize))

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		s.logger.Error("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Debug("query: %s", slog.AnyValue(query).Any())

	rows, err := s.db.Query(query, args...)
	if err != nil {
		s.logger.Error("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var materials []*material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(
			&m.UUID,
			&m.Type,
			&m.Status,
			&m.Title,
			&m.Content,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			s.logger.Error("%s: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		materials = append(materials, &m)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return materials, nil
}
