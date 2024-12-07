package constant

const (
	LogForceShutdown             = "force shutdown"
	LogShutdownCompleted         = "graceful shutdown completed"
	LogLogicsProcessingCompleted = "logic processing completed"
)

const (
	FormatLogStartShutdown  = "starting shutdown process, timeout: %s"
	FormatLogFailedShutdown = "failed to shutdown: %+v"
)
