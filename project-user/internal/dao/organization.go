package dao

import (
	"context"
	"project-user/internal/data/organization"
	"project-user/internal/database"
	"project-user/internal/database/gorms"
)

type OrganizationDao struct {
	// 不再存储连接状态
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{}
}

func (o *OrganizationDao) FindOrganizationByMemId(ctx context.Context, memId int64) ([]organization.Organization, error) {
	var orgs []organization.Organization
	conn := gorms.New()
	err := conn.Session(ctx).Where("member_id=?", memId).Find(&orgs).Error
	return orgs, err
}

func (o *OrganizationDao) SaveOrganization(conn database.DbConn, ctx context.Context, org *organization.Organization) error {
	gormConn, ok := conn.(*gorms.GormConn)
	if !ok {
		return nil
	}
	return gormConn.Tx(ctx).Create(org).Error
}
