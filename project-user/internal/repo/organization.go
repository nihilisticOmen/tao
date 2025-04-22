package repo

import (
	"context"
	"project-user/internal/data/organization"
)

type OrganizationRepo interface {
	FindOrganizationByMemId(ctx context.Context, memId int64) ([]organization.Organization, error)
	SaveOrganization(ctx context.Context, org *organization.Organization) error
}
