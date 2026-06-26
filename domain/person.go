package domain

import "time"

// Person represents an individual in the system.
type Person struct {
    ID               int       `json:"id"`
    FullName         string    `json:"full_name"`
    NationalID       *string   `json:"national_id"`
    BirthDate        *time.Time `json:"birth_date"`
    Gender           *string   `json:"gender"`
    LastSeenLocation *string   `json:"last_seen_location"`
    LastSeenDate     *time.Time `json:"last_seen_date"`
    Status           string    `json:"status"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
