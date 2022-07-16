package Go_Logger_MDC

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
	"time"
	"unsafe"
)

func SugarLog(ctx context.Context) *zap.SugaredLogger {

	logger := ConfigZap()

	loggerID := logContextInternals(logger, ctx)

	return loggerID

}

func SugarLogWithoutContext() *zap.SugaredLogger {
	logger := ConfigZap()
	return logger
}

func logContextInternals(logger *zap.SugaredLogger, ctx interface{}) *zap.SugaredLogger {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	var s []any

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				logger = logContextInternals(logger, reflectValue.Interface())
			} else {
				if reflect.TypeOf(reflectValue.Interface()) != nil {
					if reflect.TypeOf(reflectValue.Interface()).Kind() == reflect.String {
						s = append(s, reflectValue.Interface())
					} else if len(s) == 1 && reflect.TypeOf(s[0]).Kind() == reflect.String {
						s = append(s, reflectValue.Interface())
					}
				}
			}
		}
	} else {
	}

	for i := 0; i < len(s); i = i + 2 {
		logger = logger.With(s[i], s[i+1])
	}

	return logger
}

func ConfigZap() *zap.SugaredLogger {

	cfg := zap.Config{
		Encoding:    "console",                           //encode kiểu json hoặc console
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel), //chọn InfoLevel có thể log ở cả 3 level
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{ //Cấu hình logging, sẽ không có stacktracekey
			MessageKey:   "message",
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder, //Lấy dòng code bắt đầu log
			EncodeLevel:  CustomLevelEncoder,        //Format cách hiển thị level log
			EncodeTime:   SyslogTimeEncoder,         //Format hiển thị thời điểm log
		},
	}

	logger, _ := cfg.Build() //Build ra Logger
	return logger.Sugar()    //Trả về logger hoặc Sugaredlogger, ở đây ta chọn trả về Logger
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", // Log name
		MaxSize:    1,            // File content size, MB
		MaxBackups: 5,            // Maximum number of old files retained
		MaxAge:     30,           // Maximum number of days to keep old files
		Compress:   false,        // Is the file compressed
	}
	return zapcore.AddSync(lumberJackLogger)
}
