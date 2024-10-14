// Wrapper functions for slog functions
package logging

import (
	"fmt"
	"log/slog"
)

func LogDebug(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

func LogInfo(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

func LogError(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
}
