package usecase

import "sagala-todo/pkg/adapters"

type (
	TodoUsecase struct {
		Config adapters.Config
		Sql    map[string]*adapters.Sql
	}
)

func ProvideUsecase(Sql map[string]*adapters.Sql, cfg adapters.Config) *TodoUsecase {
	return &TodoUsecase{
		Config: cfg,
		Sql:    Sql,
	}
}
