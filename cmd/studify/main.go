package main

import (
	"context"
	"log/slog"
	"os"

	"go.uber.org/fx"

	"github.com/k6zma/avito-lab1/internal/application/services"
	domainRepos "github.com/k6zma/avito-lab1/internal/domain/repositories"
	"github.com/k6zma/avito-lab1/internal/infrastructure/ciphers"
	"github.com/k6zma/avito-lab1/internal/infrastructure/flags"
	"github.com/k6zma/avito-lab1/internal/infrastructure/persisters"
	infrastructureRepos "github.com/k6zma/avito-lab1/internal/infrastructure/repositories"
	"github.com/k6zma/avito-lab1/internal/presentation/tui"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() *slog.Logger {
			l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			slog.SetDefault(l)

			return l
		}),

		fx.Invoke(validators.InitValidators()),

		fx.Provide(
			func() (*flags.StudyFlags, error) {
				return flags.GetFlags()
			},

			func(cfg *flags.StudyFlags) (ciphers.Cipher, error) {
				return ciphers.NewAESGCM(cfg.CipherKey)
			},

			func(cfg *flags.StudyFlags, c ciphers.Cipher) persisters.StudentPersister {
				return persisters.NewJSONStudentPersister(cfg.ConfigPath, c)
			},

			func(p persisters.StudentPersister) (domainRepos.StudentRepository, error) {
				return infrastructureRepos.NewStudentStorageWithPersister(p)
			},

			func(repo domainRepos.StudentRepository) services.StudentServiceContract {
				return services.NewStudentService(repo)
			},
		),

		fx.Invoke(func(
			lc fx.Lifecycle,
			svc services.StudentServiceContract,
			sd fx.Shutdowner,
			log *slog.Logger,
		) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := tui.Run(svc); err != nil {
							log.Error(
								"TUI exited with error",
								"error", err,
							)
						}

						err := sd.Shutdown()
						if err != nil {
							log.Error(
								"Failed to shutdown TUI app",
								"error", err,
							)
						}
					}()

					return nil
				},
			})
		}),
	)

	app.Run()
}
