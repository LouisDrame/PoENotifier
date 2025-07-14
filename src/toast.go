package main

import (
	"github.com/gen2brain/beeep"
)

func showToast(title, message string) error {

	beeep.AppName = "PoE Notifier"
	if err := beeep.Notify(title, message, ""); err != nil {
		return err
	}
	return nil
}
