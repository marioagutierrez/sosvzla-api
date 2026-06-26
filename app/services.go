package app

import (
	"context"
	"fmt"

	"github.com/sosvzla/sosvzla.lat/domain"
)

// SearchService orchestrates the search operations.
type SearchService struct {
	personRepo   domain.PersonRepository
	reportRepo   domain.ReportRepository
	hospitalRepo domain.HospitalRepository
}

// NewSearchService creates a new SearchService.
func NewSearchService(
	personRepo domain.PersonRepository,
	reportRepo domain.ReportRepository,
	hospitalRepo domain.HospitalRepository,
) *SearchService {
	return &SearchService{
		personRepo:   personRepo,
		reportRepo:   reportRepo,
		hospitalRepo: hospitalRepo,
	}
}

// --- Person Use Cases ---

func (s *SearchService) RegisterPerson(ctx context.Context, person *domain.Person) error {
	if person.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if person.Status == "" {
		return fmt.Errorf("status is required")
	}
	return s.personRepo.Create(ctx, person)
}

func (s *SearchService) GetPersonByID(ctx context.Context, id int) (*domain.Person, error) {
	return s.personRepo.GetByID(ctx, id)
}

func (s *SearchService) GetPersonByNationalID(ctx context.Context, nationalID string) (*domain.Person, error) {
	return s.personRepo.GetByNationalID(ctx, nationalID)
}

func (s *SearchService) UpdatePerson(ctx context.Context, person *domain.Person) error {
	if person.ID == 0 {
		return fmt.Errorf("person ID is required for update")
	}
	return s.personRepo.Update(ctx, person)
}

func (s *SearchService) DeletePerson(ctx context.Context, id int) error {
	return s.personRepo.Delete(ctx, id)
}

func (s *SearchService) ListPersons(ctx context.Context, status string, limit, offset int) ([]*domain.Person, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.personRepo.List(ctx, status, limit, offset)
}

// --- Report Use Cases ---

func (s *SearchService) FileReport(ctx context.Context, report *domain.Report) error {
	if report.PersonID == 0 {
		return fmt.Errorf("person ID is required for a report")
	}
	if report.ReporterName == "" {
		return fmt.Errorf("reporter name is required")
	}
	if report.ReporterContact == "" {
		return fmt.Errorf("reporter contact is required")
	}
	if report.ReportType == "" {
		return fmt.Errorf("report type is required")
	}
	return s.reportRepo.Create(ctx, report)
}

func (s *SearchService) ListReportsByPersonID(ctx context.Context, personID int) ([]*domain.Report, error) {
	return s.reportRepo.ListByPersonID(ctx, personID)
}

// --- Hospital Use Cases ---

func (s *SearchService) RegisterHospital(ctx context.Context, hospital *domain.Hospital) error {
	if hospital.Name == "" {
		return fmt.Errorf("hospital name is required")
	}
	return s.hospitalRepo.Create(ctx, hospital)
}

func (s *SearchService) ListHospitals(ctx context.Context, limit, offset int) ([]*domain.Hospital, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.hospitalRepo.List(ctx, limit, offset)
}

func (s *SearchService) AdmitPersonToHospital(ctx context.Context, hp *domain.HospitalPerson) error {
	if hp.PersonID == 0 {
		return fmt.Errorf("person ID is required")
	}
	if hp.HospitalID == 0 {
		return fmt.Errorf("hospital ID is required")
	}

	// Add person to hospital registry
	err := s.hospitalRepo.AddPerson(ctx, hp)
	if err != nil {
		return err
	}

	// Update person status to 'in_hospital'
	person, err := s.personRepo.GetByID(ctx, hp.PersonID)
	if err != nil {
		return err
	}
	if person != nil {
		person.Status = "in_hospital"
		return s.personRepo.Update(ctx, person)
	}

	return nil
}

func (s *SearchService) ListPersonsInHospital(ctx context.Context, hospitalID int) ([]*domain.Person, error) {
	return s.hospitalRepo.ListPersons(ctx, hospitalID)
}
