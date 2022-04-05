package logr

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02 15:04:05"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

func CustomLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s := fmt.Sprintf("%5s", l.CapitalString())
	enc.AppendString(s)
}

func consoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    CustomLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeName: func(s string, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(fmt.Sprintf("[%-10v]", s))
		},
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(fmt.Sprintf("%-34v", caller.TrimmedPath()))
		},
		ConsoleSeparator: " ",
	}
}
