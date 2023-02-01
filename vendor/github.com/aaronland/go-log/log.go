package log

import (
	"fmt"
	"log"
)

// DEBUG_LEVEL is the numeric level associated with "debug" messages.
const DEBUG_LEVEL int = 10

// INFO_LEVEL is the numeric level associated with "info" messages.
const INFO_LEVEL int = 20

// WARNING_LEVEL is the numeric level associated with "warning" messages.
const WARNING_LEVEL int = 30

// ERROR_LEVEL is the numeric level associated with "error" messages.
const ERROR_LEVEL int = 40

// FATAL_LEVEL is the numeric level associated with "fatal" messages.
const FATAL_LEVEL int = 50

// DEBUG_PREFIX is the string prefix to prepend "debug" messages with.
const DEBUG_PREFIX string = "🪵" // "🦜"

// INFO_PREFIX is the string prefix to prepend "info" messages with.
const INFO_PREFIX string = "💬"

// WARNING_PREFIX is the string prefix to prepend "warning" messages with.
const WARNING_PREFIX string = "🧯"

// ERROR_PREFIX is the string prefix to prepend "error" messages with.
const ERROR_PREFIX string = "🔥"

// FATAL_PREFIX is the string prefix to prepend "fatal" messages with.
const FATAL_PREFIX string = "💥"

var minLevel int

// UnsetMinLevel resets the minimum level log level allowing all messages to be emitted.
func UnsetMinLevel() {
	minLevel = 0
}

// SetMinLevel assigns the minimum log level to 'l'. Log events with a lower level will not be emitted.
func SetMinLevel(l int) error {

	switch l {
	case DEBUG_LEVEL, INFO_LEVEL, WARNING_LEVEL, ERROR_LEVEL, FATAL_LEVEL:
		minLevel = l
	default:
		return fmt.Errorf("Invalid level")
	}

	return nil
}

// SetMinLevelStringWithPrefix assigns the minimum log level associated with 'prefix'. Log events with a lower level will not be emitted.
func SetMinLevelWithPrefix(prefix string) error {

	switch prefix {
	case DEBUG_PREFIX:
		return SetMinLevel(DEBUG_LEVEL)
	case INFO_PREFIX:
		return SetMinLevel(INFO_LEVEL)
	case WARNING_PREFIX:
		return SetMinLevel(WARNING_LEVEL)
	case ERROR_PREFIX:
		return SetMinLevel(ERROR_LEVEL)
	case FATAL_PREFIX:
		return SetMinLevel(FATAL_LEVEL)
	default:
		return fmt.Errorf("Invalid level")
	}
}

// Emit a "debug" log message.
func Debug(logger *log.Logger, msg string, args ...interface{}) {

	if minLevel > DEBUG_LEVEL {
		return
	}

	msg = fmt.Sprintf("%s %s", DEBUG_PREFIX, msg)
	logger.Printf(msg, args...)
}

// Emit a "info" log message.
func Info(logger *log.Logger, msg string, args ...interface{}) {

	if minLevel > INFO_LEVEL {
		return
	}

	msg = fmt.Sprintf("%s %s", INFO_PREFIX, msg)
	logger.Printf(msg, args...)
}

// Emit a "warning" log message.
func Warning(logger *log.Logger, msg string, args ...interface{}) {

	if minLevel > WARNING_LEVEL {
		return
	}

	msg = fmt.Sprintf("%s %s", WARNING_PREFIX, msg)
	logger.Printf(msg, args...)
}

// Emit an "error" log message.
func Error(logger *log.Logger, msg string, args ...interface{}) {

	if minLevel > ERROR_LEVEL {
		return
	}

	msg = fmt.Sprintf("%s %s", ERROR_PREFIX, msg)
	logger.Printf(msg, args...)
}

// Emit an "fatal" log message.
func Fatal(logger *log.Logger, msg string, args ...interface{}) {

	if minLevel > FATAL_LEVEL {
		return
	}

	msg = fmt.Sprintf("%s %s", FATAL_PREFIX, msg)
	logger.Fatalf(msg, args...)
}
