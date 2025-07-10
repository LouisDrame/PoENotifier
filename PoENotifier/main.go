package main

import (
	"context"
	"fmt"
	"github.com/go-faster/tail"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// Setup logging
	logger := setupLogging()
	logger.Println("Starting PoE Notifier...")

	// Will check if the config file exists, if not it will create it with the default config
	// User can edit the config file to change the patterns to match
	checkConfig()

	// TODO : support other OS than Windows
	// TODO : add path to config file to support other installations
	// For now, we will use the default Path of Exile installation path on Windows
	t := tail.File("C:\\Program Files (x86)\\Grinding Gear Games\\Path of Exile\\logs\\Client.txt", tail.Config{
		Follow:     true,       // tail -f
		BufferSize: 1024 * 128, // 128 kb for internal reader buffer

		NotifyTimeout: time.Duration(1 * time.Second),

		Location: &tail.Location{Whence: io.SeekEnd, Offset: 0},
	})
	ctx := context.Background()

	config, err := importConfig()
	if err != nil {
		logger.Printf("Error importing config: %v", err)
		return
	}
	logger.Println("Config checked and loaded")

	logger.Printf("Config imported successfully. Found %d patterns:", len(config.Patterns))
	for _, pattern := range config.Patterns {
		logger.Printf("  - Pattern: %s, Regex: %s", pattern.Name, pattern.Regex)
	}

	logger.Println("Starting to tail PoE log file...")

	if err := t.Tail(ctx, func(ctx context.Context, l *tail.Line) error {
		if matched, pattern := checkPattern(string(l.Data), config.Patterns, logger); matched {
			logger.Printf("PATTERN MATCHED: %s - Line: %s", pattern.Name, string(l.Data))
			if pattern.Toast {
				showToast(pattern.Name, pattern.Message)
			}
			if pattern.Beep {
				beep()
			}
		}
		return nil
	}); err != nil {
		logger.Printf("Fatal error in tail: %v", err)
		panic(err)
	}
}

func checkPattern(line string, patterns []Pattern, logger *log.Logger) (bool, Pattern) {
	for _, pattern := range patterns {
		cleanRegex := pattern.Regex
		if len(cleanRegex) > 2 && cleanRegex[0] == '`' && cleanRegex[len(cleanRegex)-1] == '`' {
			cleanRegex = cleanRegex[1 : len(cleanRegex)-1]
		}

		if matched, err := regexp.MatchString(cleanRegex, line); err != nil {
			logger.Printf("Regex error for pattern '%s': %v", pattern.Name, err)
		} else if matched {
			return true, pattern
		}
	}
	return false, Pattern{}
}

// setupLogging creates a logger that writes to both file and console
func setupLogging() *log.Logger {
	// Create logs directory if it doesn't exist
	user, _ := user.Current()
	logDir := path.Join(user.HomeDir, "Documents", "My Games", "Path of Exile", "Notifier", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return log.New(os.Stdout, "[PoENotifier] ", log.LstdFlags|log.Lshortfile)
	}

	// Create log file with timestamp
	logFileName := fmt.Sprintf("poe_notifier_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return log.New(os.Stdout, "[PoENotifier] ", log.LstdFlags|log.Lshortfile)
	}

	// Create a multi-writer to write to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Create logger with timestamp, filename and line number
	logger := log.New(multiWriter, "[PoENotifier] ", log.LstdFlags|log.Lshortfile)

	logger.Printf("Logging initialized. Log file: %s", logFilePath)
	return logger
}

// beep plays a system beep sound.
// Only supported on Windows at the moment.
func beep() {
	if runtime.GOOS != "windows" {
		// Windows beep
		fmt.Print("\a")
	} else {
		kernel32 := syscall.NewLazyDLL("user32.dll")
		kernel32.NewProc("MessageBeep").Call(880, 200)
	}
}
