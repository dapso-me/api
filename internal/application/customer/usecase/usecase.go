package customer_usecase

import customer "api/internal/application/customer/domain"

type uc struct {
	customerRepo customer.Repository
}

func New(customerRepo customer.Repository) *uc {
	return &uc{
		customerRepo: customerRepo,
	}
}
