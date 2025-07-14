package main

import (
	_ "embed" // for embedding icon data
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/getlantern/systray"
)

//go:embed icons/icon.ico
var IconData []byte

func initSystray() {
	// Initialize systray
	go func() {
		systray.Run(onReady, onExit)
	}()
}

func onReady() {
	// Create the systray icon and menu
	systray.SetIcon(IconData)
	systray.SetTitle("PoE Notifier")
	systray.SetTooltip("Path of Exile Notifier")

	// Add a quit item to the systray
	quitItem := systray.AddMenuItem("Quit", "Quit the application")
	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				systray.Quit()
				os.Exit(0)
				return
			}
		}
	}()

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
				handleOpenConfig(configDir)
			}
		}
	}()
}

func handleOpenConfig(confifgDir string) {
	// Depending on the OS, open the config directory
	switch runtime.GOOS {
	case "windows":
		// For Windows, use explorer to open the config directory
		cmd := exec.Command("explorer", confifgDir)
		if err := cmd.Run(); err != nil {
		}
	case "linux":
		// For Linux, use xdg-open to open the config directory
		cmd := exec.Command("xdg-open", confifgDir)
		if err := cmd.Run(); err != nil {
		}
	case "darwin":
		// For macOS, use open to open the config directory
		// Not used yet as macOS is not supported yet
		cmd := exec.Command("open", confifgDir)
		if err := cmd.Run(); err != nil {
		}
	}
}

func onExit() {
	// Cleanup code when the systray is exited
	// This can include closing log files, stopping goroutines, etc.
	// Nothing happens here for now.
	// TODO : If custom sounds are a thing someday, we should cleanup the sound resources.
}
