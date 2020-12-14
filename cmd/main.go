package main

import (
	"context"
	"errors"
	"os"

	"github.com/fatih/color"
	Logger "github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/escli/cmd/app"
)

func main() {
	if err := app.Run(os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, context.Canceled) {
			Logger.Debugln("ignore error since context is cancelled:", err)
		} else {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
}
