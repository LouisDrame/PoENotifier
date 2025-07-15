package main

import (
	_ "embed" // for embedding icon data
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/getlantern/systray"
)

//go:embed icons/icon.ico
var IconData []byte

func initSystray(logger *log.Logger) {
	// Initialize systray
	go func() {
		systray.Run(func() { onReady(logger) }, onExit)
	}()
}

func onReady(logger *log.Logger) {
	// Create the systray icon and menu
	systray.SetIcon(IconData)
	systray.SetTitle("PoE Notifier")
	systray.SetTooltip("Path of Exile Notifier")

	// Add menu items
	restartItem := systray.AddMenuItem("Restart", "Restart the application")
	go func() {
		for {
			select {
			case <-restartItem.ClickedCh:
				restartApplication(logger)
				return
			}
		}
	}()

	systray.AddSeparator()

	openConfigItem := systray.AddMenuItem("Open Config", "Open the configuration directory")
	go func() {
		for {
			select {
			case <-openConfigItem.ClickedCh:
				// Open the configuration file in the default editor
				configPath, err := getConfigPath()
				// Remove the filename from the path to open the directory
				configDir := filepath.Dir(configPath)
				if err != nil {
					return
				}
				handleOpenConfig(configDir, logger)
			}
		}
	}()

	systray.AddSeparator()

	// Add a quit item to the systray
	quitItem := systray.AddMenuItem("Quit", "Quit the application")
	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				logger.Println("Quitting application...")
				systray.Quit()
				os.Exit(0)
				return
			}
		}
	}()
}

func onExit() {
	// Cleanup code when the systray is exited
	// This can include closing log files, stopping goroutines, etc.
	// Nothing happens here for now.
	// TODO : If custom sounds are a thing someday, we should cleanup the sound resources.
}

func restartApplication(logger *log.Logger) {
	logger.Println("Restarting application...")

	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		logger.Printf("Error getting executable path: %v", err)
		return
	}

	// Start a new instance of the application
	cmd := exec.Command(executable)
	cmd.Dir = filepath.Dir(executable)

	// Start the new process
	if err := cmd.Start(); err != nil {
		logger.Printf("Error starting new instance: %v", err)
		return
	}

	logger.Println("New instance started, exiting current process...")
	// Exit the current process
	systray.Quit()
	os.Exit(0)
}

func handleOpenConfig(configDir string, logger *log.Logger) {
	logger.Printf("Opening config directory: %s", configDir)

	// Depending on the OS, open the config directory
	switch runtime.GOOS {
	case "windows":
		// For Windows, use explorer to open the config directory
		cmd := exec.Command("explorer", filepath.Join(configDir, "Notifier"))
		if err := cmd.Run(); err != nil {
			return
		}
	case "linux":
		// For Linux, use xdg-open to open the config directory
		cmd := exec.Command("xdg-open", configDir)
		if err := cmd.Run(); err != nil {
			logger.Printf("Error opening config directory on Linux: %v", err)
		}
	case "darwin":
		// For macOS, use open to open the config directory
		// Not used yet as macOS is not supported yet
		cmd := exec.Command("open", configDir)
		if err := cmd.Run(); err != nil {
			logger.Printf("Error opening config directory on macOS: %v", err)
		}
	}
}
