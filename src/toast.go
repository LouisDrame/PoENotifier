package main

import (
	"log"

	"github.com/gen2brain/beeep"
)

func showToast(title, message string, logger *log.Logger) error {

	beeep.AppName = "PoE Notifier"
	if err := beeep.Notify(title, message, ""); err != nil {
		logger.Printf("Error showing toast notification: %v", err)
		return err
	}
	return nil
}
