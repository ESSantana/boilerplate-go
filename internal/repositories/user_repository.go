package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ESSantana/boilerplate-go/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-go/internal/domain/models"
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
	row := r.conn.QueryRowContext(ctx, `
			SELECT
				*
			FROM
				user
			WHERE 
				id = ?
				AND deleted_at IS NULL 
			`, id)

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.ProviderOrigin,
		&user.ExternalID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err.Error())
		return user, fmt.Errorf("error scanning user: %s", err.Error())
	}

	return user, nil
}

func (r *userRepository) GetUserByExternalId(ctx context.Context, externalID string) (user models.User, err error) {
	row := r.conn.QueryRowContext(ctx, `
			SELECT
				*
			FROM
				user
			WHERE 
				external_id = ? 
				AND deleted_at IS NULL
			`, externalID)

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.ProviderOrigin,
		&user.ExternalID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, fmt.Errorf("error scanning user: %s", err.Error())
	}

	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) (err error) {
	_, err = r.conn.ExecContext(
		ctx,
		`
			INSERT INTO user (id, name, email, provider_origin, external_id, profile_image_url, created_at, updated_at)
			VALUES (?,?,?,?,?,?,?,?);
		`,
		user.ID, user.Name, user.Email, user.ProviderOrigin, user.ExternalID, user.ProfileImageURL, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %s", err.Error())
	}

	return nil
}
