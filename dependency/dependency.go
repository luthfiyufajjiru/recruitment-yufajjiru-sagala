//go:build wireinject
// +build wireinject

package dependency

import (
	"os"
	"sagala-todo/pkg/adapters"
	"sagala-todo/pkg/constants"
	"sagala-todo/src/delivery/v1http"
	"sagala-todo/src/usecase"

	"github.com/google/wire"
)

func InitOsSignalChannel() chan os.Signal {
	wire.Build(
		provideOsSignal,
	)
	return nil
}

func InitTodoV1HttpHandler() *v1http.V1Handler {
	wire.Build(
		provideConfiguration,
		wire.Value([]adapters.SqlConfig{
			{
				RegistryName: constants.ConnSqlDefault,
				DriverName:   constants.SqliteDriver,
			},
		}),
		provideSql,
		usecase.ProvideUsecase,
		wire.Bind(new(v1http.Usecaser), new(*usecase.TodoUsecase)),
		wire.Struct(new(v1http.V1Handler), "*"),
	)
	return nil
}
