package domain

import "time"

// Report represents a report about a missing or found person.
type Report struct {
    ID             int       `json:"id"`
    PersonID       int       `json:"person_id"`
    ReporterName   string    `json:"reporter_name"`
    ReporterContact string    `json:"reporter_contact"`
    ReportType     string    `json:"report_type"`
    Description    *string   `json:"description"`
    CreatedAt      time.Time `json:"created_at"`
}
