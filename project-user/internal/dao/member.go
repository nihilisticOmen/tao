package dao

import (
	"context"
	"gorm.io/gorm"
	"project-user/internal/data/member"
	"project-user/internal/database"
	"project-user/internal/database/gorms"
)

type MemberDao struct {
	// 不再存储连接状态
}

func NewMemberDao() *MemberDao {
	return &MemberDao{}
}

func (m *MemberDao) FindMember(ctx context.Context, account string, pwd string) (*member.Member, error) {
	var mem *member.Member
	conn := gorms.New()
	err := conn.Session(ctx).Where("account=? and password=?", account, pwd).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return mem, err
}

func (m *MemberDao) SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error {
	gormConn, ok := conn.(*gorms.GormConn)
	if !ok {
		return nil
	}
	return gormConn.Tx(ctx).Create(mem).Error
}

func (m *MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	conn := gorms.New()
	err := conn.Session(ctx).Model(&member.Member{}).Where("email=?", email).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) GetMemberByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	conn := gorms.New()
	err := conn.Session(ctx).Model(&member.Member{}).Where("account=?", account).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	conn := gorms.New()
	err := conn.Session(ctx).Model(&member.Member{}).Where("mobile=?", mobile).Count(&count).Error
	return count > 0, err
}
