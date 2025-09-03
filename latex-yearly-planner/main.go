package main

import (
	"os"
)

func main() {
	app := New()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
