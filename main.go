package main

import (
	"context"
	"fmt"
	"os"

	"github.com/d6o/portsservice/internal/app"
)

const (
	errorExitCode = 1
)

func main() {
	ctx := context.Background()

	cmd := app.NewApp()
	if err := cmd.Run(ctx); err != nil {
		fmt.Printf("Can't execute service. Error: %s\n", err)
		os.Exit(errorExitCode)
	}
}
