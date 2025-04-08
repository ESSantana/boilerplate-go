package repositories

import (
	"context"
	"database/sql"

	"github.com/application-ellas/ellas-backend/internal/domain/models"
	"github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
)

type customerRepository struct {
	conn *sql.DB
}

func newCustomerRepository(conn *sql.DB) interfaces.CustomerRepository {
	return &customerRepository{
		conn: conn,
	}
}

func (r *customerRepository) GetCustomerById(ctx context.Context, id string) (customer models.Customer, err error) {
	query := `SELECT * FROM customer WHERE id = ? AND deleted_at IS NULL`
	row := r.conn.QueryRowContext(ctx, query, id)
	return scanCustomer(row)
}

func (r *customerRepository) GetCustomerByExternalID(ctx context.Context, externalID string) (customer models.Customer, err error) {
	query := `SELECT * FROM customer WHERE external_id = ? AND deleted_at IS NULL`
	row := r.conn.QueryRowContext(ctx, query, externalID)
	return scanCustomer(row)
}

func (r *customerRepository) GetCustomerEmail(ctx context.Context, email string) (customer models.Customer, err error) {
	query := `SELECT * FROM customer WHERE email = ? AND deleted_at IS NULL`
	row := r.conn.QueryRowContext(ctx, query, email)
	return scanCustomer(row)
}

func (r *customerRepository) GetAllCustomers(ctx context.Context) (customers []models.Customer, err error) {
	query := `SELECT * FROM customer WHERE deleted_at IS NULL`
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	customers = make([]models.Customer, 0)
	for rows.Next() {
		customer, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *customerRepository) CreateCustomer(ctx context.Context, customer models.Customer) error {
	query := `
	INSERT INTO
		customer (
			id,
			name,
			nickname,
			email,
			password_hash,
			birth_date,
			phone_number,
			cpf,
			gender,
			address,
			address_number,
			address_complement,
			address_neighborhood,
			address_city,
			address_state,
			address_zip_code,
			provider_origin,
			external_id,
			profile_image_url,
			interests,
			how_heard_about_us,
			preferred_comunication_channel
		)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.conn.ExecContext(ctx, query,
		customer.ID, customer.Name, customer.Nickname, customer.Email, customer.PasswordHash, customer.BirthDate, customer.PhoneNumber,
		customer.CPF, customer.Gender, customer.Address, customer.AddressNumber, customer.AddressComplement, customer.AddressNeighborhood,
		customer.AddressCity, customer.AddressState, customer.AddressZipCode, customer.ProviderOrigin, customer.ExternalID,
		customer.ProfileImageURL, customer.Interests, customer.HowHeardAboutUs, customer.PreferredCommunicationChannel,
	)
	return err
}

func (r *customerRepository) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	query := `
	UPDATE
		customer
	SET
		name = ?,
		nickname = ?,
		email = ?,
		birth_date = ?,
		phone_number = ?,
		gender = ?,
		address = ?,
		address_number = ?,
		address_complement = ?,
		address_neighborhood = ?,
		address_city = ?,
		address_state = ?,
		address_zip_code = ?,
		profile_image_url = ?,
		interests = ?,
		preferred_comunication_channel = ?
	WHERE
		id = ?
		AND deleted_at IS NULL
	`
	_, err := r.conn.ExecContext(ctx, query,
		customer.Name, customer.Nickname, customer.Email, customer.BirthDate,
		customer.PhoneNumber, customer.Gender, customer.Address, customer.AddressNumber,
		customer.AddressComplement, customer.AddressNeighborhood, customer.AddressCity,
		customer.AddressState, customer.AddressZipCode, customer.ProfileImageURL, customer.Interests,
		customer.PreferredCommunicationChannel, customer.ID,
	)
	return err
}

func (r *customerRepository) SoftDeleteCustomer(ctx context.Context, id string) error {
	query := `UPDATE customer SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`
	_, err := r.conn.ExecContext(ctx, query, id)
	return err
}

func scanCustomer(scanner Scanner) (c models.Customer, err error) {
	err = scanner.Scan(
		&c.ID,
		&c.Name,
		&c.Nickname,
		&c.Email,
		&c.PasswordHash,
		&c.BirthDate,
		&c.PhoneNumber,
		&c.CPF,
		&c.Gender,
		&c.Address,
		&c.AddressNumber,
		&c.AddressComplement,
		&c.AddressNeighborhood,
		&c.AddressCity,
		&c.AddressState,
		&c.AddressZipCode,
		&c.ProviderOrigin,
		&c.ExternalID,
		&c.ProfileImageURL,
		&c.Interests,
		&c.HowHeardAboutUs,
		&c.PreferredCommunicationChannel,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	return c, err
}
