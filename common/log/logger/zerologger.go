package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/khorsl/miniotaur/common/constants"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	logger *zerolog.Logger
	ctx    context.Context
}

func NewZeroLogger(loggerType string, ctx context.Context) *ZeroLogger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &ZeroLogger{
		logger: &logger,
		ctx:    ctx,
	}
}

func (z *ZeroLogger) Debug(msg string, fields map[string]interface{}) {
	z.updateContext(fields)

	z.logger.Debug().Msg(msg)

}

func (z *ZeroLogger) Info(msg string, fields map[string]interface{}) {
	z.updateContext(fields)

	z.logger.Info().Msg(msg)
}

func (z *ZeroLogger) Warn(msg string, fields map[string]interface{}) {
	z.updateContext(fields)

	z.logger.Warn().Msg(msg)
}

func (z *ZeroLogger) Error(msg string, fields map[string]interface{}) {
	z.updateContext(fields)

	z.logger.Error().Msg(msg)
}

func (z *ZeroLogger) Fatal(msg string, fields map[string]interface{}) {
	z.updateContext(fields)

	z.logger.Fatal().Msg(msg)
}

func (z *ZeroLogger) addContextCommonFields(fields map[string]interface{}) {
	if z.ctx != nil {
		commonFields, ok := z.ctx.Value(constants.LoggerCommonFields).(map[string]interface{})
		if !ok {
			return
		}

		for k, v := range commonFields {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}
}

func (z *ZeroLogger) updateContext(fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	z.addContextCommonFields(fields)

	for k, v := range fields {
		z.logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			switch v := v.(type) {
			case string:
				return c.Str(k, v)
			case int:
				return c.Int(k, v)
			case bool:
				return c.Bool(k, v)
			case error:
				return c.Err(v)
			default:
				return c.Str(k, fmt.Sprintf("%+v", v))
			}
		})
	}
}
