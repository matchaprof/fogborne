package logging

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/matchaprof/fogborne/internal/core/config"
	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger

//todo rethink transaction tracking
// var requestCounter atomic.Uint64

// ContextualEntry assists in tracing related logs across a single transaction
type SessionContext struct {
	SessionID  string    // Unique ID for a player's session
	ActionID   string    // ID for a specific game action
	PlayerID   string    // Tracking the player involved
	ActionType string    // Type of action being performed
	StartTime  time.Time // When this context was created
}

func (sc *SessionContext) String() string {
	duration := time.Since(sc.StartTime).Round(time.Millisecond)

	var parts []string

	// Only include fields that are set
	if sc.SessionID != "" {
		parts = append(parts, fmt.Sprintf("sid:%s", sc.SessionID))
	}
	if sc.PlayerID != "" {
		parts = append(parts, fmt.Sprintf("pid:%s", sc.PlayerID))
	}
	if sc.ActionType != "" {
		parts = append(parts, fmt.Sprintf("type:%s", sc.ActionType))
	}

	// Always include duration last
	parts = append(parts, fmt.Sprintf("dur:%v", duration))

	return strings.Join(parts, " ¤ ")
}

// CustomFormatter defines how we want the logs to look
type CustomFormatter struct {
	TimestampFormat string // sets Timestampformat to use Go's time formatting
	ShowFullPath    bool   // This setting determines whether the logger shows the full filepath or just filename
	// startTime       time.Time // This setting will help in tracking performance issues
	ColorizeContext bool // Option to colorize the context differently
}

const (
	logSeparator  = "►►►"
	fileInfoWidth = 22
	colorReset    = "\033[0m"
)

// InitLogger initializes the logging system with the provided configuration
func InitLogger(cfg *config.LoggingConfig) error {
	Logger = logrus.New()

	customFormatter := &CustomFormatter{
		TimestampFormat: "2006-01-02 ¤ 15:04:05",
		ShowFullPath:    false,
	}

	// Set the formatter
	Logger.SetFormatter(customFormatter)

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set logging level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	Logger.SetLevel(level)

	// Sets report caller information based on config
	Logger.SetReportCaller(cfg.ReportCaller)

	return nil
}

// Format implements the logrus.Formatter interface
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var fileInfo string
	if entry.HasCaller() {
		// If ShowFullPath is false, just use the last two parts of the path
		// (e.g., "logging/logging.go" instead of full path to logger.go)
		if !f.ShowFullPath {
			parts := strings.Split(entry.Caller.File, "/")
			if len(parts) > 2 {
				fileInfo = strings.Join(parts[len(parts)-2:], "/")
			} else {
				fileInfo = entry.Caller.File
			}
		} else {
			fileInfo = entry.Caller.File
		}

		fileInfo = fmt.Sprintf("%s:%d", fileInfo, entry.Caller.Line)

		// Pad or truncate to consistent width
		if len(fileInfo) > fileInfoWidth {
			parts := strings.Split(fileInfo, ":")
			filePath := parts[0]
			lineNum := parts[1]

			pathParts := strings.Split(filePath, "/")

			importantParts := []string{pathParts[len(pathParts)-1]}

			if len(pathParts) > 1 {
				importantParts = append([]string{pathParts[len(pathParts)-2]}, importantParts...)
			}

			truncatedPath := strings.Join(importantParts, "/")

			maxPathWidth := fileInfoWidth - 3 - len(lineNum) - 1

			if len(truncatedPath) > maxPathWidth {
				truncatedPath = truncatedPath[len(truncatedPath)-maxPathWidth:]
			}

			fileInfo = fmt.Sprintf("…%s:%s", truncatedPath, lineNum)
		} else {
			// If too short, pad with spaces to reach consistent width
			fileInfo = fmt.Sprintf("%-*s", fileInfoWidth, fileInfo)
		}
	}

	// Get the color code based on log level
	var colorCode string
	switch entry.Level {
	case logrus.DebugLevel:
		colorCode = "\033[36m" // Cyan
	case logrus.InfoLevel:
		colorCode = "\033[32m" // Green
	case logrus.WarnLevel:
		colorCode = "\033[33m" // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		colorCode = "\033[31m" // Red
	}

	// Format the timestamp
	timestamp := entry.Time.Format(f.TimestampFormat)

	// Build the log message with custom formatting
	var sb strings.Builder

	// Add timestamp in brackets
	sb.WriteString(fmt.Sprintf("[%s] ", timestamp))

	// Add log level with padding and color
	level := strings.ToUpper(entry.Level.String())
	if entry.Level == logrus.InfoLevel {
		level = fmt.Sprintf("[ %-4s  ]", level)
	} else if entry.Level == logrus.WarnLevel {
		level = fmt.Sprintf("[%s]", level)
	} else {
		level = fmt.Sprintf("[ %-6s]", level) // pads level to 6 characters
	}

	// Add the colored level
	sb.WriteString(colorCode + level + colorReset)

	// Add colored separator and file info
	if fileInfo != "" {
		// Add separator with same color as log level
		sb.WriteString(fmt.Sprintf(" %s%s%s %s %s%s%s ",
			colorCode, logSeparator, colorReset,
			fileInfo,
			colorCode, logSeparator, colorReset))
	}

	//TODO redo this useless request counter
	// reqID := requestCounter.Add(1) // Thread-safe counter
	// duration := time.Since(f.startTime)
	// contextOutput := fmt.Sprintf("req:%d", reqID)

	// contextInfo := fmt.Sprintf(" (%s%s%s) ",
	// 	colorCode,
	// 	contextOutput,
	// 	colorReset)

	// sb.WriteString(contextInfo)

	// Add the log message
	sb.WriteString(entry.Message)

	// Add any fields as key=value pairs
	if len(entry.Data) > 0 {
		sb.WriteString(fmt.Sprintf(" %s►%s [", colorCode, colorReset))

		// Retrieving the keys and sorting alphabetically
		keys := make([]string, 0, len(entry.Data))
		for k := range entry.Data {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		first := true
		for _, k := range keys {
			if !first {
				// Seperator between pairs
				sb.WriteString(fmt.Sprintf(" %s¤ ", colorReset))
			}

			v := entry.Data[k]
			valueColor := getValueColor(v)
			formattedValue := formatValue(v)

			sb.WriteString(fmt.Sprintf("%s%s%s⇒%s%s%s",
				colorCode, k, colorCode, // Colored key
				valueColor,     // Start value color
				formattedValue, // Formatted value
				colorReset))    // Reset color
			first = false
		}

		sb.WriteString("]")
	}

	sb.WriteString("\n")
	return []byte(sb.String()), nil
}

/*
-----------------------------------------
-------- FORMATTING FUNCTIONS -----------
-----------------------------------------
*/

// getValueColor chooses colors based on value type
func getValueColor(value interface{}) string {
	switch value.(type) {
	case int, int32, int64, float32, float64:
		return "\033[36m" // Cyan for numbers
	case bool:
		if true {
			return "\033[32m" // Green for true
		}
		return "\033[31m" // Red for false
	case string:
		return "\033[35m" // Magenta for strings
	default:
		return "\033[33m" // Yellow for other types
	}
}

// formatValue chooses how to format values based on their type
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case time.Duration:
		// Formats time in a readable way (e.g., "1h45m")
		return v.String()
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case []string:
		return fmt.Sprintf("[ %s ]", strings.Join(v, ", "))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// LogFunc is a type that represents any logging function that takes a string message
type LogFunc func(args ...interface{})

// LogfFunc represents a formatted logging function (like Infof, Debugf)
type LogfFunc func(format string, args ...interface{})

// LogTitle creates a visual separation in logs for titles
// TODO Center the title along a width line of 80
func LogTitle(title string, logFn LogFunc) {
	if Logger == nil {
		return
	}

	logFn(fmt.Sprintf("          .•( %s )•. ", title))
}

// LogSection creates a visual separation in logs for different sections
func LogSection(title string, logFn LogFunc) {
	if Logger == nil {
		return
	}

	width := 80
	line := strings.Repeat("—", width)
	logFn(fmt.Sprint(line))
	logFn(fmt.Sprintf("                   .•( %s )•.", title))
	logFn(fmt.Sprint(line))
}

// LogSubSection creates a smaller visual separation for subsections
func LogSubSection(title string, logFn LogFunc) {
	if Logger == nil {
		return
	}

	logFn(fmt.Sprintf(" .•( %s )•. ", title))
}

/*
-----------------------------------------
----------- HELPER FUNCTIONS ------------
-----------------------------------------
*/

// Helper functions for the different log levels with structured data
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

func WithCorrelationID(id string) *logrus.Entry {
	return Logger.WithField("correlation_id", id)
}

// // StartPlayerSession creates a new session context
// func StartPlayerSession(playerID string) *SessionContext {
// 	return &SessionContext{
// 		SessionID: uuid.New().String(), // Generate unique session ID
// 		PlayerID:  playerID,
// 		StartTime: time.Now(),
// 	}
// }

// // StartGameAction adds action context to an existing session
// func (sc *SessionContext) StartGameAction(actionType string) *SessionContext {
// 	return &SessionContext{
// 		SessionID:  sc.SessionID,        // Keep the same session ID
// 		ActionID:   uuid.New().String(), // New action ID
// 		StartTime:  time.Now(),
// 		PlayerID:   sc.PlayerID,
// 		ActionType: actionType,
// 	}
// }

// Methods that mirror logrus's logging levels
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}
