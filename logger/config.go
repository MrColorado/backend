package logger

type Configuration struct {
	AppName string
	Logger  *LoggerConfiguration
}

type LoggerConfiguration struct {
	DevLogs    bool
	StackTrace bool
}
