package main

import (
	"lim/app"
)

func main() {
	server, err := app.New()
	if err != nil {
		panic(err)
	}

	if err = server.Run(); err != nil {
		panic(err)
	}
}
