package lojer

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func (l Logger) Printf(f string, a ...any) {
	l.Info(fmt.Sprintf(f, a...))
}

func NewLogger(dev ...bool) Logger {
	var (
		encoder zapcore.Encoder
		config  zapcore.EncoderConfig
	)

	if len(dev) >= 1 && dev[0] {
		config = zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(config)
	} else {
		config = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(config)
	}

	ws := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, ws, zap.InfoLevel)
	z := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return Logger{z}
}

func NewFromZap(z *zap.Logger) Logger {
	return Logger{z}
}
