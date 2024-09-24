package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/binarstrike/magic-url/config"
	_ "github.com/binarstrike/magic-url/docs"
	"github.com/binarstrike/magic-url/internal/delivery/http/controller"
	"github.com/binarstrike/magic-url/internal/delivery/http/route"
	"github.com/binarstrike/magic-url/internal/middleware"
	sess_repo "github.com/binarstrike/magic-url/internal/session/repository"
	sess_uc "github.com/binarstrike/magic-url/internal/session/usecase"
	user_repo "github.com/binarstrike/magic-url/internal/user/repository"
	user_uc "github.com/binarstrike/magic-url/internal/user/usecase"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func startApp(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, logger lojer.Logger) {
	lc.Append(fx.StartHook(func() {
		listenAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)

		logger.Info("Starting app...", zap.String("app_version", config.APP_VERSION), zap.String("listen_to", listenAddr))

		go func() {
			err := app.Listen(listenAddr)
			if err != nil {
				logger.Error("Error!! Shutting down app...", zap.Error(err))
				time.Sleep(time.Second * 2)
				os.Exit(1)
			}
		}()
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) {
		err := app.ShutdownWithContext(ctx)
		logger.Info("Shutting down app...", zap.Error(err))
		os.Exit(0)
	}))
}

func main() {
	logger := lojer.NewLogger(config.APP_ENV != "production")

	defer func() {
		if r := recover(); r != nil {
			logger.Error("App panicking...", zap.Any("panic_message", r))
			time.Sleep(time.Second * 2)
			os.Exit(1)
		}
	}()

	app := fx.New(
		fx.Provide(newConfig, newRedisClient, newValidator, newDBConn, newFiberApp),
		fx.Provide(controller.NewUserController),
		fx.Decorate(route.SetupUserRoute),
		fx.Provide(user_uc.NewUserUseCase, sess_uc.NewSessionUseCase),
		fx.Provide(user_repo.NewUserRepository, user_repo.NewRedisRepository, sess_repo.NewSessionRepository),
		fx.Provide(middleware.NewMiddleware),
		fx.Provide(func() lojer.Logger { return logger }),
		fx.Logger(logger),
		fx.Invoke(startApp),
	)

	app.Run()
}
