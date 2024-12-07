package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"playground/internal/entities"
	"playground/internal/entities/constant"
	"playground/internal/external/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AuthorizationBearerToken(tokenName string, tokenService entities.TokenService) HandlerFunc {
	if tokenService == nil {
		log.Fatal("token service is nil")
	}
	if tokenName == "" {
		log.Fatal("token name is empty")
	}
	return func(c Context) {
		token := strings.TrimPrefix(
			c.GetHeader(constant.HttpHeaderAuthorization),
			constant.HttpPrefixAuthorizationBearer,
		)
		if ok, _ := tokenService.Authorize(tokenName, token); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func RequestLimiter() HandlerFunc {
	limiter := rate.NewLimiter(2, 1)
	return func(c Context) {
		if limiter.Allow() {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}

func DefaultLogger() HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Output: os.Stdout,
	})
}

type LoggerConfig struct {
	Output       io.Writer
	SkipPaths    []string
	TimeLocation *time.Location
}

func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Output: gin.DefaultWriter,
	}
}

type logDetails struct {
	Method     string
	Path       string
	Proto      string
	StatusCode int
	Latency    time.Duration
	ClientIP   string
}

func (l *logDetails) String() string {
	return fmt.Sprintf("%s %s %s %d %s %s",
		l.Method,
		l.Path,
		l.Proto,
		l.StatusCode,
		l.Latency,
		l.ClientIP,
	)
}

type Logger struct {
	location *time.Location
}

func (l *Logger) TimeStamp(t time.Time) (timestamp string) {
	const FormatTimeStamp = "2006-01-02 15:04:05"
	return t.In(l.location).Format(FormatTimeStamp)
}

func (l *Logger) LogLevel(statusCode int) (level string) {
	if statusCode >= http.StatusBadRequest {
		return entities.WarnLevel.String()
	}
	if statusCode >= http.StatusInternalServerError {
		return entities.ErrorLevel.String()
	}
	return entities.InfoLevel.String()
}

func NewLogger(config *LoggerConfig) HandlerFunc {
	if config == nil {
		config = DefaultLoggerConfig()
	}
	var location *time.Location
	if config.TimeLocation == nil {
		location = time.Local
	}
	logger := &Logger{
		location: location,
	}
	if config.Output == nil {
		config.Output = gin.DefaultWriter
	}
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			details := logDetails{
				Method:     param.Method,
				Path:       param.Path,
				Proto:      param.Request.Proto,
				StatusCode: param.StatusCode,
				Latency:    param.Latency,
				ClientIP:   param.ClientIP,
			}
			if param.Latency > time.Minute {
				details.Latency = param.Latency.Truncate(time.Second)
			}
			return fmt.Sprintf("%s %s <web.logger> [web] \"%s\"\n",
				logger.TimeStamp(param.TimeStamp),
				logger.LogLevel(param.StatusCode),
				details.String(),
			)
		},
		Output:    config.Output,
		SkipPaths: config.SkipPaths,
	})
}
