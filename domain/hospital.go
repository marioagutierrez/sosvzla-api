package domain

import "time"

// Hospital represents a hospital in the system.
type Hospital struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Address     *string   `json:"address"`
    City        *string   `json:"city"`
    State       *string   `json:"state"`
    Country     *string   `json:"country"`
    PhoneNumber *string   `json:"phone_number"`
    CreatedAt   time.Time `json:"created_at"`
}

// HospitalPerson represents the relationship between a person and a hospital.
type HospitalPerson struct {
    ID             int       `json:"id"`
    PersonID       int       `json:"person_id"`
    HospitalID     int       `json:"hospital_id"`
    AdmissionDate  *time.Time `json:"admission_date"`
    DischargeDate  *time.Time `json:"discharge_date"`
    Notes          *string   `json:"notes"`
}
