package config

import (
	"time"

	"github.com/spf13/viper"
)

// use -ldflags flag to change value of these variables
//
// example: -ldflags="-X github.com/binarstrike/magic-url/config.APP_VERSION=1.0.4-beta"
var (
	APP_ENV     = "development"
	APP_VERSION = "0.0.1-alpha"
)

type databaseConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string `mapstructure:"db_name"`
}

type redisConfig struct {
	Host string
	Port int
}

type jwtConfig struct {
	Secret string
}

type appConfig struct {
	Host     string
	Port     int
	APIRoute string `mapstructure:"api_route"`
}

type consts struct {
	SessionCookieName     string
	JWTCookieName         string
	JWTHeaderName         string
	UserContextKey        string
	UserIdParamKey        string
	CSRFContextKey        string
	CSRFHeaderName        string
	CSRFCookieName        string
	SessionExpirationTime time.Duration
	UserExpirationTime    time.Duration
}

type Config struct {
	App      appConfig      `mapstructure:"app"`
	Database databaseConfig `mapstructure:"database"`
	JWT      jwtConfig      `mapstructure:"jwt"`
	Redis    redisConfig    `mapstructure:"redis"`
	Consts   consts
}

var config = &Config{
	Consts: consts{
		SessionCookieName:     "__Host_session_",
		CSRFCookieName:        "__Host_csrf_",
		SessionExpirationTime: time.Hour * 1,
		UserExpirationTime:    time.Minute * 30,
		JWTCookieName:         "jwt_token",
		JWTHeaderName:         "Authorization",
		UserIdParamKey:        "userId",
		UserContextKey:        "user",
		CSRFContextKey:        "csrf",
		CSRFHeaderName:        "X-Csrf-Token",
	},
}

func InitConfig(test ...bool) (*Config, error) {
	if len(test) >= 1 && test[0] {
		return config, nil
	}

	v := viper.New()

	// v := viper.NewWithOptions(
	// 	viper.WithLogger(
	// 		slog.New(
	// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 				Level: slog.LevelDebug,
	// 			}))))
	//

	v.AddConfigPath(".") // lokasi file konfigurasi
	v.AddConfigPath("config")
	v.SetConfigType("toml")   // tipe/format file konfigurasi
	v.SetConfigName("config") // nama file konfigurasi tanpa nama ekstensi
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = v.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
