package main

import (
	"os"

	"github.com/kudrykv/latex-yearly-planner/internal/app"
)

func main() {
	application := app.New()
	if err := application.Run(os.Args); err != nil {
		panic(err)
	}
}
