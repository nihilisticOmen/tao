package repo

import (
	"context"
	"project-user/internal/data/member"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (mem *member.Member, err error)
}
