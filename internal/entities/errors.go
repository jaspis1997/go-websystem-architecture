package entities

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

var (
	// 初期化に失敗した。
	// プロセスは終了すべきである。
	ErrorInitializeFailed = fmt.Errorf("initialize failed")

	// 無効、またはサポートされないコンフィグが指定された。
	// プロセスは終了すべきである。
	ErrorUnsupportedConfig = fmt.Errorf("unsupported config")
)

const (
	// リポジトリのマイグレーションに失敗した。
	// プロセスは終了すべきである。
	FormatErrorFailedMigrate = "failed to migrate on (%s=resource name)"

	// リポジトリのクローズに失敗した。
	FormatErrorFailedClose = "failed to close on (%s=resource name)"
)

type Error = *StackTraceableError

// type Recordable struct {
// 	err error
// }

// func (e Recordable) Error() string {
// 	return fmt.Sprintf(`"%s"`, e.err.Error())
// }

func newTraceId() string {
	return uuid.New().String()
}

type traceDetails struct {
	error
	caller string
}

type StackTraceableError struct {
	TraceId string
	errors  []traceDetails
}

func NewStackTrace(tid ...string) *StackTraceableError {
	var id string
	if len(tid) > 0 {
		id = tid[0]
	}
	return &StackTraceableError{TraceId: id}
}

func (e *StackTraceableError) SetIdAutomatically() *StackTraceableError {
	return &StackTraceableError{TraceId: newTraceId()}
}

func (e *StackTraceableError) Add(err error) *StackTraceableError {
	if err == nil {
		return e
	}
	_, file, line, _ := runtime.Caller(1)
	dir, file := path.Split(file)
	caller := fmt.Sprintf("%s:%d", path.Join(path.Base(dir), file), line)
	e.errors = append([]traceDetails{
		{error: err, caller: caller}}, e.errors...)
	return e
}

func (e *StackTraceableError) Addf(format string, args ...interface{}) *StackTraceableError {
	return e.Add(fmt.Errorf(format, args...))
}

// func (e *StackTraceableError) Recordable() Recordable {
// 	return Recordable{err: e}
// }

func (e *StackTraceableError) Error() string {
	var err error
	if len(e.errors) > 0 {
		err = e.errors[0]
	}
	return fmt.Sprintf("trace: %s, error: %v, more: %d",
		e.TraceId, err, len(e.errors)-1)
}

func (e *StackTraceableError) Is(target error) bool {
	for _, err := range e.errors {
		if err.error == target {
			return true
		}
	}
	return false
}

func (e *StackTraceableError) Unwrap() error {
	if len(e.errors) > 0 {
		return e.errors[0].error
	}
	return nil
}

func (e *StackTraceableError) StackTrace() string {
	const format = "#%d: <%s> %s"
	var stackTrace strings.Builder
	for index, err := range e.errors {
		fmt.Fprintf(&stackTrace, format+"\n", index, err.caller, err.Error())
	}
	return stackTrace.String()
}
