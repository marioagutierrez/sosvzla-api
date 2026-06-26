package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sosvzla/sosvzla.lat/domain"
)

// PostgresPersonRepository implements domain.PersonRepository.
type PostgresPersonRepository struct {
	db *sql.DB
}

// NewPostgresPersonRepository creates a new instance of PostgresPersonRepository.
func NewPostgresPersonRepository(db *sql.DB) *PostgresPersonRepository {
	return &PostgresPersonRepository{db: db}
}

var _ domain.PersonRepository = (*PostgresPersonRepository)(nil)

func (r *PostgresPersonRepository) Create(ctx context.Context, person *domain.Person) error {
	query := `
		INSERT INTO persons (full_name, national_id, birth_date, gender, last_seen_location, last_seen_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		person.FullName,
		person.NationalID,
		person.BirthDate,
		person.Gender,
		person.LastSeenLocation,
		person.LastSeenDate,
		person.Status,
	).Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error inserting person: %w", err)
	}
	return nil
}

func (r *PostgresPersonRepository) GetByID(ctx context.Context, id int) (*domain.Person, error) {
	query := `
		SELECT id, full_name, national_id, birth_date, gender, last_seen_location, last_seen_date, status, created_at, updated_at
		FROM persons
		WHERE id = $1
	`
	var person domain.Person
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&person.ID,
		&person.FullName,
		&person.NationalID,
		&person.BirthDate,
		&person.Gender,
		&person.LastSeenLocation,
		&person.LastSeenDate,
		&person.Status,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying person by id: %w", err)
	}
	return &person, nil
}

func (r *PostgresPersonRepository) GetByNationalID(ctx context.Context, nationalID string) (*domain.Person, error) {
	query := `
		SELECT id, full_name, national_id, birth_date, gender, last_seen_location, last_seen_date, status, created_at, updated_at
		FROM persons
		WHERE national_id = $1
	`
	var person domain.Person
	err := r.db.QueryRowContext(ctx, query, nationalID).Scan(
		&person.ID,
		&person.FullName,
		&person.NationalID,
		&person.BirthDate,
		&person.Gender,
		&person.LastSeenLocation,
		&person.LastSeenDate,
		&person.Status,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying person by national id: %w", err)
	}
	return &person, nil
}

func (r *PostgresPersonRepository) Update(ctx context.Context, person *domain.Person) error {
	query := `
		UPDATE persons
		SET full_name = $1, national_id = $2, birth_date = $3, gender = $4, last_seen_location = $5, last_seen_date = $6, status = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		person.FullName,
		person.NationalID,
		person.BirthDate,
		person.Gender,
		person.LastSeenLocation,
		person.LastSeenDate,
		person.Status,
		person.ID,
	).Scan(&person.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error updating person: %w", err)
	}
	return nil
}

func (r *PostgresPersonRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM persons WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting person: %w", err)
	}
	return nil
}

func (r *PostgresPersonRepository) List(ctx context.Context, status string, limit, offset int) ([]*domain.Person, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, full_name, national_id, birth_date, gender, last_seen_location, last_seen_date, status, created_at, updated_at
			FROM persons
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{status, limit, offset}
	} else {
		query = `
			SELECT id, full_name, national_id, birth_date, gender, last_seen_location, last_seen_date, status, created_at, updated_at
			FROM persons
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing persons: %w", err)
	}
	defer rows.Close()

	var persons []*domain.Person
	for rows.Next() {
		var person domain.Person
		err := rows.Scan(
			&person.ID,
			&person.FullName,
			&person.NationalID,
			&person.BirthDate,
			&person.Gender,
			&person.LastSeenLocation,
			&person.LastSeenDate,
			&person.Status,
			&person.CreatedAt,
			&person.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning person row: %w", err)
		}
		persons = append(persons, &person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iteration: %w", err)
	}

	return persons, nil
}

// PostgresReportRepository implements domain.ReportRepository.
type PostgresReportRepository struct {
	db *sql.DB
}

// NewPostgresReportRepository creates a new instance of PostgresReportRepository.
func NewPostgresReportRepository(db *sql.DB) *PostgresReportRepository {
	return &PostgresReportRepository{db: db}
}

var _ domain.ReportRepository = (*PostgresReportRepository)(nil)

func (r *PostgresReportRepository) Create(ctx context.Context, report *domain.Report) error {
	query := `
		INSERT INTO reports (person_id, reporter_name, reporter_contact, report_type, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		report.PersonID,
		report.ReporterName,
		report.ReporterContact,
		report.ReportType,
		report.Description,
	).Scan(&report.ID, &report.CreatedAt)

	if err != nil {
		return fmt.Errorf("error inserting report: %w", err)
	}
	return nil
}

func (r *PostgresReportRepository) GetByID(ctx context.Context, id int) (*domain.Report, error) {
	query := `
		SELECT id, person_id, reporter_name, reporter_contact, report_type, description, created_at
		FROM reports
		WHERE id = $1
	`
	var report domain.Report
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID,
		&report.PersonID,
		&report.ReporterName,
		&report.ReporterContact,
		&report.ReportType,
		&report.Description,
		&report.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying report by id: %w", err)
	}
	return &report, nil
}

func (r *PostgresReportRepository) ListByPersonID(ctx context.Context, personID int) ([]*domain.Report, error) {
	query := `
		SELECT id, person_id, reporter_name, reporter_contact, report_type, description, created_at
		FROM reports
		WHERE person_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, personID)
	if err != nil {
		return nil, fmt.Errorf("error listing reports by person id: %w", err)
	}
	defer rows.Close()

	var reports []*domain.Report
	for rows.Next() {
		var report domain.Report
		err := rows.Scan(
			&report.ID,
			&report.PersonID,
			&report.ReporterName,
			&report.ReporterContact,
			&report.ReportType,
			&report.Description,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning report row: %w", err)
		}
		reports = append(reports, &report)
	}

	return reports, nil
}

// PostgresHospitalRepository implements domain.HospitalRepository.
type PostgresHospitalRepository struct {
	db *sql.DB
}

// NewPostgresHospitalRepository creates a new instance of PostgresHospitalRepository.
func NewPostgresHospitalRepository(db *sql.DB) *PostgresHospitalRepository {
	return &PostgresHospitalRepository{db: db}
}

var _ domain.HospitalRepository = (*PostgresHospitalRepository)(nil)

func (r *PostgresHospitalRepository) Create(ctx context.Context, hospital *domain.Hospital) error {
	query := `
		INSERT INTO hospitals (name, address, city, state, country, phone_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		hospital.Name,
		hospital.Address,
		hospital.City,
		hospital.State,
		hospital.Country,
		hospital.PhoneNumber,
	).Scan(&hospital.ID, &hospital.CreatedAt)

	if err != nil {
		return fmt.Errorf("error inserting hospital: %w", err)
	}
	return nil
}

func (r *PostgresHospitalRepository) GetByID(ctx context.Context, id int) (*domain.Hospital, error) {
	query := `
		SELECT id, name, address, city, state, country, phone_number, created_at
		FROM hospitals
		WHERE id = $1
	`
	var hospital domain.Hospital
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hospital.ID,
		&hospital.Name,
		&hospital.Address,
		&hospital.City,
		&hospital.State,
		&hospital.Country,
		&hospital.PhoneNumber,
		&hospital.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying hospital by id: %w", err)
	}
	return &hospital, nil
}

func (r *PostgresHospitalRepository) List(ctx context.Context, limit, offset int) ([]*domain.Hospital, error) {
	query := `
		SELECT id, name, address, city, state, country, phone_number, created_at
		FROM hospitals
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing hospitals: %w", err)
	}
	defer rows.Close()

	var hospitals []*domain.Hospital
	for rows.Next() {
		var hospital domain.Hospital
		err := rows.Scan(
			&hospital.ID,
			&hospital.Name,
			&hospital.Address,
			&hospital.City,
			&hospital.State,
			&hospital.Country,
			&hospital.PhoneNumber,
			&hospital.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning hospital row: %w", err)
		}
		hospitals = append(hospitals, &hospital)
	}

	return hospitals, nil
}

func (r *PostgresHospitalRepository) AddPerson(ctx context.Context, hp *domain.HospitalPerson) error {
	query := `
		INSERT INTO hospital_persons (person_id, hospital_id, admission_date, discharge_date, notes)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (person_id, hospital_id) DO UPDATE
		SET admission_date = EXCLUDED.admission_date,
		    discharge_date = EXCLUDED.discharge_date,
		    notes = EXCLUDED.notes
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		hp.PersonID,
		hp.HospitalID,
		hp.AdmissionDate,
		hp.DischargeDate,
		hp.Notes,
	).Scan(&hp.ID)

	if err != nil {
		return fmt.Errorf("error adding person to hospital: %w", err)
	}
	return nil
}

func (r *PostgresHospitalRepository) ListPersons(ctx context.Context, hospitalID int) ([]*domain.Person, error) {
	query := `
		SELECT p.id, p.full_name, p.national_id, p.birth_date, p.gender, p.last_seen_location, p.last_seen_date, p.status, p.created_at, p.updated_at
		FROM persons p
		JOIN hospital_persons hp ON p.id = hp.person_id
		WHERE hp.hospital_id = $1
		ORDER BY hp.admission_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, hospitalID)
	if err != nil {
		return nil, fmt.Errorf("error listing persons in hospital: %w", err)
	}
	defer rows.Close()

	var persons []*domain.Person
	for rows.Next() {
		var person domain.Person
		err := rows.Scan(
			&person.ID,
			&person.FullName,
			&person.NationalID,
			&person.BirthDate,
			&person.Gender,
			&person.LastSeenLocation,
			&person.LastSeenDate,
			&person.Status,
			&person.CreatedAt,
			&person.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning person row: %w", err)
		}
		persons = append(persons, &person)
	}

	return persons, nil
}
