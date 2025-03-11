package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AleksandrVishniakov/jwt-auth/internal/e"
	"github.com/AleksandrVishniakov/jwt-auth/internal/repository/db"
	"github.com/AleksandrVishniakov/jwt-auth/internal/usecases"
)

type Repository struct {
	log     *slog.Logger
	db      *sql.DB
	queries *db.Queries
}

func New(
	log *slog.Logger,
	db *sql.DB,
	queries *db.Queries,
) *Repository {
	return &Repository{
		log:     log,
		db:      db,
		queries: queries,
	}
}

func (r *Repository) CreateUser(
	ctx context.Context,
	login string,
	passworHash []byte,
) (id int32, err error) {
	const src = "Repository.CreateUser"
	log := r.log.With(slog.String("src", src))
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: failed to create user: %w", src, err)
		}
	}()

	log.Debug("creating new user", slog.String("login", login))

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	q := r.queries.WithTx(tx)

	_, err = q.GetUserByLogin(ctx, login)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return 0, err
		}

		return 0, e.ErrAlreadyExists
	}

	id, err = q.CreateUser(ctx, db.CreateUserParams{
		Login:        login,
		PasswordHash: string(passworHash),
	})

	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) CreateSuperUser(
	ctx context.Context,
	login string,
	passworHash []byte,
) (id int32, err error) {
	const src = "Repository.CreateSuperUser"
	log := r.log.With(slog.String("src", src))
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: failed to create super user: %w", src, err)
		}
	}()

	log.Debug("creating new super user", slog.String("login", login))

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	q := r.queries.WithTx(tx)

	_, err = q.GetUserByLogin(ctx, login)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return 0, err
		}

		return 0, e.ErrAlreadyExists
	}

	id, err = q.CreateSuperUser(ctx, db.CreateSuperUserParams{
		Login:        login,
		PasswordHash: string(passworHash),
	})

	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetUserById(
	ctx context.Context,
	id int32,
) (user *usecases.UserModel, err error) {
	const src = "Repository.GetUserById"
	log := r.log.With(slog.String("src", src))
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: failed to fetch user: %w", src, err)
		}
	}()

	log.Debug("fetching user", slog.Int("id", int(id)))

	entity, err := r.queries.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}

		return nil, err
	}

	return &usecases.UserModel{
		ID:             id,
		Login:          entity.Login,
		PasswordHash:   entity.PasswordHash,
		Role:           entity.Alias,
		PermissionMask: entity.PermissionsMask,
	}, nil
}

func (r *Repository) GetUserByLogin(
	ctx context.Context,
	login string,
) (user *usecases.UserModel, err error) {
	const src = "Repository.GetUserByLogin"
	log := r.log.With(slog.String("src", src))
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: failed to fetch user: %w", src, err)
		}
	}()

	log.Debug("fetching user", slog.String("login", login))

	entity, err := r.queries.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	return &usecases.UserModel{
		ID:             entity.ID,
		Login:          entity.Login,
		PasswordHash:   entity.PasswordHash,
		Role:           entity.Alias,
		PermissionMask: entity.PermissionsMask,
	}, nil
}

func (r *Repository) UpsertRole(
	ctx context.Context,
	alias string,
	mask int64,
	isDefault bool,
	isSuper bool,
) (id int32, err error) {
	const src = "Repository.UpsertRole"
	log := r.log.With(slog.String("src", src))
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: failed to create role: %w", src, err)
		}
	}()

	log.Debug("upserting new role", slog.String("alias", alias))

	id, err = r.queries.UpsertRole(ctx, db.UpsertRoleParams{
		Alias:           alias,
		PermissionsMask: mask,
		IsDefault:       isDefault,
		IsSuper:         isSuper,
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateRoleById(
	ctx context.Context,
	userID int32,
	roleAlias string,
) (err error) {
	const src = "Repository.UpdateRoleById"
	log := r.log.With(slog.String("src", src))
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: failed to update role: %w", src, err)
		}
	}()

	log.Debug("updating user role", slog.Int("id", int(userID)))

	err = r.queries.UpdateRoleById(ctx, db.UpdateRoleByIdParams{
		ID:    userID,
		Alias: roleAlias,
	})
	if err != nil {
		return err
	}

	return nil
}
