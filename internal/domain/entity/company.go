package entity

import "time"

// Company represents a company entity.
type Company struct {
	ID                 int64
	Name               string
	RepresentativeName string
	PhoneNumber        string
	ZipCode            string
	Address            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
