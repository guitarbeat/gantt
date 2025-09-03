package main

import (
	"fmt"
	"os"

	"github.com/kudrykv/latex-yearly-planner/internal/app"
)

func main() {
	application := app.New()
	if err := application.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
