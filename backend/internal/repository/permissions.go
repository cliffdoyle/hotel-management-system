package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A set of strings for holding permissions. The empty struct value is for memory efficiency.
type PermissionSet map[string]struct{}

type PermissionRepository struct {
	DB *pgxpool.Pool
}

// Models holds all the repositories for our application.
type Models struct {
	Permissions PermissionRepository
	// We will add Users UserRepository here in the next issue
}

// NewModels returns a Models struct containing all initialized repositories.
func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Permissions: PermissionRepository{DB: db},
		// Users:       UserRepository{DB: db},
	}
}

// GetAllForUser retrieves all permission codes for a specific user.
func (r *PermissionRepository) GetAllForUser(userID uuid.UUID) (PermissionSet, error) {
	query := `
        SELECT p.code
        FROM permissions p
        INNER JOIN roles_permissions rp ON rp.permission_id = p.id
        INNER JOIN users_roles ur ON ur.role_id = rp.role_id
        WHERE ur.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx, query, userID)
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