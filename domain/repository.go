package domain

import "context"

// PersonRepository defines the interface for interacting with person data.
type PersonRepository interface {
    Create(ctx context.Context, person *Person) error
    GetByID(ctx context.Context, id int) (*Person, error)
    GetByNationalID(ctx context.Context, nationalID string) (*Person, error)
    Update(ctx context.Context, person *Person) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, status string, limit, offset int) ([]*Person, error)
}

// ReportRepository defines the interface for interacting with report data.
type ReportRepository interface {
    Create(ctx context.Context, report *Report) error
    GetByID(ctx context.Context, id int) (*Report, error)
    ListByPersonID(ctx context.Context, personID int) ([]*Report, error)
}

// HospitalRepository defines the interface for interacting with hospital data.
type HospitalRepository interface {
    Create(ctx context.Context, hospital *Hospital) error
    GetByID(ctx context.Context, id int) (*Hospital, error)
    List(ctx context.Context, limit, offset int) ([]*Hospital, error)
    AddPerson(ctx context.Context, hp *HospitalPerson) error
    ListPersons(ctx context.Context, hospitalID int) ([]*Person, error)
}
