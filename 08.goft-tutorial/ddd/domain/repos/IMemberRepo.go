package repos

import "goft-tutorial/ddd/domain/aggregates"

type IMemberRepo interface {
	FindByName(name string) *aggregates.Member
	CreateMember(member *aggregates.Member) error
	UpdateMember(member *aggregates.Member) error
}
