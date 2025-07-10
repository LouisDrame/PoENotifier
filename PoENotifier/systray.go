package main

import (
	_ "embed" // for embedding icon data
	"github.com/getlantern/systray"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed icons/icon.ico
var IconData []byte

func initSystray() {
	// Initialize systray
	systray.Run(onReady, onExit)
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
				cmd := exec.Command("explorer", configDir) // Change to your preferred editor if needed
				if err := cmd.Run(); err != nil {
				}
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
