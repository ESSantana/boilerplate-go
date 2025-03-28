package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/application-ellas/ellas-backend/internal/domain/models"
	"github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type userRepository struct {
	conn *sql.DB
}

func newUserRepository(conn *sql.DB) interfaces.UserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (r *userRepository) GetUserById(ctx context.Context, id string) (user models.User, err error) {
	rows, err := r.conn.QueryContext(ctx, `
			SELECT
				user.*,
				user_role.id,
				user_role.role
			FROM
				user
				JOIN user_role ON user.id = user_role.user_id
			WHERE 
				user.id = ?
				AND user.deleted_at IS NULL 
				AND user_role.deleted_at IS NULL;
			`, id)
	if err != nil {
		return user, fmt.Errorf("error scanning user: %s", err.Error())
	}

	userRoles := make([]models.UserRole, 0)

	for rows.Next() {
		var userRole models.UserRole
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProviderOrigin,
			&user.ExternalID,
			&user.ProfileImageURL,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&userRole.ID,
			&userRole.Role,
		)
		if err != nil {
			return user, fmt.Errorf("error scanning user: %s", err.Error())
		}
		userRoles = append(userRoles, userRole)
	}
	user.UserRoles = userRoles

	return user, nil
}

func (r *userRepository) GetUserByExternalId(ctx context.Context, externalID string) (user models.User, err error) {
	rows, err := r.conn.QueryContext(ctx, `
			SELECT
				user.*,
				user_role.id,
				user_role.role
			FROM
				user
				JOIN user_role ON user.id = user_role.user_id
			WHERE 
				user.external_id = ? 
				AND user.deleted_at IS NULL
				AND user_role.deleted_at IS NULL;
			`, externalID)
	if err != nil {
		return user, fmt.Errorf("error scanning user: %s", err.Error())
	}

	userRoles := make([]models.UserRole, 0)

	for rows.Next() {
		var userRole models.UserRole
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProviderOrigin,
			&user.ExternalID,
			&user.ProfileImageURL,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&userRole.ID,
			&userRole.Role,
		)
		if err != nil {
			return user, fmt.Errorf("error scanning user: %s", err.Error())
		}
		userRoles = append(userRoles, userRole)
	}
	user.UserRoles = userRoles

	return user, nil
}

func (r *userRepository) GetUserRoleByUserIdAndRole(ctx context.Context, userID, role string) (userRole models.UserRole, err error) {
	row := r.conn.QueryRowContext(ctx, `
			SELECT
				*
			FROM
				user_role
			WHERE 
				user_role.user_id = ?
				AND user_role.role = ?
				AND user_role.deleted_at IS NULL;
			`, userID, role)
	if err != nil {
		return userRole, fmt.Errorf("error scanning user role: %s", err.Error())
	}

	err = row.Scan(
		&userRole.ID,
		&userRole.UserID,
		&userRole.Role,
		&userRole.CreatedAt,
		&userRole.UpdatedAt,
		&userRole.DeletedAt,
	)
	if err != nil {
		return userRole, fmt.Errorf("error scanning user role: %s", err.Error())
	}

	return userRole, nil
}

func (r *userRepository) AppendRoleToUser(ctx context.Context, userID, role string) (err error) {
	_, err = r.conn.ExecContext(
		ctx,
		`
			INSERT INTO user_role (id, user_id, role)
			VALUES (?,?,?);
		`,
		uuid.New().String(), userID, role,
	)

	return err
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) (err error) {
	transaction, err := r.conn.BeginTx(
		ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating user: %s", err.Error())
	}

	defer transaction.Rollback()

	_, err = transaction.ExecContext(
		ctx,
		`
			INSERT INTO user (id, name, email, provider_origin, external_id, profile_image_url)
			VALUES (?,?,?,?,?,?);
		`,
		user.ID, user.Name, user.Email, user.ProviderOrigin, user.ExternalID, user.ProfileImageURL,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %s", err.Error())
	}

	_, err = transaction.ExecContext(
		ctx,
		`
			INSERT INTO user_role (id, user_id, role)
			VALUES (?,?,?);
		`,
		uuid.New().String(), user.ID, models.RoleCustomer,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %s", err.Error())
	}

	transaction.Commit()

	return nil
}
