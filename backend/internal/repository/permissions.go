package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionSet map[string]struct{}

type PermissionRepository interface {
	GetAllForUser(ctx context.Context, userID uuid.UUID) (PermissionSet, error)
}

type permissionRepository struct {
	db *pgxpool.Pool
}

func NewPermissionRepository(db *pgxpool.Pool) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetAllForUser(ctx context.Context, userID uuid.UUID) (PermissionSet, error) {
	query := `
        SELECT p.code
        FROM permissions p
        INNER JOIN roles_permissions rp ON rp.permission_id = p.id
        INNER JOIN users_roles ur ON ur.role_id = rp.role_id
        WHERE ur.user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make(PermissionSet)
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions[permission] = struct{}{}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}