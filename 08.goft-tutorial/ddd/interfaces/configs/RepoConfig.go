package configs

import (
	"goft-tutorial/ddd/domain/repos"
	"goft-tutorial/ddd/infrastructure/dao"
)

type RepoConfig struct{}

func NewRepoConfig() *RepoConfig {
	return &RepoConfig{}
}

func (r *RepoConfig) UserRepo() repos.IUserRepo {
	return &dao.UserRepo{}
}
