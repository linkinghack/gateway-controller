package log

import (
	"github.com/linkinghack/gateway-controller/config"
	"github.com/sirupsen/logrus"
)

var globalLogger *logrus.Logger

// InitGlobalLogger 使用全局配置中的日志配置初始化全局logger
func InitGlobalLogger() {
	globalLogger = logrus.New()
	logConf := config.GetGlobalConfig().LogConfig

	// set log level
	switch logConf.Level {
	case config.LogLevelDebug:
		globalLogger.SetLevel(logrus.DebugLevel)
		break
	case config.LogLevelInfo:
		globalLogger.SetLevel(logrus.InfoLevel)
		break
	case config.LogLevelWarn:
		globalLogger.SetLevel(logrus.WarnLevel)
		break
	case config.LogLevelTrace:
		globalLogger.SetLevel(logrus.TraceLevel)
		break
	default:
		globalLogger.SetLevel(logrus.DebugLevel)
	}

	switch logConf.LogFormat {
	case config.LogFormatText:
		globalLogger.SetFormatter(&logrus.TextFormatter{
			DisableColors: !logConf.Color,
			FullTimestamp: true,
		})
		break
	case config.LogFormatJson:
		globalLogger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: false,
		})
	default:
		globalLogger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
}

func GetGlobalLogger() *logrus.Logger {
	if globalLogger == nil {
		InitGlobalLogger()
	}
	return globalLogger
}

// GetSpecificLogger  从全局logger上创建一个子logger并指定一些固定参数
// 参数：
//
//	funcName: 函数名称。在log中将增加FuncName 字段
//	otherFields: 其他参数列表。数组长度应该为偶数，key,value,key,value....
func GetSpecificLogger(funcName string, otherFields ...string) *logrus.Entry {
	newLogger := globalLogger.WithField("FuncName", funcName)
	tail := len(otherFields)
	if len(otherFields)%2 != 0 {
		tail = len(otherFields) - 1
	}
	if tail >= 2 {
		additionalFields := logrus.Fields{}
		for i := 0; i <= tail/2; i += 2 {
			additionalFields[otherFields[i]] = otherFields[i+1]
		}
		newLogger = newLogger.WithFields(additionalFields)
	}

	return newLogger
}
