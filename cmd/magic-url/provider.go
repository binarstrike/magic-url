package main

import (
	"context"
	"fmt"

	"github.com/binarstrike/magic-url/config"
	"github.com/binarstrike/magic-url/pkg/lojer"
	"github.com/binarstrike/magic-url/pkg/validator"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/swagger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func newConfig() (*config.Config, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func newRedisClient(logger lojer.Logger, cfg *config.Config) (*redis.Client, error) {
	// TODO: tls? authentication?
	rc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		OnConnect: func(_ context.Context, cn *redis.Conn) error {
			logger.Info("redis connection established", zap.String("redis", cn.String()))
			return nil
		},
	})

	err := rc.Ping(context.TODO()).Err()
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func newValidator() *validator.Validator {
	return validator.NewValidator()
}

func newDBConn(logger lojer.Logger, cfg *config.Config) (*sqlx.DB, error) {
	// TODO: aktifkan fitur ssl
	dsn := fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' port='%d' sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("database connection established")

	return db, nil
}

func newFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Use(helmet.New())
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	return app
}
