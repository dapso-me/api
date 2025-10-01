package project_handler

import "github.com/google/uuid"

type addProjectReq struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Slug       string    `json:"slug"`
	Name       string    `json:"name"`
}
