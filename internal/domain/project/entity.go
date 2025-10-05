package project

import (
	"api/internal/common"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID         uuid.UUID  `json:"ID"`
	CustomerID uuid.UUID  `json:"customer_id"`
	Slug       string     `json:"slug"`
	Name       string     `json:"name"`
	RemovedAt  *time.Time `json:"removed_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func New(customerID uuid.UUID, slug, name string) (*Project, error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new project: %w", err)
	}

	return &Project{
		ID:         ID,
		CustomerID: customerID,
		Slug:       slug,
		Name:       name,
		RemovedAt:  nil,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

func (p *Project) IsOwner(customerID uuid.UUID) error {
	if p.CustomerID == customerID {
		return nil
	}

	return common.ErrForbidden
}

func (p *Project) Remove() {
	currentTime := time.Now().UTC()
	p.RemovedAt = &currentTime
}

func (p *Project) IsRemoved() bool {
	return p.RemovedAt != nil
}
